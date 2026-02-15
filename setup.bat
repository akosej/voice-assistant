@echo off
REM 🚀 Setup Script para Agente de Voz Inteligente (Windows)

setlocal enabledelayedexpansion

echo.
echo ╔════════════════════════════════════════════════════════════╗
echo ║  🎙️  SETUP - AGENTE DE VOZ INTELIGENTE                    ║
echo ╚════════════════════════════════════════════════════════════╝
echo.

REM 1. Verificar que Go está instalado
echo ✔️ Verificando que Go está instalado...
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go no está instalado. Por favor instálalo desde https://golang.org/
    pause
    exit /b 1
)
for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo %GO_VERSION%
echo.

REM 2. Descargar dependencias
echo ✔️ Descargando dependencias...
go mod download
if errorlevel 1 (
    echo ❌ Error al descargar dependencias
    pause
    exit /b 1
)
go mod tidy
echo ✅ Dependencias descargadas
echo.

REM 3. Compilar el proyecto
echo ✔️ Compilando el proyecto...
go build -o agent.exe
if errorlevel 1 (
    echo ❌ Error al compilar
    pause
    exit /b 1
)
echo ✅ Compilación exitosa
echo.

REM 4. Mostrar instrucciones
echo ╔════════════════════════════════════════════════════════════╗
echo ║  ✅ ESTABLECIMIENTO COMPLETADO                            ║
echo ╚════════════════════════════════════════════════════════════╝
echo.
echo Para iniciar el servidor, ejecuta:
echo.
echo   agent.exe
echo.
echo O con Go directamente:
echo.
echo   go run main.go
echo.
echo Luego llama al API:
echo.
echo   curl -X POST http://localhost:3000/voice ^
echo     -H "Content-Type: application/json" ^
echo     -d "{\"session_id\":\"test\",\"text\":\"Hola\"}"
echo.
pause
