# Простой шаблон Ansible-проекта

## Структура

- `ansible.cfg` — базовая конфигурация Ansible.
- `inventory/hosts.ini` — inventory с примерами хостов.
- `playbooks/site.yml` — главный playbook, который вызывает роль `common`.
- `roles/common/` — простая базовая роль.
- `requirements.txt` — Python-зависимости (ansible).

## Как использовать

1. Установите зависимости:
   ```bash
   pip install -r requirements.txt
   ```

2. Отредактируйте `inventory/hosts.ini`, указав реальные хосты, пользователя и, при необходимости, ключ:
   ```ini
   [all]
   server1 ansible_host=your_server_ip ansible_user=your_user
   ```

3. Запустите playbook:
   ```bash
   ansible-playbook playbooks/site.yml
   ```

Playbook:
- поставит базовый набор пакетов (`git`, `curl`, `vim`);
- задеплоит `/etc/motd` через шаблон `roles/common/templates/motd.j2`, используя переменную `motd_message`;
- при изменении motd дернет handler на перезапуск `sshd`.

