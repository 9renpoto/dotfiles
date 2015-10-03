#!/user/bin/sh

# http://qiita.com/usamik26/items/601f5612bd3f8a21cc41
cd /opt
sudo mkdir homebrew
sudo chown ${USER}:staff homebrew

mkdir -p ~/vimbackup

ln -s ~/src/github.com/9renpoto/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/src/github.com/9renpoto/dotfiles/.zshenv ~/.zshenv
ln -s ~/src/github.com/9renpoto/dotfiles/.zshrc ~/.zshrc
ln -s ~/src/github.com/9renpoto/dotfiles/_vimrc ~/.vimrc

ln -s ~/src/github.com/9renpoto/dotfiles/.hgrc ~/.hgrc
cp ~/src/github.com/9renpoto/dotfiles/.hgrc_auth.base ~/src/github.com/9renpoto/dotfiles/.hgrc_auth
ln -s ~/src/github.com/9renpoto/dotfiles/.hgrc_auth ~/.hgrc_auth
ln -s ~/src/github.com/9renpoto/dotfiles/.ignore ~/.hgignore_global

ln -s ~/src/github.com/9renpoto/dotfiles/.gitconfig ~/.gitconfig
ln -s ~/src/github.com/9renpoto/dotfiles/.ignore ~/.gitignore
ln -s ~/src/github.com/9renpoto/dotfiles/.ignore ~/.agignore

ln -s ~/src/github.com/9renpoto/dotfiles/vim ~/.vim
ln -s ~/src/github.com/9renpoto/dotfiles/tmux ~/.tmux

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc
