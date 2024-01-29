## Workflow

Сервис состоит из 3-х частей: 
- `HTTP Server` - HTTP-сервер, который принимает запросы от клиентов и отправляет данные о заказе по ID.
- `Consumer (Subscriber)` - подписчик, который слушает определенный канал *NATS-streaming server'а*, и сохраняет данные о заказе в базу данных.
- `Скрипт producer'а` - скрипт для демонстрации работы, который отправляет данные о заказе в канал *NATS-streaming server'а*.

![worklof diogram](./assets/workflow.png)

## Simple interface

- Простой интерфейс для работы с HTTP-сервером.

![simple interface](./assets/interface.png)

## WRK benchmark

- Бенчмарк на 30 секунд, используя 12 потоков и поддерживая 400 открытых HTTP-соединений.

![wgk benchmark](./assets/wrk.svg)

## Stack

- Golang
- NATS-streaming server
- PostgreSQL
- Docker
- Redis