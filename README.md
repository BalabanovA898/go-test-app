🚀 DevOps Lab: Automated Infrastructure Deployment

Учебный проект по автоматизации развертывания веб-приложения на языке Go с обеспечением отказоустойчивости базы данных и кэша.
📋 Описание проекта

Основная задача — создание полностью автоматизированного цикла доставки и развертывания (CI/CD). Проект включает в себя настройку инфраструктуры «с нуля» на целевой машине с использованием Ansible.

* Ключевые особенности:

  + Backend: Масштабируемое приложение на Go.

  + Database: PostgreSQL для надежного хранения данных.

  + High Availability: Redis Sentinel для мониторинга и автоматического переключения при сбоях в Redis.

  + CI/CD: GitHub Actions для автоматического запуска плейбуков при обновлении кода.

🏗 Архитектура системы

* Проект разворачивает следующие компоненты:

    + Web App: Служба systemd, запускающая бинарный файл Go.

    + PostgreSQL: Основное хранилище данных.

    + Redis Sentinel: Группа узлов Redis (Master/Slave) с арбитром (Sentinel) для обеспечения отказоустойчивости.

    + Ansible: Используется для конфигурации ОС, установки зависимостей и деплоя приложения.
    
    + Скрипт для создания бэкапов базы данных, записанный в cron.

    + yandex-disk: cli для копирования бэкапов на удаленный диск.
 
    + Стэк мониторинга: Prometheus + Graphana + Redis exporter + Node exporter для сбора базовых метрик

🛠 Технологический стек

  Infrastructure: Ansible (Playbooks, Roles, Variables).

  App: Go (Golang).

  Services: PostgreSQL, Redis Sentinel.

  Automation: GitHub Actions.

  Target OS: Ubuntu Server.

🚀 Как запустить
* Предварительные требования

  + Установленный Ansible на локальной машине.

  + SSH-доступ к целевому серверу.

    
Локальный запуск плейбука
```Bash

# Клонирование репозитория
git clone https://github.com/BalabanovA898/go-test-app
cd go-test-app/ansible

###
#Вам необходимо создать свой vault.yml с полями
#postgresql_app_password - пароль для пользователя postgresql, который будет исопльзоваться приложением
#postgres_password - паоль супер польлзователя для postgresql
#vault_docker_registry_username - имя пользователя на docker hub
#vault_docker_registry_password - пароль от docker hub
#yd__iid - id аккаунта yandex диска
#yd__passwd - пароль аккаунта yandex диска
#
#В hosts.init укажите ip Вашей виртуальной машины.
#
###

# Запуск Ansible Playbook
ansible-playbook playbooks/full-app-setup.yml -K --ask-vault-pass

⚙️ CI/CD Workflow
```
В проекте настроен GitHub Actions (.github/workflows/main.yml).

При каждом пуше в ветку main происходит автоматический запуск плейбука на удаленном сервере (через SSH).

Note: Для работы воркфлоу необходимо добавить ANSIBLE_VAULT_PASSWORD, SERVER_IP, SERVER_USER, SSH_PRIVATE_KEY, SUDO_PASSWORD в GitHub Secrets вашего репозитория.

* 📂 Структура репозитория
  + / — исходный код приложения на Go.

  + /ansible — плейбуки и роли для настройки окружения.

  + .github/workflows — конфигурация автоматизации.

  + /ansible/inventory/hosts.ini — пример файла конфигурации хостов.
