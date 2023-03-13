FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /
ADD . .
RUN go build -o bin/authz /cmd/main.go

FROM alpine:3.16.0
WORKDIR /

COPY --from=builder /bin/authz .
COPY --from=builder /config/test_config.yaml /config/config.yaml
COPY --from=builder /internal/constants/query/schemas /internal/constants/query/schemas
COPY --from=builder /platform/opa/authz.rego /platform/opa/authz.rego
COPY --from=builder /platform/opa/server/opa /platform/opa/server/opa
RUN ["apk","update"]
RUN ["apk","add","bash"]
RUN ["apk","add","lsof"]
RUN ["chmod","755","platform/opa/server/opa"]



EXPOSE 8181
ENTRYPOINT [ "./authz" ]
