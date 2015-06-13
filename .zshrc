export LANG=ja_JP.UTF-8;

# 補完機能
autoload -U compinit
compinit

# eval "$(direnv hook zsh)"

# 補完メッセージを読みやすくする
# http://qiita.com/PSP_T/items/82b080920a4241e96aed
zstyle ':completion:*' verbose yes
zstyle ':completion:*' format '%B%d%b'
zstyle ':completion:*:warnings' format 'No matches for: %d'
zstyle ':completion:*' group-name ''

# コマンド履歴
HISTFILE=~/.zsh_history
HISTSIZE=6000000
SAVEHIST=6000000
setopt hist_ignore_dups     # ignore duplication command history list
setopt share_history        # share command history data

# コマンド履歴検索
autoload history-search-end
zle -N history-beginning-search-backward-end history-search-end
zle -N history-beginning-search-forward-end history-search-end
bindkey "^P" history-beginning-search-backward-end
bindkey "^N" history-beginning-search-forward-end

# 補完候補を詰めて表示する
setopt list_packed

# 補完候補表示時などにピッピとビープ音をならないように設定
setopt nolistbeep

# カッコの対応などを自動的に補完する
setopt auto_param_keys

# エラーメッセージ本文出力に色付け
e_normal=`echo -e "¥033[0;30m"`
e_RED=`echo -e "¥033[1;31m"`
e_BLUE=`echo -e "¥033[1;36m"`

# 3秒以上かかった処理は詳細表示
REPORTTIME=3

# if [[ -n "$TMUX" ]]; then
#   cut -c3- ~/.tmux.conf | sh -s on_login
# fi

autoload -Uz colors; colors
autoload -Uz vcs_info

zstyle ':vcs_info:*' max-exports 3
zstyle ':vcs_info:*' enable hg
zstyle ':vcs_info:*' formats '[%s] %b'
precmd () { vcs_info }

function rprompt-git-current-branch {
  local name st color

  if [[ "$PWD" =~ '/\.git(/.*)?$' ]]; then
    return
  fi
    name=$(basename "$(git symbolic-ref HEAD 2> /dev/null)")
  if [[ -z $name ]]; then
    return
  fi
  st=$(git status 2> /dev/null)
  if [[ -n $(echo "$st" | grep "^nothing to") ]]; then
    color=${fg[green]}
  elif [[ -n $(echo "$st" | grep "^nothing added") ]]; then
    color=${fg[yellow]}
  elif [[ -n $(echo "$st" | grep "^# Untracked") ]]; then
    color=${fg_bold[red]}
  else
    color=${fg[red]}
  fi

  # %{...%} は囲まれた文字列がエスケープシーケンスであることを明示する
  # これをしないと右プロンプトの位置がずれる
  echo "%{$color%}[git] $name%{$reset_color%} "
}

function rprompt-mercurial-current-branch {
  local name st color

  if [[ "$PWD" =~ '/\.hg(/.*)?$' ]]; then
    return
  fi

  name="$vcs_info_msg_0_"
  # if [[ -z $name ]]; then
  #   return
  # fi
  st=$(hg status 2> /dev/null)
  if [[ -n $(echo "$st" | grep "^!") ]]; then
    color=${fg_bold[red]}
  elif [[ -n $(echo "$st" | grep "^M") ]]; then
    color=${fg_bold[red]}
  elif [[ -n $(echo "$st" | grep "^?") ]]; then
    color=${fg[yellow]}
  else
    color=${fg[green]}
  fi

  # %{...%} は囲まれた文字列がエスケープシーケンスであることを明示する
  # これをしないと右プロンプトの位置がずれる
  echo "%{$color%}$name%{$reset_color%} "
}

# プロンプトが表示されるたびにプロンプト文字列を評価、置換する
setopt prompt_subst

RPROMPT='`rprompt-mercurial-current-branch` `rprompt-git-current-branch` %~ %{$fg_bold[blue]%}${HOST}'

function make() {
    LANG=C command make "$@" 2>&1 | sed -e "s@[Ee]rror:.*@$e_RED&$e_normal@g" -e "s@cannot¥sfind.*@$e_RED&$e_normal@g" -e "s@[Ww]arning:.*@$e_BLUE&$e_normal@g"
}
function cwaf() {
    LANG=C command ./waf "$@" 2>&1 | sed -e "s@[Ee]rror:.*@$e_RED&$e_normal@g" -e "s@cannot¥sfind.*@$e_RED&$e_normal@g" -e "s@[Ww]arning:.*@$e_BLUE&$e_normal@g"
}

#色の設定
export LSCOLORS=gxfxxxxxcxxxxxxxxxgxgx
export LS_COLORS='di=01;36:ln=01;35:ex=01;32'
zstyle ':completion:*' list-colors 'di=36' 'ln=35' 'ex=32'

# http://qiita.com/STAR_ZERO/items/860b3f141e6f7e2cc3ca
printf "\033k$(hostname -s)\033\\"

# http://qiita.com/puttyo_bubu/items/0cf94ca5a764aa22827e
if [ -f ~/.zsh/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh ]; then
    source ~/.zsh/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
fi

case ${OSTYPE} in
darwin*) # Mac OS X
  function macvim () {
      if [ -d /Applications/MacVim.app ]
      then
        [ ! -f $1 ] && touch $1
        open -a MacVim $1
      else
        vim $1
      fi
  }
  alias vim='macvim'
  ;;
esac


source ~/src/github.com/9renpoto/dotfiles/zsh_local/alias

# http://qiita.com/awakia/items/1d5cd440ce58ef4fb8ae
[ -f ~/.zshrc.local ] && source ~/.zshrc.local

typeset -U path PATH
