# start-app.ps1 - Полный скрипт запуска приложения CodeAuthorityFirst
# Поддерживает Windows системы

param(
    [string]$Mode = "normal",
    [switch]$Force = $false,
    [switch]$Quiet = $false,
    [switch]$Help = $false
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
    Write-Host "Использование: start-app.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "OPTIONS:"
    Write-Host "  -Mode normal|minimal|debug   Режим запуска (по умолчанию: normal)"
    Write-Host "  -Force                       Принудительная пересборка"
    Write-Host "  -Quiet                      Тихий режим"
    Write-Host "  -Help                       Показать эту справку"
    Write-Host ""
    Write-Host "Примеры:"
    Write-Host "  start-app.ps1                # Запуск в обычном режиме"
    Write-Host "  start-app.ps1 -Mode debug   # Запуск в режиме отладки"
    Write-Host "  start-app.ps1 -Force        # Принудительная пересборка"
}

# Функция проверки зависимостей
function Test-Dependencies {
    Write-Info "Проверка зависимостей..."
    
    # Проверка Go
    if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Error "Go не установлен. Пожалуйста, установите Go 1.27 или выше."
        exit 1
    }
    
    $GoVersion = (go version).Split(' ')[2].TrimStart('go:')
    $RequiredVersion = "1.27"
    
    if ([version]$GoVersion -lt [version]$RequiredVersion) {
        Write-Error "Требуется Go $RequiredVersion или выше. Установлена версия: $GoVersion"
        exit 1
    }
    
    Write-Success "Go версии $GoVersion найден"
    
    # Проверка git
    if (Get-Command "git" -ErrorAction SilentlyContinue) {
        Write-Success "Git найден"
    } else {
        Write-Warning "Git не найден. Сборка будет без информации о версии."
    }
    
    # Проверка Python
    if (Get-Command "python" -ErrorAction SilentlyContinue -or Get-Command "python3" -ErrorAction SilentlyContinue) {
        Write-Success "Python найден"
    } else {
        Write-Warning "Python не найден. Некоторые компоненты могут не работать."
    }
}

# Функция сборки приложения
function Build-Application {
    Write-Info "Сборка приложения $APP_NAME..."
    
    # Проверка существующего бинарного файла
    if (Test-Path $APP_BINARY) {
        Write-Info "Бинарный файл существует. Проверка версии..."
        
        # Проверка версии через бинарный файл
        if (& $APP_BINARY --version 2>$null) {
            $CurrentVersion = & $APP_BINARY --version 2>$null
            Write-Info "Текущая версия: $CurrentVersion"
            
            # Проверка необходимости пересборки
            if (Test-Path ".git") {
                $GitChanges = (git status --porcelain 2>$null | Measure-Object).Count
                if ($GitChanges -gt 0) {
                    Write-Warning "Найдены изменения в Git. Рекомендуется пересобрать приложение."
                    if (-not $Force) {
                        $response = Read-Host "Пересобрать приложение? (y/N)"
                        if ($response -eq "y" -or $response -eq "Y") {
                            Remove-Item $APP_BINARY -Force
                        } else {
                            Write-Info "Используется существующий бинарный файл"
                            return $true
                        }
                    } else {
                        Remove-Item $APP_BINARY -Force
                        Write-Info "Принудительное удаление существующего бинарного файла"
                    }
                }
            }
        } else {
            Write-Warning "Существующий бинарный файл не работает. Удаление и пересборка..."
            Remove-Item $APP_BINARY -Force
        }
    }
    
    # Сборка приложения
    if (Test-Path "scripts\build.ps1") {
        Write-Info "Использование скрипта сборки..."
        & ".\scripts\build.ps1"
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Ошибка сборки через скрипт!"
            exit 1
        }
    } else {
        Write-Info "Ручная сборка..."
        
        # Получение версии из git или использование fallback
        if (Get-Command "git" -ErrorAction SilentlyContinue -and (Test-Path ".git")) {
            $GitVersion = (git describe --tags --dirty --always 2>$null) -replace "`n", "" -replace "`r", ""
            if (-not $GitVersion) { $GitVersion = "unknown" }
            $VersionFlag = "-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GitVersion'"
            Write-Info "Версия из git: $GitVersion"
        } else {
            $VersionFlag = "-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
            Write-Warning "Git не доступен, используется fallback версия"
        }
        
        # Флаги сборки
        $CommitHash = (git rev-parse HEAD 2>$null) -replace "`n", "" -replace "`r", "" -replace "`n", "" -replace "`r", ""
        if (-not $CommitHash) { $CommitHash = "unknown" }
        
        $DateTime = Get-Date -UFormat "%Y-%m-%dT%H:%M:%SZ"
        $Modified = (git diff-index --quiet HEAD -- 2>$null) -replace "`n", "" -replace "`r", ""
        if ($Modified -eq "") { $Modified = "false" }
        
        $LDFLAGS = "$VersionFlag -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit=$CommitHash'"
        $LDFLAGS = "$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.datetime=$DateTime'"
        $LDFLAGS = "$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.modified=$Modified'"
        
        Write-Info "Сборка с флагами: $LDFLAGS"
        
        # Компиляция
        go build -ldflags="$LDFLAGS" -o "$APP_BINARY" ".\main_security_demo"
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Приложение успешно собрано!"
        } else {
            Write-Error "Ошибка сборки приложения!"
            exit 1
        }
    }
    
    return $true
}

# Функция создания директорий
function New-Directories {
    Write-Info "Создание необходимых директорий..."
    
    $dirs = @($CONFIG_DIR, $AUDIT_DIR, "$CONFIG_DIR\security")
    foreach ($dir in $dirs) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
    }
    
    Write-Success "Директории созданы"
}

# Функция генерации конфигурации
function New-Config {
    Write-Info "Генерация конфигурации..."
    
    # Основная конфигурация CAF
    $ConfigFile = "$CONFIG_DIR\config.yml"
    if (-not (Test-Path $ConfigFile)) {
        @"
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
"@ | Out-File -FilePath $ConfigFile -Encoding UTF8
        Write-Success "Основная конфигурация CAF создана: $ConfigFile"
    } else {
        Write-Info "Основная конфигурация CAF уже существует"
    }
    
    # Конфигурация безопасности
    $SecurityConfigFile = "$CONFIG_DIR\security.yml"
    if (-not (Test-Path $SecurityConfigFile)) {
        @"
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
"@ | Out-File -FilePath $SecurityConfigFile -Encoding UTF8
        Write-Success "Конфигурация безопасности создана: $SecurityConfigFile"
    } else {
        Write-Info "Конфигурация безопасности уже существует"
    }
}

# Функция проверки конфигурации
function Test-Config {
    Write-Info "Валидация конфигурации..."
    
    # Проверка YAML синтаксиса через PowerShell
    try {
        $ConfigContent = Get-Content "$CONFIG_DIR\config.yml" -Raw -ErrorAction Stop
        $SecurityContent = Get-Content "$CONFIG_DIR\security.yml" -Raw -ErrorAction Stop
        
        # Простая проверка YAML структуры
        if ($ConfigContent -match "^version:" -and $SecurityContent -match "^yaml:") {
            Write-Success "Конфигурация валидна"
        } else {
            Write-Warning "Конфигурация имеет некорректную структуру"
        }
    } catch {
        Write-Error "Ошибка чтения конфигурации: $_"
        exit 1
    }
}

# Функция запуска приложения
function Start-Application {
    Write-Info "Запуск приложения $APP_NAME..."
    
    # Проверка, запущено ли уже приложение
    if (Test-Path $PID_FILE) {
        $PID = Get-Content $PID_FILE -ErrorAction SilentlyContinue
        if ($PID -and (Get-Process -Id $PID -ErrorAction SilentlyContinue)) {
            Write-Warning "Приложение уже запущено с PID: $PID"
            Write-Info "Используйте 'stop-app.ps1' для остановки"
            return $true
        } else {
            Write-Warning "PID файл существует, но процесс не найден. Удаление PID файла..."
            Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
        }
    }
    
    # Определение аргументов запуска
    $StartArgs = @()
    switch ($Mode) {
        "debug" { $StartArgs += "--debug" }
        "minimal" { $StartArgs += "--minimal" }
        default {}
    }
    
    # Запуск приложения в фоновом режиме
    $Process = Start-Process -FilePath ".\$APP_BINARY" -ArgumentList $StartArgs -RedirectStandardOutput "$LOG_FILE" -RedirectStandardError "$LOG_FILE" -NoNewWindow -PassThru
    
    if ($Process) {
        # Сохранение PID
        $Process.Id | Out-File -FilePath $PID_FILE -Encoding UTF8
        
        # Проверка запуска
        Start-Sleep -Seconds 2
        if (Get-Process -Id $Process.Id -ErrorAction SilentlyContinue) {
            Write-Success "Приложение запущено успешно!"
            Write-Info "PID: $($Process.Id)"
            Write-Info "Лог файл: $LOG_FILE"
            Write-Info "Директория аудита: $AUDIT_DIR"
            
            # Вывод информации о процессе
            "=========================================="
            "Приложение: $APP_NAME"
            "PID: $($Process.Id)"
            "Лог файл: $LOG_FILE"
            "Конфигурация: $CONFIG_DIR"
            "Аудит логи: $AUDIT_DIR"
            "=========================================="
            
            return $true
        } else {
            Write-Error "Ошибка запуска приложения!"
            if (Test-Path $LOG_FILE) {
                Write-Info "Проверьте лог файл: $LOG_FILE"
                Write-Info "Последние 10 строк лога:"
                Get-Content $LOGFile -Tail 10
            }
            Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
            exit 1
        }
    } else {
        Write-Error "Ошибка запуска приложения!"
        exit 1
    }
}

# Функция обработки завершения
function Register-Cleanup {
    $CleanupScript = {
        param($PID_FILE, $LOG_FILE, $APP_NAME)
        
        Write-Info "Остановка приложения..."
        if (Test-Path $PID_FILE) {
            $PID = Get-Content $PID_FILE -ErrorAction SilentlyContinue
            if ($PID -and (Get-Process -Id $PID -ErrorAction SilentlyContinue)) {
                Stop-Process -Id $PID -Force -ErrorAction SilentlyContinue
                Write-Info "Отправлен сигнал STOP процессу $PID"
                
                # Ожидание завершения
                for ($i = 1; $i -le 10; $i++) {
                    if (-not (Get-Process -Id $PID -ErrorAction SilentlyContinue)) {
                        break
                    }
                    Start-Sleep -Seconds 1
                }
                
                # Принудительное завершение, если необходимо
                if (Get-Process -Id $PID -ErrorAction SilentlyContinue) {
                    Write-Warning "Принудительное завершение процесса $PID"
                    Stop-Process -Id $PID -Force -ErrorAction SilentlyContinue
                }
            }
            Remove-Item $PID_FILE -Force -ErrorAction SilentlyContinue
        }
        Write-Success "Приложение остановлено"
    }
    
    # Регистрация обработчика завершения
    $job = Start-Job -ScriptBlock $CleanupScript -ArgumentList $PID_FILE, $LOG_FILE, $APP_NAME
    Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action {
        Receive-Job $job -Wait
        Remove-Job $job -Force
    } | Out-Null
}

# Основная функция
function Main {
    Write-Host "==========================================" -ForegroundColor $Colors.White
    Write-Host "Запуск приложения $APP_NAME" -ForegroundColor $Colors.White
    Write-Host "==========================================" -ForegroundColor $Colors.White
    
    # Проверка, что мы в корневой директории проекта
    if (-not (Test-Path "go.mod")) {
        Write-Error "Скрипт должен быть запущен из корневой директории проекта"
        exit 1
    }
    
    # Обработка справки
    if ($Help) {
        Write-Help
        exit 0
    }
    
    # Выполнение всех шагов
    Test-Dependencies
    $BuildResult = Build-Application
    if (-not $BuildResult) { exit 1 }
    
    New-Directories
    New-Config
    Test-Config
    Register-Cleanup
    Start-Application
    
    Write-Success "Приложение $APP_NAME полностью запущено!"
    
    # Ожидание завершения
    try {
        while ($true) { Start-Sleep -Seconds 1 }
    } catch {
        # Обработка прерывания
        Write-Info "Скрипт завершается..."
    }
}

# Обработка аргументов командной строки
if ($Help) {
    Write-Help
    exit 0
}

# Запуск основной функции
Main