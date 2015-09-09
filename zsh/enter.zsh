function do_enter() {
  if [ -n "$BUFFER" ]; then
    zle accept-line
    return 0
  fi
  if [ "$(git rev-parse --is-inside-work-tree 2> /dev/null)" = 'true' ]; then
    echo
    gss
  fi
  zle reset-prompt
  return 0
}
zle -N do_enter
bindkey '^m' do_enter
