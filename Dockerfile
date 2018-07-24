FROM golang:1.10

RUN curl -fsSL -o /usr/local/bin/dep \
  https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/fallafeljan/nereiden
WORKDIR /go/src/github.com/fallafeljan/nereiden

COPY Gopkg.toml Gopkg.lock Makefile ./
RUN make prebuild-ci

COPY *.go ./
RUN make pure-build

CMD ["./build/nereiden"]
