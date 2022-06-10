#!/bin/bash

CONFIG_FILE="config.json"

# Обновление пакетов
yes | sudo apt-get update
#yes | sudo apt-get upgrade

# Установка jq
yes | sudo apt-get install jq

# Установка mc
yes | sudo apt-get install mc

# Установка pip
yes | sudo apt-get install pip

# Установка модуля psycopg2
yes | sudo apt-get install libpq-dev
pip install psycopg2

# Установить зависимости для парсера и бота
pip3 install -r requirements.txt

sudo apt-get install python3-flask

# Установка PostgreSQL
yes | sudo apt-get install postgresql

## Установить пароль для пользователя postgres
postgres_password=$(jq -r .postgresql.superuser_password $CONFIG_FILE)
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD '$postgres_password';"

## Разрешить подключение к PostgreSQL из внешней сети
subnet=$(jq -r .postgresql.config.subnet $CONFIG_FILE)
sudo sh -c "cat >> /etc/postgresql/12/main/pg_hba.conf <<EOF
host    all             all             $subnet        md5
EOF"

## Прослушивать подключения к PostgreSQL по любому адресу
listen_addresses=\'$(jq -r .postgresql.config.listen_addresses $CONFIG_FILE)\'
search_string="#listen_addresses = 'localhost'"
replace_string="listen_addresses = $listen_addresses"
sudo sed -i "s/$search_string/$replace_string/" /etc/postgresql/12/main/postgresql.conf

## Перезагрузить PostgreSQL
sudo systemctl restart postgresql.service

# Развёртывание БД
user=$(jq -r .postgresql.connect.user.name $CONFIG_FILE)
password=$(jq -r .postgresql.connect.user.password $CONFIG_FILE)
database=$(jq -r .postgresql.connect.database $CONFIG_FILE)
host=$(jq -r .postgresql.connect.host $CONFIG_FILE)
port=$(jq -r .postgresql.connect.port $CONFIG_FILE)

## Создание пользователя БД
query="CREATE ROLE \"$user\" CREATEDB LOGIN ENCRYPTED PASSWORD '$password';"
sudo -u postgres psql -c "$query"

## Создание БД
query="CREATE DATABASE \"$database\" WITH OWNER = \"$user\";"
sudo -u postgres psql -c "$query"

## Поменять пользователя схемы public
query="ALTER SCHEMA \"public\" OWNER TO \"$user\";"
sudo -u postgres psql -d $database -c "$query"

## Развернуть все таблицы БД
psql -h $host -U $user -d $database -f deploy_database.sql
