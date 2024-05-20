FROM golang:1.22.2-alpine as builer

WORKDIR /app

COPY . /app

RUN <<EOT
mkdir -p bin
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
CGO_ENABLED=0 GOOS=linux go build -o /app/bin/wf -ldflags="-w -s" -trimpath
EOT

FROM scratch as prob

EXPOSE 8080

VOLUME /usr/local/share/etc/waveform/
VOLUME /var/log/waveform/
VOLUME /var/log/waveform/

COPY --from=builer /app/bin/wf /

ENTRYPOINT ["/wf"]