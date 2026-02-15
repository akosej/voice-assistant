# 🎙️ Agente de Voz Inteligente

Un middleware sensible, desacoplado y reutilizable para procesar comandos de voz en múltiples frameworks (Fiber, Gin, Echo, gRPC, WebSockets, CLI, Tests).

## 🏗 Arquitectura

```
/internal/voice
    ├── types.go              # Tipos base (State, Session, Response, Intent)
    ├── agent.go              # Agent middleware (la pieza central)
    ├── session_store.go      # SessionStore interface + MemoryStore
    ├── intent_detector.go    # IntentDetector interface + RuleBasedDetector
    └── executor.go           # Executor interface + DefaultExecutor
```

### 🧠 Patrón de Diseño

**Strategy + Dependency Injection**

```
┌─────────────────────────────────────────────────────────────┐
│                      AGENT                                   │
│                                                               │
│  Handle(sessionID, text) -> Response                         │
│                                                               │
└────────────┬──────────────┬──────────────┬──────────────────┘
             │              │              │
             ▼              ▼              ▼
      SessionStore  IntentDetector    Executor
      (Interface)   (Interface)       (Interface)
             │              │              │
             ▼              ▼              ▼
      MemoryStore   RuleBasedDetector DefaultExecutor
      (InMemory)    (Reglas)           (Acciones)
```

## 🚀 Inicio Rápido

### 1. Clonar o descargar el proyecto

```bash
cd a:\DevOps\Testing\agent
```

### 2. Descargar dependencias

```bash
go mod download
go mod tidy
```

### 3. Ejecutar el servidor

```bash
go run main.go
```

**Salida esperada:**
```
🚀 Servidor iniciado en http://localhost:3000
📝 Documenta tu API en http://localhost:3000/info
🧪 Prueba con: curl -X POST http://localhost:3000/voice -H 'Content-Type: application/json' -d '{"session_id":"test","text":"Hola"}'
```

## 🌐 Endpoints

### 1. **POST /voice** - Procesar comando de voz

Procesa un comando de voz/texto y retorna una respuesta.

**Request:**
```json
{
  "session_id": "usuario_123",
  "text": "analiza este texto: La inteligencia artificial es fascinante"
}
```

**Response:**
```json
{
  "reply": "📊 Análisis de Legibilidad:\n\n✅ Facilidad de Lectura (Flesch-Kincaid): 65.32\n...",
  "action": "readability_analyzed",
  "data": {
    "flesch_kincaid_ease": 65.32,
    "flesch_kincaid_grade": 7.5,
    "word_count": 8,
    "sentence_count": 1,
    "syllable_count": 20,
    "words_per_sentence": 8.0,
    "syllables_per_word": 2.5
  },
  "session_id": "usuario_123",
  "state": "idle"
}
```

### 2. **GET /info** - Información del agente

Devuelve información sobre las capacidades del agente.

**Response:**
```json
{
  "name": "Agente de Voz Inteligente",
  "version": "1.0.0",
  "capabilities": [
    "Analizar legibilidad de textos",
    "Gestionar tareas",
    "Detección de lenguaje natural",
    "Almacenamiento de sesiones"
  ],
  "endpoints": {
    "POST /voice": "Procesar comando de voz",
    "GET /health": "Verificar estado del servidor",
    "GET /info": "Información del agente",
    "GET /": "Información principal"
  }
}
```

### 3. **GET /health** - Verificar estado

Verifica si el servidor está activo.

**Response:**
```json
{
  "status": "ok",
  "message": "Agente de Voz en línea"
}
```

### 4. **GET /** - Información principal

**Response:**
```json
{
  "mensaje": "Bienvenido al Agente de Voz Inteligente",
  "info_url": "/info",
  "voice_api": "POST /voice",
  "health_check": "/health",
  "ejemplo": {
    "session_id": "usuario_123",
    "text": "Hola, analiza este texto: La inteligencia artificial es fascinante"
  }
}
```

## 🎯 Comandos Soportados

### 1. **Analizar Legibilidad**

```
Usuario: "analiza este texto: [TU TEXTO]"
Agente: Retorna métricas Flesch-Kincaid y más
```

Palabras clave detectadas:
- `analizar`, `legibilidad`, `leer`, `texto`

### 2. **Crear Tarea**

```
Usuario: "crear tarea"
Agente: "¿Cómo se llama la tarea?"
Usuario: "Mi primera tarea"
Agente: "✅ Tarea 'Mi primera tarea' creada correctamente."
```

Palabras clave detectadas:
- `crear`, `tarea`

### 3. **Listar Tareas**

```
Usuario: "mostrar mis tareas"
Agente: "📋 Tus Tareas:\n\n1. Mi primera tarea"
```

Palabras clave detectadas:
- `listar`, `ver`, `mostrar`, `tareas`

### 4. **Saludos**

```
Usuario: "Hola"
Agente: "¡Hola! Bienvenido al Agente de Voz Inteligente. ¿Cómo puedo ayudarte?"
```

Palabras clave detectadas:
- `hola`, `buenos`, `saludos`, `ey`, `hey`

### 5. **Ayuda**

```
Usuario: "ayuda"
Agente: [Muestra el menú de comandos disponibles]
```

Palabras clave detectadas:
- `ayuda`, `help`, `qué puedo`, `qué puedes`

## 📝 Ejemplos de Uso

### Con cURL

```bash
# Saludar
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_1",
    "text": "Hola"
  }'

# Analizar legibilidad
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_1",
    "text": "analiza este texto: La programación es el arte de decirle a otra persona qué quieres que el ordenador haga."
  }'

# Crear tarea
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_1",
    "text": "crear tarea"
  }'

curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_1",
    "text": "Implementar nueva funcionalidad de reportes"
  }'

# Listar tareas
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_1",
    "text": "mostrar tareas"
  }'
```

### Con Python

```python
import requests

BASE_URL = "http://localhost:3000"

def call_agent(session_id, text):
    response = requests.post(
        f"{BASE_URL}/voice",
        json={
            "session_id": session_id,
            "text": text
        }
    )
    return response.json()

# Ejemplo
resultado = call_agent("usuario_1", "Hola, analiza este texto: Python es genial")
print(resultado)
```

### Con JavaScript/Node.js

```javascript
const SESSION_ID = 'usuario_123';

async function callAgent(text) {
  const response = await fetch('http://localhost:3000/voice', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      session_id: SESSION_ID,
      text: text
    })
  });
  return response.json();
}

// Ejemplo
const result = await callAgent('Hola');
console.log(result);
```

## 🔧 Extensibilidad

### 🎯 Crear un detector de intenciones personalizado

```go
type CustomDetector struct{}

func (c *CustomDetector) Detect(text string) (*voice.Intent, error) {
    // Tu lógica de detección personalizada
    return &voice.Intent{
        Name: "custom_intent",
        Entities: map[string]string{
            "key": "value",
        },
    }, nil
}

// En main.go
// detector := &CustomDetector{}
// voiceAgent = voice.NewAgent(store, detector, executor)
```

### 🎯 Crear un executor personalizado

```go
type CustomExecutor struct{}

func (c *CustomExecutor) Execute(intent *voice.Intent, session *voice.Session) (*voice.Response, error) {
    // Tu lógica de ejecución personalizada
    return &voice.Response{
        Reply: "Respuesta personalizada",
        Action: "custom_action",
    }, nil
}

// En main.go
// executor := &CustomExecutor{}
// voiceAgent = voice.NewAgent(store, detector, executor)
```

### 🎯 Cambiar a Redis en 5 minutos

```go
type RedisStore struct {
    client *redis.Client
}

func (r *RedisStore) Get(id string) (*voice.Session, error) {
    // Implementar con Redis
}

func (r *RedisStore) Save(session *voice.Session) error {
    // Implementar con Redis
}

// En main.go
// store := &RedisStore{client: redisClient}
// voiceAgent = voice.NewAgent(store, detector, executor)
```

### 🎯 Integrar con LLM (OpenAI, Anthropic, etc.)

```go
type LLMDetector struct {
    client *openai.Client
}

func (l *LLMDetector) Detect(text string) (*voice.Intent, error) {
    // Usar OpenAI para detectar intenciones
}

// En main.go
// detector := &LLMDetector{client: openaiClient}
// voiceAgent = voice.NewAgent(store, detector, executor)
```

## 📊 Métricas de Legibilidad

El agente calcula automáticamente:

- **Flesch-Kincaid Reading Ease**: Cuán fácil es leer el texto (0-100)
- **Flesch-Kincaid Grade Level**: Nivel de grado requerido para entender
- **Word Count**: Total de palabras
- **Sentence Count**: Total de oraciones
- **Syllable Count**: Total de sílabas
- **Words/Sentence**: Promedio de palabras por oración
- **Syllables/Word**: Promedio de sílabas por palabra

## 🧪 Testing

Usa el agent directamente sin necesidad de HTTP:

```go
package main

import (
    "agent/internal/voice"
    "testing"
)

func TestAgent(t *testing.T) {
    store := voice.NewMemoryStore()
    detector := voice.NewRuleBasedDetector()
    executor := voice.NewDefaultExecutor()
    
    agent := voice.NewAgent(store, detector, executor)
    
    response, err := agent.Handle("test_user", "Hola")
    
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    
    if response.Reply != "¡Hola! Bienvenido..." {
        t.Fatalf("Respuesta inesperada: %s", response.Reply)
    }
}
```

## 🎯 Casos de Uso

✅ **Chatbots desacoplados del framework**
✅ **APIs de voz multi-framework**
✅ **CLI con inteligencia conversacional**
✅ **WebSocket handlers para chat en tiempo real**
✅ **Microservicios de procesamiento de voz**
✅ **Testing sin necesidad de HTTP**

## 📚 Estructura de Archivos Completa

```
agent/
├── main.go                          # Servidor Fiber + handlers HTTP
├── go.mod                           # Dependencias
├── go.sum                           # Checksums de dependencias
├── data/
│   └── movie_reviews.txt            # Datos de prueba
├── internal/
│   └── voice/
│       ├── types.go                 # Tipos base
│       ├── agent.go                 # Agent middleware
│       ├── session_store.go         # SessionStore y MemoryStore
│       ├── intent_detector.go       # IntentDetector y RuleBasedDetector
│       └── executor.go              # Executor y DefaultExecutor
├── cmd/
│   └── main.go                      # Alternativa no usada (archivo de referencia)
└── README.md                        # Este archivo
```

## 🚀 Próximos Pasos

1. **Integrar detección de lenguaje real** (OpenAI, Anthropic)
2. **Agregar persistencia** (PostgreSQL, MongoDB)
3. **Implementar WebSockets** para chat en tiempo real
4. **Crear UI** (Web, Mobile)
5. **Dockerizar** la aplicación
6. **Agregar autenticación** y autorización
7. **Implementar analytics** de usuario
8. **Crear system prompts** para diferentes dominios

## 📄 Licencia

MIT

---

**Hecho con ❤️ usando Go, Fiber y arquitectura limpia**
