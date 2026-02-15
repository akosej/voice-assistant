// Package voice proporciona el middleware reutilizable para el agente de voz inteligente.
package voice

import (
	"fmt"
)

// Agent agente de voz reutilizable y desacoplado del framework
type Agent struct {
	store    SessionStore
	detector IntentDetector
	executor Executor
}

// NewAgent crea un nuevo Agent
func NewAgent(store SessionStore, detector IntentDetector, executor Executor) *Agent {
	return &Agent{
		store:    store,
		detector: detector,
		executor: executor,
	}
}

// Handle procesa un mensaje de voz/texto y retorna la respuesta
// Es la función principal que puede usarse en cualquier framework
func (a *Agent) Handle(sessionID string, text string) (*Response, error) {
	// Obtener o crear sesión
	session, err := a.store.Get(sessionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener sesión: %w", err)
	}

	// Detectar intención
	intent, err := a.detector.Detect(text)
	if err != nil {
		return nil, fmt.Errorf("error al detectar intención: %w", err)
	}

	// Asegurar que existan las entidades
	if intent.Entities == nil {
		intent.Entities = make(map[string]string)
	}
	intent.Entities["raw"] = text

	// Ejecutar la acción
	response, err := a.executor.Execute(intent, session)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar: %w", err)
	}

	// Guardar sesión con el nuevo estado
	if err := a.store.Save(session); err != nil {
		return nil, fmt.Errorf("error al guardar sesión: %w", err)
	}

	return response, nil
}
