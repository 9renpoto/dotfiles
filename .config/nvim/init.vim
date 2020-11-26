set runtimepath+=$HOME/src/github.com/Shougo/dein.vim

if dein#load_state($HOME . '/.config/nvim/dein')
  let g:dein#cache_directory = $HOME . '/.cache/dein'
  call dein#begin($HOME . '/.cache/dein')

  let s:toml_dir  = $HOME . '/.config/nvim/dein/toml'
  let s:toml      = s:toml_dir . '/dein.toml'
  let s:lazy_toml = s:toml_dir . '/dein_lazy.toml'
  call dein#load_toml(s:toml,      {'lazy': 0})
  call dein#load_toml(s:lazy_toml, {'lazy': 1})

  call dein#end()
  call dein#save_state()
endif

filetype plugin indent on
syntax enable

set autoindent
set expandtab
set guifont=Ricty:h14
set number
set shiftwidth=2
set smartindent
set softtabstop=2
set tabstop=2
set showtabline=2
set termguicolors

colorscheme atom-dark-256

augroup on_save
  autocmd BufWritePre *.{js,jsx,ts,tsx,json,yml,rs} :%s/\s\+$//ge
augroup END

augroup HighlightTrailingSpaces
  autocmd!
  autocmd VimEnter,WinEnter,ColorScheme * highlight TrailingSpaces term=underline guibg=Red ctermbg=Red
  autocmd VimEnter,WinEnter * match TrailingSpaces /\s\+$/
augroup END

call dein#add('Shougo/deoplete.nvim')
if !has('nvim')
  call dein#add('roxma/nvim-yarp')
  call dein#add('roxma/vim-hug-neovim-rpc')
endif
let g:deoplete#enable_at_startup = 1

call dein#add('Shougo/neosnippet.vim')
call dein#add('Shougo/neosnippet-snippets')
call dein#add('junegunn/fzf', { 'build': './install --all', 'merged': 0 })
call dein#add('junegunn/fzf.vim', { 'depends': 'fzf' })
