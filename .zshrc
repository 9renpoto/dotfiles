setopt hist_ignore_all_dups
setopt EXTENDED_GLOB

eval "$(starship init zsh)"

[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh
