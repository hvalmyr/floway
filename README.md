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
├── docker-compose.yml
├── .env.example
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
   npm install
   npm run dev
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

Логин выдаёт JWT в httpOnly-куке (`POST /api/v1/admin/login` → `{"login": "...", "password": "..."}`), `POST /api/v1/admin/logout` куку сбрасывает. Кука привязана к origin из `FRONTEND_ORIGIN` (CORS с `credentials`). Защищены только мутирующие запросы (`POST`/`PUT`/`PATCH`/`DELETE`) — чтение (`GET`) и публичная форма `POST /api/v1/leads` остаются открытыми.

## Тесты

```bash
cd backend && go test ./...
cd frontend && npm run test
```

## Что уже сделано (фаза 1)

- Monorepo, docker-compose (Postgres + Mailhog + backend + frontend + одноразовый `migrate`), Dockerfile для backend и frontend (multi-stage), `.env.example`.
- Go-бэкенд: слоистый каркас (`config`, `httpserver`, `model`, `repository`, `service`), health-check `/healthz`, goose-миграции под все 9 сущностей ТЗ.
- Полный CRUD (repository + service + handler + unit-тесты на моках) для всех контентных сущностей: `faq_item`, `teachers`, `blog_posts`, `masterclasses`, `courses`, `course_blocks` (вложено под `/api/v1/courses/{courseId}/blocks`), `lessons` (вложено под `/api/v1/course-blocks/{blockId}/lessons`).
- `leads`: публичный `POST /api/v1/leads` (форма с сайта) + `GET/PATCH .../status/DELETE` для админки. Статус не редактируется целиком — только через отдельный `PATCH /{id}/status`.
- `admin_users`: repository + service с bcrypt-хэшированием пароля, HTTP CRUD по-прежнему сознательно не выставлен (создание — только через CLI `cmd/create-admin`).
- **Auth**: JWT в httpOnly-куке (`internal/auth.TokenManager`), `POST /api/v1/admin/login`/`logout`, chi-мидлвара `requireAdmin` защищает все мутирующие запросы (`POST`/`PUT`/`PATCH`/`DELETE`) во всех сущностях + `list`/`status`/`delete` у `leads`. Чтение и публичный `POST /api/v1/leads` остаются открытыми. CORS с `credentials` под `FRONTEND_ORIGIN`.
- Nuxt-приложение: маршруты-заглушки под карту сайта, `routeRules` (SSR для публичных страниц, SPA для `/admin`), Tailwind, composable `usePhoneMask` (маска `+7 (000) 000 00 00`) с тестами на vitest.

## Что дальше

- **Header/Footer и вёрстка страниц по дизайну** — заблокировано: нет доступа к Figma MCP (нужна авторизация через `/mcp` в интерактивной сессии). Пока используются нейтральные заглушки и плейсхолдер-токены в `frontend/app/assets/css/main.css`.
- Email/Telegram-уведомления при создании лида (`SMTP_*`/`TELEGRAM_*` в `.env.example` уже заведены, интеграции ещё нет).
- Админ-панель (UI поверх backend CRUD + auth-логин).
- Frontend пока не обращается к backend API вообще — интеграция страниц с реальными данными ещё впереди.
- Блок отзывов на главной — дизайн не предоставлен, оставлен как TODO.
