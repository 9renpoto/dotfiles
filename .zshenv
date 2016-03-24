export LC_CTYPE='ja_JP.UTF-8'

# vimを使う。
export EDITOR=vim

# vimがなくてもvimでviを起動する。
if ! type vim > /dev/null 2>&1; then
  alias vim=vi
fi

export PATH="/opt/homebrew/bin:$PATH"
export PATH="$PATH:$HOME/bin"
# export PATH="/usr/local/bin:$PATH"

# pyenv
export PYENV_VIRTUALENV_DISABLE_PROMPT=1
export PYENV_ROOT="${HOME}/.pyenv"
if [ -d "${PYENV_ROOT}" ]; then
  export PATH="${PYENV_ROOT}/shims:${PYENV_ROOT}/bin:$PATH"
  eval "$(pyenv init - zsh)"
  eval "$(pyenv virtualenv-init -)"
fi

# npm
export PATH="./node_modules/.bin:$PATH"

# ruby
# export RBENV_ROOT="${HOME}/.rbenv"
# if [ -d "${RBENV_ROOT}" ]; then
#   export PATH="$HOME/.rbenv/bin:$PATH"
#   eval "$(rbenv init - zsh)"
# fi

# perl
# export PLENV_ROOT="${HOME}/.pyenv"
# if [ -d "${PLENV_ROOT}" ]; then
#   export PATH="${PLENV_ROOT}/shims:${PLENV_ROOT}/bin:$PATH"
#   eval "$(plenv init - zsh)"
# fi

# golang
export GOPATH="$HOME"
export GOROOT=`go env GOROOT`
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"
export PATH="/opt/homebrew/go/bin:$PATH"
export GOENVTARGET="$HOME/.goenvtarget"
export PATH="$GOENVTARGET:$PATH"

# elixir
export EXENV_ROOT="${HOME}/.exenv"
if [ -d "${EXENV_ROOT}" ]; then
  export PATH="${EXENV_ROOT}/bin:$PATH"
  eval "$(exenv init - zsh)"
fi

export HOMEBREW_CASK_OPTS="--appdir=/Applications"

export NVM_DIR=~/.nvm
. $(brew --prefix nvm)/nvm.sh

# 重複したパスを登録しない。
typeset -U path
