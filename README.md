# Floway

Сайт-витрина школы флористики «Floway»: курсы, мастер-классы, блог, формы заявок. MVP без онлайн-оплаты — заявки обрабатывает менеджер вручную.

Архитектура и текущий статус описаны в плане фазы 1 разработки; см. также раздел «Что дальше» ниже.

## Стек

- **Frontend**: Nuxt 3/4 (Vue 3, Composition API, TypeScript), Tailwind CSS. Публичные страницы — SSR (для SEO), админка — client-side (`routeRules` в `nuxt.config.ts`).
- **Backend**: Go, chi (роутер) + pgx (Postgres), слоистая архитектура `handler → service → repository`, миграции — goose.
- **БД**: PostgreSQL.
- **Локальное окружение**: docker-compose (Postgres + Mailhog для проверки email-уведомлений).

## Структура репозитория

```
floway/
├── docker-compose.yml       # локальная разработка
├── docker-compose.prod.yml  # прод: образы из GHCR, Caddy (авто-HTTPS), без Mailhog
├── Justfile                 # шорткаты для docker compose и ansible
├── .env.example
├── ansible/                 # бутстрап VPS + деплой, см. ansible/README.md
├── frontend/   # Nuxt-приложение (Dockerfile внутри)
└── backend/    # Go-приложение (Dockerfile внутри)
```

## Быстрый старт (docker compose)

Поднимает всё одной командой: Postgres, Mailhog, миграции (одноразовый сервис `migrate`), backend и frontend — каждый в своём Dockerfile.

1. Скопировать переменные окружения:

   ```bash
   cp .env.example .env
   ```

2. Собрать и поднять весь стек:

   ```bash
   docker compose up -d --build
   ```

   - `migrate` применяет миграции goose и завершается; `backend` стартует только после его успешного завершения (`depends_on: condition: service_completed_successfully`).
   - Backend: http://localhost:8080/healthz
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
   docker compose up -d postgres mailhog
   ```

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

## Тесты, линт, форматирование

```bash
cd backend && go build ./... && go vet ./... && gofmt -l . && go test ./...
cd frontend && bun install && bun run lint && bun run format:check && bun run test && bun run build
```

Frontend-линт/форматтер — [oxlint](https://oxc.rs/docs/guide/usage/linter)/[oxfmt](https://oxc.rs/docs/guide/usage/formatter) (Rust-стек oxc, быстрее ESLint/Prettier на порядок). `bun run format` — исправить форматирование на месте. Backend — `gofmt` + `go vet` + [golangci-lint](https://golangci-lint.run/) (см. CI).

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

## Что уже сделано (фаза 1)

- Monorepo, docker-compose (Postgres + Mailhog + backend + frontend + одноразовый `migrate`), Dockerfile для backend и frontend (multi-stage), `.env.example`.
- Go-бэкенд: слоистый каркас (`config`, `httpserver`, `model`, `repository`, `service`), health-check `/healthz`, goose-миграции под все 9 сущностей ТЗ.
- Полный CRUD (repository + service + handler + unit-тесты на моках) для всех контентных сущностей: `faq_item`, `teachers`, `blog_posts`, `masterclasses`, `courses`, `course_blocks` (вложено под `/api/v1/courses/{courseId}/blocks`), `lessons` (вложено под `/api/v1/course-blocks/{blockId}/lessons`).
- `leads`: публичный `POST /api/v1/leads` (форма с сайта) + `GET/PATCH .../status/DELETE` для админки. Статус не редактируется целиком — только через отдельный `PATCH /{id}/status`.
- `admin_users`: repository + service с bcrypt-хэшированием пароля, HTTP CRUD по-прежнему сознательно не выставлен (создание — только через CLI `cmd/create-admin`).
- **Auth**: JWT в httpOnly-куке (`internal/auth.TokenManager`), `POST /api/v1/admin/login`/`logout`, `GET /api/v1/admin/me`, chi-мидлвара `requireAdmin` защищает все мутирующие запросы (`POST`/`PUT`/`PATCH`/`DELETE`) во всех сущностях + `list`/`status`/`delete` у `leads`. Чтение и публичный `POST /api/v1/leads` остаются открытыми. CORS с `credentials` под `FRONTEND_ORIGIN`.
- **Админ-панель (UI)**: `/admin/login`, guard-middleware `admin-auth`, layout с навигацией/logout. Полный набор CRUD-экранов на `useApi`/`useAdminAuth`/`useAdminResource`:
  - `teachers`, `blog-posts`, `masterclasses`, `faq` — плоские CRUD по образцу `teachers`
  - `courses` → `courses/[courseId]/blocks` → `course-blocks/[blockId]/lessons` — трёхуровневая вложенная навигация (ссылки в таблицах, без пунктов меню)
  - `leads` — только список, смена статуса (`select` → `PATCH .../status`) и удаление; создания нет — это публичная форма с сайта
- Nuxt-приложение: маршруты-заглушки под карту сайта, `routeRules` (SSR для публичных страниц, SPA для `/admin`), Tailwind, composable `usePhoneMask` (маска `+7 (000) 000 00 00`) с тестами на vitest.
- **CI**: GitHub Actions — lint/build/test для backend (Go) и frontend (Bun + oxlint/oxfmt), плюс сборка обоих Docker-образов. Деплой в пайплайн пока не подключён.
- **Инфраструктура**: Ansible-плейбуки для бутстрапа VPS и деплоя, `docker-compose.prod.yml` с Caddy (авто-HTTPS) — см. `ansible/`.

## Что дальше

- **Header/Footer и вёрстка страниц по дизайну** — заблокировано: нет доступа к Figma MCP (нужна авторизация через `/mcp` в интерактивной сессии). Пока используются нейтральные заглушки и плейсхолдер-токены в `frontend/app/assets/css/main.css`.
- Email/Telegram-уведомления при создании лида (`SMTP_*`/`TELEGRAM_*` в `.env.example` уже заведены, интеграции ещё нет).
- Публичный frontend пока не обращается к backend API вообще (кроме админки) — интеграция страниц с реальными данными ещё впереди.
- Блок отзывов на главной — дизайн не предоставлен, оставлен как TODO.
- Деплой в CI: сборка+пуш образов в GHCR по мерджу в `main`, затем `ansible-playbook playbooks/deploy.yml` (или упрощённый `ssh` + `docker compose pull && up -d`) — см. план в `ansible/README.md`.

## Известные технические заметки

- Все `List()`-методы репозиториев возвращают инициализированный `[]T{}`, а не `nil`-слайс — иначе Go сериализует пустой список как JSON `null`, что ломает фронтенд-код вида `items.length`. `useAdminResource`/`leads`-страница дополнительно подстрахованы через `?? []` на случай null из любого источника.
