set -x XDG_CONFIG_HOME $HOME/.config
set -x EDITOR nvim
set -x HOSTNAME hostname
set -x GOPATH $HOME
set -g GHQ_SELECTOR_OPTS "--no-sort --reverse --ansi --color bg+:13,hl:3,pointer:7"
set -U FZF_LEGACY_KEYBINDINGS 0

set -x PATH ~/.deno/bin $PATH
set -gx LC_ALL en_US.UTF-8

alias vim='/usr/local/bin/nvim'

source (brew --prefix asdf)/asdf.fish
eval (starship init fish)
