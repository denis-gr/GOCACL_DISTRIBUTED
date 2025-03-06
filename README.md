# GOCACL_DISTRIBUTED

## Описание проекта

GOCACL_DISTRIBUTED — это распределённый вычислитель арифметических выражений. Проект включает два основных компонента: `orchestrator` и `agent`. `orchestrator` управляет распределением задач, а `agent` выполняет вычисления.

## Требования

- Docker Compose
- Go 1.23.4

## Установка

1. Склонируйте репозиторий:
   ```sh
   git clone https://github.com/denis-gr/GOCACL_DISTRIBUTED
   cd GOCACL_DISTRIBUTED
   ```

2. Запустите проект с помощью Docker Compose:
   ```sh
   docker compose up
   ```

## Использование

Откройте [http://localhost/](http://localhost/), введите в поле "Арифметическое выражение" своё выражение и нажмите на кнопку "Отправить", чтобы наблюдать процесс его вычисления.

### Примеры использования API

- Добавление вычисления арифметического выражения:
  ```sh
  curl --location 'http://localhost/api/v1/calculate' --header 'Content-Type: application/json' --data '{ "expression": "2+2*2" }'
  ```

- Получение списка выражений:
  ```sh
  curl --location 'http://localhost/api/v1/expressions'
  ```

- Получение выражения по его идентификатору:
  ```sh
  curl --location 'http://localhost/api/v1/expressions/:id'
  ```

## Архитектура

Проект состоит из двух основных компонентов:
- `orchestrator`: управляет распределением задач и хранением результатов.
- `agent`: выполняет вычисления и отправляет результаты обратно `orchestrator`.

![Диаграмма взаимодействия сервисов](NoGo/diagram.png)

## Тестирование

Для запуска тестов выполните:
```sh
go test ./...
```

## Контакты

Если у вас есть вопросы или предложения, свяжитесь со мной через [Telegram](https://t.me/denisgrigoriev04).
