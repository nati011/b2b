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
    --migration_file_dir /app/migration \
    --smtp $SMTP \
    --email $EMAIL \
    --keycloak_client_secret $KEYCLOAK_CLIENT_SECRET \
    --env $ENV \
    --base_url $BASE_URL \
    --frontend_base_url $FRONTEND_URL \
    --min_compatible_client_version $CLIENT_VERSION

EXPOSE 8080