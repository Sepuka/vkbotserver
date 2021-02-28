init:
	dep ensure -v

dependencies:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

tests:
	go test ./...