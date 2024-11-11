SHELL := /bin/bash

copy-config:
	cp ./.env.example ./.env

ENV_FILE_RELATIVE_PATH = ./.env
ENV_FILE = $(shell echo "$(shell pwd)/$(ENV_FILE_RELATIVE_PATH)")

build-avs-docker:
	docker build -t goplus_avs:latest -f ./Dockerfile .

build-avs:
	cd ./avs && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./avs goplus/avs/cmd & cd ..

run-avs:
	@echo "Using env file: $(ENV_FILE)"
	@echo "Starting AVS..."
	@bash -c 'BLS_KEY_PASSWORD=$(BLS_KEY_PASSWORD) ./avs/avs start -c $(ENV_FILE)'

reg-with-avs:
	@echo "Using env file: $(ENV_FILE)"
	@bash -c ' \
		read -p "Enter ECDSA key store path: " ECDSA_KEY_STORE_PATH; \
		read -sp "Enter ECDSA key store password: " ECDSA_PASSWORD; \
		echo ""; \
		echo "Starting registration process..."; \
		ECDSA_KEY_STORE_PATH=$$ECDSA_KEY_STORE_PATH ECDSA_KEY_PASSWORD=$$ECDSA_PASSWORD BLS_KEY_PASSWORD=$$BLS_KEY_PASSWORD ./avs/avs register-with-avs -c $(ENV_FILE)'

dereg-with-avs:
	@echo "Using env file: $(ENV_FILE)"
	@bash -c ' \
		read -p "Enter ECDSA key store path: " ECDSA_KEY_STORE_PATH; \
		read -sp "Enter ECDSA key store password: " ECDSA_PASSWORD; \
		echo ""; \
		echo "Starting deregistration process..."; \
		ECDSA_KEY_STORE_PATH=$$ECDSA_KEY_STORE_PATH ECDSA_KEY_PASSWORD=$$ECDSA_PASSWORD ./avs/avs deregister-with-avs -c $(ENV_FILE)'

run-avs-docker:
	@echo "Using env file: $(ENV_FILE) ";
	export API_PORT=$(shell grep API_PORT $(ENV_FILE) | cut -d '=' -f 2) && envsubst < ./prometheus-template.yml > ./prometheus.yml
	@sudo bash -c 'CONFIG_FILE_PATH=$(ENV_FILE) BLS_KEY_PASSWORD=$(BLS_KEY_PASSWORD) docker compose -f ./docker-compose.yml --env-file=$(ENV_FILE) up -d'

stop-avs-in-docker:
	 sudo docker compose -f ./docker-compose.yml down


