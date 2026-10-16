# Комп’ютерний практикум 2. Робота з постійним сховищем і валідація вхідних даних

REST API для керування притулками та котами. Застосунок написаний на Go з використанням Gin, GORM і PostgreSQL.

## Завдання

1. Замінити тимчасове зберігання даних із першої роботи на постійне зберігання в базі даних.
2. Додати валідацію вхідних даних.
3. Використовувати ORM і шаблон «Репозиторій» для роботи з базою даних.
4. Отримувати репозиторії через ін’єкцію залежностей.
5. Використовувати окремі DTO замість передачі сутностей бази даних через REST API.
6. Валідувати вхідні DTO за допомогою аналога FluentValidation для Go.

## Реалізація

- HTTP-фреймворк: Gin.
- ORM: GORM.
- База даних: PostgreSQL.
- `CatRepoPsql` і `ShelterRepoPsql` реалізують постійне зберігання.
- PostgreSQL foreign key зв’язує кота з притулком і забороняє видаляти притулок, у якому є коти (`ON DELETE RESTRICT`).
- Транзакції з `SELECT FOR UPDATE` не дозволяють одночасним запитам переповнити притулок.
- Репозиторії передаються в handlers через конструктори та інтерфейси.
- In-memory репозиторії збережені для ізольованих тестів.
- Вхідні та вихідні DTO знаходяться в `internal/dto`.
- Перетворення між DTO і моделями знаходяться в `internal/mapper`.
- Валідація виконується Gin через `go-playground/validator`.
- Невідомі JSON-поля відхиляються.
- Технічні помилки бази даних записуються в серверний лог і не повертаються клієнту.
- Відомі помилки мають однакові відповіді: `400` для неправильного зв’язку, `404` для відсутнього ресурсу та `409` для ресурсу із залежними записами або заповненого притулку.

## Правила валідації

### Cat

- `name` — обов’язкове, від 1 до 100 символів;
- `ageMonths` — обов’язкове, від 0 до 360 місяців;
- `gender` — обов’язкове, значення `male` або `female`;
- `shelterID` — обов’язкове, більше нуля;
- кота не можна додати або перемістити до притулку без вільних місць;
- `id` призначається базою даних і не приймається у вхідному DTO.

### Shelter

- `name` — обов’язкове, від 1 до 100 символів;
- `city` — обов’язкове, від 1 до 100 символів;
- `address` — обов’язкове, від 1 до 200 символів;
- `capacity` — обов’язкова, більше нуля та не менша за поточну кількість котів;
- `id` призначається базою даних і не приймається у вхідному DTO.

## Запуск у Docker

Потрібні Docker і Docker Compose v2 або новіший.

```bash
cp .env.example .env
# Перед запуском задайте власний POSTGRES_PASSWORD у .env.
docker compose up --build -d
```

У PowerShell:

```powershell
Copy-Item .env.example .env
docker compose up --build -d
```

Сервіси:

- API: http://localhost:8080/api/v1
- Swagger UI: http://localhost:8080/swagger/index.html
- OpenAPI YAML: http://localhost:8080/openapi.yaml
- PostgreSQL з хоста: `localhost:5432`
- PostgreSQL у мережі Compose: `postgres:5432`

Для захищених операцій `POST`, `PUT` і `DELETE` потрібен заголовок:

```text
X-API-Key: demo
```

Перегляд логів:

```bash
docker compose logs -f
```

Зупинка:

```bash
docker compose down
```

Команда `docker compose down` зберігає дані в томі `postgres_data`. Не додавайте `-v`, якщо дані потрібно зберегти: цей параметр видаляє том разом із базою даних.

## Приклади запитів

Створення притулку:

```bash
curl -X POST http://localhost:8080/api/v1/shelters \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo" \
  -d '{"name":"Happy Paws","city":"Kyiv","address":"Main Street 1","capacity":50}'
```

Створення кота:

```bash
curl -X POST http://localhost:8080/api/v1/cats \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo" \
  -d '{"name":"Barsik","ageMonths":12,"gender":"male","adopted":false,"shelterID":1}'
```

Приклад помилки валідації:

```bash
curl -X POST http://localhost:8080/api/v1/cats \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo" \
  -d '{"id":10,"name":"Barsik","ageMonths":361,"gender":"unknown","shelterID":1}'
```

Відповідь:

```json
{"error":"invalid request body"}
```

## Тести

```bash
go test ./...
go vet ./...
```

## Генерація OpenAPI

Встановлення `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.6
```

Linux/macOS:

```bash
./scripts/generate-openapi.sh
```

PowerShell:

```powershell
./scripts/generate-openapi.ps1
```

Генеруються файли:

- `docs/docs.go`;
- `docs/swagger.json`;
- `docs/swagger.yaml`;
- `docs/openapi.yaml`.

## Структура

```text
cmd/api                  точка входу
internal/app             складання залежностей
internal/config          конфігурація
internal/dto             DTO REST API
internal/handlers        HTTP handlers
internal/mapper          перетворення DTO і моделей
internal/models          моделі бази даних
internal/repositories    PostgreSQL та in-memory репозиторії
internal/router          маршрути й middleware
docs                     Swagger та OpenAPI
```
