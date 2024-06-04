FROM golang:1.22.2-alpine as builer

WORKDIR /app

COPY . /app

RUN <<EOT
mkdir -p bin
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -trimpath --tags deploy -o /app/bin/wf
EOT

FROM scratch as prob

EXPOSE 8080
EXPOSE 8081

VOLUME /usr/local/share/etc/waveform/
VOLUME /var/log/waveform/
VOLUME /var/lib/waveform/

COPY --from=builer /app/bin/wf /

ENTRYPOINT ["/wf"]