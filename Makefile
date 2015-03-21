all: 	build

build:
	go build -o driver-limitlessled/driver-limitlessled
	cp ./package.json driver-limitlessled/package.json
