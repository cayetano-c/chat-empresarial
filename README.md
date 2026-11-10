# Chat Empresarial (Cliente-Servidor en Go con FastAPI)

Proyecto Integrador, Tercer Nivel
Universidad Internacional del Ecuador, Facultad de Ciencias Técnicas
Ingeniería en Sistemas de Información

## Integrantes

- Darwin Alcarraz
- José Almeida
- Cayetano Córdova

## Descripción del sistema

Muchas empresas pequeñas y equipos de trabajo se comunican por WhatsApp o Telegram, y así la empresa no controla nada de esas conversaciones: las decisiones quedan en el celular de cada persona, no se sabe quién ve la información y no hay registro de uso.

Para resolverlo proponemos un chat empresarial interno y cerrado, para una empresa con dos sedes, con tres componentes de software:

- **Cliente en Go:** la aplicación del empleado para iniciar sesión, enviar mensajes y ver el historial.
- **Servidor de chat en Go:** recibe y guarda los mensajes y los entrega al destinatario.
- **Servicio de autenticación en FastAPI:** controla quién puede entrar y durante cuánto tiempo.

Los componentes se comunican por HTTP con datos en JSON y funcionan sobre una red de dos sedes, con un router OpenWrt (Raspberry Pi) en cada una, VLAN separadas (10 Usuarios, 20 Servidores, 30 Invitados y 99 Administración), OSPF sobre un enlace Wi-Fi compartido y direccionamiento VLSM desde 10.10.0.0/16.

En el prototipo de la semana 10 están implementados el envío y recepción de mensajes y el historial por páginas. La autenticación está simulada con usuarios de prueba, a la espera del servicio FastAPI, y el registro de datos de uso (RF-05) queda pendiente.

## Alcance preliminar

Incluye: inicio y cierre de sesión, mensajes de texto entre dos usuarios, historial por páginas y registro de datos de uso.

No incluye (por ahora): envío de archivos, chats de grupo, notificaciones y registro de usuarios nuevos desde la app.

## Casos de uso principales

1. CU-01: Iniciar sesión
2. CU-02: Enviar mensaje
3. CU-03: Consultar historial
4. CU-04: Cerrar sesión

## Funcionalidades
Inicio y cierre de sesión en el cliente, con usuarios y token de prueba (simulados).
Envío de mensajes de texto entre dos usuarios registrados.
Recepción de mensajes nuevos
Historial de una conversación, por páginas de 50 mensajes y ordenado por fecha.
Guardado de los mensajes en un archivo JSON, que se recarga al reiniciar el servidor.
Validaciones y mensajes de error claros (mensaje vacío, destinatario inexistente, token inválido, servidor apagado, etc.).
IP y puerto configurables en el servidor y en el cliente.
Pruebas automáticas en los seis paquetes del proyecto.

## Requisitos de Software
Go versión 1.27.3
Una terminal en Linux o MacOS
curl

## Instalación 
git clone https://github.com/cayetano-c/chat-empresarial.git
cd chat-empresarial/app-go
go test ./...
## Ejecución del Prototipo
1. cd app-go
2. go test ./...
3. TERMINAL 1: go run ./cmd/server
4. TERMINAL 2: go run ./cmd/client
5. TERMINAL 3: go run ./cmd/client
6. Iniciar sesión con cualquiera de los 3 usuarios: "ana", "luis" o "marta"

## Dirección IP, puerto y parámetros configurables.
El servidor acepta cuatro parámetros: -host, la dirección IP en la que escucha (por defecto 127.0.0.1; con 0.0.0.0 acepta conexiones de otras máquinas); -port, el puerto (por defecto 8080); -users, el archivo con los usuarios de prueba (por defecto data/users.json); y -data, el archivo donde se guardan los mensajes (por defecto data/messages.json). El cliente acepta -host y -port (por defecto 127.0.0.1 y 8080), que deben coincidir con la dirección del servidor al que se conecta.
Ejemplo: go run ./cmd/server -host 0.0.0.0 -port 9000 y go run ./cmd/client -host 192.168.1.10 -port 9000.

## Datos de Prueba
app-go/data/users.json: tres usuarios registrados, ana, luis y marta.
En el prototipo el token de cada usuario es demo- más su nombre (por ejemplo demo-ana). No son credenciales reales: es una simulación temporal que se reemplazará por el servicio FastAPI.
## Estructura del repositorio

```
chat-empresarial/
├── README.md
├── LICENSE
├── .gitignore
├── app-go/
│   ├── cmd/
│   │   ├── server/main.go
│   │   └── client/main.go
│   ├── internal/
│   │   ├── models/        # User y Message
│   │   ├── storage/       # almacén de mensajes (JSON)
│   │   ├── auth/          # interfaz de autenticación + simulación temporal
│   │   ├── services/      # lógica del chat
│   │   ├── handlers/      # endpoints HTTP
│   │   └── apiclient/     # comunicación del cliente con el servidor
│   ├── data/
│   │   └── users.json
│   ├── go.mod
│   └── README.md
├── auth-service/
│   └── README.md
├── network/
│   ├── validation-plan.md
│   ├── addressing/
│   │   ├── addressing-plan.csv
│   │   ├── addressing-plan.md
│   │   └── README.md
│   ├── configs/
│   │   └── README.md
│   └── diagrams/
│       ├── Logical-Topology.png
│       ├── Physical-Topology.png
│       └── README.md
├── research/
│   └── instruments/
│       └── README.md
└── docs/
    ├── diagrama-casos-de-uso.png
    ├── ProyectoIntegradorAlc-Alm-Cay.pdf
    ├── gta1-phase1-network-plan.pdf
    └── (documento técnico del segundo entregable, agregar)
```
- `app-go/`: aplicación en Go (cliente y servidor de chat).
- `auth-service/`: servicio de autenticación en FastAPI.
- `network/`: diagramas, plan de direccionamiento y plan de validación de la red.
- `research/`: instrumentos de la investigación estadística.
- `docs/`: documentos entregados.


## Conceptos de POO aplicados
Structs para las entidades: User, Message, Store, MockAuth, ChatService, Handler, Client.
Constructores que garantizan un estado válido: NewUser, NewMessage, NewStore, NewMockAuth, NewChatService.
Encapsulamiento: User y Store tienen campos no exportados y se usan mediante métodos. Message tiene campos exportados porque viaja como JSON.
Métodos con receptor por puntero cuando modifican el estado (Save, SetDisplayName).
Composición: ChatService contiene un almacén y un autenticador; Handler contiene el servicio.
Interfaces pequeñas donde separan responsabilidades: MessageStore y Authenticator. Permiten cambiar la forma de guardar o de autenticar sin tocar la lógica del chat.
Manejo explícito de errores con el tipo error, sin panic.
Concurrencia: sync.Mutex en el almacén y una goroutine en el cliente.
## Entregables por asignatura

| Asignatura | Documento |
|---|---|
| Programación Orientada a Objetos | [ProyectoIntegradorAlca-Alm-Cay.pdf](docs/ProyectoIntegradorAlca-Alm-Cay.pdf): análisis y requisitos del chat |
| Tecnologías de Comunicación en Red | [GTA1_Phase_1_Integrated Project.pdf](docs/GTA1_Phase_1_Integrated%20Project.pdf): plan de la infraestructura de red |
docs/Desarrollo del Sistemas de Gestión- Fase 2 -.pdf
## Video
https://drive.google.com/file/d/1Yzkld3M9PBMLDh9iN60PnWhdQPSzVp1O/view
https://drive.google.com/file/d/1ZV4PqNrcFLGNOaluqNKmE9oRu0QNwYyP/view?usp=sharing

## Limitaciones
La autenticación es una simulación (sin contraseñas ni vencimiento de tokens); se reemplazará con FastAPI en las semanas 11 y 12.
Aún no se registran las estadísticas de uso (RF-05).
Los mensajes se guardan en un único archivo JSON que se reescribe en cada envío; es suficiente para el prototipo, no para un volumen grande.
Los mensajes nuevos se obtienen consultando cada 3 segundos, no en tiempo real.
La comunicación es HTTP sin cifrar, pensada para un entorno local o de laboratorio.
Quedan fuera del alcance los chats de grupo, el envío de archivos y el registro de usuarios desde la aplicación.
## Autores y contribuciones
## Darwin Alcarraz: 
(/client) Encargado del desarrollo y la estructura del cliente en Go. Diseñó la interfaz orientada a comandos para el intercambio de mensajes y estableció los mecanismos de conectividad y sincronización con el servidor, asegurando la estabilidad operativa del sistema mediante pruebas automatizadas (client_test.go).
## José Almeida: 
(Backend / Server) Responsable del diseño y la implementación del servidor central en Go. se encargo de el enrutador de peticiones, el ciclo de vida de las conexiones y el manejo de concurrencia mediante goroutines y canales para soportar la comunicación simultánea de múltiples clientes sin bloqueos. Asimismo colaboró con las pruebas unitarias de los manejadores de red (handlers_test.go).
## Cayetano Córdova: 
(data / storage) Responsable de la estructuración del almacenamiento local y el modelado de entidades (messages.json y users.json). Implementó los servicios de persistencia y serialización para garantizar el registro seguro de usuarios y el historial de chat, además de coordinar la integración técnica con la documentación y los entregables del repositorio.

## Créditos, licencia y uso de inteligencia artificial
Licencia: MIT (ver LICENSE).
Referencias del análisis del problema: Chui et al. (2012), McKinsey Global Institute; Guthkelch (2025), The HR Director (detalle en el documento de docs/).
Uso de IA: Se utilizó Claude (Anthropic) como asistente durante el desarrollo del prototipo de la semana 10. Claude ayudó a proponer la estructura de carpetas y la división en paquetes, generó borradores y correciones del código en Go, y explicó conceptos de Go (structs, constructores, interfaces, concurrencia). El equipo revisó el código, lo ejecutó, corrió las pruebas y probó el funcionamiento manualmente con el servidor y varios clientes.