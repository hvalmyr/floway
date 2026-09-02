# Floway

Сайт-витрина школы флористики «Floway»: курсы, мастер-классы, блог, формы заявок. MVP без онлайн-оплаты — заявки обрабатывает менеджер вручную.

Публичные страницы работают на реальных данных из backend (не моках); весь редактируемый контент — курсы, мастер-классы, блог, FAQ, преподаватели, тексты/картинки блоков — правится через админку, без деплоя. См. «Что уже сделано» и «Что дальше» ниже.

## Стек

- **Frontend**: Nuxt 3/4 (Vue 3, Composition API, TypeScript), Tailwind CSS. Публичные страницы — SSR (для SEO), админка — client-side (`routeRules` в `nuxt.config.ts`).
- **Backend**: Go, chi (роутер) + pgx (Postgres), слоистая архитектура `handler → service → repository → model`, миграции — goose, структурированное логирование — `log/slog` (JSON).
- **БД**: PostgreSQL.
- **Хранилище картинок**: [Garage](https://garagehq.deuxfleurs.fr/) — self-hosted S3-совместимое объектное хранилище; загрузка из админки (`POST /api/v1/admin/uploads`), отдача через backend (`GET /uploads/{key}`).
- **Локальное окружение**: docker-compose (Postgres + Mailhog для проверки email-уведомлений + Garage).

## Структура репозитория

```
floway/
├── docker-compose.yml       # локальная разработка
├── docker-compose.prod.yml  # прод: образы из GHCR, Caddy (авто-HTTPS), без Mailhog
├── Justfile                 # шорткаты для docker compose и ansible
├── .env.example
├── ansible/                 # бутстрап VPS + деплой, см. ansible/README.md
├── garage/                  # Dockerfile + entrypoint-бутстрап для S3-хранилища картинок
├── frontend/   # Nuxt-приложение (Dockerfile внутри)
└── backend/    # Go-приложение (Dockerfile внутри)
```

## Быстрый старт (docker compose)

Поднимает всё одной командой: Postgres, Mailhog, Garage, миграции (одноразовый сервис `migrate`), backend и frontend — каждый в своём Dockerfile.

1. Скопировать переменные окружения:

   ```bash
   cp .env.example .env
   ```

2. Собрать и поднять весь стек:

   ```bash
   docker compose up -d --build
   ```

   - `migrate` применяет миграции goose и завершается; `backend` стартует только после его успешного завершения (`depends_on: condition: service_completed_successfully`) и после того, как `garage` пройдёт healthcheck.
   - Backend: http://localhost:8080/healthz (liveness), http://localhost:8080/readyz (реальный пинг БД)
   - Frontend: http://localhost:3000
   - Mailhog UI: http://localhost:8025

3. Логи / остановка:

   ```bash
   docker compose logs -f backend frontend
   docker compose down
   ```

Порты (`HTTP_PORT`, `FRONTEND_PORT`, `POSTGRES_PORT`, `MAILHOG_*`) настраиваются в `.env`, если дефолтные заняты другими проектами.

## Локальная разработка без Docker (backend/frontend отдельно)

Полезно для быстрого перезапуска с hot-reload; Postgres и Mailhog всё равно берутся из docker-compose.

1. Поднять только инфраструктуру:

   ```bash
   docker compose up -d postgres mailhog garage
   ```

   Backend падает при старте, если не может достучаться до Postgres или до Garage (fail-fast — `pool.Ping`/`garageClient.EnsureBucket` в `cmd/api/main.go`), так что `garage` тоже нужен даже для локального запуска `go run ./cmd/api` напрямую.

2. Применить миграции backend (goose, без установки — через `go run`):

   ```bash
   cd backend
   go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 -dir migrations postgres "$DATABASE_URL" up
   ```

3. Запустить backend:

   ```bash
   cd backend
   go run ./cmd/api
   # curl http://localhost:8080/healthz
   ```

   > Garage-сервис в `docker-compose.yml` пока не публикует порт на хост — `GARAGE_ENDPOINT` для полностью standalone-запуска (backend вне Docker) настроить пока нечем. Либо гоняй backend тоже через `docker compose up -d --build backend` (сеть контейнеров решает это сама), либо добавь `ports: ["3900:3900"]` в сервис `garage` локально.

4. Запустить frontend:

   ```bash
   cd frontend
   bun install
   bun run dev
   # http://localhost:3000
   ```

## Админка: логин

Публичного эндпоинта для создания админ-аккаунтов нет (сознательно). Первый (и любой следующий) аккаунт создаётся CLI-командой:

```bash
# локально
cd backend && go run ./cmd/create-admin -login admin
# или внутри контейнера
docker compose exec backend ./create-admin -login admin
```

Пароль можно передать флагом `-password`, иначе команда спросит его в stdin.

Логин выдаёт JWT в httpOnly-куке (`POST /api/v1/admin/login` → `{"login": "...", "password": "..."}`), `POST /api/v1/admin/logout` куку сбрасывает, `GET /api/v1/admin/me` возвращает текущего пользователя (используется фронтендом для проверки сессии). Кука привязана к origin из `FRONTEND_ORIGIN` (CORS с `credentials`). Защищены только мутирующие запросы (`POST`/`PUT`/`PATCH`/`DELETE`) — чтение (`GET`) и публичная форма `POST /api/v1/leads` остаются открытыми.

UI-админка живёт на `/admin` (Nuxt, client-side): `/admin/login` — форма входа, `admin-auth` route middleware проверяет сессию через `/me` и редиректит на логин, если её нет.

## Уведомления о заявках (email/Telegram)

Каждый `POST /api/v1/leads` (публичная форма) пытается уведомить менеджера — email через SMTP и/или Telegram. Оба канала независимы и опциональны: канал включается наличием его переменных в env, отсутствие обеих пар — не ошибка, просто в логе при старте `lead notifications disabled: no SMTP or Telegram channel configured`. Сбой отправки (SMTP недоступен, неверный токен бота) никогда не валит саму заявку — она уже сохранена в БД, ошибка уведомления только логируется (`lead notification failed`, см. `internal/service/lead_service.go`).

### Email

`SMTP_HOST`/`SMTP_PORT`/`SMTP_FROM`/`SMTP_USER`/`SMTP_PASSWORD`/`NOTIFY_EMAIL_TO` в `.env` идут напрямую в backend — `docker-compose.yml` их больше не переопределяет. `SMTP_USER`/`SMTP_PASSWORD` опциональны: пусто — `internal/notify/email.go` шлёт через `smtp.SendMail` без auth (подходит для relay, доверяющего по IP/сети, либо для локальной пересылки через `sendmail`-совместимый relay на самом хосте); оба заполнены — используется `smtp.PlainAuth` (нужно для внешних провайдеров с логином, например mail.ru, Yandex).

`smtp.SendMail` умеет только STARTTLS (апгрейд соединения до TLS после обычного handshake), а не implicit TLS — для провайдеров с портом 465 (SSL с самого начала соединения, например mail.ru) это не сработает («connection refused»/обрыв на handshake). Используй порт с STARTTLS (587 у большинства провайдеров, включая mail.ru).

Для безопасного локального тестирования без реальной отправки — направь `.env` на Mailhog (поднимается вместе со всем стеком):

```bash
SMTP_HOST=mailhog
SMTP_PORT=1025
NOTIFY_EMAIL_TO=manager@floway.local
```

и пересобери backend: `docker compose up -d --build backend`. Письма смотреть в Mailhog UI — http://localhost:8025.

В проде (`docker-compose.prod.yml`, без Mailhog) `.env`/секреты деплоя должны сразу указывать на настоящий SMTP-relay.

### Telegram

1. В Telegram открой [@BotFather](https://t.me/BotFather) → `/newbot` → задай имя и username (должен оканчиваться на `bot`). BotFather пришлёт токен вида `123456789:ABCdefGHIjklMNOpqrsTUVwxyz` — это `TELEGRAM_BOT_TOKEN`.
2. Бот не может писать первым — напиши ему (или добавь в группу и напиши там) любое сообщение, иначе `chat_id` неоткуда взять.
3. Узнать `chat_id`:

   ```bash
   curl https://api.telegram.org/bot<TOKEN>/getUpdates
   ```

   В ответе — `result[].message.chat.id`. Для личного чата с ботом это твой `user id`; для группы — отрицательное число.

4. Прописать в `.env`:

   ```bash
   TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
   TELEGRAM_CHAT_ID=987654321
   ```

5. Пересобрать backend и отправить тестовую заявку:

   ```bash
   docker compose up -d --build backend
   curl -X POST http://localhost:8080/api/v1/leads \
     -H "Content-Type: application/json" \
     -d '{"name":"Тест","phone":"+79001112233","contactMethod":"telegram","source":"internet","requestType":"trial_lesson"}'
   ```

   Сообщение должно прийти в тот чат/группу, где боту писали на шаге 2.

Отсутствие `TELEGRAM_BOT_TOKEN`/`TELEGRAM_CHAT_ID` (или пустая строка хотя бы у одного) отключает канал целиком — частичная конфигурация не пытается угадать недостающее.

## Тесты, линт, форматирование

```bash
cd backend && gofmt -l . && go vet ./... && go build ./... && go test ./... -race -cover
cd backend && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.0 && "$(go env GOPATH)/bin/golangci-lint" run --timeout=5m
cd frontend && bun install && bun run lint && bun run format:check && bun run test && bun run build
```

Frontend-линт/форматтер — [oxlint](https://oxc.rs/docs/guide/usage/linter)/[oxfmt](https://oxc.rs/docs/guide/usage/formatter) (Rust-стек oxc, быстрее ESLint/Prettier на порядок). `bun run format` — исправить форматирование на месте. Backend — `gofmt` + `go vet` + [golangci-lint](https://golangci-lint.run/) (запинена версия `v2.12.0`, как в CI; не предустановлен, ставится через `go install`). golangci-lint — не опционально: `errcheck` (часть его дефолтного набора линтеров) ловит баги, которые `go vet`/`go test` пропускают (например, непроверенный `Close()`), и не раз был единственным, что реально останавливало CI.

Фронтенд полностью на [Bun](https://bun.sh) — пакетный менеджер и прод-рантайм. `bun.lock` — единственный лок-файл. Nitro собирается под `preset: "bun"` (`nuxt.config.ts`) — прод-сервер работает на `Bun.serve()`, а не на Node `http`. `frontend/Dockerfile` — оба стейджа на `oven/bun:1-alpine`, запуск через `bun run .output/server/index.mjs`.

### Git-хук перед коммитом

```bash
just install-hooks   # один раз на машину: git config core.hooksPath .githooks
```

`.githooks/pre-commit` прогоняет форматтер + линтер + тесты на каждый коммит, скоуп — по тому, что реально застейджено (`frontend/` и/или `backend/`); форматтер исправляет файлы на месте и дозастейджит их. Это лёгкое подмножество CI (без `go build -race`, `nuxt build`, `golangci-lint` — слишком медленно на каждый коммит), полная проверка всё равно остаётся в CI.

## CI

`.github/workflows/ci.yml` (GitHub Actions) на каждый push в `main` и pull request: lint + build + test backend (Go: `gofmt`, `go vet`, `golangci-lint`, `go test -race`) и frontend (Bun: `oxlint`, `oxfmt --check`, `vitest`, `nuxt build`), плюс сборка обоих Docker-образов (без пуша — деплой пока не в CI, см. `ansible/`).

## Прод-деплой

Бутстрап VPS (Ubuntu 24.04, security-хардненинг + Docker) и деплой самого
приложения (образы из GHCR + Caddy с авто-HTTPS) — через Ansible, см.
[ansible/README.md](ansible/README.md). Шорткаты в [Justfile](Justfile):
`just bootstrap-first-run`, `just deploy` и т.д.

## Что уже сделано

- Monorepo, docker-compose (Postgres + Mailhog + Garage + backend + frontend + одноразовый `migrate`), Dockerfile для backend и frontend (multi-stage), `.env.example`.
- **Backend — полный CRUD** (repository + service + handler + unit/boundary-тесты) для всех контентных сущностей: `faq_item`, `teachers`, `blog_posts`, `masterclasses`, `courses`, `course_blocks` (вложено под `/api/v1/courses/{courseId}/blocks`, принадлежность родителю реально проверяется на уровне SQL, не декоративно), `lessons` (вложено под `/api/v1/course-blocks/{blockId}/lessons`, аналогично), `features`, `about_items`, `social_links`, `page_content` (generic key-value для текстов/картинок, редактируемых без деплоя).
- `leads`: публичный `POST /api/v1/leads` (форма с сайта, rate-limit 3/час/IP) + `GET/PATCH .../status/DELETE` для админки. Каждая заявка триггерит уведомление (email через SMTP и/или Telegram — оба канала независимы и опциональны, включаются наличием `SMTP_HOST`+`NOTIFY_EMAIL_TO` / `TELEGRAM_BOT_TOKEN`+`TELEGRAM_CHAT_ID` в env; ни один не настроен — просто лог `lead notifications disabled` при старте). Best-effort: сбой уведомления не валит создание заявки, только логируется.
- **Auth**: JWT в httpOnly-куке (`internal/auth.TokenManager`, сессия 12ч), `POST /api/v1/admin/login` (rate-limit 5/мин/IP) / `logout` / `GET /api/v1/admin/me`. `logout` реально инвалидирует токен (счётчик `token_version` в БД, сверяется на каждый admin-запрос) — не просто чистит куку на клиенте. Аккаунты создаются только через CLI `cmd/create-admin`, HTTP-эндпоинта для этого нет.
- **Единая таксономия ошибок**: `internal/apperr` (`ErrValidation`/`ErrNotFound`/`ErrConflict`), репозитории транслируют driver-ошибки на границе (`translateNotFound`/`checkDeleted`), хендлеры маппят на HTTP-коды через `writeServiceError`. Внутренние ошибки (500) никогда не утекают клиенту как есть — логируются на сервере с `request_id`, клиенту уходит только этот id для корреляции.
- **Картинки**: загрузка в Garage (S3) из админки, реальный content-type определяется по байтам файла (не по заявленному клиентом заголовку), `X-Content-Type-Options: nosniff` на отдаче, автоматическая очистка старого объекта в Garage при замене картинки в `page_content`.
- **Устойчивость/безопасность**: per-IP rate limiting на login/leads, лимит размера тела запроса (10MB), `/readyz` (реальный пинг БД) отдельно от `/healthz` (liveness), graceful shutdown, fail-fast старт при недоступности БД/Garage, структурированное JSON-логирование каждого запроса (`log/slog`, корреляция по `request_id`).
- **Тесты**: unit на сервисах (моки репозиториев) + boundary-тесты на HTTP-слое (`httptest` + реальные сервисы + фейковые репозитории) — TDD-подход, проверяют реальные граничные случаи (404 вместо утечки 500, чужой родитель в nested-роуте не даёт хайджекнуть запись, неаутентифицированный запрос никогда не видит черновики блога).
- **Админ-панель (UI)**: `/admin/login`, guard-middleware `admin-auth`, layout с навигацией/logout. CRUD-экраны на `useApi`/`useAdminAuth`/`useAdminResource` для всех сущностей выше, включая `page-content` (тексты + загрузка картинок) и `about-items`/`features`.
- **Публичный frontend на реальных данных**: главная, блог (список + страница поста), курсы, мастер-классы, контакты — все берут данные из backend API, не моки.
- **CI**: GitHub Actions — lint/build/test для backend (Go: `gofmt`/`go vet`/`golangci-lint`/`go test -race`) и frontend (Bun: `oxlint`/`oxfmt`/`vitest`/`nuxt build`), плюс сборка обоих Docker-образов. Деплой в пайплайн пока не подключён.
- **Инфраструктура**: Ansible-плейбуки для бутстрапа VPS и деплоя, `docker-compose.prod.yml` с Caddy (авто-HTTPS) — см. `ansible/`.

## Что дальше

- **Блок отзывов на главной** — дизайна карточек пока нет (см. `TODO` в `frontend/app/pages/index.vue`), в макете только заголовок секции.
- **Деплой в CI** — сборка образов уже есть, но пуш в GHCR и прогон `ansible-playbook playbooks/deploy.yml` (или `ssh` + `docker compose pull && up -d`) по мерджу в `main` ещё не подключены; сейчас деплой запускается вручную, см. `ansible/README.md`.
- **page_content: нет API для создания новых ключей** — ключи заводятся только через миграции (см. «Известные технические заметки»); если понадобится добавлять редактируемые текстовые/картиночные поля без деплоя бэкенда, потребуется `POST /page-content`.
- **Локальный standalone-запуск backend без Docker требует ручной правки compose** — Garage не публикует порт на хост (см. предупреждение в разделе «Локальная разработка без Docker» выше).

## Известные технические заметки

- Все `List()`-методы репозиториев возвращают инициализированный `[]T{}`, а не `nil`-слайс — иначе Go сериализует пустой список как JSON `null`, что ломает фронтенд-код вида `items.length`. `useAdminResource`/`leads`-страница дополнительно подстрахованы через `?? []` на случай null из любого источника.
- `page_content`-ключи (`home_hero_title`, `home_hero_image` и т.д.) создаются только миграциями — `PUT /api/v1/page-content/{key}` работает лишь для уже существующего ключа, эндпоинта на создание нового нет (осознанное решение, см. `internal/repository/page_content_repository.go`). Значения-картинки — реальный контент, миграции `00010`/`00012`/`00014` намеренно не трогать при чистке «тестовых» данных; `features`/`about_items`/`social_links`, наоборот, заводятся только через админку — их миграции содержат только схему.
- Узкие per-service интерфейсы репозиториев (каждый `service.go` описывает свой собственный интерфейс, ровно под то, что ему нужно) — устоявшаяся конвенция во всех ~13 сущностях бэкенда, а не абстракция «на будущее».
- `golangci-lint` не входит в git-хук pre-commit (слишком медленно для каждого коммита) — обязательно прогонять вручную перед пушем, см. «Тесты, линт, форматирование» выше; CI без него не проходит.
