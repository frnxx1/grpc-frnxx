## Запуск
docker run -it --name user_test -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=user -e POSTGRES_DB=db  -e POSTGRES_HOST=localhost -p 5432:5432  postgres
go run main.go
