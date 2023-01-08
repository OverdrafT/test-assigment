
FROM golang:1.16-alpine

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./

COPY . .

RUN go mod download

RUN go build ./cmd/main.go 

#EXPOSE 8080

CMD [ "./main" ]