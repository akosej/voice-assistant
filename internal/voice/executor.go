// Package voice proporciona el middleware reutilizable para el agente de voz inteligente.
package voice

import (
	"fmt"
	"strings"

	"go.rtnl.ai/nlp/text"
)

// Executor interfaz para ejecutar acciones
type Executor interface {
	Execute(intent *Intent, session *Session) (*Response, error)
}

// DefaultExecutor implementación por defecto del executor
type DefaultExecutor struct {
	tasks map[string][]string // sessionID -> lista de tareas
}

// NewDefaultExecutor crea un nuevo DefaultExecutor
func NewDefaultExecutor() *DefaultExecutor {
	return &DefaultExecutor{
		tasks: make(map[string][]string),
	}
}

// Execute ejecuta la acción según la intención y estado de la sesión
func (e *DefaultExecutor) Execute(intent *Intent, session *Session) (*Response, error) {
	// Primero verificar si estamos en estados que requieren input específico
	// En estos estados, ignoramos la detección de intención y simplemente procesamos el input
	switch session.State {
	case StateCreatingTaskTitle:
		// Usuario está en proceso de crear tarea, esperando el título
		return e.handleCreateTask(intent, session)

	case StateCreatingTaskDescription:
		// Usuario está en proceso de crear tarea, esperando la descripción
		return e.handleCreateTask(intent, session)

	case StateAnalyzingText:
		// Usuario está en proceso de análisis de legibilidad, esperando el texto
		return e.handleAnalyzeReadability(intent, session)
	}

	// Si no estamos en un estado especial, detectar la intención del usuario
	switch intent.Name {

	case "greeting":
		return &Response{
			Reply:     "¡Hola! Bienvenido al Agente de Voz Inteligente. ¿Cómo puedo ayudarte?",
			Action:    "greeting_received",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	case "help":
		help := `Puedo ayudarte con:
1. **Analizar Legibilidad**: Envía un texto y te daré métricas de legibilidad (Flesch-Kincaid)
2. **Crear Tareas**: Puedo crear tareas para ti
3. **Listar Tareas**: Te muestro todas tus tareas
4. **Ayuda**: Este menú

¿Qué deseas hacer?`
		return &Response{
			Reply:     help,
			Action:    "help_provided",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	case "analyze_readability":
		return e.handleAnalyzeReadability(intent, session)

	case "create_task":
		return e.handleCreateTask(intent, session)

	case "list_tasks":
		return e.handleListTasks(intent, session)

	default:
		return &Response{
			Reply:     "No entendí el comando. Escribe 'ayuda' para ver las opciones disponibles.",
			Action:    "unknown_command",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil
	}
}

// handleAnalyzeReadability analiza la legibilidad del texto
func (e *DefaultExecutor) handleAnalyzeReadability(intent *Intent, session *Session) (*Response, error) {
	switch session.State {

	case StateIdle:
		// Pedir texto para analizar
		session.State = StateAnalyzingText
		return &Response{
			Reply:     "¿Qué texto deseas analizar? (Envía el texto completo)",
			Action:    "waiting_for_text",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	case StateAnalyzingText:
		// Analizar el texto recibido
		textToAnalyze := intent.Raw
		if textToAnalyze == "" {
			return &Response{
				Reply:     "No recibí texto. Intenta de nuevo.",
				SessionID: session.ID,
				State:     string(session.State),
			}, nil
		}

		metrics, err := e.calculateReadabilityMetrics(textToAnalyze)
		if err != nil {
			session.State = StateIdle
			return &Response{
				Reply:     fmt.Sprintf("Error al analizar: %v", err),
				SessionID: session.ID,
				State:     string(session.State),
			}, nil
		}

		session.State = StateIdle

		reply := fmt.Sprintf(`📊 Análisis de Legibilidad:

✅ Facilidad de Lectura (Flesch-Kincaid): %.2f
✅ Nivel de Grado: %.2f
✅ Total de Palabras: %d
✅ Total de Oraciones: %d
✅ Total de Sílabas: %d
✅ Promedio Palabras/Oración: %.2f
✅ Promedio Sílabas/Palabra: %.2f

Interpretación:
- Score 90-100: Muy fácil (5to grado)
- Score 60-70: Estándar (8vo grado)
- Score 0-30: Difícil (Universitario)`,
			metrics.FleschKincaidEase,
			metrics.FleschKincaidGrade,
			metrics.WordCount,
			metrics.SentenceCount,
			metrics.SyllableCount,
			metrics.WordsPerSentence,
			metrics.SyllablesPerWord,
		)

		return &Response{
			Reply:     reply,
			Action:    "readability_analyzed",
			Data:      metrics,
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	default:
		session.State = StateIdle
		return &Response{
			Reply:     "Estado desconocido. Reiniciando.",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil
	}
}

// handleCreateTask crea una nueva tarea con título y descripción
func (e *DefaultExecutor) handleCreateTask(intent *Intent, session *Session) (*Response, error) {
	switch session.State {

	case StateIdle:
		// Paso 1: Iniciar creación de tarea
		session.State = StateCreatingTaskTitle
		return &Response{
			Reply:     "Perfecto. ¿Qué título quieres ponerle a la tarea?",
			Action:    "waiting_for_task_title",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	case StateCreatingTaskTitle:
		// Paso 2: Recibir título y pedir descripción
		taskTitle := intent.Raw
		if taskTitle == "" {
			return &Response{
				Reply:     "No recibí el título. Intenta de nuevo.",
				SessionID: session.ID,
				State:     string(session.State),
			}, nil
		}

		// Guardar título temporalmente
		if session.TemporaryData == nil {
			session.TemporaryData = make(map[string]interface{})
		}
		session.TemporaryData["task_title"] = taskTitle

		// Cambiar a estado de descripción
		session.State = StateCreatingTaskDescription

		return &Response{
			Reply:     "Título guardado. Ahora, ¿qué descripción quieres agregar? (o escribe 'saltar' si no quieres descripción)",
			Action:    "waiting_for_task_description",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	case StateCreatingTaskDescription:
		// Paso 3: Recibir descripción y crear tarea
		taskDescription := intent.Raw

		// Obtener el título guardado
		taskTitle, ok := session.TemporaryData["task_title"].(string)
		if !ok {
			session.State = StateIdle
			return &Response{
				Reply:     "Error: no se encontró el título. Intenta de nuevo.",
				SessionID: session.ID,
				State:     string(session.State),
			}, nil
		}

		// Ignorar descripción si dice "saltar"
		if strings.ToLower(strings.TrimSpace(taskDescription)) == "saltar" {
			taskDescription = ""
		}

		// Guardar tarea completa
		if _, ok := e.tasks[session.ID]; !ok {
			e.tasks[session.ID] = []string{}
		}

		// Guardar con descripción si existe
		taskEntry := taskTitle
		if taskDescription != "" && strings.ToLower(strings.TrimSpace(taskDescription)) != "saltar" {
			taskEntry = fmt.Sprintf("%s - %s", taskTitle, taskDescription)
		}

		e.tasks[session.ID] = append(e.tasks[session.ID], taskEntry)

		// Limpiar datos temporales
		delete(session.TemporaryData, "task_title")

		// Volver al estado idle
		session.State = StateIdle

		return &Response{
			Reply:  fmt.Sprintf("Tarea '%s' creada correctamente.", taskTitle),
			Action: "task_created",
			Data: map[string]string{
				"title":       taskTitle,
				"description": taskDescription,
			},
			SessionID: session.ID,
			State:     string(session.State),
		}, nil

	default:
		session.State = StateIdle
		return &Response{
			Reply:     "Estado desconocido. Reiniciando.",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil
	}
}

// handleListTasks lista las tareas del usuario
func (e *DefaultExecutor) handleListTasks(intent *Intent, session *Session) (*Response, error) {
	session.State = StateIdle

	tasks, ok := e.tasks[session.ID]
	if !ok || len(tasks) == 0 {
		return &Response{
			Reply:     "No tienes tareas creadas aún.",
			Action:    "no_tasks",
			SessionID: session.ID,
			State:     string(session.State),
		}, nil
	}

	taskList := "📋 Tus Tareas:\n\n"
	for i, task := range tasks {
		taskList += fmt.Sprintf("%d. %s\n", i+1, task)
	}

	return &Response{
		Reply:     taskList,
		Action:    "tasks_listed",
		Data:      tasks,
		SessionID: session.ID,
		State:     string(session.State),
	}, nil
}

// calculateReadabilityMetrics calcula las métricas de legibilidad usando la librería NLP
func (e *DefaultExecutor) calculateReadabilityMetrics(textToAnalyze string) (*ReadabilityMetrics, error) {
	myText, err := text.New(textToAnalyze)
	if err != nil {
		return nil, err
	}

	wordCount := myText.WordCount()
	sentenceCount := myText.SentenceCount()
	syllableCount := myText.SyllableCount()

	// Evitar división por cero
	wordsPerSentence := 0.0
	if sentenceCount > 0 {
		wordsPerSentence = float64(wordCount) / float64(sentenceCount)
	}

	syllablesPerWord := 0.0
	if wordCount > 0 {
		syllablesPerWord = float64(syllableCount) / float64(wordCount)
	}

	metrics := &ReadabilityMetrics{
		FleschKincaidEase:  myText.FleschKincaidReadingEase(),
		FleschKincaidGrade: myText.FleschKincaidGradeLevel(),
		WordCount:          wordCount,
		SentenceCount:      sentenceCount,
		SyllableCount:      syllableCount,
		WordsPerSentence:   wordsPerSentence,
		SyllablesPerWord:   syllablesPerWord,
	}

	return metrics, nil
}
