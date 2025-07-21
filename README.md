# Rule hub
[![Backend CI/CD](https://github.com/NRF24l01/rulehub/actions/workflows/backend.yml/badge.svg)](https://github.com/NRF24l01/rulehub/actions/workflows/backend.yml)
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