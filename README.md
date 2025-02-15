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

## Нагрузочное тестирование
Тестирование было проведено с помощью [k6](scripts/k6_load_test.js). Результаты представлены ниже.
```
     ✗ logged in successfully
      ↳  99% — ✓ 90749 / ✗ 4
     ✓ status is not 500

     checks.........................: 99.99% 376175 out of 376179
     data_received..................: 495 MB 1.6 MB/s
     data_sent......................: 110 MB 363 kB/s
     dropped_iterations.............: 204854 675.915172/s
     http_req_blocked...............: avg=97.86µs min=2.53µs     med=8.6µs   max=348.55ms p(90)=12.7µs   p(95)=14.6µs
     http_req_connecting............: avg=86.42µs min=0s         med=0s      max=348.26ms p(90)=0s       p(95)=0s
     http_req_duration..............: avg=37.83ms min=12.11ms    med=25.35ms max=5.03s    p(90)=58.94ms  p(95)=82.64ms
       { expected_response:true }...: avg=37.69ms min=12.11ms    med=25.33ms max=5.03s    p(90)=58.6ms   p(95)=82.12ms
     http_req_failed................: 0.94%  3570 out of 376179
     http_req_receiving.............: avg=1.51ms  min=-1851631ns med=80.5µs  max=3.54s    p(90)=813.64µs p(95)=1.87ms
     http_req_sending...............: avg=53.5µs  min=8.03µs     med=31.27µs max=15.31ms  p(90)=101.08µs p(95)=134.6µs
     http_req_tls_handshaking.......: avg=0s      min=0s         med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=36.26ms min=11.94ms    med=25.02ms max=2.41s    p(90)=57.14ms  p(95)=79.95ms
     http_reqs......................: 376179 1241.201506/s
     iteration_duration.............: avg=3.15s   min=16.51ms    med=3.11s   max=8.13s    p(90)=3.22s    p(95)=3.29s
     iterations.....................: 95146  313.933948/s
     vus............................: 46     min=46               max=1000
     vus_max........................: 1000   min=447              max=1000
```

## Вопросы и решения
1. Во время работы над интеграционными тестами понадобилось быть уверенным в доступности API. С этой целью добавил маршрут `/health`, возвращающий `200 OK`.