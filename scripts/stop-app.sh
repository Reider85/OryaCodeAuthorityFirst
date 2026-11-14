#!/bin/bash

# stop-app.sh - Скрипт полной остановки приложения CodeAuthorityFirst
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
STATE_FILE="$CONFIG_DIR/app.state"

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

# Функция поиска процессов приложения
find_processes() {
    local pids=()
    
    # Поиск по имени бинарного файла
    if [ -f "$APP_BINARY" ]; then
        local binary_pids=$(pgrep -f "$APP_BINARY" || true)
        if [ -n "$binary_pids" ]; then
            while read -r pid; do
                pids+=("$pid")
            done <<< "$binary_pids"
        fi
    fi
    
    # Поиск по PID файлу
    if [ -f "$PID_FILE" ]; then
        local pid_from_file=$(cat "$PID_FILE" 2>/dev/null || true)
        if [ -n "$pid_from_file" ] && ps -p "$pid_from_file" > /dev/null 2>&1; then
            # Проверяем, что PID еще не в списке
            local found=false
            for pid in "${pids[@]}"; do
                if [ "$pid" = "$pid_from_file" ]; then
                    found=true
                    break
                fi
            done
            if [ "$found" = false ]; then
                pids+=("$pid_from_file")
            fi
        else
            log_warning "PID файл существует, но процесс не найден. Удаление PID файла..."
            rm -f "$PID_FILE"
        fi
    fi
    
    # Поиск по ключевым словам в командной строке
    local keyword_pids=$(pgrep -f "CodeAuthorityFirst\|security.*demo" || true)
    if [ -n "$keyword_pids" ]; then
        while read -r pid; do
            local found=false
            for existing_pid in "${pids[@]}"; do
                if [ "$pid" = "$existing_pid" ]; then
                    found=true
                    break
                fi
            done
            if [ "$found" = false ]; then
                pids+=("$pid")
            fi
        done <<< "$keyword_pids"
    fi
    
    echo "${pids[@]}"
}

# Функция остановки процесса
stop_process() {
    local pid=$1
    local force=$2
    
    log_info "Остановка процесса $pid..."
    
    if [ "$force" = true ]; then
        # Принудительное завершение
        kill -9 "$pid" 2>/dev/null
        log_info "Отправлен сигнал KILL процессу $pid"
    else
        # Мягкое завершение
        kill "$pid" 2>/dev/null
        log_info "Отправлен сигнал TERM процессу $pid"
        
        # Ожидание завершения
        for i in {1..10}; do
            if ! ps -p "$pid" > /dev/null 2>&1; then
                break
            fi
            sleep 1
        done
        
        # Принудительное завершение, если необходимо
        if ps -p "$pid" > /dev/null 2>&1; then
            log_warning "Принудительное завершение процесса $pid"
            kill -9 "$pid" 2>/dev/null
        fi
    fi
    
    # Проверка результата
    if ps -p "$pid" > /dev/null 2>&1; then
        log_error "Не удалось остановить процесс $pid"
        return 1
    else
        log_success "Процесс $pid остановлен"
        return 0
    fi
}

# Функция сохранения состояния
save_state() {
    log_info "Сохранение состояния приложения..."
    
    # Создание директории состояния
    mkdir -p "$CONFIG_DIR/state"
    
    local timestamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
    local state_file="$CONFIG_DIR/state/stop_$timestamp.json"
    
    local state='{
        "stopped_at": "'"$timestamp"'",
        "app_name": "'"$APP_NAME"'",
        "audit_dir": "'"$AUDIT_DIR"'",
        "config_dir": "'"$CONFIG_DIR"'",
        "processes": []
    }'
    
    # Добавление информации о процессах
    local pids=($(find_processes))
    for pid in "${pids[@]}"; do
        if ps -p "$pid" > /dev/null 2>&1; then
            local process_info=$(ps -p "$pid" -o pid=,ppid=,cmd= 2>/dev/null | tail -n 1)
            local cmd_info=$(ps -p "$pid" -o args= 2>/dev/null | tail -n 1)
            
            state=$(echo "$state" | jq --arg pid "$pid" --arg ppid "$(ps -p "$pid" -o ppid= 2>/dev/null)" --arg cmd "$cmd_info" \
                '.processes += [{"pid": $pid, "ppid": $ppid, "cmd": $cmd}]' 2>/dev/null || echo "$state")
        fi
    done
    
    echo "$state" | jq '.' 2>/dev/null > "$state_file" || echo "$state" > "$state_file"
    
    log_success "Состояние сохранено: $state_file"
}

# Функция очистки
cleanup() {
    log_info "Очистка временных файлов..."
    
    # Удаление PID файла
    if [ -f "$PID_FILE" ]; then
        rm -f "$PID_FILE"
        log_info "PID файл удален"
    fi
    
    # Очистка временных файлов приложения
    local temp_files=$(find . -name "*.tmp" -o -name "*.temp" -o -name "*~" 2>/dev/null | head -20)
    if [ -n "$temp_files" ]; then
        echo "$temp_files" | xargs rm -f 2>/dev/null || true
        log_info "Временные файлы очищены"
    fi
    
    # Очистка кэша (если есть)
    if [ -d "$CONFIG_DIR/cache" ]; then
        rm -rf "$CONFIG_DIR/cache"/* 2>/dev/null || true
        log_info "Кэш очищен"
    fi
    
    # Очистка логов старше 7 дней
    if [ -d "$AUDIT_DIR" ]; then
        find "$AUDIT_DIR" -name "*.log" -mtime +7 -delete 2>/dev/null || true
        log_info "Старые логи удалены"
    fi
    
    log_success "Очистка завершена"
}

# Функция проверки зависимостей
check_dependencies() {
    # Проверка ps
    if ! command -v ps &> /dev/null; then
        log_error "Команда ps не доступна"
        exit 1
    fi
    
    # Проверка kill
    if ! command -v kill &> /dev/null; then
        log_error "Команда kill не доступна"
        exit 1
    fi
    
    # Проверка jq (для JSON обработки)
    if ! command -v jq &> /dev/null; then
        log_warning "jq не установлен. JSON обработка будет упрощена."
    fi
}

# Функция проверки, запущено ли приложение
is_running() {
    local pids=($(find_processes))
    [ ${#pids[@]} -gt 0 ]
}

# Основная функция
main() {
    echo "=========================================="
    echo "Остановка приложения $APP_NAME"
    echo "=========================================="
    
    # Проверка, что мы в корневой директории проекта
    if [ ! -f "go.mod" ]; then
        log_error "Скрипт должен быть запущен из корневой директории проекта"
        exit 1
    fi
    
    # Проверка зависимостей
    check_dependencies
    
    # Проверка, запущено ли приложение
    if ! is_running; then
        log_info "Приложение $APP_NAME не запущено"
        
        # Проверка наличия PID файла
        if [ -f "$PID_FILE" ]; then
            log_warning "Найден PID файл, но процесс не активен. Удаление PID файла..."
            rm -f "$PID_FILE"
        fi
        
        log_success "Приложение уже остановлено"
        exit 0
    fi
    
    # Получение списка процессов
    local pids=($(find_processes))
    log_info "Найдено процессов: ${#pids[@]}"
    
    for pid in "${pids[@]}"; do
        log_info "Найден процесс: $pid"
    done
    
    # Подтверждение остановки
    echo ""
    read -p "Остановить все найденные процессы? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Отмена остановки"
        exit 0
    fi
    
    # Остановка всех процессов
    local all_stopped=true
    for pid in "${pids[@]}"; do
        if ! stop_process "$pid" false; then
            all_stopped=false
        fi
    done
    
    if [ "$all_stopped" = false ]; then
        log_warning "Некоторые процессы не удалось остановить мягко"
        echo ""
        read -p "Принудительно остановить оставшиеся процессы? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            for pid in "${pids[@]}"; do
                if ps -p "$pid" > /dev/null 2>&1; then
                    stop_process "$pid" true
                fi
            done
        fi
    fi
    
    # Сохранение состояния
    save_state
    
    # Очистка
    cleanup
    
    log_success "Приложение $APP_NAME полностью остановлено!"
    
    # Вывод информации о состоянии
    echo ""
    echo "=========================================="
    echo "Информация о остановке:"
    echo "Приложение: $APP_NAME"
    echo "Конфигурация: $CONFIG_DIR"
    echo "Лог файл: $LOG_FILE"
    echo "Аудит логи: $AUDIT_DIR"
    echo "Состояние: $CONFIG_DIR/state/"
    echo "=========================================="
}

# Запуск основной функции
main "$@"