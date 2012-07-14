set nocompatible
filetype off

set rtp+=~/vimfiles/vundle.git/        #vundleのディレクトリ
call vundle#rc()
Bundle 'Shougo/neocomplcache'        #Bundle...は使用するプラグインを書く。詳細はguthubのREADMEが詳しい。
Bundle 'Shougo/unite.vim'
Bundle 'thinca/vim-ref'
Bundle 'thinca/vim-quickrun'
filetype plugin indent on     " required!
