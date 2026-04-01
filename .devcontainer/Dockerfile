FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

ARG DEBIAN_FRONTEND=noninteractive

# hadolint ignore=DL3008
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        locales \
 && rm -rf /var/lib/apt/lists/*

RUN locale-gen en_US.UTF-8 \
 && update-locale LANG=en_US.UTF-8
