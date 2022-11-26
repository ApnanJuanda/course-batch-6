# Get Golang image from docker hub
FROM golang:1.19

# Define the working directory
WORKDIR /app

# Copy all files to container
COPY . .

# Install all dependencies
RUN go mod tidy

# Build the go project to binary
RUN go build -o /app/exercise-api app/main.go

# Expose port 1234, then localhost can access the container 
EXPOSE 1234

# Run the binary build of go project
CMD [ "/app/exercise-api" ]