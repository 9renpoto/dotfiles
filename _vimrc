" plugin

source ~/github/dotfiles/_vimrc_bundle
source ~/github/dotfiles/_vimrc_neocomplcache

"全角スペース
highlight JpSpace cterm=underline ctermfg=7 guifg=7
au BufRead,BufNew * match JpSpace /■/

" 全角スペースの表示
highlight ZenkakuSpace cterm=underline ctermbg=red guibg=#666666
au BufWinEnter * let w:m3 = matchadd("ZenkakuSpace", '　')
au WinEnter * let w:m3 = matchadd("ZenkakuSpace", '　')

"文末空白の表示
highlight WhitespaceEOL ctermbg=red guibg=red
match WhitespaceEOL /\s\+$/
autocmd WinEnter * match WhitespaceEOL /\s\+$/

command! Utf8 e ++enc=utf-8
command! Euc e ++enc=euc-jp
command! Sjis e ++enc=cp932
command! Jis e ++enc=iso-2022-jp
command! WUtf8 w ++enc=utf-8 | e
command! WEuc w ++enc=euc-jp | e
command! WSjis w ++enc=cp932 | e
command! WJis w ++enc=iso-2022-jp | ecommand! Utf8 e ++enc=utf-8
command! Euc e ++enc=euc-jp
command! Sjis e ++enc=cp932
command! Jis e ++enc=iso-2022-jp
command! WUtf8 w ++enc=utf-8 | e
command! WEuc w ++enc=euc-jp | e
command! WSjis w ++enc=cp932 | e
command! WJis w ++enc=iso-2022-jp | e

" ESCでIMEを確実にOFF http://nobeans.hatenablog.com/entry/20090211/1234326782
inoremap <ESC> <ESC>:set iminsert=0<CR>

" http://qiita.com/knt45/items/9717e30ca6a0f1fdad0f
set guioptions-=r
set guioptions-=R
set guioptions-=l
set guioptions-=L

set encoding=utf-8
"新しい行のインデントを現在行と同じにする
set autoindent
"ファイル保存ダイアログの初期ディレクトリをバッファファイル位置に設定
set browsedir=buffer
"Vi互換をオフ
set nocompatible
"スワップファイル用のディレクトリ
set directory=$HOME/vimbackup
"タブの代わりに空白文字を挿入する
set expandtab
"変更中のファイルでも、保存しないで他のファイルを表示
set hidden
"インクリメンタルサーチを行う
set incsearch
"タブ文字、行末など不可視文字を表示する
" set list
"listで表示される文字のフォーマットを指定する
" set listchars=tab:>\ ,extends:<
"行番号を表示する
set number
"閉じ括弧が入力されたとき、対応する括弧を表示する
set showmatch
"検索時に大文字を含んでいたら大/小を区別
set smartcase
"新しい行を作ったときに高度な自動インデントを行う
set smartindent

" 保存時に行末の空白を除去する
autocmd BufWritePre * :%s/\s\+$//ge
" 保存時にtabをスペースに変換する
" autocmd BufWritePre * :%s/\t/  /ge

"シフト移動幅
set shiftwidth=4
"ファイル内の <Tab> が対応する空白の数
au BufNewFile,BufRead *.haml set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.rb set tabstop=2 softtabstop=2 shiftwidth=2
au BufNewFile,BufRead *.cs set tabstop=4 softtabstop=4 shiftwidth=4
au BufNewFile,BufRead *.md set tabstop=2 softtabstop=2 shiftwidth=2

"行頭の余白内で Tab を打ち込むと、'shiftwidth' の数だけインデントする。
set smarttab
"カーソルを行頭、行末で止まらないようにする
set whichwrap=b,s,h,l,<,>,[,]
"検索をファイルの先頭へループしない
set nowrapscan
"検索結果文字列のハイライトを有効にする
set hlsearch

" カーソル行をハイライト
set cursorline
" 縦行ハイライト
" set cursorcolumn
augroup cch
        autocmd! cch
            autocmd WinLeave * set nocursorline
                autocmd WinEnter,BufRead * set cursorline
            augroup END
            :hi clear CursorLine
            :hi CursorLine gui=underline
            hi CursorLine ctermbg=black guibg=black

" set guifont=MyFont\ for\ Powerline
let g:Powerline_symbols = 'fancy'

" http://qiita.com/methane/items/4905f40e4772afec3e60
" :Fmt などで gofmt の代わりに goimports を使う
let g:gofmt_command = 'goimports'

filetype off
filetype plugin indent off
" Go に付属の plugin と gocode を有効にする
if $GOROOT != ''
    set rtp+=$GOROOT/misc/vim
endif

exe "set rtp+=".globpath($GOPATH, "src/github.com/nsf/gocode/vim")
set completeopt=menu,preview

filetype plugin indent on

" 保存時に :Fmt する
au BufWritePre *.go Fmt
au BufNewFile,BufRead *.go set sw=4 noexpandtab ts=4
au FileType go compiler go

syntax on
