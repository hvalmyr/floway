# Локальная разработка — см. README.md.
# Продакшен-инфраструктура (Ansible) — см. ansible/README.md.

set shell := ["bash", "-uc"]

ansible_dir := "ansible"

# --- Ansible ---

# Установить коллекции Ansible (community.docker, community.general, ansible.posix)
ansible-deps:
    ansible-galaxy collection install -r {{ansible_dir}}/requirements.yml

# Первый прогон бутстрапа против чистой VPS (пока нет пользователя deploy — коннект от root)
bootstrap-first-run:
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml -e ansible_user=root --vault-password-file .vault_pass

# Повторный/идемпотентный прогон бутстрапа (уже под пользователем deploy)
bootstrap:
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml --vault-password-file .vault_pass

# Деплой приложения (образы должны быть уже собраны и запушены в GHCR)
deploy:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --vault-password-file .vault_pass

# Деплой конкретного тега образа вместо latest, например: just deploy-tag sha-abc1234
deploy-tag tag:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --vault-password-file .vault_pass -e image_tag={{tag}}

# Синтаксис-проверка обоих плейбуков
ansible-check:
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml --syntax-check
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --syntax-check

# Редактировать зашифрованные секреты
vault-edit:
    cd {{ansible_dir}} && ansible-vault edit --vault-password-file .vault_pass inventory/group_vars/floway_prod/vault.yml

# --- decor-сайт (flo-way.ru) — тот же хост, отдельная inventory-группа ---

# Деплой decor-сайта (образы должны быть уже собраны и запушены в GHCR)
deploy-decor:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --vault-password-file .vault_pass -e target_group=floway_decor_prod

# Деплой конкретного тега образа decor-сайта вместо latest
deploy-decor-tag tag:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --vault-password-file .vault_pass -e target_group=floway_decor_prod -e image_tag={{tag}}

# Редактировать зашифрованные секреты decor-сайта
vault-edit-decor:
    cd {{ansible_dir}} && ansible-vault edit --vault-password-file .vault_pass inventory/group_vars/floway_decor_prod/vault.yml

# Отключить decor-сайт на несезон — контейнеры остановлены, БД/файлы/образы
# сохраняются (см. TZ п. 4 "Отключение на несезон"). Школьный сайт и общий
# Caddy не затрагиваются.
decor-stop:
    cd {{ansible_dir}} && ansible floway_decor_prod -i inventory/production.yml -m shell -a "docker compose stop chdir=/opt/floway-decor" --vault-password-file .vault_pass

# Включить decor-сайт обратно перед следующим сезоном
decor-start:
    cd {{ansible_dir}} && ansible floway_decor_prod -i inventory/production.yml -m shell -a "docker compose start chdir=/opt/floway-decor" --vault-password-file .vault_pass

# --- Локальная разработка (см. README.md за подробностями) ---

dev-up:
    docker compose up -d --build

dev-down:
    docker compose down

dev-test:
    cd backend && go test ./...
    cd frontend && bun run test

frontend-lint:
    cd frontend && bun run lint && bun run format:check

# Включить git-хуки репозитория (.githooks/) — прогоняют форматтер/линтер/тесты
# перед каждым коммитом, скоуп по тому, что реально застейджено (frontend/backend).
install-hooks:
    git config core.hooksPath .githooks
    @echo "OK: git config core.hooksPath .githooks"
