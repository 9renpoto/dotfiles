#!/user/bin/sh

cd ~/

mkdir ~/vimbackup
mkdir ~/github
mkdir ~/bitbucket
mkdir ~/.zsh
git clone git://github.com/zsh-users/zsh-syntax-highlighting.git ~/.zsh/zsh-syntax-highlighting

ln -s ~/dotfiles/_bashrc ~/.bashrc
ln -s ~/dotfiles/_tmux.conf ~/.tmux.conf
ln -s ~/dotfiles/_zshrc ~/.zshrc

ln -s ~/dotfiles/_vimrc ~/.vimrc

# http://qiita.com/tatsuya6502@github/items/a356db103f3b654d6f27
# ln -s ~/dotfiles/_kerlrc ~/.kerlrc

