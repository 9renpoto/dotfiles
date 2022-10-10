-- This file can be loaded by calling `lua require('plugins')` from your init.vim

-- Only required if you have packer configured as `opt`
vim.cmd([[packadd packer.nvim]])

return require("packer").startup(function(use)
    -- Packer can manage itself
    use("wbthomason/packer.nvim")
    use("wakatime/vim-wakatime")

    -- Lazy loading:
    -- Load on specific commands
    use({ "tpope/vim-dispatch", opt = true, cmd = { "Dispatch", "Make", "Focus", "Start" } })

    -- Load on an autocommand event
    use({ "andymass/vim-matchup", event = "VimEnter" })

    -- Load on a combination of conditions: specific filetypes or commands
    -- Also run code after load (see the "config" key)
    use({
        "dense-analysis/ale",
        ft = { "sh", "zsh", "bash", "c", "cpp", "cmake", "html", "markdown", "racket", "vim", "tex" },
        cmd = "ALEEnable",
        config = "vim.cmd[[ALEEnable]]",
    })

    use({
        "nvim-telescope/telescope.nvim",
        tag = "0.1.0",
        requires = { { "nvim-lua/plenary.nvim" } },
    })

    -- use("hrsh7th/nvim-cmp")
    -- use("hrsh7th/cmp-nvim-lsp")

    use({
        "nvim-lualine/lualine.nvim",
        requires = { "kyazdani42/nvim-web-devicons", opt = true },
    })
    local builtin = require("telescope.builtin")
    vim.keymap.set("n", "ff", builtin.find_files, {})
    vim.keymap.set("n", "fg", builtin.live_grep, {})
    vim.keymap.set("n", "fb", builtin.buffers, {})
    vim.keymap.set("n", "fh", builtin.help_tags, {})

    -- Use dependency and run lua function after load
    -- use({
    -- 	"lewis6991/gitsigns.nvim",
    -- 	requires = { "nvim-lua/plenary.nvim" },
    -- 	config = function()
    -- 		require("gitsigns").setup()
    -- 	end,
    -- })

    -- You can specify multiple plugins in a single call
    use({ "tjdevries/colorbuddy.vim", { "nvim-treesitter/nvim-treesitter", opt = true } })

    -- You can alias plugin names
    use({ "dracula/vim", as = "dracula" })
end)
