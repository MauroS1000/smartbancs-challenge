# Declaración de Uso de Inteligencia Artificial

Como parte de los lineamientos de transparencia de la evaluación técnica para el backend de SmartBancs, se detalla a continuación el uso de asistentes de Inteligencia Artificial Generativa durante el desarrollo del proyecto.

## 1. Herramientas Utilizadas

*   **Gemini (Google):** Utilizado como asistente de programación (pair-programming), consultor de arquitectura de software y revisor de redacción técnica.

## 2. Áreas de Aplicación

La Inteligencia Artificial fue empleada exclusivamente como un acelerador de desarrollo en las siguientes áreas:

*   **Generación y Refactorización de Código:** Apoyo en la estructuración idiomática de los paquetes en Go, resolución de dependencias, y estructuración de la lógica de concurrencia (*goroutines*). Asimismo, se utilizó para la generación del *boilerplate* del script de transformación de datos (ETL) en Python.
*   **Diseño de Arquitectura:** Validación teórica de los patrones de microservicios utilizados (como el Patrón Outbox) y estrategias de mitigación de *Deadlocks* en bases de datos relacionales bajo alta concurrencia.
*   **Documentación Técnica:** Revisión ortográfica, estructuración de formato Markdown y consolidación de ideas técnicas para la redacción de los documentos formales (README, Documento de Arquitectura e Informe Post Mortem).

## 3. Justificación y Propiedad

El uso de estas herramientas se limitó a potenciar la productividad y validar conceptos técnicos. Todas las decisiones de diseño, la configuración de la infraestructura (Docker), la lógica de negocio central y la ejecución de las pruebas fueron dirigidas, validadas y certificadas manualmente por el desarrollador para garantizar el cumplimiento estricto de los requisitos funcionales y no funcionales del reto.
