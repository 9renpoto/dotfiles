set -x XDG_CONFIG_HOME $HOME/.config
set -x EDITOR nvim
set -x HOSTNAME hostname
set -x GOPATH $HOME

set -x PATH /opt/homebrew/bin $PATH
set -gx LC_ALL en_US.UTF-8

alias vim='/opt/homebrew/bin/nvim'

if test -e $HOME/src/bitbucket.org/9renpoto/ssh/.keychain
  . $HOME/src/bitbucket.org/9renpoto/ssh/.keychain
end
source /opt/homebrew/opt/asdf/asdf.fish
