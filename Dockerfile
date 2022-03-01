FROM alpine

RUN apk update  && apk add --no-cache ca-certificates

CMD ["/bin/sh", "-c", "./test_assignment"]

EXPOSE 8081

WORKDIR /test_assignment

COPY . /test_assignment/

RUN chmod +x test_assignment
