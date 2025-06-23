NAME = socialNetwork

IN-PORT=8080

# Port to connect to the webapp
OUT-PORT=9090

all: get-keys run-frontend run-backend

# docker:
# 	@echo "\033[1m\033[92mGetting docker ready for first use\nPlease Wait...\033[0m"
# 	@curl -fsSL https://get.docker.com/rootless 2>/dev/null | sh >/dev/null 2>&1
# 	@echo "\033[1m\033[92mCopy paste the following command to start docker in rootless mode:\033[0m"
# 	export PATH=$(HOME)/bin:$$PATH
# 	export DOCKER_HOST=unix://$(XDG_RUNTIME_DIR)/docker.sock

# $(IMAGE_NAME): build run status
# 	@echo "\033[1m\033[92mListing what's inside the container\033[0m"
# 	@docker exec -it $(NAME) /bin/bash -c "ls -l && exit"
# 	@echo "\033[1m\033[92mWebapp running at: \033[0m\033[92mhttps://localhost:$(OUT-PORT)\033[0m"

# build:
# 	@echo "\033[1m\033[92mBuilding Docker image...\033[0m"
# 	@docker image build -f Dockerfile -t $(IMAGE_NAME) . 

# run:
# 	@exec "clear"
# 	@echo "\033[1m\033[92mRunning the Container\033[0m"
# 	docker container run -p $(OUT-PORT):$(IN-PORT) --detach --name $(NAME) $(IMAGE_NAME)

run-frontend:
	@cd ./front-end && npm install && npm run dev &

run-backend:
	@cd ./backend && go run .

get-keys:
	@if [ ! -f ./private/.env ]; then \
		echo "Keys file not found, generating..."; \
		cd ./backend && python3 init.py; \
	else \
		echo "Keys file already exists, skipping generation."; \
	fi

clean:
	rm -fr ./private

dockerClean:
	@echo "\033[1m\033[92mRemoving all containers\033[0m"
	@if [ -n "$$(docker container ls -aq)" ]; then \
		 docker container rm -f $$(docker container ls -aq) > /dev/null 2>&1; \
	else \
		echo "\033[1m\033[93mNo containers to remove\033[0m"; \
	fi

	@echo "\033[1m\033[92mRemoving all unused images\033[0m"
	@if [ -n "$$(docker image ls -aq)" ]; then \
		 docker image prune -a -f > /dev/null 2>&1; \
	else \
		echo "\033[1m\033[93mNo images to prune\033[0m"; \
	fi

fclean: dockerClean clean
	

re: clean all

.PHONY: all run-frontend run-backend

