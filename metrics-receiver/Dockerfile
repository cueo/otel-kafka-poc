FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
RUN go env -w GO111MODULE=on

# Copy the source code. Note the slash at the end, as explained in
ADD . .

# Build
RUN go build -o /receiver

# Run
CMD ["/receiver"]
