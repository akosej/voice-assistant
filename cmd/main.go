package main

import (
	"agent/internal/voice"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var voiceAgent *voice.Agent

// Inicializar el agente de voz con sus dependencias
func initVoiceAgent() {
	store := voice.NewMemoryStore()
	detector := voice.NewRuleBasedDetector()
	executor := voice.NewDefaultExecutor()

	voiceAgent = voice.NewAgent(store, detector, executor)

	log.Println("✅ Agente de voz inicializado correctamente")
}

// VoiceRequest es una solicitud al agente de voz.
type VoiceRequest struct {
	SessionID string `json:"session_id" form:"session_id"`
	Text      string `json:"text" form:"text"`
}

// VoiceErrorResponse es una respuesta de error.
type VoiceErrorResponse struct {
	Error string `json:"error"`
}

// HandleVoice maneja las solicitudes de voz
func HandleVoice(c *fiber.Ctx) error {
	var req VoiceRequest

	// Intentar parsear como JSON primero, luego como form data
	if err := c.BodyParser(&req); err != nil {
		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(VoiceErrorResponse{
				Error: "solicitud inválida: necesitas session_id y text",
			})
		}
	}

	// Validar campos requeridos
	if req.SessionID == "" || req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(VoiceErrorResponse{
			Error: "session_id y text son requeridos",
		})
	}

	// Procesar con el agente
	response, err := voiceAgent.Handle(req.SessionID, req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(VoiceErrorResponse{
			Error: fmt.Sprintf("error procesando solicitud: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleHealth verifica si el servidor está activo
func HandleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Agente de Voz en línea",
	})
}

// HandleInfo devuelve información del agente
func HandleInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":    "Agente de Voz Inteligente",
		"version": "1.0.0",
		"capabilities": []string{
			"Analizar legibilidad de textos",
			"Gestionar tareas",
			"Detección de lenguaje natural",
			"Almacenamiento de sesiones",
		},
		"endpoints": fiber.Map{
			"POST /voice": "Procesar comando de voz",
			"GET /health": "Verificar estado del servidor",
			"GET /info":   "Información del agente",
			"GET /":       "Información principal",
		},
	})
}

func main() {
	// Inicializar el agente
	initVoiceAgent()

	// Crear la aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName: "Agente de Voz Inteligente v1.0.0",
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Rutas
	app.Get("/health", HandleHealth)
	app.Get("/info", HandleInfo)
	app.Post("/voice", HandleVoice)

	// Ruta raíz
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"mensaje":      "Bienvenido al Agente de Voz Inteligente",
			"info_url":     "/info",
			"voice_api":    "POST /voice",
			"health_check": "/health",
			"ejemplo": fiber.Map{
				"session_id": "usuario_123",
				"text":       "Hola, analiza este texto: La inteligencia artificial es fascinante",
			},
		})
	})

	// Iniciar servidor en puerto 3000
	port := ":3000"
	log.Printf("🚀 Servidor iniciado en http://localhost%s\n", port)
	log.Printf("📝 Documenta tu API en http://localhost%s/info\n", port)
	log.Printf("🧪 Prueba con: curl -X POST http://localhost%s/voice -H 'Content-Type: application/json' -d '{\"session_id\":\"test\",\"text\":\"Hola\"}'\n", port)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
