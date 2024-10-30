
FROM golang:1.22

COPY . /src
RUN  cd /src && go build -o ./mock_secware/cmd/mock_secware ./mock_secware/cmd/main.go
CMD ["/src/mock_secware/cmd/mock_secware"]
