-- lua/strudel-nvim/init.lua

local M = {}

local core = require("strudel-nvim.core")

-- Default keymap configuration
local default_keymaps = {
	start = "<leader>ss",
	stop = "<leader>sq",
	eval_line = "<leader>se",
	connect = "<leader>sc",
}

function M.setup(opts)
	opts = opts or {}
	
	-- Merge user keymaps with defaults
	local keymaps = vim.tbl_deep_extend("force", default_keymaps, opts.keymaps or {})
	
	-- Setup core functionality
	core.setup(opts)
	
	-- Create user commands
	vim.api.nvim_create_user_command("StrudelStart", function()
		core.start_backend()
	end, { desc = "Start the Strudel backend server and connect to it" })

	vim.api.nvim_create_user_command("StrudelStop", function()
		core.stop_backend()
	end, { desc = "Stop the Strudel backend server and disconnect" })

	vim.api.nvim_create_user_command("StrudelEval", function()
		core.eval_line()
	end, { desc = "Evaluate the current line with Strudel" })

	vim.api.nvim_create_user_command("StrudelConnect", function()
		core.connect_to_backend()
	end, { desc = "Connect to an already running Strudel backend" })
	
	-- Setup keymaps only if they're not disabled
	if opts.disable_keymaps ~= true then
		if keymaps.start and keymaps.start ~= "" then
			vim.keymap.set("n", keymaps.start, "<Cmd>StrudelStart<CR>", {
				noremap = true,
				silent = true,
				desc = "Strudel: Start backend",
			})
		end
		
		if keymaps.stop and keymaps.stop ~= "" then
			vim.keymap.set("n", keymaps.stop, "<Cmd>StrudelStop<CR>", {
				noremap = true,
				silent = true,
				desc = "Strudel: Stop backend",
			})
		end
		
		if keymaps.eval_line and keymaps.eval_line ~= "" then
			vim.keymap.set("n", keymaps.eval_line, "<Cmd>StrudelEval<CR>", {
				noremap = true,
				silent = true,
				desc = "Strudel: Evaluate current line",
			})
		end
		
		if keymaps.connect and keymaps.connect ~= "" then
			vim.keymap.set("n", keymaps.connect, "<Cmd>StrudelConnect<CR>", {
				noremap = true,
				silent = true,
				desc = "Strudel: Connect to backend",
			})
		end
	end

	vim.notify("Strudel plugin loaded!", vim.log.levels.INFO)
end

return M
