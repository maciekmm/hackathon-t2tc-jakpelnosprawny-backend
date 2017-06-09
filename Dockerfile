FROM golang:alpine

WORKDIR /go/src/github.com/maciekmm/t2tc-backend
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 4000

CMD ["go-wrapper", "run"]