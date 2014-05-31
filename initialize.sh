#!/user/bin/sh

cd ~/

mkdir ~/vimbackup
mkdir ~/github
mkdir ~/bitbucket

mkdir ~/.zsh
git clone git://github.com/zsh-users/zsh-syntax-highlighting.git ~/.zsh/zsh-syntax-highlighting

mkdir -p ~/.vim/bundle
# git clone git://github.com/Shougo/neobundle.vim ~/.vim/bundle/neobundle.vim

mkdir ~/.tmux
# git clone git://github.com/erikw/tmux-powerline.git ~/.tmux/tmux-powerline

ln -s ~/github/dotfiles/_bash_profile ~/.bash_profile
ln -s ~/github/dotfiles/_bashrc ~/.bashrc
ln -s ~/github/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/github/dotfiles/_zshenv ~/.zshenv
ln -s ~/github/dotfiles/_zshrc ~/.zshrc
ln -s ~/github/dotfiles/_vimrc ~/.vimrc

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc


