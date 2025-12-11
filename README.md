# 🤖 AI Browser Agent

Автономный AI-агент для управления веб-браузером на базе Claude AI.

## ✨ Возможности

- 🌐 **Автономное управление браузером** - навигация, клики, ввод текста
- 🧠 **Claude AI** - принятие решений на основе анализа страницы
- 🔄 **SmartClick** - умный клик с 4 fallback-стратегиями
- 📑 **Управление вкладками** - автоматическое обнаружение и переключение
- 🔒 **Security Layer** - подтверждение опасных операций (оплата, удаление)
- 📸 **Скриншоты** - автоматическое сохранение для отладки

## 🏗️ Архитектура

```
internal/
├── app/           # CLI приложение, DI контейнер
├── agent/         # Логика агента, выполнение действий
├── ai/            # Claude AI клиент, парсинг ответов
├── browser/       # Rod браузер, DOM, действия
├── security/      # Проверка безопасности, подтверждения
├── config/        # Конфигурация через env
├── domain/        # Доменные модели
└── llm/           # Абстракция LLM провайдеров
```

## 🚀 Быстрый старт

### Требования

- Go 1.24+
- Chrome/Chromium (устанавливается автоматически через Rod)

### Установка

```bash
git clone https://github.com/Daniil-Sakharov/BrowserAgent.git
cd BrowserAgent

# Установить зависимости
go mod download

# Скопировать конфиг
cp .env.example .env

# Добавить API ключ
nano .env  # ANTHROPIC_API_KEY=sk-ant-...

# Собрать
go build -o bin/agent ./cmd/agent
```

### Запуск

```bash
# Интерактивный режим
./bin/agent run

# С задачей
./bin/agent run "Открой google.com и найди информацию о Go"
```

## ⚙️ Конфигурация

Файл `.env`:

```env
# API
ANTHROPIC_API_KEY=sk-ant-xxx
ANTHROPIC_MODEL=claude-sonnet-4-5-20250929

# Браузер
BROWSER_HEADLESS=false      # true для Docker
BROWSER_TIMEOUT=30

# Агент
AGENT_MAX_STEPS=30
AGENT_INTERACTIVE=true

# Безопасность
SECURITY_ENABLED=true
SECURITY_AUTO_CONFIRM=false  # true = не спрашивать подтверждение
```

## 🐳 Docker

```bash
# Сборка
docker-compose -f deploy/compose/docker-compose.yml build

# Запуск
ANTHROPIC_API_KEY=sk-ant-xxx docker-compose -f deploy/compose/docker-compose.yml up
```

## 🛠️ Разработка

```bash
# Форматирование
go fmt ./...

# Линтинг
golangci-lint run ./...

# Или через Task
task format
task lint
task check  # format + lint
```

## 📝 Примеры

### Заказ еды
```bash
./bin/agent run "Закажи шаурму на Яндекс Еда"
```

### Поиск вакансий
```bash
./bin/agent run "Найди вакансии Go-разработчика на hh.ru"
```

### Работа с почтой
```bash
./bin/agent run "Открой яндекс почту и удали спам"
```

## 🔒 Безопасность

Security Layer автоматически обнаруживает опасные операции:

| Уровень | Действия | Поведение |
|---------|----------|-----------|
| **Critical** | Оплата, удаление аккаунта | Красное предупреждение + подтверждение |
| **High** | Удаление данных, отправка писем | Предупреждение + подтверждение |
| **Medium** | Отправка форм | Подтверждение |

При попытке оплаты появится диалог:
```
💳 ФИНАНСОВАЯ ОПЕРАЦИЯ - ПОДТВЕРДИТЕ ОПЛАТУ
┌────────────────────────────────────────┐
│ Действие: Клик                         │
│ Селектор: text:Pay                     │
│ Риск: КРИТИЧЕСКИЙ                      │
└────────────────────────────────────────┘
  Выполнить
▶ Отменить
```

## 📁 Структура проекта

```
.
├── cmd/agent/          # Точка входа
├── internal/           # Внутренние пакеты
├── pkg/                # Публичные пакеты (logger, closer)
├── deploy/
│   ├── docker/         # Dockerfile, .env.docker
│   └── compose/        # docker-compose.yml
├── screenshots/        # Скриншоты агента
├── logs/               # Логи
├── .env.example        # Шаблон конфигурации
├── .golangci.yml       # Конфиг линтера
└── Taskfile.yml        # Task runner
```

## 📄 Лицензия

MIT

## 👤 Автор

Daniil Sakharov
