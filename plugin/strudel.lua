-- Strudel plugin entry point

if vim.g.loaded_strudel then
  return
end
vim.g.loaded_strudel = 1

-- The plugin now requires explicit setup() call
-- All commands and keymaps are configured in the setup() function
-- This ensures users have full control over configuration

