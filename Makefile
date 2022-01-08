tests:
	go test ./...

mocks:
	go get github.com/vektra/mockery/v2/.../
	mockery --all --dir api --output api/mocks