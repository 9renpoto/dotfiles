#!/bin/sh

# http://qiita.com/usamik26/items/601f5612bd3f8a21cc41
cd /opt
sudo mkdir -p homebrew
sudo chown ${USER}:staff homebrew

mkdir -p ~/vimbackup

ln -s ~/src/github.com/9renpoto/dotfiles/.agignore ~/.agignore
ln -s ~/src/github.com/9renpoto/dotfiles/.asdf ~/.asdf
ln -s ~/src/github.com/9renpoto/dotfiles/.editorconfig ~/.editorconfig
ln -s ~/src/github.com/9renpoto/dotfiles/.gitconfig ~/.gitconfig
ln -s ~/src/github.com/9renpoto/dotfiles/.hyper.js ~/.hyper.js
ln -s ~/src/github.com/9renpoto/dotfiles/.ignore ~/.gitignore
ln -s ~/src/github.com/9renpoto/dotfiles/.tmux ~/.tmux
ln -s ~/src/github.com/9renpoto/dotfiles/.tmux.conf ~/.tmux.conf
ln -s ~/src/github.com/9renpoto/dotfiles/.tool-versions ~/.tool-versions
ln -s ~/src/github.com/9renpoto/dotfiles/.vimrc ~/.vimrc
ln -s ~/src/github.com/9renpoto/dotfiles/.zshenv ~/.zshenv
ln -s ~/src/github.com/9renpoto/dotfiles/.zshrc ~/.zshrc
ln -s ~/src/github.com/9renpoto/dotfiles/vim ~/.vim

asdf plugin-add nodejs https://github.com/asdf-vm/asdf-nodejs.git
