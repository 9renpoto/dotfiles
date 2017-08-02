set -x XDG_CONFIG_HOME $HOME/.config
set -x EDITOR nvim
set -x HOSTNAME hostname

set -x PATH /opt/homebrew/bin $PATH

set -x PATH $HOME/.pyenv/bin $PATH
. (pyenv init - | psub)

if test -e $HOME/src/bitbucket.org/9renpoto/ssh/.keychain
  . $HOME/src/bitbucket.org/9renpoto/ssh/.keychain
end
