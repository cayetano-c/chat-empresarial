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

## Alcance preliminar

Incluye: inicio y cierre de sesión, mensajes de texto entre dos usuarios, historial por páginas y registro de datos de uso.

No incluye (por ahora): envío de archivos, chats de grupo, notificaciones y registro de usuarios nuevos desde la app.

## Casos de uso principales

1. CU-01: Iniciar sesión
2. CU-02: Enviar mensaje
3. CU-03: Consultar historial
4. CU-04: Cerrar sesión

## Estructura del repositorio

```
chat-empresarial/
├── README.md
├── LICENSE
├── app-go/
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
    └── gta1-phase1-network-plan.pdf
```

- `app-go/`: aplicación en Go (cliente y servidor de chat).
- `auth-service/`: servicio de autenticación en FastAPI.
- `network/`: diagramas, plan de direccionamiento y plan de validación de la red.
- `research/`: instrumentos de la investigación estadística.
- `docs/`: documentos entregados.

## Entregables por asignatura

| Asignatura | Documento |
|---|---|
| Programación Orientada a Objetos | [docs/ProyectoIntegradorAlc-Alm-Cay.pdf](docs/ProyectoIntegradorAlc-Alm-Cay.pdf): análisis y requisitos del chat |
| Tecnologías de Comunicación en Red | [docs/gta1-phase1-network-plan.pdf](docs/gta1-phase1-network-plan.pdf): plan de la infraestructura de red |

## Video

[pegar enlace del video]