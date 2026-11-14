# stop-app.ps1 - Скрипт полной остановки приложения CodeAuthorityFirst
# Поддерживает Windows системы

param(
    [switch]$Force = $false,
    [switch]$Clean = $false,
    [switch]$Quiet = $false,
    [switch]$Help = $false,
    [switch]$SaveState = $true
)

# Настройка вывода
$ErrorActionPreference = "Stop"

# Константы
$APP_NAME = "CodeAuthorityFirst"
$APP_BINARY = "main_security_demo.exe"
$CONFIG_DIR = ".caf"
$AUDIT_DIR = "audit_logs"
$PID_FILE = "$CONFIG_DIR\app.pid"
$LOG_FILE = "$CONFIG_DIR\app.log"
$STATE_DIR = "$CONFIG_DIR\state"

# Цвета для вывода
$Colors = @{
    Red = [ConsoleColor]::Red
    Green = [ConsoleColor]::Green
    Yellow = [ConsoleColor]::Yellow
    Blue = [ConsoleColor]::Blue
    White = [ConsoleColor]::White
}

# Функции вывода
function Write-Info($message) {
    if (-not $Quiet) {
        Write-Host "[$(Get-Date -Format 'HH:mm:ss')] [INFO] $message" -ForegroundColor $Colors.Blue
    }
}

function Write-Success($message) {
    if (-not $Quiet) {
        Write-Host "[$(Get-Date -Format 'HH:mm:ss')] [SUCCESS] $message" -ForegroundColor $Colors.Green
    }
}

function Write-Warning($message) {
    if (-not $Quiet) {
        Write-Host "[$(Get-Date -Format 'HH:mm:ss')] [WARNING] $message" -ForegroundColor $Colors.Yellow
    }
}

function Write-Error($message) {
    Write-Host "[$(Get-Date -Format 'HH:mm:ss')] [ERROR] $message" -ForegroundColor $Colors.Red
}

function Write-Help() {
    Write-Host "Использование: stop-app.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:"
    Write-Host "  -Force         Принудительная остановка"
    Write-Host "  -Clean         Очистка временных файлов"
    Write-Host "  -Quiet         Тихий режим"
    Write-Host "  -SaveState    Сохранять состояние (по умолчанию: true)"
    Write-Host "  -Help          Показать эту справку"
    Write-Host ""
    Write-Host "Примеры:"
    Write-Host "  stop-app.ps1                # Остановка приложения"
    Write-Host "  stop-app.ps1 -Force         # Принудительная остановка"
    Write-Host "  stop-app.ps1 -Clean         # Остановка с очисткой"
    Write-Host "  stop-app.ps1 -Force -Clean  # Принудительная остановка с очисткой"
}

# Функция поиска процессов приложения
function Find-Processes {
    $processes = @()
    
    # Поиск по имени бинарного файла
    if (Test-Path $APP_BINARY) {
        $binaryProcesses = Get-Process -Name $APP_BINARY -ErrorAction SilentlyContinue
        foreach ($proc in $binaryProcesses) {
            $processes += $proc
        }
    }
    
    # Поиск по PID файлу
    if (Test-Path $PID_FILE) {
        $pidFromFile = Get-Content $PID_FILE -ErrorAction SilentlyContinue
        if ($pidFromFile -and (Get-Process -Id $pidFromFile -ErrorAction SilentlyContinue)) {
            $found = $false
            foreach ($proc in $processes) {
                if ($proc.Id -eq $pidFromFile) {
                    $found = $true
                    break
                }
            }
            if (-not $found) {
                $processes += (Get-Process -Id $pidFromFile -ErrorAction SilentlyContinue)
            }
        } else {
            Write-Warning "PID файл существует, но процесс не найден. Удаление PID файла..."
            Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
        }
    }
    
    # Поиск по ключевым словам в командной строке
    $keywordProcesses = Get-Process | Where-Object {
        $_.ProcessName -match "CodeAuthorityFirst" -or 
        $_.CommandLine -match "CodeAuthorityFirst" -or
        $_.CommandLine -match "security.*demo"
    }
    
    foreach ($proc in $keywordProcesses) {
        $found = $false
        foreach ($existingProc in $processes) {
            if ($existingProc.Id -eq $proc.Id) {
                $found = $true
                break
            }
        }
        if (-not $found) {
            $processes += $proc
        }
    }
    
    return $processes
}

# Функция остановки процесса
function Stop-ProcessWrapper {
    param($Process, $Force)
    
    $pid = $Process.Id
    $processName = $Process.ProcessName
    
    Write-Info "Остановка процесса $pid ($processName)..."
    
    if ($Force) {
        # Принудительное завершение
        Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
        Write-Info "Отправлен сигнал KILL процессу $pid"
    } else {
        # Мягкое завершение
        Stop-Process -Id $pid -ErrorAction SilentlyContinue
        Write-Info "Отправлен сигнал STOP процессу $pid"
        
        # Ожидание завершения
        for ($i = 1; $i -le 10; $i++) {
            if (-not (Get-Process -Id $pid -ErrorAction SilentlyContinue)) {
                break
            }
            Start-Sleep -Seconds 1
        }
        
        # Принудительное завершение, если необходимо
        if (Get-Process -Id $pid -ErrorAction SilentlyContinue) {
            Write-Warning "Принудительное завершение процесса $pid"
            Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
        }
    }
    
    # Проверка результата
    if (Get-Process -Id $pid -ErrorAction SilentlyContinue) {
        Write-Error "Не удалось остановить процесс $pid"
        return $false
    } else {
        Write-Success "Процесс $pid остановлен"
        return $true
    }
}

# Функция сохранения состояния
function Save-State {
    param($Processes)
    
    if (-not $SaveState) {
        Write-Info "Сохранение состояния отключено"
        return
    }
    
    Write-Info "Сохранение состояния приложения..."
    
    # Создание директории состояния
    if (-not (Test-Path $STATE_DIR)) {
        New-Item -ItemType Directory -Path $STATE_DIR -Force | Out-Null
    }
    
    $timestamp = Get-Date -UFormat "%Y-%m-%dT%H:%M:%SZ"
    $stateFile = "$STATE_DIR\stop_$timestamp.json"
    
    $state = @{
        stopped_at = $timestamp
        app_name = $APP_NAME
        audit_dir = $AUDIT_DIR
        config_dir = $CONFIG_DIR
        processes = @()
    }
    
    # Добавление информации о процессах
    foreach ($proc in $Processes) {
        $processInfo = @{
            pid = $proc.Id
            ppid = $proc.ProcessName
            name = $proc.ProcessName
            cpu = $proc.CPU
            memory = $proc.WorkingSet / 1MB
            cmd = $proc.CommandLine -join " "
        }
        $state.processes += $processInfo
    }
    
    $state | ConvertTo-Json -Depth 3 | Out-File -FilePath $stateFile -Encoding UTF8
    
    Write-Success "Состояние сохранено: $stateFile"
}

# Функция очистки
function Cleanup {
    Write-Info "Очистка временных файлов..."
    
    # Удаление PID файла
    if (Test-Path $PID_FILE) {
        Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
        Write-Info "PID файл удален"
    }
    
    # Очистка временных файлов приложения
    $tempFiles = Get-ChildItem -Path . -Name "*.tmp", "*.temp", "*~" -ErrorAction SilentlyContinue | Select-Object -First 20
    if ($tempFiles) {
        $tempFiles | Remove-Item -Force -ErrorAction SilentlyContinue
        Write-Info "Временные файлы очищены"
    }
    
    # Очистка кэша (если есть)
    $cacheDir = "$CONFIG_DIR\cache"
    if (Test-Path $cacheDir) {
        Get-ChildItem -Path $cacheDir -File | Remove-Item -Force -ErrorAction SilentlyContinue
        Write-Info "Кэш очищен"
    }
    
    # Очистка логов старше 7 дней
    if (Test-Path $AUDIT_DIR) {
        $oldLogs = Get-ChildItem -Path $AUDIT_DIR -Name "*.log" -ErrorAction SilentlyContinue | 
                   Where-Object { (Get-Item $_).LastWriteTime -lt (Get-Date).AddDays(-7) }
        if ($oldLogs) {
            $oldLogs | Remove-Item -Force -ErrorAction SilentlyContinue
            Write-Info "Старые логи удалены"
        }
    }
    
    Write-Success "Очистка завершена"
}

# Функция проверки зависимостей
function Test-Dependencies {
    # Проверка Get-Process
    if (-not (Get-Command "Get-Process" -ErrorAction SilentlyContinue)) {
        Write-Error "Команда Get-Process не доступна"
        exit 1
    }
    
    # Проверка Stop-Process
    if (-not (Get-Command "Stop-Process" -ErrorAction SilentlyContinue)) {
        Write-Error "Команда Stop-Process не доступна"
        exit 1
    }
}

# Функция проверки, запущено ли приложение
function Is-Running {
    $processes = Find-Processes
    return ($processes.Count -gt 0)
}

# Основная функция
function Main {
    Write-Host "==========================================" -ForegroundColor $Colors.White
    Write-Host "Остановка приложения $APP_NAME" -ForegroundColor $Colors.White
    Write-Host "==========================================" -ForegroundColor $Colors.White
    
    # Обработка справки
    if ($Help) {
        Write-Help
        exit 0
    }
    
    # Проверка, что мы в корневой директории проекта
    if (-not (Test-Path "go.mod")) {
        Write-Error "Скрипт должен быть запущен из корневой директории проекта"
        exit 1
    }
    
    # Проверка зависимостей
    Test-Dependencies
    
    # Проверка, запущено ли приложение
    if (-not (Is-Running)) {
        Write-Info "Приложение $APP_NAME не запущено"
        
        # Проверка наличия PID файла
        if (Test-Path $PID_FILE) {
            Write-Warning "Найден PID файл, но процесс не активен. Удаление PID файла..."
            Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
        }
        
        Write-Success "Приложение уже остановлено"
        exit 0
    }
    
    # Получение списка процессов
    $processes = Find-Processes
    Write-Info "Найдено процессов: $($processes.Count)"
    
    foreach ($proc in $processes) {
        Write-Info "Найден процесс: $($proc.Id) - $($proc.ProcessName)"
    }
    
    # Подтверждение остановки
    Write-Host ""
    $response = Read-Host "Остановить все найденные процессы? (y/N)"
    if ($response -ne "y" -and $response -ne "Y") {
        Write-Info "Отмена остановки"
        exit 0
    }
    
    # Остановка всех процессов
    $allStopped = $true
    foreach ($proc in $processes) {
        if (-not (Stop-ProcessWrapper -Process $proc -Force:$Force)) {
            $allStopped = $false
        }
    }
    
    if (-not $allStopped) {
        Write-Warning "Некоторые процессы не удалось остановить мягко"
        Write-Host ""
        $response = Read-Host "Принудительно остановить оставшиеся процессы? (y/N)"
        if ($response -eq "y" -or $response -eq "Y") {
            foreach ($proc in $processes) {
                if (Get-Process -Id $proc.Id -ErrorAction SilentlyContinue) {
                    Stop-ProcessWrapper -Process $proc -Force:$true
                }
            }
        }
    }
    
    # Сохранение состояния
    Save-State -Processes $processes
    
    # Очистка
    if ($Clean) {
        Cleanup
    }
    
    Write-Success "Приложение $APP_NAME полностью остановлено!"
    
    # Вывод информации о состоянии
    Write-Host ""
    Write-Host "==========================================" -ForegroundColor $Colors.White
    Write-Host "Информация об остановке:" -ForegroundColor $Colors.White
    Write-Host "Приложение: $APP_NAME" -ForegroundColor $Colors.White
    Write-Host "Конфигурация: $CONFIG_DIR" -ForegroundColor $Colors.White
    Write-Host "Лог файл: $LOG_FILE" -ForegroundColor $Colors.White
    Write-Host "Аудит логи: $AUDIT_DIR" -ForegroundColor $Colors.White
    Write-Host "Состояние: $STATE_DIR\" -ForegroundColor $Colors.White
    Write-Host "==========================================" -ForegroundColor $Colors.White
}

# Обработка аргументов командной строки
if ($Help) {
    Write-Help
    exit 0
}

# Запуск основной функции
Main