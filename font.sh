#!/bin/sh

brew tap sanemat/font
brew install ricty

cmd=$(brew info ricty | grep "Ricty\*.ttf" | sed -e "s/.*\(cp -f.*\)/\1/")
echo $cmd
eval $cmd

#git clone https://github.com/Lokaltog/powerline.git ~/.powerline
#fontforge -script $HOME/.powerline/font/fontpatcher.py ~/Library/Fonts/Ricty-Regular.ttf

git clone https://github.com/Lokaltog/vim-powerline ~/.vim-powerline
fontforge -lang=py -script ~/.vim-powerline/fontpatcher/fontpatcher ~/Library/Fonts/Ricty-Regular.ttf

mv -f *.ttf ~/Library/Fonts/

# wait a minute
fc-cache -vf
