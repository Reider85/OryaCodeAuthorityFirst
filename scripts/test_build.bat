@echo off
REM test_build.bat - Test build script for CodeAuthorityFirst
REM This script tests the build process with version embedding

cd /d "%~dp0\.."

echo 🔨 Testing build with version embedding...

REM Test build with hardcoded version
go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=test-version'" -o caf_test.exe ./cmd/caf

if %errorlevel% equ 0 (
    echo ✅ Build successful!
    echo 📋 Binary created: caf_test.exe
    
    REM Test the binary
    echo 🧪 Testing binary...
    caf_test.exe --version
    
    echo 🎉 Test complete!
) else (
    echo ❌ Build failed!
    exit /b 1
)

REM Clean up
del caf_test.exe 2>nul