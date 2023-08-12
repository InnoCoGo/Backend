# InnoCoGo
The carpooling/ridesharing application designed for residents of city Innopolis to organize economical and fast trips, find fellow travelers, and communicate seamlessly.

### Demo
https://github.com/InnoCoTravel/Backend/assets/47076924/6b2a4b26-b279-4a9e-b82e-f79bb282e916

#### Function list of MVP *(Telegram-webapp / Website)* 
1. Telegram authentication
2. Trip filtration
3. Trip creation/deletion and joining/leaving
4. Accepting/rejecting join request to the trip
5. Notification
6. Automatic translation of trip description

#### Function list of Final Product (Without Telegram)

7. The list above+
8. User rating/commenting
9. Automatic group chat setup
10. Voice message
11. Fraud message detection ML
12. Distributed deployment (through kafka global message queue, unified message delivery, the system can be expanded horizontally)

### Backend
The backend part of the application satisfies all the functions listed above, with the exception of "ML for detecting fraud messages".

Backend Technologies and Frameworks

- Gin web framework
- Long connection WebSocket (gorilla)
- viper for configuration management
- Makefile writing
- PostgreSQL
- Kafka
- Docker, docker compose


### Frontend
Currently [frontend](https://github.com/InnoCoTravel/Frontend ) the project includes only the MVP. The final product of this project is under development.

## How to Run
#### LibreTranslate
1. Clone https://github.com/LibreTranslate/LibreTranslate/tree/main into a directory
2. Change the last line in run.sh from:
```
docker run -ti --rm -p $LT_PORT:$LT_PORT $DB_VOLUME -v lt-local:/home/libretranslate/.local libretranslate/libretranslate ${ARGS[@]}`
```
to:
```
docker run --name libretranslate -ti --rm --network=innocogo -p $LT_PORT:$LT_PORT $DB_VOLUME -v lt-local:/home/libretranslate/.local libretranslate/libretranslate ${ARGS[@]}
```
3. Create a shared network for docker containers using the command
```docker network create innocogo```
4. Run LibreTranslate only in Russian and English using the command
``(./run.sh --load-only=ru,en )&``
#### Telegram webhook
5. All self-signed certificates for the IP address are generated via:
```
openssl req -newkey rsa:2048 -sha256 -nodes -keyout PRIVATE.key -x509 -days 365 -out PUBLIC.pem -subj "/C=US/ST=State/L=City/O=pinkyhi/CN={HERE_IP}"
```
6. The endpoint for webhooks is registered in the telegram via:
```
curl -F "url=https://{HERE_IP}:8443/telegram_endpoint" -F "certificate=@PUBLIC.pem" https://api.telegram.org/bot{HERE_BOT_TOKEN}/setWebhook?secret_token={HERE_TG_SECRET_TOKEN_FROM_docker-compose.yml}
```
#### Docker
7. Fill the environment variables below

```
KAFKA_TOPIC=
KAFKA_HOSTS=localhost:9092

TG_BOT_TOKEN=
CERT_FILE=path/PUBLIC.pem
PKEY_FILE=path/PRIVATE.key
HOST=0.0.0.0
TG_SECRET_TOKEN=
PORT=
BACKEND_SECRET_TOKEN=
BACKEND_URL=
PERSISTENT_FOLDER=

DB_PASSWORD=
DB_HOST=
DB_PORT=5432
BOT_TOKEN=
TG_BOT_URL=
BACKEND_SECRET_TOKEN=
TRANSLATE_URL=http://libretranslate:5000
TRANSLATE_API_KEY=

POSTGRES_PASSWORD=
```
8. copy the docker-compose.yml and migrate_up.sh (ensure to modify its permissions to make it executable) files into the directory
9. run:
```
docker compose -f docker-compose.yml down -v
docker compose stop
docker compose rm -f
docker compose pull
docker compose up -d

./migrate_up.sh
```
