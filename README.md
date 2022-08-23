# Портфолио по разработке backend на Golang

## Clone Reddit
Представляет собой сайт, предоставляющий функционал добавления постов на опредленные темы, а также
комментириев к ним и их оценки.

## Для запуска клона реддита
Нужно зайти backend/reddicloneDB и поднять докером БД (MySQL + MongoDB):

    docker-compose up

Удалются через команды:

    docker rm $(docker ps -a -q)
    docker volume prune -f

Потом в backend/redditclone/cmd запустить main.go (изначально на 8080 порту):

    go run main.go

Открыть в браузере можно по адресу:

    http://localhost:8080/

Также для основных функций были написаны тесты (в папка pkg/handlers, pkg/posts, pkg/user), их можно запустить:

    go test -v

## Микросервис
Код находится в microservice/service.go, в admin.go - логика администрирования, чтобы получать логи с сервиса, данные
передается по протоколу protobuf.
Для данной функции есть тесты (service_test.go), которые можно запустить через:

    go test -v -race