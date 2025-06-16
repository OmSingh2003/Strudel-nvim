-- lua/strudel-nvim/init.lua

local M = {}

local core = require("strudel-nvim.core")

-- This is the main setup function for the plugin.
-- It will be called by  plugin manager.
function M.setup(opts)
	-- Create user commands
	vim.api.nvim_create_user_command("StrudelStart", function()
		core.start_backend()
	end, { desc = "Start the Go backend and connect to it." })

	vim.api.nvim_create_user_command("StrudelStop", function()
		core.stop_backend()
	end, { desc = "Disconnect from and stop the Go backend." })

	vim.api.nvim_create_user_command("StrudelEvalLine", function()
		core.eval_line()
	end, { desc = "Evaluate the current line with Strudel." })

	-- Create a keymap
	-- Pressing <Space>e will  run the command.
	vim.keymap.set("n", "<leader>e", "<Cmd>StrudelEvalLine<CR>", {
		noremap = true,
		silent = true,
		desc = "Strudel: Evaluate current line",
	})

	vim.notify("Strudel plugin loaded!", vim.log.levels.INFO)
end

return M
