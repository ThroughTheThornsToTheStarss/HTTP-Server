# HTTP Server (amoCRM → Unisender)

Микросервис на Go для хранения аккаунтов/интеграций и фоновой обработки событий по контактам из **amoCRM** через очередь **beanstalkd**.

## Стек и компоненты
- Go (HTTP сервер + воркеры)
- MySQL (хранение данных)
- beanstalkd (очередь задач, tube: `sync_contacts`)
- gRPC (отключение аккаунта)
- Docker Compose (поднимает всё локально)
- phpMyAdmin (смотреть БД)

## Конфигурация (.env)

Создай файл `.env` в корне проекта. 
Пример:

```env
# HTTP / gRPC
HTTP_PORT=8080
GRPC_PORT=8091

# MySQL (в Docker-хост: fullstack-mysql)
db_host=fullstack-mysql
db_port=3306
db_user=app_user
db_password=app_password
db_name=amoCRM_http_server

# Beanstalkd
BEANSTALK_ADDR=beanstalkd:11300

# amoCRM OAuth
AMO_BASE_DOMAIN=example.amocrm.ru
AMO_CLIENT_ID=xxxx
AMO_CLIENT_SECRET=xxxx
AMO_REDIRECT_URI=http://localhost:8080/amo/oauth/callback

Запуск (Docker)
Сборка и старт:
docker compose up -d --build
docker compose ps
Логи:
docker compose logs -f app
docker compose logs -f workers
Остановить с удалением volume БД:
docker compose down -v

Доступные сервисы
HTTP API: http://localhost:8080
gRPC: localhost:8091
phpMyAdmin: http://localhost:9090

проверка, что всё работает

Создать аккаунт

curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"id":1,"referer":"example.amocrm.ru","is_active":true}'

обновление контакта

curl -X POST http://localhost:8080/webhooks/contacts \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "account_id=1" \
  --data-urlencode "contacts[update][0][id]=111" \
  --data-urlencode "contacts[update][0][name]=Bob"
  
Проверить результат
посмотри логи воркеров: docker compose logs -f workers
открой phpMyAdmin и проверь таблицу sync_histories

HTTP API 

Список аккаунтов:
curl http://localhost:8080/accounts

Отключить аккаунт (is_active=false):
curl -X DELETE "http://localhost:8080/accounts?account_id=1" -i

Создать/обновить интеграцию:

curl -X POST http://localhost:8080/integrations \
  -H "Content-Type: application/json" \
  -d '{"account_id":1,"client_id":"cid","secret_key":"sk","redirect_url":"http://localhost:8080/amo/oauth/callback"}'

Сохранить Unisender key и поставить initial_sync job:
curl -X POST http://localhost:8080/integrations/unisender \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "account_id=1" \
  --data-urlencode "unisender_key=YOUR_KEY"

amoCRM OAuth
GET /amo/auth/start — редирект на авторизацию amoCRM
GET /amo/oauth/callback?code=...&referer=... — обмен кода на токены и сохранение аккаунта

Очередь и воркеры
Tube: sync_contacts
Виды задач (kind):
initial_sync
webhook_upsert
webhook_delete

В Docker воркеры стартуют автоматически.
Режимы:
--mode all — все задачи
--mode webhook — только webhook-задачи
--mode init — только initial_sync
