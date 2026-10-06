# chat-empresarial

Proyecto Integrador, Fase 1-2
Universidad Internacional del Ecuador, Facultad de Ciencias Técnicas
Programación Orientada a Objetos (03-SIN-A)
Docente: Ing. Ivan Galo Reyes Chacon

## Integrantes

- Darwin Alcarraz
- José Almeida
- Cayetano Córdova

## Descripción del sistema

Muchas empresas pequeñas y equipos de trabajo se comunican por WhatsApp o Telegram, y así la empresa no controla nada de esas conversaciones: las decisiones quedan en el celular de cada persona, no se sabe quién ve la información y no hay registro de uso.

Para resolverlo proponemos un chat empresarial interno y cerrado, con tres componentes:

- *Cliente en Go:* la aplicación del empleado para iniciar sesión, enviar mensajes y ver el historial.
- *Servidor de chat en Go:* recibe y guarda los mensajes y los entrega al destinatario.
- *Servicio de autenticación en FastAPI:* crea y revisa el token de sesión y controla cuándo empieza y termina cada sesión.

Los tres se comunican por HTTP con datos en JSON. Cada acción se guarda con usuario, fecha y hora, y de ahí salen datos medibles como mensajes por usuario, sesiones por día y tiempo promedio de sesión.

## Alcance preliminar

Incluye: inicio y cierre de sesión con token, mensajes de texto entre dos usuarios, historial por páginas y registro de datos de uso.

No incluye (por ahora): envío de archivos, chats de grupo, notificaciones, registro de usuarios nuevos desde la app y servicios externos.

## Casos de uso principales

1. CU-01: Iniciar sesión
2. CU-02: Enviar mensaje
3. CU-03: Consultar historial
4. CU-04: Cerrar sesión

## Estructura del repositorio


chat-empresarial/
├── README.md
├── docs/
│   ├── documento-fase1.pdf
│   └── diagrama-casos-de-uso.png
├── client-go/
│   └── README.md
└── auth-fastapi/
    └── README.md


- docs/: documento de análisis y especificación de requisitos (PDF) y diagrama de casos de uso.
- client-go/: cliente del chat en Go.
- auth-fastapi/: servicio de autenticación en FastAPI.

## Video