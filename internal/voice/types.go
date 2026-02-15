// Package voice proporciona el middleware reutilizable para el agente de voz inteligente.
package voice

// State representa el estado de una sesión del usuario.
type State string

const (
	StateIdle                  State = "idle"
	StateCreatingTaskTitle     State = "creating_task_title"
	StateCreatingTaskDescription State = "creating_task_description"
	StateAnalyzingText         State = "analyzing_text"
	StateWaitingForConfirm     State = "waiting_for_confirm"
)

// Session almacena información sobre la sesión de un usuario.
type Session struct {
	ID            string
	State         State
	Context       map[string]string
	TemporaryData map[string]interface{}
}

// Response es la respuesta que retorna el agente.
type Response struct {
	Reply     string      `json:"reply"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	SessionID string      `json:"session_id"`
	State     string      `json:"state"`
}

// Intent representa la intención detectada del usuario.
type Intent struct {
	Name     string
	Entities map[string]string
	Raw      string
}

// ReadabilityMetrics contiene las métricas de análisis de legibilidad.
type ReadabilityMetrics struct {
	FleschKincaidEase  float64 `json:"flesch_kincaid_ease"`
	FleschKincaidGrade float64 `json:"flesch_kincaid_grade"`
	WordCount          int     `json:"word_count"`
	SentenceCount      int     `json:"sentence_count"`
	SyllableCount      int     `json:"syllable_count"`
	WordsPerSentence   float64 `json:"words_per_sentence"`
	SyllablesPerWord   float64 `json:"syllables_per_word"`
}
