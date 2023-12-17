FROM golang:alpine

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . /app

RUN go build -o main ./main.go

EXPOSE 3242 3242
CMD [ "/app/main" ]