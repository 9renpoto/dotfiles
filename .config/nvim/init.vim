set runtimepath+=$HOME/src/github.com/Shougo/dein.vim

if dein#load_state($HOME . '/src/github/9renpoto/dotfiles/.config/nvim/dein')
  call dein#begin($HOME . '/src/github/9renpoto/dotfiles/.cache/dein')

  call dein#add($HOME . '/src/github.com/Shougo/dein.vim')
  call dein#add('Shougo/neocomplete.vim')
  call dein#add('gosukiwi/vim-atom-dark')

  call dein#end()
  call dein#save_state()
endif

filetype plugin indent on
syntax enable

set expandtab
set tabstop=2
set shiftwidth=2
set softtabstop=2
set autoindent
set smartindent
set number

set guifont=Ricty:h14

set showtabline=2

set termguicolors

colorscheme atom-dark-256

