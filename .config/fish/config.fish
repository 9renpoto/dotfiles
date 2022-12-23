set -x XDG_CONFIG_HOME $HOME/.config
set -x EDITOR nvim
set -x HOSTNAME hostname
set -x GOPATH $HOME

set -x PATH ~/.deno/bin $PATH
set -gx LC_ALL en_US.UTF-8

eval (starship init fish)
