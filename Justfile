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
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml -e ansible_user=root --ask-vault-pass

# Повторный/идемпотентный прогон бутстрапа (уже под пользователем deploy)
bootstrap:
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml --ask-vault-pass

# Деплой приложения (образы должны быть уже собраны и запушены в GHCR)
deploy:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --ask-vault-pass

# Деплой конкретного тега образа вместо latest, например: just deploy-tag sha-abc1234
deploy-tag tag:
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --ask-vault-pass -e image_tag={{tag}}

# Синтаксис-проверка обоих плейбуков
ansible-check:
    cd {{ansible_dir}} && ansible-playbook playbooks/bootstrap.yml --syntax-check
    cd {{ansible_dir}} && ansible-playbook playbooks/deploy.yml --syntax-check

# Редактировать зашифрованные секреты
vault-edit:
    cd {{ansible_dir}} && ansible-vault edit group_vars/floway_prod/vault.yml

# --- Локальная разработка (см. README.md за подробностями) ---

dev-up:
    docker compose up -d --build

dev-down:
    docker compose down

dev-test:
    cd backend && go test ./...
    cd frontend && npm run test
