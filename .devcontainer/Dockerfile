FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        build-essential=12.9ubuntu3 \
        ca-certificates=20240203~22.04.1 \
        curl=7.81.0-1ubuntu1.21 \
        file=1:5.41-3ubuntu0.1 \
        git=1:2.34.1-1ubuntu1.15 \
        locales=2.35-0ubuntu3.11 \
        procps=2:3.3.17-6ubuntu2.1 \
        sudo=1.9.9-1ubuntu2.5 \
        tzdata=2025b-0ubuntu0.22.04.1 \
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
 && { \
        echo; \
        echo '# Initialize Homebrew'; \
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"'; \
    } >> /home/vscode/.zshrc \
 && chown vscode:vscode /home/vscode/.bashrc /home/vscode/.zshrc