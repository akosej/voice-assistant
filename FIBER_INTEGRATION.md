# 🌐 Integración con Fiber

## ¿Por qué Fiber?

Fiber es más rápido y consum menos memoria que Gin para proyectos pequeños a medianos:

| Métrica | Gin | Fiber | Ventaja |
|---------|-----|-------|---------|
| Throughput | ~32k req/s | ~48k req/s | ⚡ Fiber 50% más rápido |
| Memoria | 9MB | 2MB | 💾 Fiber 4x menos |
| Tiempo de startup | 150ms | 50ms | 🚀 Fiber 3x más rápido |
| Tamaño binario | 12MB | 6MB | 📦 Fiber 2x más pequeño |

---

## Estructura actual (main.go)

```go
// 1. Inicializar el agente (una sola vez)
var voiceAgent *voice.Agent

func initVoiceAgent() {
    store := voice.NewMemoryStore()
    detector := voice.NewRuleBasedDetector()
    executor := voice.NewDefaultExecutor()
    voiceAgent = voice.NewAgent(store, detector, executor)
}

// 2. Handler HTTP que usa el agente
func HandleVoice(c *fiber.Ctx) error {
    var req VoiceRequest
    c.BodyParser(&req)
    
    response, err := voiceAgent.Handle(req.SessionID, req.Text)
    
    return c.JSON(response)
}

// 3. Crear app Fiber y registrar handlers
func main() {
    initVoiceAgent()
    
    app := fiber.New()
    app.Post("/voice", HandleVoice)
    app.Get("/health", HandleHealth)
    
    app.Listen(":3000")
}
```

---

## API Endpoints

### 1. POST /voice - Procesar comando

**Request:**
```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "usuario_123",
    "text": "Hola, analiza este texto"
  }'
```

**Response:**
```json
{
  "reply": "¡Hola! Bienvenido al Agente de Voz Inteligente...",
  "action": "greeting_received",
  "session_id": "usuario_123",
  "state": "idle"
}
```

---

### 2. GET /health - Verificar salud

**Request:**
```bash
curl http://localhost:3000/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Agente de Voz en línea"
}
```

---

### 3. GET /info - Información del agente

**Request:**
```bash
curl http://localhost:3000/info
```

**Response:**
```json
{
  "name": "Agente de Voz Inteligente",
  "version": "1.0.0",
  "capabilities": [
    "Analizar legibilidad de textos",
    "Gestionar tareas",
    "Detección de lenguaje natural"
  ]
}
```

---

## Middlewares Configurados

### 1. Logger Middleware
Registra todas las solicitudes HTTP:
```
GET   /health   200   1.23 ms
POST  /voice    200   45.2 ms
GET   /info     200   0.98 ms
```

### 2. CORS Middleware
Permitir solicitudes desde cualquier origen:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
```

---

## Ejemplos de Uso

### Con cURL

```bash
# 1. Saludo
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"Hola"}'

# 2. Solicitar ayuda
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"ayuda"}'

# 3. Analizar texto
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"analizar: Python es genial"}'

# 4. Crear tarea
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"crear tarea"}'

curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"Implementar nueva feature"}'

# 5. Listar tareas
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"test","text":"mostrar tareas"}'
```

### Con JavaScript/Fetch

```javascript
const API_URL = 'http://localhost:3000';

async function callAgent(sessionId, text) {
  const response = await fetch(`${API_URL}/voice`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      session_id: sessionId,
      text: text
    })
  });
  return response.json();
}

// Uso
const result = await callAgent('user123', 'Hola');
console.log(result.reply);
```

### Con Axios (Node.js/Vue/React)

```javascript
import axios from 'axios';

const agent = axios.create({
  baseURL: 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json'
  }
});

async function sendMessage(sessionId, text) {
  try {
    const { data } = await agent.post('/voice', {
      session_id: sessionId,
      text: text
    });
    return data;
  } catch (error) {
    console.error('Error:', error.message);
  }
}

// Uso
const response = await sendMessage('user123', 'analizar este texto');
console.log(response.reply);
```

### Con Python/Requests

```python
import requests

BASE_URL = 'http://localhost:3000'
SESSION_ID = 'user123'

def call_agent(text):
    response = requests.post(
        f'{BASE_URL}/voice',
        json={
            'session_id': SESSION_ID,
            'text': text
        }
    )
    return response.json()

# Uso
result = call_agent('Hola')
print(result['reply'])
```

---

## Manejo de Errores

### Error 400 - Bad Request
```json
{
  "error": "session_id y text son requeridos"
}
```

Causas comunes:
- Falta `session_id`
- Falta `text`
- JSON inválido

### Error 500 - Internal Server Error
```json
{
  "error": "error procesando solicitud: ..."
}
```

Causas comunes:
- El agente retornó un error
- Problema al guardar sesión
- Problema al detectar intención

---

## Optimizaciones para Producción

### 1. Connection Pool
```go
app := fiber.New(fiber.Config{
    Concurrency: 256 * 1024,  // Max conexiones simultáneas
    ReadBufferSize: 4096,      // Buffer para leer
})
```

### 2. Compresión
```go
app.Use(compress.New(compress.Config{
    Level: compress.LevelBestBalance,
}))
```

### 3. Rate Limiting
```go
app.Use(limiter.New(limiter.Config{
    Max:        100,           // Max 100 solicitudes
    Expiration: 1 * time.Minute,
}))
```

### 4. Request ID (Tracing)
```go
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Request-ID", uuid.New().String())
    return c.Next()
})
```

---

## Deployment

### Docker

```dockerfile
FROM golang:1.25 as builder
WORKDIR /build
COPY . .
RUN go build -o agent

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/agent .
COPY data/ ./data/
EXPOSE 3000
CMD ["./agent"]
```

```bash
docker build -t voice-agent:1.0 .
docker run -p 3000:3000 voice-agent:1.0
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: voice-agent
spec:
  replicas: 3
  selector:
    matchLabels:
      app: voice-agent
  template:
    metadata:
      labels:
        app: voice-agent
    spec:
      containers:
      - name: voice-agent
        image: voice-agent:1.0
        ports:
        - containerPort: 3000
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5
```

---

## Performance Tips

1. **Reutilizar Agent**: Se inicializa una sola vez
2. **Connection pooling**: Fiber lo maneja automáticamente
3. **Async processing**: Puedes procesar en background
4. **Caching**: Guardar resultados de análisis frecuentes
5. **Redis backend**: Para escalar a múltiples instancias

---

**Fiber es rápido ⚡, mantenible 🔧, y fácil de escalar 📈**
