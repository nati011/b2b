FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/

CMD ./main \
    --port $PORT \
    --keycloak_base_url $KEYCLOAK_BASE_URL \
    --keycloak_user_name $KEYCLOAK_USER_NAME \
    --keycloak_password $KEYCLOAK_PASSWORD \
    --keycloak_realm $KEYCLOAK_REALM \
    --keycloak_client_id $KEYCLOAK_CLIENT_ID \
    --keycloak_application_realm $KEYCLOAK_APPLICATION_REALM \
    --db $DB_URL \
    --migration_file_dir /app \
    --smtp $SMTP \
    --email $EMAIL \
    --keycloak_client_secret $KEYCLOAK_CLIENT_SECRET

EXPOSE 8080