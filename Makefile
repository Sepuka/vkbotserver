git tests:
	go test ./...

tidy:
	go mod tidy

mocks:
	go get github.com/vektra/mockery/v2/.../
	mockery --all --dir api --output api/mocks