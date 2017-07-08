#!/bin/sh

# http://qiita.com/usamik26/items/601f5612bd3f8a21cc41
sudo mkdir -p /opt
cd /opt
sudo mkdir -p homebrew
sudo chown ${USER}:staff homebrew

ln -s ~/src/github.com/9renpoto/dotfiles/.agignore ~/.agignore
ln -s ~/src/github.com/9renpoto/dotfiles/.config ~/.config
ln -s ~/src/github.com/9renpoto/dotfiles/.editorconfig ~/.editorconfig
ln -s ~/src/github.com/9renpoto/dotfiles/.gitconfig ~/.gitconfig
ln -s ~/src/github.com/9renpoto/dotfiles/.hyper.js ~/.hyper.js
ln -s ~/src/github.com/9renpoto/dotfiles/.ignore ~/.gitignore
ln -s ~/src/github.com/9renpoto/dotfiles/.tmux ~/.tmux
ln -s ~/src/github.com/9renpoto/dotfiles/.tmux.conf ~/.tmux.conf
