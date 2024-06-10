FROM golang:latest
RUN apt-get update && \
    apt-get install -y vim && \
    cd /usr/share && mkdir nvim && cd nvim && \
    curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim.appimage && chmod u+x nvim.appimage && \
    ln -s nvim.appimage /usr/local/bin/nvim
