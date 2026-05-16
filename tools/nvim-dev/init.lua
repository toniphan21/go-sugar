------------------------------------------------------------------------------------------------
--- Vanilla lazy.nvim one file config
------------------------------------------------------------------------------------------------
-- Bootstrap lazy.nvim
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not (vim.uv or vim.loop).fs_stat(lazypath) then
  local lazyrepo = "https://github.com/folke/lazy.nvim.git"
  local out = vim.fn.system({ "git", "clone", "--filter=blob:none", "--branch=stable", lazyrepo, lazypath })
  if vim.v.shell_error ~= 0 then
    vim.api.nvim_echo({
      { "Failed to clone lazy.nvim:\n", "ErrorMsg" },
      { out, "WarningMsg" },
      { "\nPress any key to exit..." },
    }, true, {})
    vim.fn.getchar()
    os.exit(1)
  end
end
vim.opt.rtp:prepend(lazypath)

-- Make sure to setup `mapleader` and `maplocalleader` before
-- loading lazy.nvim so that mappings are correct.
-- This is also a good place to setup other settings (vim.opt)
vim.g.mapleader = " "
vim.g.maplocalleader = "\\"

-- Setup lazy.nvim
require("lazy").setup({
  spec = {
    { "LazyVim/LazyVim", import = "lazyvim.plugins" },
    -- add your plugins here
  },
  -- Configure any other settings here. See the documentation for more details.
  -- colorscheme that will be used when installing plugins.
  install = { colorscheme = { "habamax" } },
  -- automatically check for plugin updates
  checker = { enabled = true },
})


------------------------------------------------------------------------------------------------
--- GO SUGAR configuration
------------------------------------------------------------------------------------------------

vim.filetype.add({
  extension = { gos = "gos" },
})

vim.treesitter.language.register("go", "gos")

 --vim.lsp.config("gopls", {
 --  cmd = { "gopls" },
 --  filetypes = { "go", "gomod", "gowork", "gotmpl", "gos" },
 --  root_markers = { "go.mod", ".git" },
 --})
 --
 --vim.lsp.enable("gopls")

local gos_bin = vim.fn.getenv("GO_SUGAR_BIN")
local gos_log = vim.fn.getenv("GO_SUGAR_LOG")

vim.lsp.config("gos_lsp", {
  cmd = {gos_bin, "lsp", "-log", gos_log},
  filetypes = { "gos", "go", "gomod", "gowork", "gotmpl" },
  root_markers = { "go.mod", ".git" },
})

vim.lsp.enable("gos_lsp")
