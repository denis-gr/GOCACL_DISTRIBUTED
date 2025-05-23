# Very Tired Calculator  (GOCACL_DISTRIBUTED)
![Coverage](https://img.shields.io/badge/Coverage-68.5%25-yellow)

![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?logo=githubactions&logoColor=white&style=flat)
![Bootstrap](https://img.shields.io/badge/bootstrap-%238511FA.svg?logo=bootstrap&logoColor=white&style=flat)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?&logo=go&logoColor=whitee&style=flat)
[![Go](https://github.com/denis-gr/GOCACL_DISTRIBUTED/actions/workflows/go.yml/badge.svg)](https://github.com/denis-gr/GOCACL_DISTRIBUTED/actions/workflows/go.yml)

Этот ридми пишет очень уставший Денис, но несмотря на усталость, я буду рад, если тебе понравится мой калькулятор. Если будут какие-то проблемы, то пиши в тг [@denisgrigoriev04](https://t.me/denisgrigoriev04) или в [Issues](https://github.com/denis-gr/GOCACL_DISTRIBUTED/issues). Я ещё написал веб-интерфейс, поэтому если захочешь протестить, то просто запусти его, используя docker compose.

Этот калькулятор очень устал, поэтому выполняет арифметические опрации очень медлено. Поэтому он состоит из нескольких компоненов `orchestrator` и `agent`. `orchestrator` управляет распределением задач, а `agent` выполняет вычисления. Чтобы ускорить выполнения операция стоит запускать несколько `agent`, каждый агент также может считать несколько операций одновременно.


![Вид сайта](NoGo/image.png)


## Использование

1. Склонируйте репозиторий:
   ```sh
   git clone https://github.com/denis-gr/GOCACL_DISTRIBUTED
   cd GOCACL_DISTRIBUTED
   ```

2. Запустите проект с помощью Docker Compose:
   ```sh
   docker compose up
   ```

3. Сайт будет доступен по адресу [http://localhost/](http://localhost/)


### Примеры использования API (командная строка Linux)

- Регистрация пользователя:
  ```sh
  curl --location 'http://localhost/api/v1/register' --header 'Content-Type: application/json' --data '{ "login": "user", "password": "password" }'
  ```

- Вход пользователя, здесь вы получите JWT токен:
  ```sh
  curl --location 'http://localhost/api/v1/login' --header 'Content-Type: application/json' --data '{ "login": "user", "password": "password" }'
  ```

- Добавление вычисления арифметического выражения (требуется аутентификация):
  ```sh
  curl --location 'http://localhost/api/v1/calculate' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <JWT_TOKEN>' \
  --data '{ "expression": "2+2*2" }'
  ```

- Получение списка выражений (требуется аутентификация):
  ```sh
  curl --location 'http://localhost/api/v1/expressions' \
  --header 'Authorization: Bearer <JWT_TOKEN>'
  ```

- Получение выражения по его идентификатору (требуется аутентификация):
  ```sh
  curl --location 'http://localhost/api/v1/expressions/:id' \
  --header 'Authorization: Bearer <JWT_TOKEN>'
  ```

## Архитектура


![Диаграмма взаимодействия сервисов](NoGo/diagram.png)

Проект состоит из двух основных компонентов:
- `orchestrator`: управляет распределением задач и хранением результатов.
- `agent`: выполняет вычисления и отправляет результаты обратно `orchestrator`.


## Тестирование

Для запуска тестов выполните:
```sh
go test ./...
```
