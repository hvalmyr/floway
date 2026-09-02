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
2. **Переменные** (`inventory/group_vars/floway_prod/vars.yml`) — эта
   директория лежит рядом с файлом inventory намеренно: `group_vars`
   подхватывается автоматически, только если он сосед inventory-файла или
   плейбука, а не просто где-то в `ansible/`.
   - `deploy_user_ssh_public_key` — твой публичный SSH-ключ (не секрет).
   - `floway_domain` — домен сайта. **Заранее направь A/AAAA-запись на IP
     VPS** — без этого Caddy не сможет получить сертификат Let's Encrypt.
   - `caddy_acme_email` — email для уведомлений Let's Encrypt.
   - `github_owner` — владелец GHCR-образов (по умолчанию `hvalmyr`, под
     текущий `origin`).
3. **Секреты** (vault):

   ```bash
   cp inventory/group_vars/floway_prod/vault.yml.example inventory/group_vars/floway_prod/vault.yml
   $EDITOR inventory/group_vars/floway_prod/vault.yml   # реальные пароли/токены
   ansible-vault encrypt inventory/group_vars/floway_prod/vault.yml
   ```

   Зашифрованный `vault.yml` коммитить можно и нужно — секреты живут в git,
   но нечитаемы без пароля. Пароль от vault храни отдельно (менеджер паролей
   или секрет в CI), никогда не в репозитории.

4. **Vault-пароль как файл** (чтобы не вводить его на каждый запуск):
   создай `ansible/.vault_pass` с паролем внутри (одна строка, без
   переноса на конце — `printf '%s' 'твой-пароль' > ansible/.vault_pass`),
   и `chmod 600 ansible/.vault_pass`. Файл в `.gitignore`, коммитить нельзя.
   `just`-рецепты (`bootstrap-first-run`, `bootstrap`, `deploy`, `vault-edit`)
   используют `--vault-password-file .vault_pass` вместо интерактивного
   `--ask-vault-pass`.

## Бутстрап VPS

Первый прогон против чистой машины — до создания `deploy`-пользователя
подключаемся под `root` (или тем пользователем, что дал провайдер):

```bash
just bootstrap-first-run
# эквивалент: ansible-playbook playbooks/bootstrap.yml -e ansible_user=root --vault-password-file .vault_pass
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

## Сборка и публикация образов

Это делает CI (`.github/workflows/ci.yml`, джоба `docker-build`) на каждый пуш
в `main`: собирает `floway-backend`/`floway-frontend` и пушит в GHCR двумя
тегами — `sha-<короткий-sha>` и `latest`. На PR образы только собираются
(проверка, что Dockerfile жив), но не публикуются.

Руками, если понадобится:

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

## Деплой из CI

Джоба `deploy` в `.github/workflows/ci.yml` гоняет ровно тот же
`playbooks/deploy.yml`, что и `just deploy`, только из раннера GitHub Actions.
Порядок: пуш в `main` → зелёные линты/тесты/сборки → `docker-build` публикует
образы в GHCR → `deploy` раскатывает свежий `sha-`тег на VPS → смоук-проверка
`https://<floway_domain>/healthz`.

Гоняется именно Ansible, а не `ssh + docker compose pull`, чтобы изменения
`docker-compose.prod.yml`, `Caddyfile.j2` и `env.j2` в git доезжали до прода
тем же деплоем, что и код, — иначе конфиг и образы разъезжаются.

Ручной запуск — Actions → CI → Run workflow: пустой `image_tag` пересобирает и
катит текущий коммит, заполненный (например `sha-abc1234`) катит уже
опубликованный тег, ничего не пересобирая и не двигая `latest`, — то есть
откат на предыдущий образ.

### Подготовка (один раз)

1. Сгенерировать отдельный SSH-ключ для CI (без пароля — раннеру некому его
   вводить):

   ```bash
   ssh-keygen -t ed25519 -C "github-actions@floway" -N "" -f ~/.ssh/floway_ci_deploy
   ```

2. Публичную часть (`~/.ssh/floway_ci_deploy.pub`) добавить третьей строкой в
   `deploy_user_ssh_public_key` в `inventory/group_vars/floway_prod/vars.yml`
   и раскатать на VPS: `just bootstrap`.

3. Собрать отпечаток хоста для `known_hosts` (в `ansible.cfg` включён
   `host_key_checking = True`, поэтому в CI он нужен явно):

   ```bash
   ssh-keyscan -H 79.143.29.84
   ```

4. Завести секреты репозитория (Settings → Secrets and variables → Actions):

   | Секрет | Что внутри |
   |---|---|
   | `DEPLOY_SSH_PRIVATE_KEY` | содержимое `~/.ssh/floway_ci_deploy` целиком, включая строки `-----BEGIN/END OPENSSH PRIVATE KEY-----` |
   | `DEPLOY_SSH_KNOWN_HOSTS` | вывод `ssh-keyscan -H <IP VPS>` |
   | `ANSIBLE_VAULT_PASSWORD` | пароль от `ansible-vault` (то же, что лежит в `ansible/.vault_pass`) |

   `GITHUB_TOKEN` для пуша в GHCR заводить не нужно — он встроенный, джобе
   выдано `packages: write`.

Ключ и пароль vault пишутся в раннере во временные файлы и живут только время
джобы; `.vault_pass` удаляется по `trap ... EXIT`, даже если плейбук упал.

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
