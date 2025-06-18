-- lua/strudel-nvim/core.lua

local M = {}

local tcp_client = nil
local timer = nil
local backend_job_id = nil

local config = {
	host = "127.0.0.1",
	port = 8080,
	reconnect_delay = 3000, -- 3 seconds
	auto_reconnect = true, -- New option: automatically try to reconnect if backend disconnects
}

-- Simple XOR function for masking (replaces bit.bxor)
local function xor(a, b)
	local result = 0
	local bit_pos = 1
	while a > 0 or b > 0 do
		local bit_a = a % 2
		local bit_b = b % 2
		if bit_a ~= bit_b then
			result = result + bit_pos
		end
		a = math.floor(a / 2)
		b = math.floor(b / 2)
		bit_pos = bit_pos * 2
	end
	return result
end

-- Simple WebSocket frame creation (for text frames)
local function create_websocket_frame(data)
	local payload = data or ""
	local payload_len = #payload
	local frame = {}

	-- FIN=1, RSV=000, Opcode=0001 (text frame)
	frame[1] = string.char(0x81)

	-- Mask=1, Payload length
	if payload_len < 126 then
		frame[2] = string.char(0x80 + payload_len)
	elseif payload_len < 65536 then
		frame[2] = string.char(0x80 + 126)
		frame[3] = string.char(math.floor(payload_len / 256))
		frame[4] = string.char(payload_len % 256)
	else
		-- For larger payloads, we'd need 8-byte length
		vim.notify("Strudel: Payload too large", vim.log.levels.ERROR)
		return nil
	end

	-- Masking key (4 random bytes)
	local mask = {}
	for i = 1, 4 do
		mask[i] = math.random(0, 255)
		frame[#frame + 1] = string.char(mask[i])
	end

	-- Masked payload
	for i = 1, payload_len do
		local byte = string.byte(payload, i)
		local masked_byte = xor(byte, mask[((i - 1) % 4) + 1])
		frame[#frame + 1] = string.char(masked_byte)
	end

	return table.concat(frame)
end

-- Simple WebSocket handshake
local function create_websocket_handshake()
	local key = "dGhlIHNhbXBsZSBub25jZQ==" -- Base64 encoded key
	local handshake = {
		"GET /ws HTTP/1.1",
		"Host: " .. config.host .. ":" .. config.port,
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Key: " .. key,
		"Sec-WebSocket-Version: 13",
		"",
		"",
	}
	return table.concat(handshake, "\r\n")
end

-- Function to connect to the Go WebSocket server
local function connect_to_backend()
	if tcp_client then
		vim.notify("Strudel: Already connected to backend", vim.log.levels.WARN)
		return
	end

	tcp_client = vim.loop.new_tcp()
	if not tcp_client then
		vim.notify("Strudel: Failed to create TCP client", vim.log.levels.ERROR)
		return
	end

	-- Connect to the server
	tcp_client:connect(config.host, config.port, function(err)
		if err then
			vim.schedule(function()
				vim.notify("Strudel: Connection failed: " .. err, vim.log.levels.ERROR)
				tcp_client:close()
				tcp_client = nil
				-- Schedule reconnection
				schedule_reconnect()
			end)
			return
		end

		-- Send WebSocket handshake
		local handshake = create_websocket_handshake()
		tcp_client:write(handshake, function(write_err)
			if write_err then
				vim.schedule(function()
					vim.notify("Strudel: Handshake failed: " .. write_err, vim.log.levels.ERROR)
					tcp_client:close()
					tcp_client = nil
				end)
				return
			end

			vim.schedule(function()
				vim.notify("Strudel: Connected to backend!", vim.log.levels.INFO)
			end)
		end)

		-- Start reading from the server
		tcp_client:read_start(function(read_err, data)
			if read_err then
				vim.schedule(function()
					vim.notify("Strudel: Read error: " .. read_err, vim.log.levels.ERROR)
					if tcp_client then
						tcp_client:close()
						tcp_client = nil
					end
					schedule_reconnect()
				end)
				return
			end

			if data then
				vim.schedule(function()
					-- For now, just display the raw data
					-- In a full WebSocket implementation, you'd parse the frame here
					if data:find("HTTP/1.1 101") then
						-- Handshake response, ignore
						return
					end

					-- Simple frame parsing for text messages
					-- This is a simplified version - a full implementation would be more robust
					local msg = data:gsub("[\r\n]", "")
					if #msg > 0 and not msg:find("Sec%-WebSocket") then
						vim.notify("Strudel Log: " .. msg, vim.log.levels.INFO)
					end
				end)
			else
				-- Connection closed
				vim.schedule(function()
					vim.notify("Strudel: Connection closed by server", vim.log.levels.WARN)
					if tcp_client then
						tcp_client:close()
						tcp_client = nil
					end
					schedule_reconnect()
				end)
			end
		end)
	end)
end

-- Schedule reconnection attempt
local function schedule_reconnect()
	if not config.auto_reconnect then -- Check new auto_reconnect config
		vim.notify("Strudel: Auto-reconnect is disabled.", vim.log.levels.INFO)
		return
	end

	if timer then
		timer:close()
	end

	timer = vim.loop.new_timer()
	timer:start(
		config.reconnect_delay,
		0,
		vim.schedule_wrap(function()
			vim.notify("Strudel: Attempting to reconnect...", vim.log.levels.INFO)
			connect_to_backend()
			timer:close()
			timer = nil
		end)
	)
end

function M.start_backend()
	if backend_job_id then
		vim.notify("Strudel: Backend is already running.", vim.log.levels.INFO)
	else
		vim.notify("Strudel: Starting Go backend process...", vim.log.levels.INFO)
		-- Determine the path to the backend directory relative to the plugin file
		-- Use debug.getinfo to get the current file path more reliably
		local info = debug.getinfo(1, 'S')
		local script_path = info.source:match('^@(.*)$')
		local plugin_dir = vim.fn.fnamemodify(script_path, ':h:h:h')
		local backend_script_path = plugin_dir .. "/backend/start.sh"
		local backend_dir = plugin_dir .. "/backend"
		
		-- Backend paths determined successfully

		backend_job_id = vim.fn.jobstart({ backend_script_path }, {
			cwd = backend_dir, -- Set the current working directory for the job
			detach = true, -- Run the process in the background, detached from Neovim
			on_exit = function(job_id, data, event)
				vim.schedule(function()
					if backend_job_id == job_id then -- Ensure it's the job we started
						vim.notify("Strudel: Backend process exited. Code: " .. data, vim.log.levels.INFO)
						backend_job_id = nil
						-- Optionally, attempt to reconnect or notify if it was an unexpected exit
						if config.auto_reconnect then
							schedule_reconnect()
						end
					end
				end)
			end,
			on_stdout = vim.schedule_wrap(function(_, data, _)
				for _, line in ipairs(data) do
					vim.notify("Strudel Backend Output: " .. line, vim.log.levels.DEBUG)
				end
			end),
			on_stderr = vim.schedule_wrap(function(_, data, _)
				for _, line in ipairs(data) do
					vim.notify("Strudel Backend Error: " .. line, vim.log.levels.ERROR)
				end
			end),
		})

		if backend_job_id == 0 or backend_job_id == -1 then
			vim.notify(
				"Strudel: Failed to start backend process. Check if '"
					.. backend_script_path
					.. "' is executable and Go is installed.",
				vim.log.levels.ERROR
			)
			backend_job_id = nil
		else
			vim.notify(
				"Strudel: Backend process started (Job ID: " .. backend_job_id .. "). Attempting to connect...",
				vim.log.levels.INFO
			)
			-- Give a short moment for the backend to start before attempting connection
			vim.loop.new_timer():start(
				1000,
				0,
				vim.schedule_wrap(function()
					connect_to_backend()
				end)
			)
		end
	end
end

function M.stop_backend()
	if tcp_client then
		tcp_client:close()
		tcp_client = nil
		vim.notify("Strudel: Disconnected from backend.", vim.log.levels.INFO)
	else
		vim.notify("Strudel: Not connected.", vim.log.levels.WARN)
	end

	if timer then
		timer:close()
		timer = nil
	end

	if backend_job_id then
		vim.notify("Strudel: Stopping backend process (Job ID: " .. backend_job_id .. ")...", vim.log.levels.INFO)
		vim.fn.jobstop(backend_job_id) -- Terminate the background process
		backend_job_id = nil
		vim.notify("Strudel: Backend process stopped.", vim.log.levels.INFO)
	else
		vim.notify("Strudel: No backend process running to stop.", vim.log.levels.WARN)
	end
end

function M.eval_line()
	if not tcp_client then
		vim.notify("Strudel: Not connected to backend. Run :StrudelStart", vim.log.levels.ERROR)
		return
	end

	-- Get the content of the current line
	local line = vim.api.nvim_get_current_line()

	-- Create a WebSocket frame and send it
	local frame = create_websocket_frame(line)
	if not frame then
		vim.notify("Strudel: Failed to create WebSocket frame", vim.log.levels.ERROR)
		return
	end

	tcp_client:write(frame, function(err)
		if err then
			vim.schedule(function()
				vim.notify("Strudel: Failed to send code: " .. err, vim.log.levels.ERROR)
			end)
		else
			vim.schedule(function()
				vim.notify("Strudel: Evaluated line.", vim.log.levels.INFO)
			end)
		end
	end)
end

-- Function to configure the plugin
function M.setup(opts)
	if opts then
		if opts.host then
			config.host = opts.host
		end
		if opts.port then
			config.port = opts.port
		end
		if opts.reconnect_delay then
			config.reconnect_delay = opts.reconnect_delay
		end
		if opts.auto_reconnect ~= nil then -- New configuration option
			config.auto_reconnect = opts.auto_reconnect
		end
	end
end

-- Expose the connect function for manual reconnection
function M.connect_to_backend()
	connect_to_backend()
end

return M
