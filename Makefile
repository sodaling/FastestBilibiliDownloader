APP_NAME=downloader

build: ## Build the container
	docker build -t $(APP_NAME) .

build-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME) .
run:
	docker run -it  -v ${CURDIR}/download:/download $(APP_NAME)

