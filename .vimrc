" http://rbtnn.hateblo.jp/entry/2014/11/30/174749
set encoding=utf-8
scriptencoding utf-8

augroup vimrc
  autocmd!
augroup END

source ~/src/github.com/9renpoto/dotfiles/vim/neobundle
source ~/src/github.com/9renpoto/dotfiles/vim/neocomplete

if has('gui_macvim')
  set noundofile
endif

highlight JpSpace cterm=underline ctermfg=7 guifg=7
au BufRead,BufNew * match JpSpace /■/

highlight ZenkakuSpace cterm=underline ctermbg=red guibg=#666666

augroup space
  autocmd BufWinEnter * let w:m3 = matchadd("ZenkakuSpace", '　')
  autocmd WinEnter * let w:m3 = matchadd("ZenkakuSpace", '　')
augroup END

highlight WhitespaceEOL ctermbg=red guibg=red
match WhitespaceEOL /\s\+$/

augroup EndSpace
  autocmd WinEnter * match WhitespaceEOL /\s\+$/
augroup END

" http://qiita.com/knt45/items/9717e30ca6a0f1fdad0f
set guioptions-=r
set guioptions-=R
set guioptions-=l
set guioptions-=L

set termencoding=utf-8
set fileformats=unix,mac,dos
set ambiwidth=double
set autoindent
set browsedir=buffer
set backup
set backupdir=$HOME/vimbackup
set directory=$HOME/vimbackup
set expandtab
set hidden
set incsearch
set number
set showmatch
set smartcase
set smartindent

augroup bufWritePre
  autocmd BufWritePre * :%s/\s\+$//ge
augroup END

set shiftwidth=2
set smarttab

au BufNewFile,BufRead *.yml set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.html set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.coffee set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.js set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.rb set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.cs set tabstop=4 softtabstop=4 shiftwidth=4
au BufNewFile,BufRead *.md set tabstop=2 softtabstop=2 shiftwidth=2

set whichwrap=b,s,h,l,<,>,[,]
set nowrapscan
set hlsearch

set cursorline
augroup cch
  autocmd! cch
  autocmd WinLeave * set nocursorline
  autocmd WinEnter,BufRead * set cursorline
augroup END
:hi clear CursorLine
:hi CursorLine gui=underline
hi CursorLine ctermbg=black guibg=black

filetype off
filetype plugin indent off
set completeopt=menu,preview

" http://kazuph.github.io/presentation/yapc_vim_2013_github/#/46
" perl
augroup filetypedetect
  au BufNewFile,BufRead *.psgi    setf perl
  au BufNewFile,BufRead *.t       setf perl
  au BufNewFile,BufRead cpanfile  setf perl
augroup END

filetype plugin indent on

syntax on
