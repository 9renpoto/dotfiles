set encoding=utf-8
scriptencoding utf-8

" Note: Skip initialization for vim-tiny or vim-small.
if 0 | endif

if has('vim_starting')
  " Required:
  set runtimepath+=~/.vim/bundle/neobundle.vim/
endif

" Required:
call neobundle#begin(expand('~/.vim/bundle/'))

" Let NeoBundle manage NeoBundle
" Required:
NeoBundleFetch 'Shougo/neobundle.vim'
NeoBundle 'git://github.com/elixir-lang/vim-elixir.git'
NeoBundle 'git://github.com/scrooloose/nerdtree.git'
NeoBundle 'git://github.com/Shougo/neocomplete.git'
NeoBundle 'git://github.com/Shougo/unite-outline'
NeoBundle 'git://github.com/Shougo/unite.vim.git'
NeoBundle 'git://github.com/Shougo/vim-vcs.git'
NeoBundle 'git://github.com/Shougo/vimfiler.git'
NeoBundle 'git://github.com/thinca/vim-localrc.git'
NeoBundle 'JavaScript-syntax'
NeoBundle 'kchmck/vim-coffee-script'
NeoBundle 'nginx.vim'
NeoBundle 'rking/ag.vim'
NeoBundle 'Townk/vim-autoclose'
NeoBundle "mhinz/vim-signify" " http://deris.hatenablog.jp/entry/2013/12/15/235606
NeoBundle 'osyo-manga/vim-over'  " http://qiita.com/PSP_T/items/3a1af1301ee197b32a8a

" http://blog.glidenote.com/blog/2012/04/01/new-gist.vim/
NeoBundle 'mattn/gist-vim'
NeoBundle 'mattn/webapi-vim'
let g:gist_show_privates = 1
let g:gist_post_private = 1

NeoBundle 'scrooloose/syntastic'
let g:syntastic_javascript_checkers = ['eslint']

NeoBundle 'git://github.com/nathanaelkane/vim-indent-guides.git'  " http://chiiiiiiiii.hatenablog.com/entry/2012/12/02/102815'
" let g:indent_guides_auto_colors=0
" let g:indent_guides_color_change_percent=30
" let g:indent_guides_enable_on_vim_startup=1
" let g:indent_guides_guide_size=1
" let g:indent_guides_start_level=1
" autocmd VimEnter,Colorscheme * :hi IndentGuidesOdd  guibg=#262626 ctermbg=black
" autocmd VimEnter,Colorscheme * :hi IndentGuidesEven guibg=#3c3c3c ctermbg=black

NeoBundle 'tpope/vim-surround'
set laststatus=2
set t_Co=256

" http://yuku-tech.hatenablog.com/entry/20110427/1303868482
NeoBundle 'git://github.com/tpope/vim-fugitive.git'

NeoBundle 'itchyny/lightline.vim'
let g:lightline = {
    \ 'colorscheme': 'wombat',
    \ 'mode_map': { 'c': 'NORMAL' },
    \ 'active': {
    \   'left': [ [ 'mode', 'paste' ], [ 'pyenv' ], [ 'fugitive', 'filename' ] ]
    \ },
    \ 'component_function': {
    \   'modified': 'MyModified',
    \   'readonly': 'MyReadonly',
    \   'fugitive': 'MyFugitive',
    \   'filename': 'MyFilename',
    \   'fileformat': 'MyFileformat',
    \   'filetype': 'MyFiletype',
    \   'fileencoding': 'MyFileencoding',
    \   'mode': 'MyMode',
    \   'pyenv': 'pyenv#statusline#component',
    \ },
    \ 'separator': { 'left': '', 'right': '' },
    \ 'subseparator': { 'left': '', 'right': '' }
    \ }

function! g:MyModified()
  return &ft =~? 'help\|vimfiler\|gundo' ? '' : &modified ? '+' : &modifiable ? '' : '-'
endfunction

function! g:MyReadonly()
  return &ft !~? 'help\|vimfiler\|gundo' && &readonly ? '⭤' : ''
endfunction

function! g:MyFilename()
  return ('' !=? g:MyReadonly() ? g:MyReadonly() . ' ' : '') .
        \ (&ft ==? 'vimfiler' ? g:vimfiler#get_status_string() :
        \  &ft ==? 'unite' ? g:unite#get_status_string() :
        \  &ft ==? 'vimshell' ? g:vimshell#get_status_string() :
        \ '' !=? expand('%:t') ? expand('%:t') : '[No Name]') .
        \ ('' !=? g:MyModified() ? ' ' . g:MyModified() : '')
endfunction

function! g:MyFugitive()
  if &ft !~? 'vimfiler\|gundo' && exists('*fugitive#head')
    let l:_ = g:fugitive#head()
    return strlen(l:_) ? ''.l:_ : ''
  endif
  return ''
endfunction

function! g:MyFileformat()
  return winwidth('.') > 70 ? &fileformat : ''
endfunction

function! g:MyFiletype()
  return winwidth('.') > 70 ? (strlen(&filetype) ? &filetype : 'no ft') : ''
endfunction

function! g:MyFileencoding()
  return winwidth('.') > 70 ? (strlen(&fenc) ? &fenc : &enc) : ''
endfunction

function! g:MyMode()
  return winwidth('.') > 60 ? g:lightline#mode() : ''
endfunction

NeoBundle "tyru/caw.vim.git"  " http://ichyo.jp/blog/2014/03/14/how-to-comment-out-with-vim/
nmap <Leader>c <Plug>(caw:i:toggle)
vmap <Leader>c <Plug>(caw:i:toggle)

NeoBundle 'tomasr/molokai'

call neobundle#end()

" カラースキーム
let scheme = 'molokai'
augroup guicolorscheme
  autocmd!
  execute 'autocmd GUIEnter * colorscheme' scheme
augroup END
execute 'colorscheme' scheme

" Required:
filetype plugin indent on

" If there are uninstalled bundles found on startup,
" this will conveniently prompt you to install them.
NeoBundleCheck