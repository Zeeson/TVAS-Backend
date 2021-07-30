FROM golang:alpine

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp -r /build/html .
RUN cp /build/main .

# Environment Variables
ENV DB_HOST="127.0.0.1" \
    DB_DRIVER="postgres" \
    DB_USER="dbadmin" \
    DB_PASSWORD="password" \
    DB_NAME="truvest_identity_management" \
    DB_PORT=5432 \
    APP_PROTOCOL="https" \
    APP_HOST="localhost" \
    APP_PORT=9191 \
    SMTP_HOST="smtp.gmail.com" \
    SMTP_PORT=587 \
    API_SECRET="9c5b761a9bf8f74da84c57101669b8042b860ca5871a15c4e729758ad1dee58d562ddf8fa303ada1f3e0278fdff4af9a72f57ce0514be98f6a959daabda926f5" \
    SYSTEM_EMAIL="system@gmail.com" \
    SYSTEM_EMAIL_PASSWORD="change_me" \
    SYSTEM_ADMIN_USERNAME="admin" \
    SYSTEM_ADMIN_FIRSTNAME="admin" \
    SYSTEM_ADMIN_LASTNAME="user" \
    SYSTEM_ADMIN_PASSWORD="password" \
    ALLOWED_ORIGINS="*" \
    ACCESS_TOKEN_EXPIRY_IN_MILLISECOND=3600000

# Export necessary port
EXPOSE 9191

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait
RUN chmod +x /wait

# Command to run when starting the container
CMD /wait && /dist/main