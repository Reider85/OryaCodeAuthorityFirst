#!/bin/bash

# start-app.sh - Полный скрипт запуска приложения CodeAuthorityFirst
# Поддерживает Linux/macOS системы

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

APP_NAME="CodeAuthorityFirst"
APP_BINARY="main_security_demo"
CONFIG_DIR=".caf"
AUDIT_DIR="audit_logs"
PID_FILE="$CONFIG_DIR/app.pid"
LOG_FILE="$CONFIG_DIR/app.log"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции вывода
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Функция проверки зависимостей
check_dependencies() {
    log_info "Проверка зависимостей..."
    
    # Проверка Go
    if ! command -v go &> /dev/null; then
        log_error "Go не установлен. Пожалуйста, установите Go 1.27 или выше."
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    REQUIRED_VERSION="1.27"
    
    if ! [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
        log_error "Требуется Go $REQUIRED_VERSION или выше. Установлена версия: $GO_VERSION"
        exit 1
    fi
    
    log_success "Go версии $GO_VERSION найден"
    
    # Проверка git (для сборки с версионированием)
    if command -v git &> /dev/null; then
        log_success "Git найден"
    else
        log_warning "Git не найден. Сборка будет без информации о версии."
    fi
    
    # Проверка Python (для некоторых компонентов)
    if command -v python3 &> /dev/null; then
        log_success "Python3 найден"
    else
        log_warning "Python3 не найден. Некоторые компоненты могут не работать."
    fi
}

# Функция сборки приложения
build_application() {
    log_info "Сборка приложения $APP_NAME..."
    
    # Проверка существующего бинарного файла
    if [ -f "$APP_BINARY" ]; then
        log_info "Бинарный файл существует. Проверка версии..."
        
        # Проверка версии через бинарный файл
        if ./$APP_BINARY --version &> /dev/null; then
            CURRENT_VERSION=$(./$APP_BINARY --version 2>/dev/null || echo "unknown")
            log_info "Текущая версия: $CURRENT_VERSION"
            
            # Проверка необходимости пересборки
            if [ -f ".git" ]; then
                GIT_CHANGED=$(git status --porcelain 2>/dev/null | wc -l)
                if [ "$GIT_CHANGED" -gt 0 ]; then
                    log_warning "Найдены изменения в Git. Рекомендуется пересобрать приложение."
                    read -p "Пересобрать приложение? (y/N): " -n 1 -r
                    echo
                    if [[ $REPLY =~ ^[Yy]$ ]]; then
                        rm -f "$APP_BINARY"
                    else
                        log_info "Используется существующий бинарный файл"
                        return 0
                    fi
                fi
            fi
        else
            log_warning "Существующий бинарный файл не работает. Удаление и пересборка..."
            rm -f "$APP_BINARY"
        fi
    fi
    
    # Сборка приложения
    if [ -f "scripts/build.sh" ]; then
        log_info "Использование скрипта сборки..."
        bash scripts/build.sh
    else
        log_info "Ручная сборка..."
        
        # Получение версии из git или использование fallback
        if command -v git &> /dev/null && [ -d .git ]; then
            GIT_VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo "unknown")
            VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GIT_VERSION'"
            log_info "Версия из git: $GIT_VERSION"
        else
            VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
            log_warning "Git не доступен, используется fallback версия"
        fi
        
        # Флаги сборки
        LDFLAGS="$VERSION_FLAG -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit=$(git rev-parse HEAD 2>/dev/null || echo unknown)'"
        LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.datetime=$(date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo unknown)'"
        LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.modified=$(git diff-index --quiet HEAD -- 2>/dev/null || echo true)'"
        
        log_info "Сборка с флагами: $LDFLAGS"
        
        # Компиляция
        go build -ldflags="$LDFLAGS" -o "$APP_BINARY" ./main_security_demo
        
        if [ $? -eq 0 ]; then
            log_success "Приложение успешно собрано!"
        else
            log_error "Ошибка сборки приложения!"
            exit 1
        fi
    fi
}

# Функция создания директорий
create_directories() {
    log_info "Создание необходимых директорий..."
    
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$AUDIT_DIR"
    mkdir -p "$CONFIG_DIR/security"
    
    log_success "Директории созданы"
}

# Функция генерации конфигурации
generate_config() {
    log_info "Генерация конфигурации..."
    
    # Основная конфигурация CAF
    if [ ! -f "$CONFIG_DIR/config.yml" ]; then
        cat > "$CONFIG_DIR/config.yml" << EOF
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
EOF
        log_success "Основная конфигурация CAF создана: $CONFIG_DIR/config.yml"
    else
        log_info "Основная конфигурация CAF уже существует"
    fi
    
    # Конфигурация безопасности
    if [ ! -f "$CONFIG_DIR/security.yml" ]; then
        cat > "$CONFIG_DIR/security.yml" << EOF
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
EOF
        log_success "Конфигурация безопасности создана: $CONFIG_DIR/security.yml"
    else
        log_info "Конфигурация безопасности уже существует"
    fi
}

# Функция проверки конфигурации
validate_config() {
    log_info "Валидация конфигурации..."
    
    # Проверка YAML синтаксиса
    if command -v python3 &> /dev/null; then
        python3 -c "
import yaml
import sys
try:
    with open('$CONFIG_DIR/config.yml', 'r') as f:
        yaml.safe_load(f)
    with open('$CONFIG_DIR/security.yml', 'r') as f:
        yaml.safe_load(f)
    print('✓ YAML синтаксис корректен')
except Exception as e:
    print(f'✗ Ошибка YAML синтаксиса: {e}')
    sys.exit(1)
" 2>/dev/null || log_warning "Python3 недоступен для проверки YAML синтаксиса"
    fi
    
    log_success "Конфигурация валидна"
}

# Функция запуска приложения
start_application() {
    log_info "Запуск приложения $APP_NAME..."
    
    # Проверка, запущено ли уже приложение
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            log_warning "Приложение уже запущено с PID: $PID"
            log_info "Используйте './stop-app.sh' для остановки"
            return 0
        else
            log_warning "PID файл существует, но процесс не найден. Удаление PID файла..."
            rm -f "$PID_FILE"
        fi
    fi
    
    # Запуск приложения в фоновом режиме
    nohup ./$APP_BINARY > "$LOG_FILE" 2>&1 &
    PID=$!
    
    # Сохранение PID
    echo "$PID" > "$PID_FILE"
    
    # Проверка запуска
    sleep 2
    if ps -p "$PID" > /dev/null 2>&1; then
        log_success "Приложение запущено успешно!"
        log_info "PID: $PID"
        log_info "Лог файл: $LOG_FILE"
        log_info "Директория аудита: $AUDIT_DIR"
        
        # Вывод информации о процессе
        echo "=========================================="
        echo "Приложение: $APP_NAME"
        echo "PID: $PID"
        echo "Лог файл: $LOG_FILE"
        echo "Конфигурация: $CONFIG_DIR"
        echo "Аудит логи: $AUDIT_DIR"
        echo "=========================================="
    else
        log_error "Ошибка запуска приложения!"
        if [ -f "$LOG_FILE" ]; then
            log_info "Проверьте лог файл: $LOG_FILE"
            echo "Последние 10 строк лога:"
            tail -n 10 "$LOG_FILE"
        fi
        rm -f "$PID_FILE"
        exit 1
    fi
}

# Функция настройки обработчика сигналов
setup_signal_handlers() {
    # Функция очистки при выходе
    cleanup() {
        log_info "Остановка приложения..."
        if [ -f "$PID_FILE" ]; then
            PID=$(cat "$PID_FILE")
            if ps -p "$PID" > /dev/null 2>&1; then
                kill "$PID" 2>/dev/null
                log_info "Отправлен сигнал STOP процессу $PID"
                
                # Ожидание завершения
                for i in {1..10}; do
                    if ! ps -p "$PID" > /dev/null 2>&1; then
                        break
                    fi
                    sleep 1
                done
                
                # Принудительное завершение, если необходимо
                if ps -p "$PID" > /dev/null 2>&1; then
                    log_warning "Принудительное завершение процесса $PID"
                    kill -9 "$PID" 2>/dev/null
                fi
            fi
            rm -f "$PID_FILE"
        fi
        log_success "Приложение остановлено"
        exit 0
    }
    
    # Настройка обработчиков сигналов
    trap cleanup SIGINT SIGTERM
}

# Основная функция
main() {
    echo "=========================================="
    echo "Запуск приложения $APP_NAME"
    echo "=========================================="
    
    # Проверка, что мы в корневой директории проекта
    if [ ! -f "go.mod" ]; then
        log_error "Скрипт должен быть запущен из корневой директории проекта"
        exit 1
    fi
    
    # Выполнение всех шагов
    check_dependencies
    build_application
    create_directories
    generate_config
    validate_config
    setup_signal_handlers
    start_application
    
    log_success "Приложение $APP_NAME полностью запущено!"
    
    # Ожидание завершения
    while true; do
        sleep 1
    done
}

# Запуск основной функции
main "$@"