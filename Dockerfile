FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /
ADD . .
RUN go build -o bin/authz /cmd/main.go

FROM alpine:3.16.0
WORKDIR /

COPY --from=builder /bin/authz .
COPY --from=builder /config/test_config.yaml /config/config.yaml
COPY --from=builder /internal/constants/query/schemas /internal/constant/query/schemas
COPY --from=builder /platform/opa/authz.rego /platform/opa/authz.rego


EXPOSE 5184
ENTRYPOINT [ "./authz" ]
