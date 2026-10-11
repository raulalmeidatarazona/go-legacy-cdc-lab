# Guion de vídeo: Go Legacy CDC Lab

Duración objetivo: 3–5 minutos. Grabar con el README y el terminal visibles. El ejemplo usa datos sintéticos; ejecutar desde la raíz del repo.

## Guion principal (3–5 min)

### 0:00–0:25 · Gancho

**En cámara:**

> ¿Cómo conectas un sistema legado con servicios nuevos sin pedirle al escritor que haga una segunda escritura a un broker? En este laboratorio enseño una respuesta pequeña y ejecutable: una outbox transaccional en MySQL y un lector de binlog escrito en Go.

**Frase clave en inglés:** “Capture committed changes without coupling the legacy writer to a broker.”

### 0:25–1:05 · Problema y alcance

**En cámara, señalando el pitch y el diagrama del README:**

> Imaginemos un sistema de inventario que ya escribe en MySQL. Consultar la tabla cada pocos segundos genera polling aunque no haya cambios. Hacer que el mismo proceso escriba en MySQL y en un broker crea dos destinos que pueden fallar por separado.
>
> Aquí el escritor inserta un movimiento. Un trigger actualiza el inventario y añade un evento a `outbox_events` dentro de la misma transacción. Después, un proceso en Go lee el cambio comprometido del row binlog y lo representa como CloudEvent.
>
> Es una demo acotada, no una plataforma CDC lista para producción.

**Frase clave en inglés:** “The database transaction is the atomicity boundary.”

### 1:05–2:15 · Demo paso a paso

**Compartir pantalla. En el terminal, desde la raíz del repositorio, ejecutar:**

```bash
make up && make demo
```

**Narración mientras arranca:**

> Este comando levanta MySQL 8.4 con Docker Compose y ejecuta el demo. La base se inicializa con un SKU sintético, `DEMO-001`, y una cantidad de cien.

**Cuando aparezca el JSON en stdout, señalar:**

> El programa primero lee la posición actual del binlog y luego inserta un movimiento de menos uno. El trigger actualiza el saldo e inserta la fila de outbox. Como el lector empieza desde esa posición, captura el cambio que acaba de ocurrir.
>
> Esta salida es un CloudEvent: el `id` identifica el evento, `type` indica `inventory.stock-moved.v1`, `subject` es el SKU y `data` contiene el movimiento. Los IDs y el número de secuencia pueden cambiar en cada ejecución.

**Enseñar brevemente los tests y ejecutar:**

```bash
make test
```

> Aquí están los tests de conversión de eventos y validación del payload. También existe una prueba de integración que verifica rollback con MySQL local:

```bash
INTEGRATION=1 go test ./...
```

> Esa prueba requiere haber ejecutado `make up`.

### 2:15–3:45 · Decisiones técnicas

**Decisión 1 — outbox en la misma transacción:**

> El movimiento, el cambio de saldo y el evento se escriben en la misma transacción. Si la restricción de saldo no negativo rechaza el cambio, también se revierte la fila de movimiento y no queda un evento huérfano. El coste es trabajo adicional en la escritura y dependencia del trigger.

**Frase clave en inglés:** “One committed source transaction produces one outbox row.”

**Decisión 2 — leer el binlog:**

> Para esta demo, leer cambios del binlog evita consultar repetidamente la outbox. A cambio, el proceso necesita permisos de replicación y el binlog debe conservarse el tiempo suficiente. Eso desplaza parte de la complejidad a la operación.

**Frase clave en inglés:** “The binlog trades polling queries for replication and retention concerns.”

**Decisión 3 — honestidad sobre entrega:**

> El evento tiene un ID estable generado al insertarlo, pero este ejecutable solo lee desde la posición actual, imprime un evento y termina. No guarda checkpoint GTID, no publica en un broker y no tiene consumidor idempotente. Por tanto, todavía no ofrece recuperación tras reinicio ni entrega exactamente una vez.

**Frase clave en inglés:** “An event ID is useful, but it does not make delivery idempotent by itself.”

### 3:45–4:15 · Cierre

**En cámara:**

> Lo que quería demostrar aquí es cómo planteo una integración entre un legado y servicios nuevos: definir con precisión la transacción, elegir dónde capturar los cambios y dejar visible qué garantías existen y cuáles faltan. El siguiente paso sería persistir un checkpoint, añadir publicación con confirmación y probar duplicados en una proyección idempotente.
>
> El código y los límites están en el README. Gracias por verlo.

**Frase final en inglés:** “Small runnable slice, explicit guarantees, clear next steps.”

## Versión de 60 segundos

> ¿Cómo conectamos un MySQL legado con nuevos consumidores sin añadir una doble escritura a un broker? Este laboratorio muestra una pieza ejecutable: al insertar un movimiento, un trigger actualiza el saldo y escribe una fila de outbox en la misma transacción. Si la operación falla, MySQL revierte todo.
>
> Luego, un proceso en Go lee la nueva fila del row binlog y la convierte a CloudEvent JSON. Para probarlo: `make up && make demo`. El evento aparece en el terminal con su ID, tipo, SKU y datos del movimiento.
>
> La decisión reduce polling, pero añade dependencia del binlog y permisos de replicación. Y el límite es importante: esta demo no guarda checkpoints, no publica en un broker y no implementa idempotencia del consumidor. Es una base concreta para discutir fiabilidad y los próximos pasos, no una promesa de CDC lista para producción.

**Frase clave en inglés:** “The source write is atomic; durable delivery is still future work.”
