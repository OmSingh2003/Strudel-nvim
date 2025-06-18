-- Strudel plugin commands

if vim.g.loaded_strudel then
  return
end
vim.g.loaded_strudel = 1

-- Create user commands
vim.api.nvim_create_user_command('StrudelStart', function()
  require('strudel-nvim.core').start_backend()
end, {
  desc = 'Start the Strudel backend server and connect to it'
})

vim.api.nvim_create_user_command('StrudelStop', function()
  require('strudel-nvim.core').stop_backend()
end, {
  desc = 'Stop the Strudel backend server and disconnect'
})

vim.api.nvim_create_user_command('StrudelEval', function()
  require('strudel-nvim.core').eval_line()
end, {
  desc = 'Evaluate the current line with Strudel'
})

vim.api.nvim_create_user_command('StrudelConnect', function()
  local core = require('strudel-nvim.core')
  core.connect_to_backend()
end, {
  desc = 'Connect to an already running Strudel backend'
})

-- Default keymaps
vim.keymap.set('n', '<leader>ss', ':StrudelStart<CR>', { desc = 'Start Strudel backend' })
vim.keymap.set('n', '<leader>sq', ':StrudelStop<CR>', { desc = 'Stop Strudel backend' })
vim.keymap.set('n', '<leader>se', ':StrudelEval<CR>', { desc = 'Evaluate current line with Strudel' })

print("Strudel commands loaded! Use :StrudelStart to begin.")

