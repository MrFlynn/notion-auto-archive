FROM alpine:latest as base

RUN adduser -u 10000 -H -D ${PROJECT_NAME}

FROM scratch

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /etc/passwd /etc/passwd

COPY ${PROJECT_NAME} /bin/${PROJECT_NAME}

USER ${PROJECT_NAME}

ENTRYPOINT [ "/bin/${PROJECT_NAME}" ]