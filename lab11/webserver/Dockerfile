FROM golang:1.15-alpine AS build
EXPOSE 8000
WORKDIR /src/
COPY webserver.go ./ 
COPY go.mod ./
COPY go.sum ./
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/mongo/options
RUN CGO_ENABLED=0 go build  -o /bin/server webserver.go 

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]
