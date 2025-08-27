---

## Estado Actual y Logros Recientes

**Objetivo General:** Refactorizar la aplicación para alinearla con los patrones del proyecto de referencia `hatmax/blog`, mejorar la consistencia del código y solucionar bugs críticos.

### Logros de la Sesión:

*   **Solucionado: Error Crítico en el Arranque (Seeder Bug)**
    *   **Síntoma:** La aplicación no arrancaba debido a un error `could not find name shortID` durante la fase de seeding de la base de datos.
    *   **Causa Raíz:** Se descubrió una inconsistencia en el repositorio (`internal/repo/sqlite/ssg.go`). El método `CreateLayout` usaba `NamedExecContext` de `sqlx`, que requiere que los campos del struct sean exportados (públicos). Sin embargo, los campos comunes y de auditoría (`shortID`, `createdBy`, etc.) eran privados.
    *   **Solución:** Se estandarizó el método `CreateLayout` para que use `ExecContext`, pasando los parámetros por orden, al igual que los otros métodos `Create` del repositorio. Esto eliminó la dependencia de los nombres de los campos del struct y solucionó el error de arranque.

*   **Refactorización Masiva de Modelos (`auth` y `ssg`)**
    *   Para soportar la solución anterior y mejorar la consistencia, se refactorizaron **todas** las entidades de los paquetes `auth` y `ssg`.
    *   **Patrón Aplicado:**
        1.  Los campos comunes (`ShortID`) y de auditoría (`CreatedBy`, `CreatedAt`, etc.) ahora son **exportados** (empiezan con mayúscula) para ser accesibles por `sqlx`.
        2.  Se asignó el tag `db:"<column_name>"` a todos estos campos para un mapeo correcto con la base de datos.
        3.  Se asignó el tag `json:"-"` a todos estos campos para asegurar que **no** sean expuestos en las respuestas de la API, cumpliendo con los requisitos de visibilidad.
    *   Se adoptó la convención de usar `Kind` (en lugar de `Type` o `ResourceType`) para el campo que define el tipo de un `Resource`, basándose en el proyecto de referencia `hatmax/blog`.

*   **Solucionado: `POST` en BFF para Crear Contenido**
    *   **Síntoma:** El formulario para crear un nuevo contenido fallaba al guardar.
    *   **Causa Raíz:** El `APIHandler` (`internal/feat/ssg/apihandler.go`) recibía correctamente los datos del BFF, pero los descartaba y creaba un nuevo objeto `Content` solo con el `Heading` y el `Body`, perdiendo `UserID`, `SectionID`, etc.
    *   **Solución:** Se modificó el `APIHandler` para que utilice directamente el objeto `Content` decodificado del cuerpo de la petición, asegurando que todos los datos se pasen correctamente a la capa de servicio.

### Nueva Arquitectura BFF (Backend-for-Frontend)

*   Se ha introducido una nueva struct `BFF` en `internal/feat/ssg/bff.go` que reemplaza gradualmente la funcionalidad de `WebHandler`.
*   **Diferencia Clave:** Mientras que `WebHandler` opera directamente contra la capa de servicio (`Service`), el `BFF` actúa como un cliente HTTP para el `APIHandler` del mismo feature.
*   **Motivación:** Este desacoplamiento permite que en el futuro, el BFF pueda ser extraído a su propio microservicio escalable e independiente.

### Próximos Pasos:

*   Continuar la migración de la lógica de `WebHandler` a la nueva arquitectura `BFF`.
*   Verificar el funcionamiento de las operaciones de `Update` y `Delete` a través del flujo BFF -> API.

---