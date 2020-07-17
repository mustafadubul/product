.PHONY: \
		build \
		test \
		start \
		stop 

SHELL := /bin/bash

# Run Unit tests
test:
	docker-compose -f docker-compose.test.yml up --build

# Builds the Docker container
build:
	docker build -t mustafadubul/product:latest .

# Starts an instance of the product service 
start:
	docker-compose up 

# Stops the product service instance
stop:
	docker-compose down


