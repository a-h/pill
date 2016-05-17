FROM golang:1.6
ENV GOPATH /go
RUN go get gopkg.in/mgo.v2 && go get github.com/gorilla/mux && \
 go get github.com/gorilla/sessions && \
 go get github.com/gorilla/securecookie && \
 go get github.com/gorilla/context
COPY . /go/src/github.com/a-h/pill
WORKDIR /go/src/github.com/a-h/pill
RUN go get -d -v ./...
WORKDIR /go/src/github.com/a-h/pill/httpservice/main
RUN go build github.com/a-h/pill/httpservice/main
CMD ./main
EXPOSE 8080
