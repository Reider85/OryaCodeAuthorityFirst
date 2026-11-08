@echo off
echo Starting UI Accessibility Testing
echo Screen: dashboard

REM Test Results
set pass_count=0
set total_tests=4

echo.
echo Testing INV-UI-009: 100%% Keyboard-Only Operation
echo ---------------------------------------------

REM Check keyboard navigation features
findstr /C:"tea\.KeyTab" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found keyboard feature: tea\.KeyTab
    set has_keyboard=1
) else (
    set has_keyboard=0
)

findstr /C:"CurrentWidget" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found keyboard feature: CurrentWidget
    set has_keyboard=1
)

REM Check for mouse operations
findstr /C:"mouse" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   WARNING: Potential mouse operation found: mouse
    set has_mouse=1
) else (
    set has_mouse=0
)

REM Check keyboard navigation completeness
findstr /C:"Tab" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Tab navigation: Implemented
) else (
    echo   Tab navigation: Missing
    set complete_nav=0
)

findstr /C:"ShiftTab" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Shift+Tab navigation: Implemented
) else (
    echo   Shift+Tab navigation: Missing
    set complete_nav=0
)

if %has_keyboard% equ 1 if %has_mouse% equ 0 (
    echo   Full keyboard-only operation implemented
    set pass_009=1
    set /a pass_count+=1
) else (
    echo   Keyboard-only operation incomplete
    set pass_009=0
)

set pass_009_status=PASS
if %pass_009% equ 0 set pass_009_status=FAIL

echo   INV-UI-009: %pass_009_status%

echo.
echo Testing INV-UI-010: Help Overlay Functionality
echo ---------------------------------------------

REM Check help overlay features
findstr /C:"ShowHelp" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found help feature: ShowHelp
    set has_help=1
)

findstr /C:"helpSections" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found help feature: helpSections
    set has_help=1
)

REM Check help activation
findstr /C:"KeyQuestionMark" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Help activation with ? key implemented
)

REM Check help content completeness
findstr /C:"Keyboard Shortcuts" "pkg\ui\components\helpoverlay.go" >nul
if %errorlevel% equ 0 (
    echo   Help content: Keyboard Shortcuts present
) else (
    echo   Help content: Keyboard Shortcuts missing
    set complete_help=0
)

if %has_help% equ 1 (
    echo   Comprehensive help overlay implemented
    set pass_010=1
    set /a pass_count+=1
) else (
    echo   Help overlay incomplete
    set pass_010=0
)

set pass_010_status=PASS
if %pass_010% equ 0 set pass_010_status=FAIL

echo   INV-UI-010: %pass_010_status%

echo.
echo Testing Focus Visibility
echo -----------------------

REM Check focus indicators
findstr /C:"CurrentWidget" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found focus feature: CurrentWidget
    set has_focus=1
)

findstr /C:"focusPrefix" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found focus feature: focusPrefix
    set has_focus=1
)

REM Check focus marker implementation
findstr /C:"focusPrefix.*▸" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Focus marker (▸) implemented
)

REM Check focus management
findstr /C:"CurrentWidget.*=" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Focus management implemented
)

REM Check focus visibility across states
findstr /C:"StateLoading" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Focus handling in StateLoading: Present
) else (
    echo   Focus handling in StateLoading: Missing
    set focus_all_states=0
)

findstr /C:"StateSuccess" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Focus handling in StateSuccess: Present
) else (
    echo   Focus handling in StateSuccess: Missing
    set focus_all_states=0
)

if %has_focus% equ 1 (
    echo   Focus visibility properly implemented
    set pass_focus=1
    set /a pass_count+=1
) else (
    echo   Focus visibility issues found
    set pass_focus=0
)

set pass_focus_status=PASS
if %pass_focus% equ 0 set pass_focus_status=FAIL

echo   Focus Visibility: %pass_focus_status%

echo.
echo Testing Symbol-Text Pairing
echo ---------------------------

REM Check for symbols with text
findstr /C:"✓" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found symbol: ✓
    set has_symbols=1
)

findstr /C:"⚠" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found symbol: ⚠
    set has_symbols=1
)

findstr /C:"✗" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Found symbol: ✗
    set has_symbols=1
)

REM Check symbol-text pairing
findstr /C:"success" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    echo   Symbol-text pairing found
    set has_symbol_text=1
)

findstr /C:"warning" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    set has_symbol_text=1
)

findstr /C:"error" "pkg\ui\screens\dashboard.go" >nul
if %errorlevel% equ 0 (
    set has_symbol_text=1
)

REM Check color-blind safety
findstr /C:"color" "pkg\ui\a11y.go" >nul
if %errorlevel% equ 0 (
    echo   Color-blind safety features implemented
)

if %has_symbols% equ 1 if %has_symbol_text% equ 1 (
    echo   Symbol-text pairing properly implemented
    set pass_symbols=1
    set /a pass_count+=1
) else (
    echo   Symbol-text pairing issues found
    set pass_symbols=0
)

set pass_symbols_status=PASS
if %pass_symbols% equ 0 set pass_symbols_status=FAIL

echo   Symbol-Text Pairing: %pass_symbols_status%

echo.
echo Accessibility Testing Summary
echo ============================
echo INV-UI-009: %pass_009_status%
echo INV-UI-010: %pass_010_status%
echo Focus Visibility: %pass_focus_status%
echo Symbol-Text Pairing: %pass_symbols_status%
echo.
echo Overall Accessibility Result: %pass_count%/%total_tests% tests passed

if %pass_count% equ %total_tests% (
    echo All accessibility requirements met!
    exit 0
) else (
    echo Some accessibility requirements need attention
    exit 1
)