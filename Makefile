NAME = socialNetwork

IN-PORT=8080

# Port to connect to the webapp
OUT-PORT=9090

all: get-keys run-frontend run-backend
	@echo "\033[1m\033[92mAll services are running!\033[0m"
	@echo "\033[1m\033[92mWebapp running at: \033[0m\033[92mhttps://localhost:8080\033[0m"
	@echo "\033[1m\033[92m--------------------------------------------------------\033[0m"

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
	@cd ./front-end && npm install > '/dev/null' 2>&1 && npm run dev > '/dev/null' 2>&1 &
	@echo "\033[1m\033[92mFrontend service is running!\033[0m"

run-backend:
	@cd ./backend && go run .  > '/dev/null' 2>&1 &
	@echo "\033[1m\033[92mBackend service is running!\033[0m"

get-keys:
	@if [ ! -f ./private/.env ]; then \
		echo "Keys file not found, generating..."; \
		cd ./backend && python3 init.py; \
		clear ;\
		echo "\033[1m\033[92m--------------------------------------------------------\033[0m"; \
		echo "\033[1m\033[92mKeys file generated successfully!\033[0m"; \
	else \
		echo "\033[1m\033[93m--------------------------------------------------------\033[0m"; \
		echo "\033[1m\033[93mKeys file already exists, skipping generation.\033[0m"; \
	fi


clean:
	@if fuser -n tcp 3000 > /dev/null 2>&1; then \
		echo "Port 3000 is in use. Cleaning..."; \
		fuser -n tcp 3000 -k 2>/dev/null || true; \
		rm -rf ./private; \
	fi

	@if fuser -n tcp 8080 > /dev/null 2>&1; then \
		echo "Port 8080 is in use. Cleaning..."; \
		fuser -n tcp 8080 -k 2>/dev/null || true; \
		rm -rf ./private; \
	fi
	@rm -fr ./private
	clear

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

