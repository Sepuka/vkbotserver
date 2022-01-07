tests:
	go test ./...

mocks:
	mockery -all -dir api -output api/mocks