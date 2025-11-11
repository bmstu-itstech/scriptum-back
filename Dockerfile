FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/http/http.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache \
    bash \
    git \
    curl \
    build-base \
    libffi-dev \
    openssl-dev \
    bzip2-dev \
    readline-dev \
    sqlite-dev \
    zlib-dev \
    xz-dev

RUN git clone https://github.com/pyenv/pyenv.git ~/.pyenv

ENV PYENV_ROOT="/root/.pyenv"
ENV PATH="$PYENV_ROOT/bin:$PATH"

RUN echo 'eval "$(pyenv init -)"' >> ~/.bashrc

COPY --from=builder /app/app .

CMD ["./app"]
