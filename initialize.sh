#!/user/bin/sh

cd ~/

mkdir -p ~/vimbackup
mkdir -p ~/github
mkdir -p ~/bitbucket

ln -s ~/github/dotfiles/_bash_profile ~/.bash_profile
ln -s ~/github/dotfiles/_bashrc ~/.bashrc
ln -s ~/github/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/github/dotfiles/_zshenv ~/.zshenv
ln -s ~/github/dotfiles/_zshrc ~/.zshrc
ln -s ~/github/dotfiles/_vimrc ~/.vimrc

ln -s ~/github/dotfiles/vim ~/.vim
ln -s ~/github/dotfiles/zsh ~/.zsh
ln -s ~/github/dotfiles/tmux ~/.tmux

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc
