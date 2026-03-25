[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)
[![CI Lint](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-lint.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-lint.yml)
[![CI Build](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-build.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-build.yml)
[![CI Tests](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-tests.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-tests.yml)
[![CI Coverage](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-coverage.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-X1ag/actions/workflows/ci-coverage.yml)

# Room Booking Service

Тестовое задание на backend-стажировку: сервис бронирования переговорок с ролями `admin` и `user`.

Сервис позволяет:
- администратору создавать переговорки;
- администратору один раз задавать расписание доступности переговорки;
- пользователю получать доступные слоты по переговорке и дате;
- пользователю создавать и отменять свои брони;
- пользователю опционально создавать ссылку на конференцию при бронировании;
- администратору смотреть список всех броней с пагинацией;
- пользователю смотреть список своих будущих броней.

Все даты и время хранятся и отдаются в UTC.

## Содержание

- [Запуск проекта](#запуск-проекта)
- [Переменные окружения](#переменные-окружения)
- [Что реализовано](#что-реализовано)
- [Стек](#стек)
- [CI И Quality Gates](#ci-и-quality-gates)
- [Архитектура](#архитектура)
- [Принятые решения](#принятые-решения)
- [Модель данных](#модель-данных)
- [Smoke test](#smoke-test)
- [Тесты](#тесты)
- [Структура проекта](#структура-проекта-1)
- [Дополнительные задания](#сознательно-принятые-упрощения)

## Запуск проекта

### 1. Подготовка env

```bash
cp .env.example .env
```

Этот шаг опционален.
`docker-compose.yaml` уже содержит дефолтные значения для всех используемых переменных, поэтому проект можно запускать без дополнительной настройки.
`.env` нужен только если хочется переопределить значения под свою среду.

### 2 Запуск через Docker Compose или Makefile

```bash
docker compose up --build
```

Но также можно запустить проект с помощью команды Makefile
```bash
make up
```

После запуска сервис доступен по адресу:

```text
http://localhost:8080
```

### 3. Полный сброс локальной БД

Если нужно начать с чистого состояния:

```bash
docker compose down -v
docker compose up --build
```

## Переменные окружения

Основные переменные уже описаны в [`.env.example`](./.env.example).
При этом у приложения и `docker-compose.yaml` есть встроенные дефолтные значения, так что `.env` не является обязательным.

Самые важные:
- `HTTP_PORT`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`
- `MIGRATIONS_PATH`
- `JWT_SECRET`
- `DUMMY_ADMIN_ID`
- `DUMMY_USER_ID`

## Что реализовано

Обязательная часть API:
- `POST /dummyLogin`
- `GET /_info`
- `GET /rooms/list`
- `POST /rooms/create`
- `POST /rooms/{roomId}/schedule/create`
- `GET /rooms/{roomId}/slots/list`
- `POST /bookings/create`
- `POST /bookings/{bookingId}/cancel`
- `GET /bookings/my`
- `GET /bookings/list`

## Стек

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- Docker Compose
- JWT

## CI И Quality Gates

В GitHub Actions настроены отдельные workflow для проверки качества:
- `CI Lint` запускает `golangci-lint` по конфигурации из `.golangci.yaml`;
- `CI Build` проверяет, что проект собирается;
- `CI Tests` запускает тесты;
- `CI Coverage` считает покрытие и проверяет порог не ниже `40%`.

Локально для линтера и тестов можно использовать `Makefile`:
- `make lint`
- `make test-cover`

## Архитектура

Проект построен по слоям:
- `handlers` отвечают за HTTP, парсинг запроса и маппинг ошибок в HTTP-коды;
- `usecase` содержит бизнес-логику;
- `repository` отвечает только за работу с PostgreSQL;
- `middleware` отвечает за JWT-аутентификацию и проверку роли.

Основной поток зависимости:

`handler -> usecase -> repository`

Это позволяет не размазывать бизнес-логику по HTTP-слою и проще тестировать доменные сценарии.

## Принятые решения

### Структура проекта

Отдельно для себя я принял решение строить проект по layered-подходу. Я посмотрел несколько вариантов организации backend-проекта и остановился на разделении на `handlers`, `usecase` и `repository`. Такой подход показался мне самым понятным для этого задания, потому что он помогает не смешивать HTTP-логику, бизнес-правила и работу с базой данных.

### Генерация слотов

Выбрана ленивая стратегия генерации слотов:
- расписание хранится как правило доступности переговорки;
- при запросе `GET /rooms/{roomId}/slots/list?date=...` сервис:
- 1) проверяет, что переговорка существует;
- 2) получает расписание переговорки;
- 3) проверяет, применимо ли расписание к указанной дате;
- 4) генерирует 30-минутные слоты на эту дату;
- 5) сохраняет их в БД;
- 6) затем возвращает только свободные слоты.

Почему выбран этот подход:
- в задании самый нагруженный endpoint это получение доступных слотов;
- пользователи почти всегда работают в пределах ближайших 7 дней;
- нет смысла заранее генерировать большой объем слотов на далекое будущее;
- слоты при этом получают стабильные UUID, что позволяет бронировать их по `slotId`.

Повторная генерация на ту же дату безопасна(операция идемпотентна):
- на таблице `slots` есть уникальность по `(room_id, start_at, end_at)`;
- вставка выполняется через `ON CONFLICT DO NOTHING`.

### Защита от двойного бронирования

На уровне бизнес-логики перед созданием брони я проверяю, что у слота нет активной брони.

Дополнительно защита есть на уровне БД:
- создан уникальный partial index `idx_bookings_active_slot`
- он разрешает только одну запись со `status = 'active'` на один `slot_id`

Это защищает от гонок при параллельных запросах на один и тот же слот.

### Опциональная ссылка на конференцию

Чтобы не зашивать генерацию ссылки прямо в `booking/usecase`, я вынес это в маленький интерфейс `ConferenceService`.
У него только одна ответственность: вернуть ссылку на конференцию для `slotId` и `userId`.

В приложении подключена мок-реализация этого сервиса:
- она имитирует внешний conference service;
- генерирует ссылку формата `https://meet.example.com/{uuid}`;
- не требует отдельного HTTP-сервиса и не усложняет локальный запуск.

Принятые решения по сбоям:
- если conference service недоступен или вернул ошибку, бронь не создается вообще
- если conference service успешно ответил, но запись брони в БД потом упала, клиент получает ошибку, а созданная внешняя ссылка может остаться висячей
- отдельную компенсацию удаления ссылки я не добавлял, потому что для тестового задания это уже заметно усложняет решение.

Почему выбрал именно так:
- usecase не знает, как именно строится ссылка, а только оркестрирует шаги;
- если пользователь явно попросил conference link, система не создает "успешную" бронь без ссылки.
### Запуск линтера

Добавил проверку кода линтеров в Github Actions

Почему выбрал именно так:
- В крупных компаниях очень важно чтобы люди писали похожий и чистый код
- Это уменьшает кол-во ошибок и приводит к единому code-style

### Идемпотентная отмена брони

Отмена брони реализована как идемпотентная операция:
- если бронь уже `cancelled`, повторный вызов не должен приводить к ошибке;
- клиент получает актуальное состояние брони.

### Dummy users

Для `dummyLogin` используются фиксированные UUID:
- `admin`: `11111111-1111-1111-1111-111111111111`
- `user`: `22222222-2222-2222-2222-222222222222`

Эти пользователи добавляются автоматически миграцией `000002_seed_dummy_users.up.sql`, когда приложение стартует и применяет миграции.

Команда `make seed` нужна для загрузки дополнительных тестовых данных из `seed/test_data.sql`:
- тестовых переговорок;
- расписаний для них.

Команду можно запускать так:
```bash
make seed
```

Это решение пришлось принять отдельно, потому что сам `dummyLogin` выдает JWT с фиксированными UUID, но в базе данных таких пользователей изначально нет. Из-за этого создание брони падало по foreign key на таблицу `users`. Отдельная seed-миграция сделала сценарий с тестовыми токенами полностью рабочим и согласованным с моделью данных.

## Модель данных

Основные таблицы:
- `users`
- `rooms`
- `schedules`
- `schedule_days`
- `slots`
- `bookings`

Ключевые ограничения:
- одно расписание на одну переговорку;
- слот не может пересекаться сам с собой за счет фиксированного разбиения по 30 минут;
- только одна активная бронь на слот;
- статус брони: `active` или `cancelled`.


## Smoke test

Минимальный сценарий ручной проверки(все прошло успешно):
- получить токен `admin` через `/dummyLogin`
- создать переговорку
- создать расписание
- получить токен `user` через `/dummyLogin`
- получить список доступных слотов
- создать бронь
- повторно попытаться забронировать тот же слот и получить `409`
- получить `/bookings/my`
- отменить бронь
- повторно вызвать `cancel`
- получить `/bookings/list` под `admin`

## Тесты

E2E-тесты сами поднимают только HTTP server через `httptest`.
PostgreSQL нужно поднять отдельно перед запуском тестов, например через Docker Compose или `Makefile`.

В GitHub Actions отдельно настроен workflow на покрытие.
Он считает `coverage.out` и падает, если общее покрытие становится ниже `40%`.

Локальный запуск тестов:

```bash
go test ./...
```

но лучше запускать тесты через Makefile:
```bash
make test-cover
```

Вывод команды `make test-cover`:
```text
$ make test-cover
docker compose up -d --wait postgres
[+] up 1/1
 ✔ Container room-booking-postgres Healthy                                                                        5.9s
        test-backend-1-X1ag/cmd/room-booking            coverage: 0.0% of statements
        test-backend-1-X1ag/internal/app                coverage: 0.0% of statements
ok      test-backend-1-X1ag/internal/auth       0.584s  coverage: 1.7% of statements in ./...
ok      test-backend-1-X1ag/internal/booking    0.800s  coverage: 5.7% of statements in ./...
?       test-backend-1-X1ag/internal/clock      [no test files]
        test-backend-1-X1ag/internal/conference         coverage: 0.0% of statements
        test-backend-1-X1ag/internal/config             coverage: 0.0% of statements
?       test-backend-1-X1ag/internal/http/dto   [no test files]
        test-backend-1-X1ag/internal/http/handlers              coverage: 0.0% of statements
        test-backend-1-X1ag/internal/http/middleware            coverage: 0.0% of statements
        test-backend-1-X1ag/internal/http/response              coverage: 0.0% of statements
        test-backend-1-X1ag/internal/logger             coverage: 0.0% of statements
        test-backend-1-X1ag/internal/repository/postgres                coverage: 0.0% of statements
ok      test-backend-1-X1ag/internal/room       0.498s  coverage: 3.1% of statements in ./...
ok      test-backend-1-X1ag/internal/schedule   0.921s  coverage: 4.9% of statements in ./...
ok      test-backend-1-X1ag/internal/slot       0.709s  coverage: 4.9% of statements in ./...
ok      test-backend-1-X1ag/internal/user       1.528s  coverage: 2.9% of statements in ./...
ok      test-backend-1-X1ag/tests/e2e   1.115s  coverage: 47.5% of statements in ./...
total:                                                                          (statements)            56.6%
```

## Структура проекта

Ключевые директории:
- `cmd/room-booking` — точка входа приложения
- `internal/app` - запуск и инициализация приложения
- `internal/http/handlers` — HTTP-обработчики
- `internal/http/middleware` — JWT/auth/role middleware
- `internal/room` — логика переговорок
- `internal/schedule` — логика расписаний
- `internal/slot` — генерация и получение слотов
- `internal/booking` — создание, отмена и получение броней
- `internal/repository/postgres` — PostgreSQL-репозитории
- `migrations` — SQL-миграции

## Дополнительные задания 

Я не успел реализовать все дополнительные задания. 
Список сделанных доп. заданий:
- Регистрацию и логин по email и паролю
- CI в Github Actions 
- Конфигурация линтера (.golangci.yaml)
- Makefile
- Опциональное создание ссылки на конференцию при бронировании.

Не реализованы дополнительные задачи:
- нагрузочное тестирование;
- swagger codegen.

Если бы времени было больше, следующим шагом я бы:
- создал openapi спецификацию 
- сделал нагрузочное

Все время было потрачено на основные эндпоинты и только после того как я убедился, что основной функционал работает я пошел делать доп. задания. 
