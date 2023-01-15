
FROM golang:1.19-alpine

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./

COPY . .

RUN go mod tidy

RUN go build ./cmd/main.go 

#EXPOSE 8080

CMD [ "./main" ]