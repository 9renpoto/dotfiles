FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        file \
        git \
        locales \
        procps \
        sudo \
        tzdata \
 && rm -rf /var/lib/apt/lists/*

RUN locale-gen en_US.UTF-8 \
 && update-locale LANG=en_US.UTF-8

RUN curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh -o /tmp/install.sh \
 && chmod +x /tmp/install.sh \
 && su -l vscode -c 'NONINTERACTIVE=1 /tmp/install.sh' \
 && rm -f /tmp/install.sh

RUN mkdir -p /etc/profile.d \
 && printf 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"\n' > /etc/profile.d/homebrew.sh \
 && chmod 0644 /etc/profile.d/homebrew.sh \
 && printf '\n# Initialize Homebrew\n' >> /home/vscode/.bashrc \
 && printf 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"\n' >> /home/vscode/.bashrc \
 && printf '\n# Initialize Homebrew\n' >> /home/vscode/.zshrc \
 && printf 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"\n' >> /home/vscode/.zshrc \
 && chown vscode:vscode /home/vscode/.bashrc /home/vscode/.zshrc
