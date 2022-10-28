FROM alpine:3.16

RUN apk add --no-cache wget restic unzip

RUN wget https://releases.hashicorp.com/vault/1.12.0/vault_1.12.0_linux_amd64.zip --quiet && \
    unzip vault_1.12.0_linux_amd64.zip && \
    rm vault_1.12.0_linux_amd64.zip && \
    mv /vault /usr/bin/vault
