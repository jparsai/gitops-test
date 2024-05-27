FROM golang:1.22-alpine

WORKDIR /gitops-test

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE='/gitops-test/.cache/go-build'

RUN go build main.go
#RUN mkdir -p /.cache
#RUN chmod -R 777 /.cache

CMD ["go", "run", "main.go"]
