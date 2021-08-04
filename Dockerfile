# Initialize image and work directory
FROM golang:1.16-alpine
WORKDIR /app

# Copy app files
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY *.html ./
COPY *.env ./

# Run mod and build commands
RUN go mod download
RUN go build -o /leafylink

# Expose web port and run listener
EXPOSE 8080
CMD [ "/leafylink" ]