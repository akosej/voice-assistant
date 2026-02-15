#!/usr/bin/env python3
"""
Prueba interactiva del nuevo flujo de creación de tareas.
Asegúrate de que el servidor esté corriendo en http://localhost:3000
"""

import requests
import json
import sys

BASE_URL = "http://localhost:3000"


def call_agent(session_id, text):
    """Llama al agente con un texto."""
    try:
        response = requests.post(
            f"{BASE_URL}/voice",
            json={
                "session_id": session_id,
                "text": text
            },
            timeout=5
        )
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"❌ Error: {e}")
        return None


def print_response(response, show_raw=False):
    """Imprime la respuesta de forma legible."""
    if not response:
        return

    print(f"\n{'─' * 70}")
    print(f"💬 Agente: {response.get('reply', 'Sin respuesta')}")
    print(f"⚡ Acción: {response.get('action', 'N/A')}")
    print(f"🔄 Estado: {response.get('state', 'N/A')}")

    if response.get('data'):
        print(f"📊 Datos: {json.dumps(response['data'], indent=2, ensure_ascii=False)}")

    if show_raw:
        print(f"\n🔍 Respuesta completa:")
        print(json.dumps(response, indent=2, ensure_ascii=False))

    print(f"{'─' * 70}\n")


def test_create_task_workflow():
    """Prueba completa del flujo de creación de tareas."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║           🎯 PRUEBA: CREAR TAREA CON TÍTULO Y DESCRIPCIÓN        ║
╚════════════════════════════════════════════════════════════════════╝
    """)

    session_id = "test_task_flow"

    print("Paso 1️⃣ : Usuario dice 'crear tarea'\n")
    print(">>> crear tarea")
    response1 = call_agent(session_id, "crear tarea")
    print_response(response1)

    # Esperar entrada del usuario
    print("Paso 2️⃣ : Usuario envía el TÍTULO\n")
    title = input(">>> ").strip()
    if not title:
        print("❌ No enviaste nada")
        return

    response2 = call_agent(session_id, title)
    print_response(response2)

    print("Paso 3️⃣ : Usuario envía la DESCRIPCIÓN (o 'saltar')\n")
    description = input(">>> ").strip()

    response3 = call_agent(session_id, description or "saltar")
    print_response(response3)

    print("\nPaso 4️⃣ : Verificar que la tarea se creó\n")
    print(">>> mostrar tareas")
    response4 = call_agent(session_id, "mostrar tareas")
    print_response(response4)


def test_multiple_tasks():
    """Crear varias tareas y listarlas."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║            🎯 PRUEBA: CREAR MÚLTIPLES TAREAS                      ║
╚════════════════════════════════════════════════════════════════════╝
    """)

    session_id = "test_multiple_tasks"
    tasks = [
        ("Aprender Go", "Completar el tutorial oficial"),
        ("Escribir documentación", "Documentar el proyecto"),
        ("Hacer ejercicio", "saltar"),  # Sin descripción
    ]

    for i, (title, description) in enumerate(tasks, 1):
        print(f"\n📝 Crear tarea {i}: '{title}'")

        # Iniciar creación
        call_agent(session_id, "crear tarea")

        # Enviar título
        response = call_agent(session_id, title)
        print(f"✅ Título: {response.get('action')}")

        # Enviar descripción
        response = call_agent(session_id, description)
        print(f"✅ Tarea creada: {response.get('action')}")

    # Listar todas
    print("\n\n📋 Listando todas las tareas:")
    print(">>> mostrar tareas")
    response = call_agent(session_id, "mostrar tareas")
    print_response(response)


def test_edge_cases():
    """Prueba casos especiales."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║            🧪 PRUEBA: CASOS ESPECIALES                            ║
╚════════════════════════════════════════════════════════════════════╝
    """)

    print("\n1️⃣ Descripción vacía (presionar Enter sin escribir):")
    session1 = "test_edge_empty_description"
    call_agent(session1, "crear tarea")
    call_agent(session1, "Mi tarea")
    response = call_agent(session1, "")
    print_response(response)

    print("\n2️⃣ Descripción con saltar:")
    session2 = "test_edge_skip"
    call_agent(session2, "crear tarea")
    call_agent(session2, "Otra tarea")
    response = call_agent(session2, "saltar")
    print_response(response)

    print("\n3️⃣ Descripción muy larga:")
    session3 = "test_edge_long_desc"
    call_agent(session3, "crear tarea")
    call_agent(session3, "Proyecto grande")
    long_desc = "Esta es una descripción muy larga que contiene múltiples líneas de texto " * 5
    response = call_agent(session3, long_desc)
    print_response(response)

    print("\n4️⃣ Caracteres especiales:")
    session4 = "test_edge_special_chars"
    call_agent(session4, "crear tarea")
    call_agent(session4, "Tarea con 🎉 emojis 🚀")
    response = call_agent(session4, "Descripción con símbolos: @#$%^&*()")
    print_response(response)


def test_rapid_requests():
    """Prueba múltiples sesiones simultáneamente."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║        🚀 PRUEBA: MÚLTIPLES USUARIOS AL MISMO TIEMPO              ║
╚════════════════════════════════════════════════════════════════════╝
    """)

    usuarios = ["alice", "bob", "charlie", "diana"]

    for usuario in usuarios:
        print(f"\n👤 Usuario: {usuario}")

        # Crear tarea
        call_agent(usuario, "crear tarea")
        response = call_agent(usuario, f"Tarea de {usuario}")
        response = call_agent(usuario, f"Descripción para {usuario}")

        print(f"   ✅ {response.get('action')}")


def print_menu():
    """Muestra el menú de opciones."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║     🎙️  PRUEBAS DEL FLUJO DE CREACIÓN DE TAREAS (MEJORADO)       ║
╚════════════════════════════════════════════════════════════════════╝

Selecciona una prueba:

1. ✨ Flujo normal: Crear tarea con título y descripción (INTERACTIVO)
2. 📋 Crear múltiples tareas
3. 🧪 Casos especiales (vacío, special chars, etc.)
4. 🚀 Múltiples usuarios al mismo tiempo
5. 🔍 Prueba manual (ingresa comandos libremente)
0. ❌ Salir

Opción: """)


def manual_mode():
    """Modo manual para ingresar comandos."""
    print("""
╔════════════════════════════════════════════════════════════════════╗
║               🔍 MODO MANUAL - INGRESA COMANDOS                   ║
╚════════════════════════════════════════════════════════════════════╝

Escribe 'salir' para terminar
Usa una sesión consistente para las pruebas (ej: 'test', 'usuario1')
    """)

    while True:
        print("\n" + "─" * 70)
        session_id = input("Session ID: ").strip()
        if session_id.lower() == "salir":
            break

        text = input("Comando: ").strip()
        if text.lower() == "salir":
            break

        response = call_agent(session_id, text)
        print_response(response, show_raw=False)


def main():
    """Función principal."""
    while True:
        print_menu()

        try:
            opcion = input().strip()

            if opcion == "1":
                test_create_task_workflow()
            elif opcion == "2":
                test_multiple_tasks()
            elif opcion == "3":
                test_edge_cases()
            elif opcion == "4":
                test_rapid_requests()
            elif opcion == "5":
                manual_mode()
            elif opcion == "0":
                print("\n👋 Hasta luego!")
                sys.exit(0)
            else:
                print("\n❌ Opción no válida. Intenta de nuevo.")

            input("\nPresiona Enter para continuar...")

        except KeyboardInterrupt:
            print("\n\n⚠️ Interrumpido por el usuario")
            sys.exit(0)
        except Exception as e:
            print(f"\n❌ Error: {e}")


if __name__ == "__main__":
    # Verificar conexión al servidor
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=2)
        response.raise_for_status()
        print("✅ Servidor conectado\n")
    except requests.exceptions.RequestException as e:
        print(f"""
❌ No se puede conectar al servidor en {BASE_URL}

Asegúrate de que el servidor esté corriendo:
  $ go run main.go

Error: {e}
        """)
        sys.exit(1)

    main()
