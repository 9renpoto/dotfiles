#!/user/bin/sh

cd ~/

mkdir -p ~/vimbackup
mkdir -p ~/github/9renpoto
mkdir -p ~/bitbucket/9renpoto

ln -s ~/github/9renpoto/dotfiles/_bash_profile ~/.bash_profile
ln -s ~/github/9renpoto/dotfiles/_bashrc ~/.bashrc
ln -s ~/github/9renpoto/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/github/9renpoto/dotfiles/_zshenv ~/.zshenv
ln -s ~/github/9renpoto/dotfiles/_zshrc ~/.zshrc
ln -s ~/github/9renpoto/dotfiles/_vimrc ~/.vimrc

ln -s ~/github/9renpoto/dotfiles/vim ~/.vim
ln -s ~/github/9renpoto/dotfiles/zsh ~/.zsh
ln -s ~/github/9renpoto/dotfiles/tmux ~/.tmux

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc
