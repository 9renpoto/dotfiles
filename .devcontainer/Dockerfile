FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

ARG DEBIAN_FRONTEND=noninteractive

# hadolint ignore=DL3008
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        locales \
 && rm -rf /var/lib/apt/lists/*

RUN locale-gen en_US.UTF-8 \
 && update-locale LANG=en_US.UTF-8

RUN curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh -o /tmp/install.sh \
 && chmod +x /tmp/install.sh \
 && su -l vscode -c 'NONINTERACTIVE=1 /tmp/install.sh' \
 && rm -f /tmp/install.sh

# hadolint ignore=SC2016
RUN mkdir -p /etc/profile.d \
 && echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' > /etc/profile.d/homebrew.sh \
 && chmod 0644 /etc/profile.d/homebrew.sh \
 && { \
        echo; \
        echo '# Initialize Homebrew'; \
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"'; \
    } >> /home/vscode/.bashrc \
 && chown vscode:vscode /home/vscode/.bashrc
