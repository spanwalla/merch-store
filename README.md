# Внутренний магазин мерча
Более подробное описание можно найти [здесь](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md).

## Настройка и запуск
1. Склонируйте репозиторий.
```
git clone https://github.com/spanwalla/merch-store
cd merch-store
```
2. Создайте файл `.env` в корневом каталоге проекта по образцу [.env.example](.env.example).
3. Для запуска контейнеров выполните команду:
```
make compose-up
```
ИЛИ `docker-compose up --build -d`.

4. Для остановки используйте команду:
```
make compose-down
```
ИЛИ `docker-compose down --remove-orphans`.

5. Для запуска юнит-тестов можно использовать команды:
```
make test
```
ИЛИ `go test -v './internal/...'`.

6. Для запуска интеграционных тестов используйте:
```
make compose-up-integration-test
```
ИЛИ `docker-compose --profile tests up --build --abort-on-container-exit --exit-code-from integration`.

После запуска сервис будет доступен по адресу `http://localhost:8080`.
