# ShitChat

**A crude and fairly crappy chat server and client, made for an assignment in TTM4100.**

### Installation
~~~
go get github.com/hdhauk/ShitChat
go get github.com/asaskevich/govalidator
cd $GOPATH/src/github.com/hdhauk/ShitChat/sc-client
go install .
cd $GOPATH/src/github.com/hdhauk/ShitChat/sc-server
go install .
~~~

### How to run
1. Start the server: `sc-server -port <your port>`
2. Start the client: `sc-client -server <ip:port to your server>`
