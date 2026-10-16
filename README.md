# Отчёт по лабораторной работе

## CI/CD для Go-приложения: GitHub Actions + Docker

**Выполнил:** Потапов Иван 28ИПО8481
**Репозиторий:** https://github.com/agazsaz/my-go-app
**Дата:** 16.10.2026

---

## 1. Цель работы

Настроить полный CI/CD-пайплайн для простого Go-приложения: автоматический линтинг, тестирование с покрытием и сборку Docker-образа через GitHub Actions. Освоить локальную сборку и запуск Go-приложения через Docker.

---

## 2. Используемые инструменты

- **Go 1.22** — язык программирования
- **Docker Desktop** — контейнеризация
- **GitHub Actions** — CI/CD-платформа
- **golangci-lint** — статический анализатор кода
- **Git Bash** — командная оболочка

---

## 3. Структура проекта

```
my-go-app/
├── .github/
│   └── workflows/
│       └── ci.yml
├── main.go
├── sum.go
├── sum_test.go
├── Dockerfile
├── go.mod
└── README.md
```

| Файл | Назначение |
|---|---|
| `main.go` | Точка входа приложения |
| `sum.go` | Функция `Sum(a, b int) int` |
| `sum_test.go` | Юнит-тесты для функции `Sum` |
| `Dockerfile` | Многослойная сборка Docker-образа |
| `ci.yml` | Конфигурация GitHub Actions |
| `go.mod` | Манифест Go-модуля |

---

## 4. Исходный код

### main.go

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello from Go app!")
	fmt.Println("2 + 3 =", Sum(2, 3))
}
```

### sum.go

```go
package main

func Sum(a, b int) int {
	return a + b
}
```

### sum_test.go

```go
package main

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{-1, 1, 0},
		{0, 0, 0},
	}
	for _, tt := range tests {
		result := Sum(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("Sum(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}
```

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o my-app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/my-app .
CMD ["./my-app"]
```

### .github/workflows/ci.yml

```yaml
name: CI for Go App

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  test:
    name: Lint & Test
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

      - name: Run tests with coverage
        run: go test -v -coverprofile=coverage.out ./...

      - name: Upload coverage report
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.out

  docker-build:
    name: Build Docker Image (no push)
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v4
      - name: Build Docker image
        run: docker build -t my-go-app:test .
```

---

## 5. Ход работы

### 5.1. Создание репозитория

На GitHub создан публичный репозиторий `my-go-app` с файлом `README.md`. Репозиторий клонирован локально.

### 5.2. Инициализация Go-модуля

Так как Go локально не установлен, все операции выполнялись через Docker-контейнер:

```bash
docker run --rm -v "${PWD}:/app" -w /app golang:1.22-alpine go mod init my-go-app
```

Создан файл `go.mod`:

```
module my-go-app

go 1.22
```

### 5.3. Запуск юнит-тестов

```bash
docker run --rm -v "${PWD}:/app" -w /app golang:1.22-alpine go test ./...
```

**Результат:** `ok  my-go-app  ...` — все три теста пройдены.

<img width="951" height="81" alt="изображение" src="https://github.com/user-attachments/assets/00533e27-1376-4784-9f9b-4b10da3d5fc1" />


### 5.4. Локальная сборка Docker-образа

```bash
docker build -t my-go-app:latest .
docker run --rm my-go-app:latest
```
<img width="1225" height="721" alt="изображение" src="https://github.com/user-attachments/assets/3fa445d1-93c8-4d91-86a7-71215119d30a" />

**Результат:**

```
Hello from Go app!
2 + 3 = 5
```

<img width="485" height="108" alt="изображение" src="https://github.com/user-attachments/assets/15e17e57-9c0a-4556-828b-d58457688cfe" />


### 5.5. Пуш в GitHub и запуск CI

```bash
git add .
git commit -m "Fix sum.go and replace CI with reference workflow"
git push -u origin main
```

Автоматически запустился workflow **CI for Go App**.

<img width="1205" height="660" alt="изображение" src="https://github.com/user-attachments/assets/59c573b3-351a-4a7d-9fdd-ac8f30f2a862" />


---

## 6. Результаты

### 6.1. GitHub Actions — успешный запуск

<img width="951" height="81" alt="yml 1 part" src="https://github.com/user-attachments/assets/bbec8b09-8d31-4d9e-9833-476bd5348c0b" />


**Параметры запуска:**

- Статус: **Success**
- Общее время: **48 секунд**
- Артефактов: **1** (`coverage-report`)

**Джобы:**

- ✅ `Lint & Test` — 28 секунд
  - Checkout кода
  - Установка Go 1.22
  - Скачивание зависимостей
  - Запуск `golangci-lint`
  - Запуск тестов с профилем покрытия
  - Загрузка артефакта `coverage.out`
- ✅ `Build Docker Image (no push)` — 14 секунд
  - Checkout кода
  - Сборка Docker-образа `my-go-app:test`
  - Запуск возможен только после успешного прохождения `Lint & Test` (зависимость `needs: test`)

### 6.2. Список файлов репозитория
<img width="299" height="340" alt="изображение" src="https://github.com/user-attachments/assets/7d34a87e-eae6-49d7-9dba-2b8952f43767" />


### 6.3. Локальный запуск контейнера

[СКРИНШОТ: вывод docker run --rm my-go-app:latest — Hello from Go app! и 2 + 3 = 5]

Программа корректно выводит:

```
<img width="220" height="48" alt="изображение" src="https://github.com/user-attachments/assets/c320b2ae-e099-4ce2-885a-617fa92c900a" />

```

---

## 7. Выводы

В ходе лабораторной работы:

1. Освоен полный цикл CI/CD для Go-приложения с использованием GitHub Actions.
2. Настроены два последовательных job'а: тестирование и сборка Docker-образа.
3. Подключён статический анализатор `golangci-lint` для контроля качества кода.
4. Реализован сбор покрытия тестами с публикацией артефакта.
5. Освоена многослойная сборка Docker-образа: builder на `golang:1.22-alpine`, финальный образ на `alpine:latest` (минимальный размер).
6. Проверена работа пайплайна — все job'ы завершились со статусом **Success**.
7. Локально подтверждена работоспособность приложения через `docker run`.

Все поставленные задачи выполнены в полном объёме.
