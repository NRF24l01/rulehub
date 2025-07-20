# Rule hub
*Например для вашего майнкрафт сервера*
*Просто хранилка маркдауна без возможности регистрации*

## Запуск бека
### DEV
**Нужно:**
- Постгрес
- Минио
- Golang >= 1.24

**Запуск**
```shell
cd backend
go mod download
go run main.go
```

### Тесты
**Запуск**
```shell
cd backend/tests
go test
```