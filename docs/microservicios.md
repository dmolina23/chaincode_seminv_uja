# Implementación de microservicios para la implementación del lado del cliente
La arquitectura de microservicios es una muy buena elección para integrar múltiples proyectos basados en blockchain, ya que permite escalabilidad, modularidad y facilita la interoperabilidad entre los distintos sistemas.

Antes de ver las recomendaciones sobre este tipo de implementación, vamos a definir cuándo y cuándo no debemos utilizarla:

### Cuándo usar Microservicios
* Si los proyectos necesitan independencia y escalabilidad.
* Si cada sistema tiene lógica y requisitos muy distintos.
* Si planean integrar más proyectos a futuro.

### Cuándo NO usar Microservicios
* Si el equipo no tiene experiencia en manejar sistemas distribuidos.
* Si los sistemas están demasiado interconectados y la separación no es clara.
* Si el tráfico es bajo y no hay necesidad de escalabilidad inmediata.


Conociendo esto, yo utilizaría **microservicios** tal y como se ha propuesto.


## Recomendaciones para el Desarrollo de Microservicios
Para construir nuestro sistema, sugiero un stack tecnológico moderno y robusto:

### Elección de Lenguaje y Frameworks
Dependiendo del equipo, se pueden utilizar varios lenguajes:

* **Node.js (Express / NestJS)**: Ideal para aplicaciones rápidas y escalables. Nos proporciona una mejor integración con el desarrollo del FrontEnd (que se realizará con Angular y Node)
* **Go (Gin / Fiber)**: Excelente rendimiento y compatibilidad con Hyperledger Fabric.
* **Python (FastAPI / Flask)**: Buenas herramientas para integración con Machine Learning y análisis de datos.

Por facilidades de integración la aplicación web y móvil, yo utilizaría NodeJS con Express. Aunque el desarrollo de FastAPI con Python puede ser mucho más sencillo.


###  Comunicación entre Microservicios
Para que los microservicios interactúen correctamente, podríamos utilizar:

* **gRPC** → Comunicación eficiente y rápida con soporte para múltiples lenguajes.
* **RabbitMQ** / Kafka → Sistemas de mensajería para eventos asincrónicos.
* **REST** / GraphQL → API HTTP para interacciones con clientes externos.

En este caso, si los microservicios deben responder en tiempo real, gRPC sería la mejor opción, mientras que si manejan eventos, Kafka o RabbitMQ pueden ser más eficientes.

Personalmente, con la única que estoy familiarizado es con la creación una API utilizando REST y GraphQL.


### Orquestación y Despliegue
Para gestionar los microservicios en producción:

* **Docker + Kubernetes (K8s)**: Para contenerización y escalabilidad automática.
* **Docker Compose**: Para entornos de desarrollo local.
* **CI/CD (GitHub Actions, GitLab CI/CD, Jenkins)**: Para automatizar despliegues y pruebas.


### Bases de Datos
Cada microservicio puede necesitar bases de datos diferentes según su función:

* **MongoDB / PostgreSQL** → Para almacenar datos generales del sistema.
* **IPFS / Hyperledger Fabric Ledger** → Para el almacenamiento descentralizado y auditable de datos en blockchain.
* **Redis** → Para mejorar la velocidad de lectura y caché en consultas frecuentes.


### Seguridad y Gestión de Identidades
Si se está manejando una implementación de SSI (Self-Sovereign Identity), es clave considerar:

* **Hyperledger Indy / Aries / Ursa** → Para identidad descentralizada.
* **OAuth 2.0 / OpenID Connect** → Para autenticación y autorización de usuarios externos.
* **JWT (JSON Web Tokens)** → Para autenticación entre microservicios.

Creo recordar que en el seminario se propuso utilizar JWT.


## Posibl Alternativa a los Microservicios: Arquitectura Modular en Monolito
Si el equipo no está preparado para manejar microservicios, una arquitectura modular monolítica también podría funcionar.
* Se pueden dividir los proyectos en módulos dentro de un mismo backend, en lugar de hacer microservicios separados.
* Menos sobrecarga de infraestructura y más simple de desarrollar y desplegar.

Si el equipo ya tiene experiencia con microservicios, adelante con esa opción. Pero si es la primera vez trabajando en una arquitectura distribuida, puede ser más eficiente iniciar con un monolito bien modularizado.

## Resumen de herramientas recomendadas
*  Lenguajes: Node.js (NestJS), Go (Gin), Python (FastAPI)
*  Comunicación: gRPC / Kafka / RabbitMQ
*  Orquestación: Docker + Kubernetes
*  Base de datos: MongoDB, PostgreSQL, IPFS, Redis
*  Seguridad: Hyperledger Indy, OAuth 2.0, JWT

