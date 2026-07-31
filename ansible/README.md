# Ansible: инфраструктура Floway

Два плейбука:

- **`playbooks/bootstrap.yml`** — приводит чистую VPS (Ubuntu 24.04) в
  безопасное рабочее состояние: пользователь для деплоя, Docker, своп,
  автообновления безопасности, fail2ban, ufw, ssh-хардненинг. Идемпотентен —
  можно и нужно гонять повторно (например, после ручных правок на машине).
- **`playbooks/deploy.yml`** — раскладывает и запускает само приложение
  (`docker-compose.prod.yml` + Caddy) поверх уже забутстрапленной машины.

Архитектура прода: **Caddy** (авто-HTTPS по домену, единственная точка входа
на 80/443) → проксирует `/api/*` и `/healthz` на `backend`, всё остальное на
`frontend`. `postgres`/`backend`/`frontend` не публикуют порты наружу —
только внутренняя docker-сеть. Backend/frontend — готовые образы из GHCR
(`ghcr.io/<owner>/floway-backend`, `floway-frontend`), собранные в CI, а не на
самой VPS.

## Предварительные требования

```bash
just ansible-deps   # ansible-galaxy collection install -r requirements.yml
```

Нужны коллекции `community.docker`, `community.general`, `ansible.posix`
(перечислены в `requirements.yml`).

## Настройка перед первым запуском

1. **Inventory** (`inventory/production.yml`): впиши белый IP VPS в
   `ansible_host`.
2. **Переменные** (`group_vars/floway_prod/vars.yml`):
   - `deploy_user_ssh_public_key` — твой публичный SSH-ключ (не секрет).
   - `floway_domain` — домен сайта. **Заранее направь A/AAAA-запись на IP
     VPS** — без этого Caddy не сможет получить сертификат Let's Encrypt.
   - `caddy_acme_email` — email для уведомлений Let's Encrypt.
   - `github_owner` — владелец GHCR-образов (по умолчанию `hvalmyr`, под
     текущий `origin`).
3. **Секреты** (vault):

   ```bash
   cp group_vars/floway_prod/vault.yml.example group_vars/floway_prod/vault.yml
   $EDITOR group_vars/floway_prod/vault.yml   # реальные пароли/токены
   ansible-vault encrypt group_vars/floway_prod/vault.yml
   ```

   Зашифрованный `vault.yml` коммитить можно и нужно — секреты живут в git,
   но нечитаемы без пароля. Пароль от vault храни отдельно (менеджер паролей
   или секрет в CI), никогда не в репозитории.

## Бутстрап VPS

Первый прогон против чистой машины — до создания `deploy`-пользователя
подключаемся под `root` (или тем пользователем, что дал провайдер):

```bash
just bootstrap-first-run
# эквивалент: ansible-playbook playbooks/bootstrap.yml -e ansible_user=root --ask-vault-pass
```

**Перед тем как закрывать эту сессию** — открой отдельный терминал и
проверь, что новый пользователь реально работает:

```bash
ssh deploy@<IP>
```

Роль `ssh_hardening` идёт последней в плейбуке и намеренно: сначала
создаётся `deploy`-пользователь с рабочим ключом и sudo, и только потом
отключаются root-логин и пароль. Конфиг `sshd` проверяется командой `sshd -t`
_до_ перезапуска демона — если проверка не проходит, плейбук падает, а
старый (рабочий) `sshd` остаётся жив. Но если ключ `deploy`-пользователя
всё-таки неверный, а хардненинг уже применился — доступ можно потерять.
Отсюда и обязательная проверка в отдельном терминале. У большинства VPS-хостеров
есть веб-консоль (VNC/serial) как аварийный доступ, если что-то пошло не так.

Повторные прогоны — уже под `deploy`:

```bash
just bootstrap
```

## Сборка и публикация образов (пока вручную, до GitHub Actions)

```bash
docker buildx build --push -t ghcr.io/<owner>/floway-backend:latest ./backend
docker buildx build --push -t ghcr.io/<owner>/floway-frontend:latest ./frontend
```

Если пакеты в GHCR приватные — либо сделай их публичными (Settings пакета →
Change visibility, самый простой вариант для сайта без закрытого кода), либо
поставь `ghcr_login_required: true` в `vars.yml` и заполни
`vault_ghcr_username`/`vault_ghcr_token` (GitHub PAT с правом `read:packages`)
в vault.

## Деплой приложения

```bash
just deploy
# или конкретный тег: just deploy-tag sha-abc1234
```

Что делает: копирует `docker-compose.prod.yml` (из корня репозитория) и
`Caddyfile`/`.env` (шаблонизированные из vars+vault) на VPS в
`/opt/floway/`, логинится в GHCR (если нужно), `docker compose pull && up -d`,
чистит неиспользуемые образы.

## Планы по CI/CD (когда дойдут руки до GitHub Actions)

Предлагаемый пайплайн:

1. На пуш в `main` — джоба собирает и пушит `floway-backend`/`floway-frontend`
   в GHCR с тегом `sha-<short-sha>` (и обновляет `latest`).
2. Джоба деплоя — либо `ansible-playbook playbooks/deploy.yml` прямо в раннере
   GitHub Actions (нужен SSH-ключ и пароль vault как секреты Action), либо
   более простой вариант — деплой без Ansible на этом шаге: `ssh` на VPS и
   `docker compose pull && up -d` (Ansible-плейбук уже разложил
   docker-compose.yml/Caddyfile/.env один раз при бутстрапе/первом деплое,
   для рутинных обновлений образов полноценный Ansible-прогон избыточен).
3. Пароль vault хранить как секрет GitHub Actions (`ANSIBLE_VAULT_PASSWORD`),
   передавать через `--vault-password-file <(echo "$ANSIBLE_VAULT_PASSWORD")`.

## Что не входит в эти плейбуки (но стоит сделать по мере роста прод-нагрузки)

- **Бэкапы Postgres** — `postgres_data` живёт в docker-volume на одной VPS,
  без внешних бэкапов. Минимум — cron с `pg_dump` + отправкой в S3-совместимое
  хранилище (Backblaze B2 / Selectel S3 и т.п.).
- **Мониторинг/алерты** — сейчас только `docker compose logs`. Для начала
  достаточно uptime-проверки (UptimeRobot/Healthchecks.io на `/healthz`).
- **Zero-downtime деплой** — `docker compose up -d` пересоздаёт контейнеры с
  секундами простоя. При текущей нагрузке сайта-витрины это не критично;
  если станет важно — можно посмотреть в сторону blue-green через два набора
  контейнеров и переключение Caddy, но это не обязательно.
