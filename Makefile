PROGRAM_NAME=vkbotserver

init:
	dep ensure -v

build:
	go build -o $(PROGRAM_NAME)