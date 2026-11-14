# Конфигурационные параметры CodeAuthorityFirst

## Обзор

Документация описывает все конфигурационные параметры приложения CodeAuthorityFirst. Приложение использует два основных файла конфигурации:

1. **`.caf/config.yml`** - Основная конфигурация CAF (Code Authority Framework)
2. **`.caf/security.yml`** - Конфигурация безопасности и системы аудита

## Структура конфигурации

### Основная конфигурация (.caf/config.yml)

```yaml
version: "0.4.0-triz2"
triviality_threshold_loc: 50
ui:
  theme: "dark"
  language: "ru"
lsp:
  strict: true
security:
  signed_adapters_required: false
  sandbox_parser: true
  audit_log: true
```

### Конфигурация безопасности (.caf/security.yml)

```yaml
# Конфигурация безопасности CodeAuthorityFirst
yaml:
  sandbox_parser: true
  max_nesting_depth: 10
  max_document_size: 1048576  # 1MB
  allow_custom_tags: false
  validation_mode: "strict"
  required_fields: []
  blocked_patterns: []

secrets:
  enabled: true
  max_line_length: 1000
  min_secret_length: 8
  max_secret_length: 512
  severity_threshold: 0.7
  enable_heuristics: true
  enable_entropy_check: true
  enable_context_check: true
  ignore_files:
    - ".git"
    - ".gitignore"
    - "node_modules"
  custom_rules: []

audit:
  enabled: true
  log_to_console: true
  log_to_file: true
  log_level: "info"
  max_file_size: 10485760  # 10MB
  max_files: 5
  retain_days: 30
  output_dir: "audit_logs"
  enable_json: true
  enable_structured: true

authentication:
  signed_adapters_required: false
  token:
    min_length: 16
    max_length: 256
    allowed_characters: "a-zA-Z0-9-_."
    expiration: 86400  # 24 часа
  session:
    timeout: 86400  # 24 часа
    max_concurrent: 5
    encryption: true
    cleanup_interval: 3600  # 1 час
```

## Подробное описание параметров

### Основная конфигурация (.caf/config.yml)

#### version
- **Тип**: String
- **Описание**: Версия CAF фреймворка
- **Значение по умолчанию**: `"0.4.0-triz2"`
- **Валидация**: Должен соответствовать формату версии

#### triviality_threshold_loc
- **Тип**: Integer
- **Описание**: Порог тривиальности для строк кода (Lines of Code)
- **Значение по умолчанию**: `50`
- **Валидация**: Должен быть положительным целым числом

#### ui
- **Тип**: Object
- **Описание**: Настройки пользовательского интерфейса

##### ui.theme
- **Тип**: String
- **Описание**: Тема интерфейса
- **Допустимые значения**: `"dark"`, `"light"`, `"auto"`
- **Значение по умолчанию**: `"dark"`

##### ui.language
- **Тип**: String
- **Описание**: Язык интерфейса
- **Допустимые значения**: `"ru"`, `"en"`, `"es"`, `"fr"`
- **Значение по умолчанию**: `"ru"`

#### lsp
- **Тип**: Object
- **Описание**: Настройки Language Server Protocol

##### lsp.strict
- **Тип**: Boolean
- **Описание**: Строгий режим валидации LSP
- **Значение по умолчанию**: `true`

#### security
- **Тип**: Object
- **Описание**: Основные настройки безопасности

##### security.signed_adapters_required
- **Тип**: Boolean
- **Описание**: Требуются ли подписанные адаптеры
- **Значение по умолчанию**: `false`

##### security.sandbox_parser
- **Тип**: Boolean
- **Описание**: Использовать песочницу для парсинга
- **Значение по умолчанию**: `true`

##### security.audit_log
- **Тип**: Boolean
- **Описание**: Включить аудит логирование
- **Значение по умолчанию**: `true`

### Конфигурация безопасности (.caf/security.yml)

#### yaml
- **Тип**: Object
- **Описание**: Настройки YAML безопасности

##### yaml.sandbox_parser
- **Тип**: Boolean
- **Описание**: Использовать песочницу для парсинга YAML
- **Значение по умолчанию**: `true`
- **Рекомендации**: Всегда включайте в продакшене

##### yaml.max_nesting_depth
- **Тип**: Integer
- **Описание**: Максимальная глубина вложенности YAML
- **Значение по умолчанию**: `10`
- **Валидация**: Должен быть в диапазоне 1-100
- **Рекомендации**: Увеличьте для сложных конфигураций

##### yaml.max_document_size
- **Тип**: Integer
- **Описание**: Максимальный размер YAML документа в байтах
- **Значение по умолчанию**: `1048576` (1MB)
- **Валидация**: Должен быть в диапазоне 1024-104857600

##### yaml.allow_custom_tags
- **Тип**: Boolean
- **Описание**: Разрешить кастомные YAML теги
- **Значение по умолчанию**: `false`
- **Рекомендации**: Отключите для безопасности

##### yaml.validation_mode
- **Тип**: String
- **Описание**: Режим валидации YAML
- **Допустимые значения**: `"basic"`, `"strict"`, `"paranoid"`
- **Значение по умолчанию**: `"strict"`
- **Рекомендации**: Используйте `"paranoid"` для высоких требований безопасности

##### yaml.required_fields
- **Тип**: Array[String]
- **Описание**: Обязательные поля в YAML документах
- **Значение по умолчанию**: `[]`
- **Пример**: `["name", "version", "config"]`

##### yaml.blocked_patterns
- **Тип**: Array[String]
- **Описание**: Блокированные паттерны в YAML содержимом
- **Значение по умолчанию**: `[]`
- **Пример**: `["password", "secret", "token"]`

#### secrets
- **Тип**: Object
- **Описание**: Настройки сканирования секретов

##### secrets.enabled
- **Тип**: Boolean
- **Описание**: Включить сканирование секретов
- **Значение по умолчанию**: `true`
- **Рекомендации**: Всегда включайте в продакшене

##### secrets.max_line_length
- **Тип**: Integer
- **Описание**: Максимальная длина строки для сканирования
- **Значение по умолчанию**: `1000`
- **Валидация**: Должен быть в диапазоне 100-10000

##### secrets.min_secret_length
- **Тип**: Integer
- **Описание**: Минимальная длина секрета
- **Значение по умолчанию**: `8`
- **Валидация**: Должен быть в диапазоне 4-1024

##### secrets.max_secret_length
- **Тип**: Integer
- **Описание**: Максимальная длина секрета
- **Значение по умолчанию**: `512`
- **Валидация**: Должен быть в диапазоне 8-2048

##### secrets.severity_threshold
- **Тип**: Float
- **Описание**: Порог серьезности для обнаружения секретов
- **Значение по умолчанию**: `0.7`
- **Валидация**: Должен быть в диапазоне 0.0-1.0
- **Рекомендации**: Уменьшите для более строгого сканирования

##### secrets.enable_heuristics
- **Тип**: Boolean
- **Описание**: Включить эвристическое сканирование
- **Значение по умолчанию**: `true`
- **Рекомендации**: Включите для лучшего обнаружения

##### secrets.enable_entropy_check
- **Тип**: Boolean
- **Описание**: Включить проверку энтропии
- **Значение по умолчанию**: `true`
- **Рекомендации**: Включите для обнаружения случайных строк

##### secrets.enable_context_check
- **Тип**: Boolean
- **Описание**: Включить контекстную проверку
- **Значение по умолчанию**: `true`
- **Рекомендации**: Включите для точного обнаружения

##### secrets.ignore_files
- **Тип**: Array[String]
- **Описание**: Игнорируемые файлы при сканировании
- **Значение по умолчанию**: `[".git", ".gitignore", "node_modules"]`
- **Пример**: `["*.test.js", "config/local.yml"]`

##### secrets.custom_rules
- **Тип**: Array[Object]
- **Описание**: Кастомные правила обнаружения секретов
- **Значение по умолчанию**: `[]`
- **Структура**:
  ```yaml
  - name: "custom_api_key"
    pattern: "API_KEY_[A-Z0-9]{20,}"
    severity: 0.8
    description: "Custom API key pattern"
  ```

#### audit
- **Тип**: Object
- **Описание**: Настройки аудита и логирования

##### audit.enabled
- **Тип**: Boolean
- **Описание**: Включить аудит логирование
- **Значение по умолчанию**: `true`
- **Рекомендации**: Всегда включайте в продакшене

##### audit.log_to_console
- **Тип**: Boolean
- **Описание**: Логировать в консоль
- **Значение по умолчанию**: `true`
- **Рекомендации**: Отключите в продакшене для производительности

##### audit.log_to_file
- **Тип**: Boolean
- **Описание**: Логировать в файл
- **Значение по умолчанию**: `true`
- **Рекомендации**: Всегда включайте

##### audit.log_level
- **Тип**: String
- **Описание**: Уровень логирования
- **Допустимые значения**: `"debug"`, `"info"`, `"warning"`, `"error"`, `"critical"`
- **Значение по умолчанию**: `"info"`
- **Рекомендации**: Используйте `"warning"` или выше в продакшене

##### audit.max_file_size
- **Тип**: Integer
- **Описание**: Максимальный размер лог файла в байтах
- **Значение по умолчанию**: `10485760` (10MB)
- **Валидация**: Должен быть в диапазоне 1024-104857600

##### audit.max_files
- **Тип**: Integer
- **Описание**: Максимальное количество лог файлов
- **Значение по умолчанию**: `5`
- **Валидация**: Должен быть в диапазоне 1-100

##### audit.retain_days
- **Тип**: Integer
- **Описание**: Количество дней хранения логов
- **Значение по умолчанию**: `30`
- **Валидация**: Должен быть в диапазоне 1-365

##### audit.output_dir
- **Тип**: String
- **Описание**: Директория для вывода логов
- **Значение по умолчанию**: `"audit_logs"`
- **Рекомендации**: Используйте отдельную файловую систему для логов

##### audit.enable_json
- **Тип**: Boolean
- **Описание**: Включить JSON формат логирования
- **Значение по умолчанию**: `true`
- **Рекомендации**: Включите для автоматической обработки

##### audit.enable_structured
- **Тип**: Boolean
- **Описание**: Включить структурированное логирование
- **Значение по умолчанию**: `true`
- **Рекомендации**: Включите для лучшей аналитики

#### authentication
- **Тип**: Object
- **Описание**: Настройки аутентификации

##### authentication.signed_adapters_required
- **Тип**: Boolean
- **Описание**: Требуются ли подписанные адаптеры
- **Значение по умолчанию**: `false`
- **Рекомендации**: Включите в продакшене

##### authentication.token
- **Тип**: Object
- **Описание**: Настройки токенов

###### authentication.token.min_length
- **Тип**: Integer
- **Описание**: Минимальная длина токена
- **Значение по умолчанию**: `16`
- **Валидация**: Должен быть в диапазоне 8-1024

###### authentication.token.max_length
- **Тип**: Integer
- **Описание**: Максимальная длина токена
- **Значение по умолчанию**: `256`
- **Валидация**: Должен быть >= min_length

###### authentication.token.allowed_characters
- **Тип**: String
- **Описание**: Разрешенные символы в токене
- **Значение по умолчанию**: `"a-zA-Z0-9-_."`
- **Рекомендации**: Используйте строгие правила для безопасности

###### authentication.token.expiration
- **Тип**: Integer
- **Описание**: Время жизни токена в секундах
- **Значение по умолчанию**: `86400` (24 часа)
- **Валидация**: Должен быть в диапазоне 60-2592000

##### authentication.session
- **Тип**: Object
- **Описание**: Настройки сессий

###### authentication.session.timeout
- **Тип**: Integer
- **Описание**: Таймаут сессии в секундах
- **Значение по умолчанию**: `86400` (24 часа)
- **Валидация**: Должен быть в диапазоне 60-2592000

###### authentication.session.max_concurrent
- **Тип**: Integer
- **Описание**: Максимальное количество одновременных сессий
- **Значение по умолчанию**: `5`
- **Валидация**: Должен быть в диапазоне 1-100

###### authentication.session.encryption
- **Тип**: Boolean
- **Описание**: Шифрование сессий
- **Значение по умолчанию**: `true`
- **Рекомендации**: Всегда включайте

###### authentication.session.cleanup_interval
- **Тип**: Integer
- **Описание**: Интервал очистки сессий в секундах
- **Значение по умолчанию**: `3600` (1 час)
- **Валидация**: Должен быть в диапазоне 60-86400

## Примеры конфигураций

### Минимальная конфигурация для разработки
```yaml
# .caf/config.yml
version: "0.4.0-triz2"
triviality_threshold_loc: 50
ui:
  theme: "dark"
  language: "ru"
lsp:
  strict: false
security:
  signed_adapters_required: false
  sandbox_parser: true
  audit_log: true

# .caf/security.yml
yaml:
  sandbox_parser: true
  max_nesting_depth: 10
  max_document_size: 1048576
  allow_custom_tags: false
  validation_mode: "basic"
  required_fields: []
  blocked_patterns: []

secrets:
  enabled: true
  max_line_length: 1000
  min_secret_length: 8
  max_secret_length: 512
  severity_threshold: 0.5
  enable_heuristics: true
  enable_entropy_check: true
  enable_context_check: true
  ignore_files: [".git", ".gitignore", "node_modules"]
  custom_rules: []

audit:
  enabled: true
  log_to_console: true
  log_to_file: true
  log_level: "debug"
  max_file_size: 10485760
  max_files: 3
  retain_days: 7
  output_dir: "audit_logs"
  enable_json: true
  enable_structured: true

authentication:
  signed_adapters_required: false
  token:
    min_length: 16
    max_length: 256
    allowed_characters: "a-zA-Z0-9-_."
    expiration: 3600
  session:
    timeout: 3600
    max_concurrent: 3
    encryption: true
    cleanup_interval: 1800
```

### Максимальная безопасность для продакшена
```yaml
# .caf/config.yml
version: "0.4.0-triz2"
triviality_threshold_loc: 100
ui:
  theme: "dark"
  language: "ru"
lsp:
  strict: true
security:
  signed_adapters_required: true
  sandbox_parser: true
  audit_log: true

# .caf/security.yml
yaml:
  sandbox_parser: true
  max_nesting_depth: 5
  max_document_size: 524288  # 512KB
  allow_custom_tags: false
  validation_mode: "paranoid"
  required_fields: ["name", "version", "environment"]
  blocked_patterns: ["password", "secret", "token", "key"]

secrets:
  enabled: true
  max_line_length: 500
  min_secret_length: 12
  max_secret_length: 256
  severity_threshold: 0.9
  enable_heuristics: true
  enable_entropy_check: true
  enable_context_check: true
  ignore_files: [".git", ".gitignore", "node_modules", "tests"]
  custom_rules:
    - name: "aws_key"
      pattern: "AKIA[0-9A-Z]{16}"
      severity: 1.0
      description: "AWS Access Key"
    - name: "github_token"
      pattern: "ghp_[a-zA-Z0-9]{36}"
      severity: 1.0
      description: "GitHub Personal Access Token"

audit:
  enabled: true
  log_to_console: false
  log_to_file: true
  log_level: "warning"
  max_file_size: 5242880  # 5MB
  max_files: 10
  retain_days: 90
  output_dir: "/var/log/orya/audit"
  enable_json: true
  enable_structured: true

authentication:
  signed_adapters_required: true
  token:
    min_length: 32
    max_length: 128
    allowed_characters: "a-zA-Z0-9-_."
    expiration: 1800  # 30 минут
  session:
    timeout: 3600  # 1 час
    max_concurrent: 1
    encryption: true
    cleanup_interval: 300  # 5 минут
```

## Валидация конфигурации

Приложение автоматически валидирует конфигурационные файлы при запуске. Возможные ошибки:

### Ошибки YAML синтаксиса
```
Error: Failed to parse YAML: yaml: line 10: found character that cannot start any token
```

### Ошибки валидации параметров
```
Error: Invalid security configuration: secrets.severity_threshold must be between 0.0 and 1.0
```

### Ошибки типов данных
```
Error: Type mismatch for parameter 'triviality_threshold_loc': expected integer, got string
```

## Управление конфигурацией

### Автоматическая генерация
Скрипты запуска автоматически создают конфигурационные файлы с настройками по умолчанию, если они отсутствуют.

### Ручная настройка
Вы можете вручную отредактировать файлы конфигурации в директории `.caf/`.

### Переопределение конфигурации
Конфигурация может быть переопределена через переменные окружения:
```bash
export CAF_CONFIG_PATH="/custom/path/to/config.yml"
export CAF_SECURITY_CONFIG_PATH="/custom/path/to/security.yml"
```

### Горячее обновление конфигурации
Некоторые параметры могут быть обновлены без перезапуска приложения:
- Уровни логирования
- Пороги серьезности
- Настройки аудита

Для этого используйте команду:
```bash
./caf config reload
```

## Мониторинг и логирование

### Логи приложения
Логи приложения сохраняются в:
- Консоль (если включено)
- Файлы в директории `audit_logs/`
- Системный журнал (если настроено)

### Структура логов
```
audit_logs/
├── app_2024-01-15.log          # Основной лог приложения
├── security_2024-01-15.log     # Логи безопасности
├── audit_2024-01-15.json       # Структурированные логи в JSON
└── state/
    ├── app_2024-01-15.json     # Состояние приложения
    └── stop_2024-01-15.json    # Логи остановки
```

### Анализ логов
Для анализа логов можно использовать:
```bash
# Просмотр последних логов
tail -f audit_logs/app_*.log

# Поиск ошибок
grep "ERROR" audit_logs/app_*.log

# Анализ JSON логов
jq '.' audit_logs/audit_2024-01-15.json | less
```

## Безопасность конфигурации

### Защита конфигурационных файлов
```
# Установка правильных прав доступа
chmod 600 .caf/config.yml
chmod 600 .caf/security.yml

# Шифрование конфигурации (опционально)
openssl enc -aes-256-cbc -salt -in .caf/config.yml -out .caf/config.yml.enc
```

### Управление секретами
Не храните секреты в конфигурационных файлах. Используйте:
- Переменные окружения
- Системы управления секретами (Vault, AWS Secrets Manager)
- Файлы с ограниченным доступом

### Регулярный аудит
Регулярно проверяйте конфигурацию на соответствие политикам безопасности:
```bash
# Аудит конфигурации
./caf security audit-config

# Проверка секретов
./caf security scan-config
```

## Тестирование конфигурации

### Валидация конфигурации
```bash
# Проверка синтаксиса
./caf config validate

# Тестирование конфигурации
./caf config test

# Dry-run запуск
./caf start --dry-run
```

### Интеграционное тестирование
```bash
# Тестирование безопасности
./caf security test

# Тестирование аудита
./caf audit test

# Тестирование аутентификации
./caf auth test
```

## Часто задаваемые вопросы

### Q: Как изменить путь к конфигурационным файлам?
A: Используйте переменные окружения `CAF_CONFIG_PATH` и `CAF_SECURITY_CONFIG_PATH`.

### Q: Как включить подробное логирование?
A: Установите `audit.log_level: "debug"` в конфигурации безопасности.

### Q: Как увеличить лимит размера YAML документов?
A: Измените параметр `yaml.max_document_size` в конфигурации безопасности.

### Q: Как добавить кастомное правило обнаружения секретов?
A: Добавьте запись в секцию `secrets.custom_rules` в конфигурации безопасности.

### Q: Как настроить ротацию логов?
A: Используйте параметры `audit.max_file_size` и `audit.max_files` для автоматической ротации.

## Версионирование конфигурации

Каждое изменение конфигурации должно сопровождаться:
1. Обновлением версии в `config.yml`
2. Документированием изменений
3. Тестированием новой конфигурации
4. Резервным копированием старой конфигурации

Пример:
```yaml
# Версия конфигурации
config_version: "1.2.0"
last_updated: "2024-01-15T10:00:00Z"
changes:
  - "Увеличен лимит размера YAML документов до 2MB"
  - "Добавлено новое правило обнаружения API ключей"
```

## Заключение

Правильная конфигурация безопасности является критически важной для работы приложения CodeAuthorityFirst. Регулярно обновляйте конфигурацию в соответствии с изменяющимися требованиями безопасности и следуйте рекомендациям, изложенным в этой документации.