run:
	cp config-dev.json config.json
	DEBUG=1 go run main.go

deploy:
	cp config-prod.json config.json
	gcloud app deploy --promote