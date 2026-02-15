# 🏗️ ARQUITECTURA DETALLADA

## Patrón: Strategy + Dependency Injection

Este proyecto implementa un middleware desacoplado utilizando dos patrones de diseño clave:

### 1. **Strategy Pattern**
Permite cambiar el comportamiento en tiempo de ejecución inyectando diferentes implementaciones.

```go
type IntentDetector interface {
    Detect(text string) (*Intent, error)
}

// Puedes cambiar entre:
detector1 := &RuleBasedDetector{}          // Basado en reglas
detector2 := &OpenAIDetector{}             // OpenAI
detector3 := &HuggingFaceDetector{}        // HuggingFace
```

### 2. **Dependency Injection**
Las dependencias se inyectan en el constructor, no se crean internamente.

```go
type Agent struct {
    store    SessionStore
    detector IntentDetector
    executor Executor
}

// Inyección en constructor
agent := NewAgent(store, detector, executor)
```

---

## Flujo de Procesamiento

```
┌─────────────────────────────────────────────────────────────────────┐
│                      SOLICITUD HTTP (Fiber)                         │
│                                                                      │
│  POST /voice {session_id: "user_123", text: "Hola"}                │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Agent.Handle()                               │
├─────────────────────────────────────────────────────────────────────┤
│ 1️⃣ Obtener o crear sesión                                           │
│    SessionStore.Get(sessionID)                                      │
│    Resultado: Session{ID, State, Context}                          │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 2️⃣ Detectar intención                                               │
│    IntentDetector.Detect(text)                                      │
│    Resultado: Intent{Name, Entities}                               │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 3️⃣ Ejecutar acción                                                  │
│    Executor.Execute(intent, session)                                │
│    - Lee el estado actual de la sesión                             │
│    - Cambia el estado si es necesario                              │
│    - Calcula respuesta basada en intención y estado                │
│    Resultado: Response{Reply, Action, Data}                        │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 4️⃣ Guardar sesión actualizada                                       │
│    SessionStore.Save(session)                                       │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      RESPUESTA HTTP (JSON)                          │
│                                                                      │
│  {                                                                  │
│    "reply": "...",                                                 │
│    "action": "greeting_received",                                  │
│    "session_id": "user_123",                                       │
│    "state": "idle"                                                 │
│  }                                                                  │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Máquina de Estados

El agente utiliza **máquinas de estados** para mantener contexto en conversaciones multi-turno.

```
                    ┌─────────────┐
                    │   IDLE      │
                    │ (Inicial)   │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
   ┌─────────────┐  ┌──────────────┐  ┌───────────────┐
   │ CREATING    │  │ ANALYZING    │  │ WAITING FOR   │
   │ TASK        │  │ TEXT         │  │ CONFIRM       │
   │             │  │              │  │               │
   └─────────────┘  └──────────────┘  └───────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   IDLE      │
                    │ (Respuesta) │
                    └─────────────┘
```

**Ejemplo: Crear Tarea**

```
Usuario: "crear tarea"
←──────────────────────────────────┐
State: IDLE                         │
Executor: Cambia a CREATING_TASK    │
Response: "¿Cómo se llama?"        │
                                    │
Usuario: "Mi tarea importante"
←──────────────────────────────────┤
State: CREATING_TASK                │
Executor: Guarda tarea              │
Executor: Cambia a IDLE            │
Response: "✅ Tarea creada"        │
```

---

## Componentes Principales

### 1. **Agent (agent.go)**

```go
type Agent struct {
    store    SessionStore      // Almacenar sesiones
    detector IntentDetector    // Detectar intenciones
    executor Executor          // Ejecutar acciones
}

// El método principal
func (a *Agent) Handle(sessionID, text string) (*Response, error)
```

**Responsabilidades:**
- Orquestar el flujo
- Guardar sesiones
- Propagar errores

### 2. **SessionStore (session_store.go)**

```go
type SessionStore interface {
    Get(id string) (*Session, error)
    Save(session *Session) error
    Delete(id string) error
}
```

**Implementaciones disponibles:**
- `MemoryStore` (en memoria, para desarrollo)
- Puedes crear: `RedisStore`, `PostgresStore`, `MongoStore`

### 3. **IntentDetector (intent_detector.go)**

```go
type IntentDetector interface {
    Detect(text string) (*Intent, error)
}
```

**Implementaciones disponibles:**
- `RuleBasedDetector` (basado en palabras clave)
- Puedes crear: `OpenAIDetector`, `HuggingFaceDetector`

### 4. **Executor (executor.go)**

```go
type Executor interface {
    Execute(intent *Intent, session *Session) (*Response, error)
}
```

**Implementaciones disponibles:**
- `DefaultExecutor` (maneja: saludo, ayuda, análisis, tareas)
- Puedes crear: `CustomExecutor` con tu lógica

### 5. **Types (types.go)**

Define estructuras base:
- `State` - Estados de la sesión
- `Session` - Datos de sesión
- `Intent` - Información detectada
- `Response` - Respuesta del agente
- `ReadabilityMetrics` - Métricas de análisis

---

## Cómo Extender

### Crear un nuevo Detector

```go
// Paso 1: Crear estructura
type LLMDetector struct {
    apiKey string
    client *openai.Client
}

// Paso 2: Implementar interface
func (l *LLMDetector) Detect(text string) (*Intent, error) {
    // Usar LLM para detectar
    prompt := fmt.Sprintf("Detecta la intención: %s", text)
    result, _ := l.client.CreateCompletion(prompt)
    
    return &Intent{
        Name: parseName(result),
        Entities: parseEntities(result),
    }, nil
}

// Paso 3: Usar en main.go
detector := &LLMDetector{apiKey: "sk-..."}
agent := voice.NewAgent(store, detector, executor)
```

### Crear un nuevo Executor

```go
// Paso 1: Crear estructura
type CustomExecutor struct {
    db database.DB
}

// Paso 2: Implementar interface
func (c *CustomExecutor) Execute(intent *Intent, session *Session) (*Response, error) {
    switch intent.Name {
    case "custom_intent":
        // Tu lógica personalizada
        return &Response{Reply: "Custom!"}, nil
    }
    return nil, nil
}

// Paso 3: Usar en main.go
executor := &CustomExecutor{db: myDB}
agent := voice.NewAgent(store, detector, executor)
```

### Cambiar a Redis

```go
// Paso 1: Crear RedisStore
type RedisStore struct {
    client *redis.Client
}

func (r *RedisStore) Get(id string) (*Session, error) {
    data, _ := r.client.Get(id)
    session := &Session{}
    json.Unmarshal(data, session)
    return session, nil
}

func (r *RedisStore) Save(session *Session) error {
    data, _ := json.Marshal(session)
    return r.client.Set(session.ID, data)
}

// Paso 2: Usar en main.go
redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
store := &RedisStore{client: redisClient}
agent := voice.NewAgent(store, detector, executor)
```

---

## Por Qué Esta Arquitectura

### ✅ **Mantenible**
- Código limpio y organizado
- Fácil de entender
- Separación de responsabilidades

### ✅ **Testeable**
- Las interfaces facilitan los mocks
- Puedes testear sin HTTP
- Tests unitarios simples

```go
func TestAgent(t *testing.T) {
    // Mock del SessionStore
    store := &MockStore{}
    
    // Mock del IntentDetector
    detector := &MockDetector{}
    
    // Agent con mocks
    agent := voice.NewAgent(store, detector, executor)
    
    response, _ := agent.Handle("user", "hola")
    assert.Equal(t, "¡Hola!...", response.Reply)
}
```

### ✅ **Extensible**
- Agregar nuevos detectores
- Agregar nuevos ejecutores
- Agregar nuevos stores
- Sin tocar el código existente

### ✅ **Desacoplado**
- El agente no conoce de HTTP/Fiber
- Puede usarse en CLI, gRPC, WebSocket, etc.
- El framework es solo un "envoltorio"

### ✅ **Escalable**
- Fácil pasar de MemoryStore a Redis
- Fácil cambiar LLM
- Fácil multiplicar instancias

---

## Relaciones entre Componentes

```
┌─────────────────────────────────────────────────────────┐
│                                                          │
│              ┌────────────────────┐                      │
│              │    Agent           │                      │
│              │  (Orquestrador)    │                      │
│              └────────────────────┘                      │
│                     │                                    │
│      ┌──────────────┼──────────────┐                     │
│      │              │              │                     │
│      ▼              ▼              ▼                     │
│  ┌───────────┐ ┌──────────┐ ┌──────────┐               │
│  │ Store     │ │Detector  │ │ Executor │               │
│  │(Interface)│ │(Interface) │(Interface)│               │
│  └───────────┘ └──────────┘ └──────────┘               │
│      │              │              │                     │
│      ▼              ▼              ▼                     │
│  ┌───────────┐ ┌──────────┐ ┌──────────┐               │
│  │Memory     │ │RuleBase  │ │ Default  │               │
│  │Store      │ │Detector  │ │ Executor │               │
│  └───────────┘ └──────────┘ └──────────┘               │
│                                                          │
│  Inyección en constructor → Fácil cambiar               │
│  Interfaces → Fácil extender                            │
│  Sin lock-in → Puedes usar en cualquier framework       │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Flujo Completo: Ejemplo

**Usuario envía:** `"crear tarea"`

```
1. HTTP Request (Fiber)
   POST /voice {session_id: "user1", text: "crear tarea"}
   
2. Agent.Handle("user1", "crear tarea")

3. SessionStore.Get("user1")
   Returns: Session{ID: "user1", State: StateIdle, Context: {}}

4. IntentDetector.Detect("crear tarea")
   Returns: Intent{Name: "create_task", Entities: {type: "task"}}

5. Executor.Execute(intent, session)
   - Lee: session.State == StateIdle
   - Cambia: session.State = StateCreatingTask
   - Retorna: Response{Reply: "¿Cómo se llama?", ...}

6. SessionStore.Save(session)
   Guarda: Session{State: StateCreatingTask}

7. HTTP Response (Fiber)
   200 {reply: "¿Cómo se llama?", state: "creating_task"}

---

**Usuario envía:** `"Mi primer proyecto"`

1. HTTP Request (Fiber)
   POST /voice {session_id: "user1", text: "Mi primer proyecto"}

2. Agent.Handle("user1", "Mi primer proyecto")

3. SessionStore.Get("user1")
   Returns: Session{ID: "user1", State: StateCreatingTask} <- ¡Cambió!

4. IntentDetector.Detect("Mi primer proyecto")
   Returns: Intent{Name: "unknown", Entities: {}}

5. Executor.Execute(intent, session)
   - Lee: session.State == StateCreatingTask <- ¡Toma acción diferente!
   - Guarda tarea: "Mi primer proyecto"
   - Cambia: session.State = StateIdle
   - Retorna: Response{Reply: "✅ Tarea creada", Action: "task_created"}

6. SessionStore.Save(session)
   Guarda: Session{State: StateIdle}

7. HTTP Response (Fiber)
   200 {reply: "✅ Tarea creada", state: "idle"}
```

---

## Ventajas sobre Arquitectura Monolítica

### ❌ **Sin estos patrones:**
```go
func HandleVoice(c *gin.Context) {
    if contains(text, "hola") {
        // Saludo
    } else if contains(text, "crear tarea") {
        // Crear tarea
        if session["state"] == "idle" {
            // Primera vez
        } else if session["state"] == "creating_task" {
            // Segunda vez
        }
        // ... 50 líneas más de ifs anidados
    }
    // Imposible testear sin HTTP
    // Imposible reutilizar en CLI
    // Imposible cambiar de framework
}
```

### ✅ **Con estos patrones:**
```go
response := agent.Handle(sessionID, text)  // ¡Limpio!

// Testeable
response, _ := agent.Handle("user", "hola")
assert.Equal(t, "¡Hola!...", response.Reply)

// Reutilizable
// CLI: response := agent.Handle(sessionID, userInput)
// gRPC: response := agent.Handle(sessionID, protoMsg)
// WebSocket: response := agent.Handle(sessionID, wsMsg)
```

---

**Conclusión:** Esta arquitectura es enterprise-grade, mantenible, testeable y escalable. 🚀
