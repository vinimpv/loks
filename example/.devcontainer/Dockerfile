FROM mcr.microsoft.com/devcontainers/base:ubuntu

ENV container=docker

COPY --chmod=755 ./bootstrap.sh ./bootstrap.sh

ENTRYPOINT ["./bootstrap.sh"]
