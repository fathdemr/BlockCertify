FROM alpine:3.20 AS alpine
RUN apk add -U --no-cache ca-certificates tzdata && \
    mkdir -p /tmp/uploads && \
    chmod 777 /tmp/uploads

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine /tmp /tmp
ENV ZONEINFO=/usr/share/zoneinfo

LABEL org.opencontainers.image.title="BlockCertify Api" \
      org.opencontainers.image.description="BlockCertify Api" \
      org.opencontainers.image.authors="Fatih Demir <fath.demmr@gmail.com>"

# Non-root olarak çalış (scratch'ta numeric UID kullan)
USER 65532:65532


# Buradaki binary adı ile Makefile'daki çıktı adı uyumlu olmalı (aşağıya bak)
COPY dist/linux/api api

EXPOSE 5075


ENTRYPOINT ["./api"]