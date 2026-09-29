# Roadmap de 40 horas: un portfolio que abre conversaciones

## Regla comercial

Este repositorio debe ayudar a conseguir entrevistas y contratos, no convertirse en un proyecto infinito. Publica la primera versión ejecutable hoy y enlázala en propuestas relevantes. Dedica el bloque más fresco de cada día a prospectar contratos; el portfolio mejora la conversión de esa prospección. Un repositorio excelente no garantiza que una empresa omita su prueba técnica.

## Qué enseñar al comprador

**Posicionamiento:** ingeniero principal de backend en Go y modernización de sistemas con experiencia real en WMS, pagos y plataformas con dinero transaccionado. El demo enseña un patrón concreto: extraer cambios de un sistema MySQL legado, estandarizarlos y preparar una proyección segura frente a duplicados.

**Prueba pública:** código ejecutable, escenario de fallo, tests, métricas reproducibles, decisión documentada, vídeo breve. La experiencia comercial se cuenta por separado con cifras que puedan publicarse o con autorización del cliente. No afirmar que estos clientes usan exactamente el código de este repositorio.

## Secuencia de 40 horas

| Bloque | Horas | Resultado verificable | Criterio para pasar |
| --- | ---: | --- | --- |
| 1. Publicar la base | 4 | Repo público con README, diagrama, demo y CI | Un tercero ejecuta `make up && make demo` |
| 2. GTID y reinicio | 8 | Checkpoint duradero y reproducción tras caída | Test de reinicio sin pérdida |
| 3. Broker y confirmación | 7 | RabbitMQ con confirmaciones y reintentos | Broker parado no avanza checkpoint |
| 4. Consumidor idempotente | 7 | Proyección con clave única de evento y API simple | Entregar dos veces aplica una sola vez |
| 5. Fallos y observabilidad | 6 | Inyección de caídas, backlog, métricas, runbook | Se puede explicar y reproducir cada fallo |
| 6. Medición y presentación | 5 | Benchmark honesto, diagrama final, vídeo de 90 s | Números acompañados de entorno y comandos |
| 7. Distribución comercial | 3 | Perfil GitHub fijado y caso de estudio enlazado en ofertas | Al menos 10 propuestas relevantes usan el enlace |

Total: **40 horas**. Si aparece una entrevista o un cliente, se interrumpe este plan sin culpa: la finalidad es conseguir ingresos.

## Revisión de calidad antes de presentarlo como proyecto avanzado

- El README separa claramente lo implementado de lo planificado.
- Hay tests de rollback, duplicado, caída entre publish/checkpoint, broker no disponible y reanudación.
- El benchmark incluye hardware, versión de MySQL, tamaño de datos, comandos y p50/p95.
- La demo se reproduce desde cero y el CI la ejecuta.
- No aparecen secretos, nombres ni código de clientes.
- Cada afirmación del CV puede defenderse oralmente; las cifras privadas se atribuyen como tales.

## Mensaje de propuesta corto

> Soy ingeniero backend y de sistemas, disponible como contractor B2B desde Malta. He liderado sistemas de inventario y pagos, y trabajo con Go, MySQL y arquitecturas de integración. Aquí hay una demo ejecutable de cómo extraería cambios de un legado con outbox + binlog, con escenarios de fallo y tests: [enlace al repositorio]. ¿Tenéis un problema parecido de modernización, integraciones o fiabilidad?

Personaliza la primera frase al problema real de cada empresa. No envíes el mismo texto masivamente ni gastes los 35 InMails en contactos sin encaje.
