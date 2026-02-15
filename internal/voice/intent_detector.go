// Package voice proporciona el middleware reutilizable para el agente de voz inteligente.
package voice

import (
	"strings"
)

// IntentDetector interfaz para detectar intenciones
type IntentDetector interface {
	Detect(text string) (*Intent, error)
}

// RuleBasedDetector detector de intenciones basado en reglas
type RuleBasedDetector struct{}

// NewRuleBasedDetector crea un nuevo RuleBasedDetector
func NewRuleBasedDetector() *RuleBasedDetector {
	return &RuleBasedDetector{}
}

// Detect detecta la intención del usuario basado en palabras clave
func (r *RuleBasedDetector) Detect(text string) (*Intent, error) {
	textLower := strings.ToLower(text)
	text = strings.TrimSpace(text)

	// Detectar intención de analizar texto
	if contains(textLower, []string{"analizar", "legibilidad", "leer", "texto"}) {
		return &Intent{
			Name: "analyze_readability",
			Entities: map[string]string{
				"type": "readability",
			},
			Raw: text,
		}, nil
	}

	// Detectar intención de crear tarea
	if contains(textLower, []string{"crear", "tarea"}) {
		return &Intent{
			Name: "create_task",
			Entities: map[string]string{
				"type": "task",
			},
			Raw: text,
		}, nil
	}

	// Detectar intención de listar tareas
	if contains(textLower, []string{"listar", "ver", "mostrar", "tareas"}) {
		return &Intent{
			Name: "list_tasks",
			Entities: map[string]string{
				"type": "task",
			},
			Raw: text,
		}, nil
	}

	// Detectar intención de saludo
	if contains(textLower, []string{"hola", "buenos", "saludos", "ey", "hey"}) {
		return &Intent{
			Name: "greeting",
			Entities: map[string]string{
				"type": "greeting",
			},
			Raw: text,
		}, nil
	}

	// Detectar intención de ayuda
	if contains(textLower, []string{"ayuda", "help", "qué puedo", "qué puedes"}) {
		return &Intent{
			Name: "help",
			Entities: map[string]string{
				"type": "help",
			},
			Raw: text,
		}, nil
	}

	// Intención desconocida
	return &Intent{
		Name: "unknown",
		Entities: map[string]string{
			"type": "unknown",
		},
		Raw: text,
	}, nil
}

// contains verifica si el texto contiene alguna de las palabras clave
func contains(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
