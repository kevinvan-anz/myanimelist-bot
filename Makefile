# Variables
IMAGE_NAME := my-mal-app
CONTAINER_NAME := my-mal-container
PORT := 8080

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the Docker container
run:
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

# Stop and remove the Docker container
stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

# Clean up (remove Docker image)
clean:
	docker rmi $(IMAGE_NAME)

# Run the MAL application in the Docker container
start:
	docker exec -it $(CONTAINER_NAME) ./map-app

# Build, run, and start the MAL application in one command
add: build run start
