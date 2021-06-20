FROM alpine:latest as base

RUN adduser -u 10000 -H -D notion-auto-archive

FROM scratch

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /etc/passwd /etc/passwd

COPY notion-auto-archive /bin/notion-auto-archive

USER notion-auto-archive

ENTRYPOINT [ "/bin/notion-auto-archive" ]
CMD [ "-config=/config.yml" ]