FROM alpine:3.18

ARG APPVER
RUN arch="$(apk --print-arch)"; case "$arch" in 'x86_64') arch="amd64"; ;; 'armhf') arch="armv6"; ;; 'armv7') arch="armv7"; ;; 'aarch64') arch="arm64"; ;; *) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; esac \
    && wget "https://github.com/lollipopkit/nano-db/releases/download/v${APPVER}/ndb_${APPVER}_linux_$arch.tar.gz" \
    && tar -xvf "ndb_${APPVER}_linux_$arch.tar.gz" \
    && rm "ndb_${APPVER}_linux_$arch.tar.gz" \
    && mv ndb /usr/bin \
    && chmod +x /usr/bin/ndb

ENTRYPOINT ["/usr/bin/ndb", "serve"]