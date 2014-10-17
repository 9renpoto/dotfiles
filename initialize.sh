#!/user/bin/sh

mkdir -p ~/vimbackup

ln -s ~/src/github.com/9renpoto/dotfiles/_bash_profile ~/.bash_profile
ln -s ~/src/github.com/9renpoto/dotfiles/_bashrc ~/.bashrc
ln -s ~/src/github.com/9renpoto/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/src/github.com/9renpoto/dotfiles/_zshenv ~/.zshenv
ln -s ~/src/github.com/9renpoto/dotfiles/_zshrc ~/.zshrc
ln -s ~/src/github.com/9renpoto/dotfiles/_vimrc ~/.vimrc

ln -s ~/src/github.com/9renpoto/dotfiles/vim ~/.vim
ln -s ~/src/github.com/9renpoto/dotfiles/zsh ~/.zsh
ln -s ~/src/github.com/9renpoto/dotfiles/tmux ~/.tmux

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc
