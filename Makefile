APP = gocrawl
VERSION = 0.0.1

test:
	go test -v -race github.com/brewkode/gocrawl

setup:
	go get github.com/parnurzeal/gorequest
	go get github.com/PuerkitoBio/goquery
	go get github.com/hashicorp/go-multierror

build:
	go build -o ${APP} .

clean:
	rm -f ${APP}

install: build
	sudo install -d /usr/local/bin
	sudo install -c ${APP} /usr/local/bin/${APP}

uninstall:
	sudo rm /usr/local/bin/${APP}
