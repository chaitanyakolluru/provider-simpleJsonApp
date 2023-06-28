FROM docker.nexus.heb.tools/golang:1.20.2-alpine AS golang

# build 
RUN mkdir -p /build/provider-simpleJsonApp
COPY cmd/ /build/provider-simpleJsonApp/cmd/
COPY apis/ /build/provider-simpleJsonApp/apis/
COPY internal/ /build/provider-simpleJsonApp/internal/
COPY go.* /build/provider-simpleJsonApp/

WORKDIR /build/provider-simpleJsonApp
RUN go build -o provider-simpleJsonApp ./cmd/provider

FROM docker.nexus.heb.tools/alpine:3.16.0

RUN apk update && apk add --no-cache openssl \
    && rm -rf /var/cache/apk/*

COPY --from=golang /build/provider-simpleJsonApp/provider-simpleJsonApp /usr/local/bin/

ENTRYPOINT [ "/bin/sh", "-c", "provider-simpleJsonApp" ]
