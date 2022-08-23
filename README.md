# Портфолио по разработке backend на Golang

## Для запуска клона реддита
Нужно зайти backend/reddicloneDB и поднять докером БД
Потом в backend/redditclone/cmd запустить main.go (изначально на 8080 порту)

## Микросервис
Код находится в microservice/service.go, в admin.go - логика администрирования, чтобы получать логи с сервиса.
Для данной функции есть тесты, которые можно запустить через go test -v -race.