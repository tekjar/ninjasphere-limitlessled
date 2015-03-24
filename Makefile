host:
	go build -o driver-limitlessled/driver-limitlessled
	cp ./package.json driver-limitlessled/package.json
	cp ./milights.xml driver-limitlessled/milights.xml

target:
	GOOS=linux GOARCH=arm go build -o driver-limitlessled/driver-limitlessled
	cp ./package.json driver-limitlessled/package.json
	cp ./milights.xml driver-limitlessled/milights.xml



#./{your-driver-name} --mqtt.host=ninjasphere.local [Stop the driver running in the sphere first]
