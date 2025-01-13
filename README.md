# Favorites

## Описание

**Favorites API** — это API для обращения с сущностями типа Favorite, написанное при помощи фреймворка Go Gin,
в качестве базы данных используется PostgreSQL.

## Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/VanillaFroggy/favorites.git
cd favorites
```

Настройте переменные окружения:

Создайте файл .env в папке `deploy/` проекта и укажите следующие переменные:

```bash
PORT=8080
DATABASE_USER=user
DATABASE_PASSWORD=password
DATABASE_NAME=favorites
DATABASE_URL=postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@db:5432/${DATABASE_NAME}?sslmode=disable
```

## Развёртывание в Docker
Приложение разворачивается через Docker Compose.

Убедитесь, что у вас установлены Docker и Docker Compose.

Соберите и запустите контейнеры:

```bash
docker-compose up --build
```
Это создаст два контейнера:

- db: для PostgreSQL
- app: для API

Сервер будет доступен по адресу: http://localhost:8080.

## Эндпоинты

Для ознакомления с эндпоинтами через `swagger`, запустите приложение, и перейдите по ссылке:
http://localhost:8080/docs/index.html

## Структура проекта

```bash
favorites/
│
├── cmd/
│   └── favorites/
│       └── main.go                           # Входная точка приложения
│
├── config/
│   └── config.go                             # Конфигурация базы данных
│
├── deploy/
│   ├── .env                                  # Файл с переменными окружения
│   ├── .env.example                          # Файл с примером переменных окружения
│   ├── docker-compose.yml                    # Файл с иструкциями для сборки и развёртывания
│   │                                         # приложения с базой данных в Docker-контейнерах
│   └── Dockerfile                            # Файл с инструкциями для сборки
│
├── docs/                                     # Папка со сгенерированной документацией swagger
│
├── internal/
│   ├── db/
│   │   ├── migrations/                       # Папка с миграциями в БД
│   │   ├── db.go                             # Файл с функциями подключения к БД
│   │   └── migrate.go                        # Файл с функциями применения миграций к БД
│   ├── handlers/
│   │   ├── dto/                              # Папка с сущностями тел запросов или ответов
│   │   │   └── create_favorite_request.go    # Тело запроса для создания сущности БД
│   │   └── favorite_handler.go               # Файл с регистрацией и описания поведения эндпоинтов
│   ├── models/
│   │   └── favorite/                         # Папка с сущностями по тегу favorite
│   │       ├── enums.go                      # Перечисления по тегу favorite
│   │       └── favorite                      # Сущность Favorite
│   └── repository/
│       └── favorite_repo.go                  # Файл с методами для взаимодействия с БД
│
├── tests/                                    # Тесты
│   └── integration                           # Интеграционные тесты
│       └── integration_test.go               # Интеграционный тест по тегу favorite
│
├── go.mod                                    # Файл go-модуля с зависимостями
└── README.md                                 # Документация
```

## Тестирование

Для запуска тестов выполните команду до сборки docker-контейнера:

```bash
go test -v ./tests/...
```

Тесты находятся в директории `tests/` и покрывают основные эндпоинты и функционал приложения.