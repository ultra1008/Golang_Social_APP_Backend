FROM golang:alpine as builder

# ENV
ENV APPDIR $GOPATH/src/github.com/niklod/highload-social-network

# FS
RUN mkdir -p ${APPDIR}
WORKDIR ${APPDIR}

# Скачивание модулей
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копирование зависимостей
COPY templates /build/templates
COPY static /build/static

COPY . .

RUN go build -o /build/hsn ./cmd/highload-social-network 

FROM alpine:3.7

COPY --from=builder /build/hsn /hsn
COPY --from=builder /build/templates templates
COPY --from=builder /build/static static

CMD ["./hsn"]