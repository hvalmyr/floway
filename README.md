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
├── frontend/   # Nuxt-приложение
└── backend/    # Go-приложение
```

## Быстрый старт

1. Скопировать переменные окружения:

   ```bash
   cp .env.example .env
   ```

2. Поднять Postgres и Mailhog:

   ```bash
   docker compose up -d
   ```

3. Применить миграции backend (goose, без установки — через `go run`):

   ```bash
   cd backend
   go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 -dir migrations postgres "$DATABASE_URL" up
   ```

4. Запустить backend:

   ```bash
   cd backend
   go run ./cmd/api
   # curl http://localhost:8080/healthz
   ```

5. Запустить frontend:

   ```bash
   cd frontend
   npm install
   npm run dev
   # http://localhost:3000
   ```

## Тесты

```bash
cd backend && go test ./...
cd frontend && npm run test
```

## Что уже сделано (фаза 1)

- Monorepo, docker-compose (Postgres + Mailhog), `.env.example`.
- Go-бэкенд: слоистый каркас (`config`, `httpserver`, `model`, `repository`, `service`), health-check `/healthz`, goose-миграции под все 9 сущностей ТЗ.
- Полный вертикальный срез на примере `faq_item` (repository + service + unit-тесты на моках) — шаблон для остальных сущностей.
- Nuxt-приложение: маршруты-заглушки под карту сайта, `routeRules` (SSR для публичных страниц, SPA для `/admin`), Tailwind, composable `usePhoneMask` (маска `+7 (000) 000 00 00`) с тестами на vitest.

## Что дальше

- **Header/Footer и вёрстка страниц по дизайну** — заблокировано: нет доступа к Figma MCP (нужна авторизация через `/mcp` в интерактивной сессии). Пока используются нейтральные заглушки и плейсхолдер-токены в `frontend/app/assets/css/main.css`.
- CRUD для остальных сущностей (courses, course_blocks, lessons, masterclasses, teachers, blog_posts, leads, admin_users) — по шаблону `faq_item`.
- Форма заявки (lead) + email/Telegram-уведомления.
- Админ-панель (UI поверх backend CRUD).
- Блок отзывов на главной — дизайн не предоставлен, оставлен как TODO.
