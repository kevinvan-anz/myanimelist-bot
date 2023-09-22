# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy your MAL application source code into the container
COPY . .

# Install any required dependencies (if applicable)
# RUN go get -d -v ./...

# Build your MAL application
RUN go build -o mal-app

# Expose the port your MAL application will listen on
EXPOSE 8080

# Define the command to run your MAL application
CMD ["./mal-app"]