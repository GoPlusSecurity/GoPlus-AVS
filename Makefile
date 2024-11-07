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
	sudo bash -c './avs/avs start -c $(ENV_FILE)'

reg-with-avs:
	@echo "Using env file: $(ENV_FILE)"
	sudo bash -c './avs/avs register-with-avs -c $(ENV_FILE)'

dereg-with-avs:
	@echo "Using env file: $(ENV_FILE)"
	sudo bash -c './avs/avs deregister-with-avs -c $(ENV_FILE)'

run-avs-docker:
	@echo "Using env file: $(ENV_FILE)"
	export API_PORT=$(shell grep API_PORT $(ENV_FILE) | cut -d '=' -f 2) && envsubst < ./prometheus-template.yml > ./prometheus.yml
	sudo bash -c 'CONFIG_FILE_PATH=$(ENV_FILE) docker compose -f ./docker-compose.yml --env-file=$(ENV_FILE) up -d'

stop-avs-in-docker:
	 sudo docker compose -f ./docker-compose.yml down


