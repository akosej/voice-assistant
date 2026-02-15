# 🎯 Flujo de Creación de Tareas (Mejorado)

## Ahora es más intuitivo ✨

### **Antes (No intuitivo):**
```
Usuario: "crear tarea"
Agente:  "¿Cómo se llama la tarea?"
Usuario: "crear tarea: Mi primera tarea"  ❌ Tenia que repetir el comando
Agente:  "✅ Tarea creada"
```

### **Ahora (Intuitivo):**
```
Usuario: "crear tarea"
Agente:  "📝 Perfecto. ¿Qué título quieres ponerle a la tarea?"

Usuario: "Implementar autenticación"
Agente:  "✅ Título guardado. Ahora, ¿qué descripción quieres agregar? (o escribe 'saltar' si no quieres descripción)"

Usuario: "Implementar OAuth2 en la aplicación"
Agente:  "✅ Tarea 'Implementar autenticación' creada correctamente."
```

---

## 🔄 Máquina de Estados de Tareas

```
                    ┌──────────────────┐
                    │   IDLE           │
                    │ (Estado inicial) │
                    └────────┬─────────┘
                             │
              Usuario dice:   │  "crear tarea"
                             ▼
                    ┌──────────────────────────────────┐
                    │ CREATING_TASK_TITLE              │
                    │ Agente pide: "¿Qué título?"     │
                    └────────┬─────────────────────────┘
                             │
              Usuario envía:  │  Título (ej: "Mi tarea")
                             ▼
                    ┌──────────────────────────────────┐
                    │ CREATING_TASK_DESCRIPTION        │
                    │ Agente pide: "¿Qué descripción?" │
                    └────────┬─────────────────────────┘
                             │
                    ┌────────┴─────────────────────────┐
                    │                                  │
        Usuario envía:                        Usuario escribe:
        "Descripción"                          "saltar"
                    │                                  │
                    ▼                                  ▼
          Descripción guardada              Sin descripción
                    │                                  │
                    └────────┬─────────────────────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │   IDLE           │
                    │ Tarea creada    │
                    └──────────────────┘
```

---

## 📝 Ejemplos de Uso

### Ejemplo 1: Crear tarea con descripción

```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_1","text":"crear tarea"}'

# Respuesta:
# {
#   "reply": "📝 Perfecto. ¿Qué título quieres ponerle a la tarea?",
#   "action": "waiting_for_task_title",
#   "state": "creating_task_title"
# }
```

```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_1","text":"Aprender Go"}'

# Respuesta:
# {
#   "reply": "✅ Título guardado. Ahora, ¿qué descripción quieres agregar?",
#   "action": "waiting_for_task_description",
#   "state": "creating_task_description"
# }
```

```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_1","text":"Aprender las bases de Go y crear mi primer proyecto"}'

# Respuesta:
# {
#   "reply": "✅ Tarea 'Aprender Go' creada correctamente.",
#   "action": "task_created",
#   "data": {
#     "title": "Aprender Go",
#     "description": "Aprender las bases de Go y crear mi primer proyecto"
#   },
#   "state": "idle"
# }
```

### Ejemplo 2: Crear tarea sin descripción (usar "saltar")

```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_2","text":"crear tarea"}'

# Agente pide título...

curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_2","text":"Hacer ejercicio"}'

# Agente pide descripción...

curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_2","text":"saltar"}'

# Respuesta:
# {
#   "reply": "✅ Tarea 'Hacer ejercicio' creada correctamente.",
#   "action": "task_created",
#   "data": {
#     "title": "Hacer ejercicio",
#     "description": ""
#   },
#   "state": "idle"
# }
```

### Ejemplo 3: Listar tareas después de crear

```bash
curl -X POST http://localhost:3000/voice \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"usuario_1","text":"mostrar tareas"}'

# Respuesta:
# {
#   "reply": "📋 Tus Tareas:\n\n1. Aprender Go - Aprender las bases de Go y crear mi primer proyecto",
#   "action": "tasks_listed",
#   "data": [
#     "Aprender Go - Aprender las bases de Go y crear mi primer proyecto"
#   ]
# }
```

---

## 🧪 Test con Python

```python
import requests

BASE_URL = "http://localhost:3000"
SESSION_ID = "usuario_test"

def call_agent(text):
    response = requests.post(
        f"{BASE_URL}/voice",
        json={"session_id": SESSION_ID, "text": text}
    )
    return response.json()

# Paso 1: Iniciar creación de tarea
response1 = call_agent("crear tarea")
print(f"Paso 1: {response1['reply']}")

# Paso 2: Enviar título
response2 = call_agent("Estudiar Python")
print(f"Paso 2: {response2['reply']}")

# Paso 3: Enviar descripción
response3 = call_agent("Completar ejercicios de DataFrames en pandas")
print(f"Paso 3: {response3['reply']}")

# Paso 4: Listar tareas
response4 = call_agent("mostrar tareas")
print(f"Paso 4: {response4['reply']}")
```

---

## 💾 Estructura de Datos Permanencia

Cuando creas una tarea con título y descripción:

**Almacenamiento en MemoryStore:**
```
Usuario: "usuario_1"
  Sesión:
    - ID: "usuario_1"
    - State: "idle" (después de crear)
    - TemporaryData: {} (se limpian después de crear)
    - Tasks:
      [
        "Aprender Go - Aprender las bases de Go y crear mi primer proyecto",
        "Hacer ejercicio"
      ]
```

---

## 🔍 Detalles Técnicos

### Estados Disponibles (en `types.go`):
- `StateIdle` - Esperando comando
- `StateCreatingTaskTitle` - Esperando título
- `StateCreatingTaskDescription` - Esperando descripción
- `StateAnalyzingText` - Analizando legibilidad
- `StateWaitingForConfirm` - Esperando confirmación

### Flujo en `executor.go`:
1. **StateIdle + Intent "create_task"**
   - Cambia a `StateCreatingTaskTitle`
   - Pide el título

2. **StateCreatingTaskTitle + Intent.Raw (texto)**
   - Guarda título en `session.TemporaryData["task_title"]`
   - Cambia a `StateCreatingTaskDescription`
   - Pide descripción

3. **StateCreatingTaskDescription + Intent.Raw (texto)**
   - Si es "saltar" → sin descripción
   - Si es texto regular → guarda descripción
   - Crea la tarea en `e.tasks[sessionID]`
   - Limpia datos temporales
   - Cambia a `StateIdle`
   - Devuelve confirmación

---

## 🚀 Próximas Mejoras Sugeridas

- [ ] Editar tareas existentes
- [ ] Marcar tareas como completadas
- [ ] Eliminar tareas
- [ ] Prioridad en tareas
- [ ] Fecha límite para tareas
- [ ] Categorías de tareas
