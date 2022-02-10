# The base go-image
FROM golang:1.16.3
LABEL name="forum"
LABEL description="Docker"
LABEL authors="mus11110; zhangir11;"
LABEL release-date="05.02.2022"
# Create a directory for the app
RUN mkdir /forum
 
# Copy all files from the current directory to the app directory
COPY . /forum
 
# Set working directory
WORKDIR /forum
 
# Run command as described:
RUN go mod download
# go build will build an executable file named main in the current directory
RUN go build -o main ./cmd/mainProg/main.go 
 
# Run the main executable
CMD [ "/forum/main" ]