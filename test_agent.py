"""
Ejemplos de uso del Agente de Voz desde Python.
Asegúrate de que el servidor esté corriendo en http://localhost:3000
"""

import requests
import json
import time

BASE_URL = "http://localhost:3000"
SESSION_ID = "python_cliente"


def call_agent(text, session_id=SESSION_ID):
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


def print_response(response):
    """Imprime la respuesta de forma legible."""
    if not response:
        return

    print("\n" + "=" * 80)
    print(f"📌 Sesión: {response.get('session_id')}")
    print(f"🔄 Estado: {response.get('state')}")
    print(f"⚡ Acción: {response.get('action', 'N/A')}")
    print("\n💬 Respuesta:")
    print(response.get('reply', 'Sin respuesta'))

    if response.get('data'):
        print("\n📊 Datos:")
        if isinstance(response['data'], dict):
            for key, value in response['data'].items():
                if isinstance(value, float):
                    print(f"   {key}: {value:.2f}")
                else:
                    print(f"   {key}: {value}")
        else:
            print(f"   {response['data']}")

    print("=" * 80 + "\n")


def test_greeting():
    """Prueba 1: Saludar al agente."""
    print("\n🧪 Test 1: Saludar al agente")
    response = call_agent("Hola")
    print_response(response)


def test_help():
    """Prueba 2: Pedir ayuda."""
    print("\n🧪 Test 2: Pedir ayuda")
    response = call_agent("ayuda")
    print_response(response)


def test_readability_analysis():
    """Prueba 3: Analizar legibilidad."""
    print("\n🧪 Test 3: Analizar legibilidad")
    
    # Paso 1: Iniciar análisis
    print("Paso 1: Solicitar análisis de legibilidad")
    response1 = call_agent("analiza este texto")
    print_response(response1)

    # Paso 2: Enviar el texto a analizar
    time.sleep(1)
    print("Paso 2: Enviar texto para analizar")
    texto = "La programación es el arte de escribir código que otro humano pueda entender. El computador es secundario. Es importante escribir código limpio, legible y mantenible."
    response2 = call_agent(texto)
    print_response(response2)


def test_task_creation():
    """Prueba 4: Crear tareas."""
    print("\n🧪 Test 4: Crear tareas")

    # Paso 1: Iniciar creación de tarea
    print("Paso 1: Solicitar creación de tarea")
    response1 = call_agent("crear tarea")
    print_response(response1)

    # Paso 2: Enviar nombre de la tarea
    time.sleep(1)
    print("Paso 2: Enviar nombre de la tarea")
    response2 = call_agent("Implementar autenticación OAuth2")
    print_response(response2)


def test_list_tasks():
    """Prueba 5: Listar tareas."""
    print("\n🧪 Test 5: Listar tareas")
    response = call_agent("mostrar mis tareas")
    print_response(response)


def test_complex_text_analysis():
    """Prueba 6: Análisis de texto complejo."""
    print("\n🧪 Test 6: Análisis de texto más complejo")

    textos = [
        "Python es simple.",
        "La inteligencia artificial es una rama de la informática que busca crear máquinas inteligentes capaces de realizar tareas que normalmente requieren inteligencia humana.",
        "El aprendizaje automático es un subcampo de la inteligencia artificial que se enfoca en desarrollar algoritmos y sistemas que pueden aprender y mejorar su desempeño basado en experiencias previas."
    ]

    for i, texto in enumerate(textos, 1):
        print(f"\nAnalizando texto #{i}:")
        # Iniciar análisis
        call_agent("analizar legibilidad")
        time.sleep(0.5)
        # Enviar texto
        response = call_agent(texto)
        print_response(response)
        time.sleep(1)


def test_unknown_command():
    """Prueba 7: Comando desconocido."""
    print("\n🧪 Test 7: Comando desconocido")
    response = call_agent("klsdjf asdfjk lkjsdf")
    print_response(response)


def test_multiple_sessions():
    """Prueba 8: Múltiples sesiones."""
    print("\n🧪 Test 8: Múltiples sesiones")

    usuarios = ["usuario_1", "usuario_2", "usuario_3"]

    for usuario in usuarios:
        print(f"\nCreando tarea para {usuario}:")
        # Iniciar creación
        call_agent("crear tarea", session_id=usuario)
        time.sleep(0.5)
        # Enviar nombre
        response = call_agent(f"Tarea del {usuario}", session_id=usuario)
        print_response(response)
        time.sleep(1)

    # Listar tareas de cada usuario
    for usuario in usuarios:
        print(f"\nTareas de {usuario}:")
        response = call_agent("listar tareas", session_id=usuario)
        print_response(response)


def test_health_check():
    """Prueba 9: Verificar salud del servidor."""
    print("\n🧪 Test 9: Verificar salud del servidor")
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        response.raise_for_status()
        data = response.json()
        print(f"✅ Servidor: {data['status'].upper()}")
        print(f"📝 Mensaje: {data['message']}")
    except requests.exceptions.RequestException as e:
        print(f"❌ Error: {e}")


def test_agent_info():
    """Prueba 10: Información del agente."""
    print("\n🧪 Test 10: Información del agente")
    try:
        response = requests.get(f"{BASE_URL}/info", timeout=5)
        response.raise_for_status()
        data = response.json()
        print(f"🤖 Agente: {data['name']} v{data['version']}")
        print(f"\n📋 Capacidades:")
        for cap in data['capabilities']:
            print(f"   ✓ {cap}")
        print(f"\n🔌 Endpoints:")
        for endpoint, desc in data['endpoints'].items():
            print(f"   {endpoint}: {desc}")
    except requests.exceptions.RequestException as e:
        print(f"❌ Error: {e}")


def run_all_tests():
    """Ejecuta todas las pruebas."""
    print("\n" + "=" * 80)
    print("🚀 INICIANDO SUITE DE PRUEBAS DEL AGENTE DE VOZ")
    print("=" * 80)

    try:
        test_health_check()
        test_agent_info()
        test_greeting()
        test_help()
        test_unknown_command()
        test_readability_analysis()
        test_task_creation()
        test_list_tasks()
        test_complex_text_analysis()
        test_multiple_sessions()

        print("\n" + "=" * 80)
        print("✅ TODAS LAS PRUEBAS COMPLETADAS")
        print("=" * 80)

    except KeyboardInterrupt:
        print("\n\n⚠️ Pruebas interrumpidas por el usuario")


if __name__ == "__main__":
    print("""
    ╔═══════════════════════════════════════════════════════════════╗
    ║  🎙️  CLIENTE DE PRUEBAS - AGENTE DE VOZ INTELIGENTE          ║
    ╚═══════════════════════════════════════════════════════════════╝
    
    Asegúrate de que el servidor esté corriendo:
    $ go run main.go
    
    El servidor debe estar escuchando en: http://localhost:3000
    """)

    input("Presiona Enter para comenzar las pruebas...")
    run_all_tests()
