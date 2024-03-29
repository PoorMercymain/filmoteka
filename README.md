# filmoteka
[![codecov](https://codecov.io/gh/PoorMercymain/filmoteka/graph/badge.svg?token=COEUBP3510)](https://codecov.io/gh/PoorMercymain/filmoteka)
# Быстрый запуск (Docker)
Клонируем репозиторий, переходим в его корневую папку, переименовываем `.env.example` в `.env`, после чего в терминале прописываем `docker-compose up`
Ждем, пока пройдет healthcheck постгреса и сервис напишет в логи, что он слушает на заданном порту (по умолчанию `8080`)

# Swagger
Документация swagger сгенерирована из комментариев-аннотаций. Чтобы получить доступ к Swagger UI, после запуска сервиса нужно обратиться к `/swagger/` (по умолчанию `http://localhost:8080/swagger/`)

# Миграции
Файлы миграций находятся в директории `migrations` корневой папки. Для их применения при запуске сервис использует `golang-migrate`

# Эндпойнты
У сервиса присутствуют следующие эндпойнты:</br>
`POST /actor` - добавить актера в БД</br>
`PUT /actor/{id}` - обновить актера</br>
`DELETE /actor/{id}` - удалить актера из БД</br>
`GET /actors` - получить список актеров с соответствующими им фильмами</br>
</br>
`POST /film` - добавить фильм в БД</br>
`PUT /film/{id}` - обновить фильм</br>
`DELETE /film/{id}` - удалить фильм из БД</br>
`GET /films` - получить список фильмов с возможностью сортировки по различным полям</br>
`GET /films/search` - найти фильм по фрагменту названия и/или фрагменту имени актера</br>
</br>
`POST /register` - зарегистрироваться в сервисе</br>
`POST /login` - получить токен авторизации</br>
Подробнее они расписаны в Swagger
