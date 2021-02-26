FROM golang:1.14.3-alpine as build

ENV CGO_ENABLED 0

WORKDIR /go/src/app
ADD . /go/src/app

RUN go build -o /go/bin/app

FROM gcr.io/distroless/static-debian10
COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/assets /assets
CMD ["/app"]