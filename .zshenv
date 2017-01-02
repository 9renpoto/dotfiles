export LC_CTYPE='ja_JP.UTF-8'

# vimを使う。
export EDITOR=nvim

export HOMEBREW_CASK_OPTS="--appdir=/Applications"
export HOSTNAME=`hostname`
export XDG_CONFIG_HOME="~/.config"

PATH="/opt/homebrew/bin:$PATH"
PATH="$HOME/bin:$PATH"

# pyenv
export PYENV_VIRTUALENV_DISABLE_PROMPT=1
export PYENV_ROOT="${HOME}/.pyenv"
if [ -d "${PYENV_ROOT}" ]; then
  PATH="${PYENV_ROOT}/shims:${PYENV_ROOT}/bin:$PATH"
  eval "$(pyenv init - zsh)"
  eval "$(pyenv virtualenv-init -)"
fi

# golang
export GOPATH="$HOME"
export GOROOT=`go env GOROOT`
PATH="$GOROOT/bin:$GOPATH/bin:$PATH"
PATH="/opt/homebrew/go/bin:$PATH"
export GOENVTARGET="$HOME/.goenvtarget"
PATH="$GOENVTARGET:$PATH"

# npm
PATH="./node_modules/.bin:$HOME/src/github.com/9renpoto/dotfiles/node_modules/.bin:$PATH"

export PATH

# 重複したパスを登録しない。
typeset -U path
