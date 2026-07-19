# steam-collector

Сервис собирает страницы предметов Steam Market, приводит данные к domain-модели и публикует snapshot-ы в Kafka-compatible broker.

## Структура

```text
cmd/steam-collector         composition root
internal/config             чтение env-настроек
internal/domain/market      чистые domain-типы
internal/usecase/collector  сценарий сбора и публикации snapshot-ов
internal/adapters/steam     Steam Market HTTP adapter
internal/adapters/kafka     Kafka publisher adapter на franz-go
internal/adapters/scheduler запуск usecase по расписанию
```

## Конфигурация

Шаблон переменных окружения лежит в `.env.example`.

Пока сервис читает только реальные env-переменные через `envconfig`; `.env` файл автоматически не загружается. Позже добавим Makefile или loader для локального запуска.

Минимально важные параметры:

```text
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=raw.market.prices.steam
SCHEDULER_INTERVAL=10s
SCHEDULER_RATE_LIMIT_PAUSE=10m
SCHEDULER_JITTER=10s
SCHEDULER_PAGE_SIZE=30
```

Steam adapter перед первым запросом может сделать bootstrap-запрос на `STEAM_HOME_URL`, чтобы получить cookies в `http.CookieJar`. Это включается переменной `STEAM_BOOTSTRAP_COOKIES=true`.

При HTTP `429 Too Many Requests` adapter возвращает rate-limit ошибку, а scheduler ждет `SCHEDULER_RATE_LIMIT_PAUSE` перед следующим запросом.

## Инфраструктура

Локально поднимаем Kafka-compatible broker Redpanda и UI:

```powershell
docker compose up -d
```

UI будет доступен по адресу:

```text
http://localhost:8080
```

## Проверка

```powershell
go test ./...
```

## Запуск

Для реального запуска нужен Kafka/Redpanda broker, потому что Kafka publisher проверяет produce connection при старте.

```powershell
go run ./cmd/steam-collector
```
