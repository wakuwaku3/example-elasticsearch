FROM mcr.microsoft.com/vscode/devcontainers/universal:focal

# codespace ユーザーに /go 以下の権限を付与
RUN chown -R codespace:codespace /go

WORKDIR /home/codespace
USER codespace
