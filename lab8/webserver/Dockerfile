FROM golang:1.14-alpine AS build
EXPOSE 8050
WORKDIR /src/
COPY webserver.go ./ 
RUN apk add git
RUN go get go.mongodb.org/mongo-driver/bson &&\
    go get go.mongodb.org/mongo-driver/mongo &&\
    go get go.mongodb.org/mongo-driver/mongo/options &&\
    CGO_ENABLED=0 go build  -o /bin/server webserver.go 

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]


