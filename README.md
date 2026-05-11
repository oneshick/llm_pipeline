# LLM Pipeline — CSV → Groq API → JSON

Пайплайн на Go для автоматического извлечения характеристик товаров из текстовых описаний через Groq API (LLaMA 3).

---

## Структура проекта

```
llm_pipeline/
├── .env                  # конфигурация (API ключ и настройки)
├── products.csv          # входные данные
├── results.json          # результат (создаётся при запуске)
├── go.mod                # модуль Go
├── main.go               # точка входа
├── config/
│   └── config.go         # чтение .env и переменных окружения
├── orm/
│   └── models.go         # структуры данных
├── reader/
│   └── csv_reader.go     # чтение CSV-файла
├── llm/
│   └── groq_client.go    # HTTP-клиент к Groq API
├── processor/
│   └── processor.go      # формирование промпта, парсинг ответа, retry
└── writer/
    └── json_writer.go    # сохранение результата в JSON
```

---

## Как работает пайплайн

```
products.csv
     │
     ▼
 reader/        — читает CSV, возвращает []Product{ID, Description}
     │
     ▼
 processor/     — формирует промпт для каждого товара
     │
     ▼
 llm/           — отправляет POST на api.groq.com, получает JSON-ответ
     │
     ▼
 processor/     — парсит JSON из ответа LLM в структуру ProductFeatures
     │
     ▼
 writer/        — сохраняет итоговый массив с метаданными в results.json
```

---

## Быстрый старт

### 1. Установите Go

Скачайте и установите с [go.dev/dl](https://go.dev/dl/). Проверьте:

```powershell
go version
```

### 2. Получите API ключ Groq

Зарегистрируйтесь на [console.groq.com](https://console.groq.com/keys) и создайте ключ. Ключ начинается на `gsk_...`

### 3. Вставьте ключ в .env

Откройте файл `.env` и замените:

```
GROQ_API_KEY=gsk_ВАШ_КЛЮЧ_СЮДА
```

### 4. Запустите

```powershell
cd llm_pipeline

go mod tidy
go build -o pipeline.exe .
.\pipeline.exe
```

---

## Формат входного CSV

Файл `products.csv` должен содержать два столбца:

```csv
id,description
1,"Название товара. Цена: 9 990 руб. Описание характеристик..."
2,"Другой товар. Цена: 4 500 руб. ..."
```

---

## Формат выходного JSON

```json
{
  "generated_at": "2026-05-09T14:22:01Z",
  "model": "llama-3.1-8b-instant",
  "total_items": 5,
  "processed": 5,
  "failed": 0,
  "products": [
    {
      "id": "1",
      "name": "Apple iPhone 15 Pro 256GB Natural Titanium",
      "brand": "Apple",
      "category": "Смартфоны",
      "price_rub": 89990,
      "currency": "RUB",
      "key_specs": [
        "Чип A17 Pro",
        "Камера 48 МП",
        "Экран 6.1\" Super Retina XDR"
      ]
    }
  ]
}
```

---

## Настройки (.env)

| Переменная       | По умолчанию           | Описание                                  |
|------------------|------------------------|-------------------------------------------|
| `GROQ_API_KEY`   | —                      | API ключ (обязательно)                    |
| `GROQ_MODEL`     | `llama-3.1-8b-instant` | Модель Groq                               |
| `GROQ_MAX_TOKENS`| `512`                  | Максимум токенов в ответе                 |
| `CSV_PATH`       | `products.csv`         | Путь к входному файлу                     |
| `OUT_PATH`       | `results.json`         | Путь к выходному файлу                    |
| `DELAY_MS`       | `300`                  | Пауза между запросами (защита rate limit) |

---

## Возможные ошибки

| Ошибка | Причина | Решение |
|--------|---------|---------|
| `package llm_pipeline/config is not in std` | Нет `go.mod` | Выполните `go mod init llm_pipeline` |
| `API ключ не задан` | Пустой `.env` | Вставьте ключ в `.env` |
| `Rate limit reached` | Исчерпан дневной лимит токенов | Подождите до завтра или смените модель |
| `EOF` | Нет интернета или неверный URL | Проверьте соединение |
| `products.csv: no such file` | Файл не в папке | Положите `products.csv` рядом с `main.go` |
