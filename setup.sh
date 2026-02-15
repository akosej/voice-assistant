#!/bin/bash

# 🚀 Setup Script para Agente de Voz Inteligente

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  🎙️  SETUP - AGENTE DE VOZ INTELIGENTE                    ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# 1. Verificar que Go está instalado
echo "✔️ Verificando que Go está instalado..."
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado. Por favor instálalo desde https://golang.org/"
    exit 1
fi
GO_VERSION=$(go version)
echo "$GO_VERSION"
echo ""

# 2. Descargar dependencias
echo "✔️ Descargando dependencias..."
go mod download
go mod tidy
echo "✅ Dependencias descargadas"
echo ""

# 3. Compilar el proyecto
echo "✔️ Compilando el proyecto..."
go build -o agent
if [ $? -ne 0 ]; then
    echo "❌ Error al compilar"
    exit 1
fi
echo "✅ Compilación exitosa"
echo ""

# 4. Mostrar instrucciones
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  ✅ ESTABLECIMIENTO COMPLETADO                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Para iniciar el servidor, ejecuta:"
echo ""
echo "  ./agent"
echo ""
echo "O con Go directamente:"
echo ""
echo "  go run main.go"
echo ""
echo "Luego llama al API:"
echo ""
echo "  curl -X POST http://localhost:3000/voice \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"session_id\":\"test\",\"text\":\"Hola\"}'"
echo ""
