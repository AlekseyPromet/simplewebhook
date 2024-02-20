# simplewebhook
Сервис - простой  webhook

Необходимо написать программу(сервис) реализующим функционал отправки запросров на сторонний сервис(реализация webhook).

Пример конфигурации

url: https://some.domain
requests:
  amount: 1000
  per_second: 10

Отправляем POST запросом на url запросы с body = { "iteration": ${number} } где number порядковый номер итерации из requests.amount. Реализация должна учитывать количество отправляемых запросов request.per_second. Написать тест для предложенной реализации.

Исползуемый инстументарий
- go 1.22
- https://github.com/uber-go/zap
- https://github.com/spf13/cobra
- https://github.com/uber-go/fx
- https://github.com/go-resty/resty

Тестирование сервиса:

- запустить тестовый клиент 
    $ go run ./tests/client/main.go

- запустить сервер
    $ go mod download
    $ go run main.go

- запрос на создание webhook для localhost:8090
    curl --location 'localhost:8088/invoke' \
    --header 'Content-Type: application/json' \
    --data '{
    "url": "http://localhost:8090",
    "requests": {
        "amount": 1000,
        "per_seconds": 1
        }
    }'

- запрос на получение правила по ключу
    curl --location 'localhost:8088/get' \
    --header 'Content-Type: application/json' \
    --data '{
        "key": "d32b1d54-3b09-40e9-b14e-66174d7ad5c2"
    }'