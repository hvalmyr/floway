# Floway frontend (Nuxt)

Пакетный менеджер и раннер — [Bun](https://bun.sh). См. также корневой [README](../README.md).

## Установка зависимостей

```bash
bun install
```

## Дев-сервер

```bash
bun run dev
# http://localhost:3000
```

## Продакшен-сборка

```bash
bun run build
bun run preview   # локально проверить прод-сборку
```

## Линт, форматирование, тесты

```bash
bun run lint          # oxlint
bun run format:check  # oxfmt --check
bun run format         # oxfmt (исправить на месте)
bun run test           # vitest
```

Подробнее про [Nuxt](https://nuxt.com/docs/getting-started/introduction).
