# Wildberries internship Level 0

***

Тестовое задание WBTech

В сервисе:
<ol>
    <li>Подключение и подписка на канал в nats-streaming</li>
    <li>Полученные данные писать в Postgres</li>
    <li>Так же полученные данные сохранить in memory в сервисе (Кеш)</li>
    <li>В случае падения сервиса восстанавливать Кеш из Postgres</li>
    <li>Поднять http сервер и выдавать данные по id из кеша</li>
</ol>

Как поднять:
<ol>
    <li>make up - поднять контейнер для бд (docker-compose up)
    <li>make server - поднять сервер + nats jetstream
    <li>make client - запустить клиента
</ol>
***

### Stack:
* Go
* Postgres
* Docker-Compose
* Nats (jetstream)
* Vegeta (стресс тесты)
