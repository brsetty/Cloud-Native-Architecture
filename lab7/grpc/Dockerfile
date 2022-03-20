FROM golang:1.15-alpine AS build
EXPOSE 8050
WORKDIR /src/
COPY movieapi movieapi/
COPY movieserver movieserver/
COPY go.mod ./
COPY go.sum ./
RUN go get ./... &&\
    CGO_ENABLED=0 go build  -o /bin/server movieserver/server.go 

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]

