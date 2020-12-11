FROM golang:1.15 as builder

WORKDIR /workdir
COPY . .
RUN go build

FROM alpine:3.11

LABEL maintainer="Gary Kim <gary@garykim.dev>"

RUN apk add --no-cache ca-certificates libc6-compat libstdc++
COPY --from=builder /workdir/cooperdiscord .

CMD ["./cooperdiscord", "start"]
