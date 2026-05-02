# Golang_PricingEngine

## Qué hace este proyecto

Es un **servicio de motor de precios** para escenarios de envío de dinero (transferencias internacionales u online). Recibe por HTTP el contexto de una cotización — país y moneda de origen, país y moneda de destino, partner comercial, opcionalmente estado de origen, estilo de precio (por ejemplo menor comisión, entrega rápida, mejor tasa), forma de recepción del beneficiario, pagador, montos y tipo de forma de pago del envío — y devuelve:

- **Comisiones** desglosadas y disponibilidad por método de pago del remitente (ACH, tarjeta débito, tarjeta crédito).
- **Tipo de cambio** aplicable (tasa, base, monto destino estimado) cuando hay monto origen.
- **Comisión y tipo de cambio juntos** en una sola respuesta cuando aplica.
- **Opciones de estilo de precio y métodos de entrega** válidos para la ruta origen → destino (según reglas en base de datos).
- **Monto en origen** calculado a partir de un monto destino deseado (con validación frente a bandas).
- **Tipos de cambio por pagador** en destino, incluyendo soporte para ciudad/estado destino cuando se usa la variante avanzada del procedimiento.

La autoridad numérica y de reglas complejas vive en **SQL Server**, mediante **procedimientos almacenados** en el esquema `dbo`. Esta aplicación **no** incorpora un ORM ni lee tablas de pricing directamente para esos flujos: actúa como **adaptador HTTP → SP → respuesta JSON**, aplicando después del resultado del SP un conjunto acotado de reglas (filtros por estado de origen y código `ALL`, exclusión de ciertas filas por tipo de forma de pago, selección de la mejor tasa disponible, redondeos y ensamblado de DTOs).

Dos operaciones del dominio expuestas como rutas históricas (`GetBestRatesByCountry`, `GetCountriesAndCurrencies`) **no tienen implementación** en este código: el servicio responde con **501** hasta que exista una capa de datos adicional que las cubra.

---

## Arquitectura

La aplicación sigue una separación por **capas internas** (`internal/`), con dependencias hacia dentro (el dominio no depende del framework HTTP).

```text
                    ┌─────────────────┐
                    │   HTTP (Chi)    │  internal/httpserver
                    │ rutas + query   │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │    validate     │  reglas de entrada por operación
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │ pricing.Service │  orquestación y reglas posteriores al SP
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   repo (SQL)    │  EXEC de SP, filas → mapas tipados luego
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   SQL Server    │
                    └─────────────────┘

     internal/config  ───► lectura de entorno en arranque (cmd/server)
     internal/models  ───► estructuras de petición/respuesta JSON
```

- **`cmd/server`**: composición del programa — carga configuración, abre el pool de conexiones, comprueba conectividad (opcional, con log si falla), registra el router y gestiona el apagado ante señales del sistema operativo.
- **`internal/httpserver`**: superficie REST mínima (`GET`), parseo de query string compatible con nombres tipo `PartnerId` o `partnerId`, y traducción de errores de dominio a códigos HTTP (200 / 204 / 400 / 501).
- **`internal/validate`**: validaciones declarativas por operación (campos obligatorios, rangos de ids), antes de tocar la base.
- **`internal/pricing`**: **casos de uso** del motor — invoca al repositorio con los parámetros correctos, elige variante V1 o V2 del SP según configuración por partner, aplica filtros y construye los modelos de salida (incluidos métodos de pago y agrupación por pagador donde corresponde).
- **`internal/repo`**: única frontera con SQL Server (`database/sql` + driver oficial Microsoft); cada método público corresponde a un SP concreto; los resultados se normalizan desde columnas con distinto casing.
- **`internal/config`**: políticas operativas sin código — por ejemplo qué partners usan los SP `_v2`, cómo se elige la opción de tasa base (`BASE_RATE_OPTION` o listas país/moneda para simulador o tesorería), y flags para desactivar rutas que solo tienen implementación vía SP (`DELIVERY_METHODS_USE_SP`, etc.).
- **`internal/models`**: contratos JSON (principalmente **camelCase** en respuestas, convención estándar de `encoding/json` en Go).

No hay inyección de dependencias pesada: el **árbol de objetos** se construye en `main` (config → repo → servicio → servidor).

---

## API HTTP

Escucha en `HTTP_ADDR` (por defecto `:8080`).

| Ruta | Función |
|------|---------|
| `GET /health` | Comprobación de proceso vivo. |
| `GET /PricingEngine/GetDeliveryMethodAndStyleOptions` | Estilos y métodos de entrega disponibles. |
| `GET /PricingEngine/GetFees` | Comisiones y métodos de pago. |
| `GET /PricingEngine/GetExchangeRates` | Mejor tipo de cambio según criterios del servicio. |
| `GET /PricingEngine/GetFeesExchangeRates` | Comisión + tipo de cambio combinados. |
| `GET /PricingEngine/GetOriginAmount` | Monto origen desde monto destino. |
| `GET /PricingEngine/GetExchangeRatesByPayer` | Tasas por pagador (con ramas V1/V2 según config). |
| `GET /PricingEngine/GetBestRatesByCountry` | Sin implementación → **501**. |
| `GET /PricingEngine/GetCountriesAndCurrencies` | Sin implementación → **501**. |

**200** con JSON, **204** si no hay resultado de negocio, **400** con `{ "error": "..." }`, **501** cuando la funcionalidad no está construida en este binario.

Parámetros habituales (query): origen y destino geográficos y monetarios, `PartnerId`, `StyleId`, `ReceivingOptionId`, `PayerId`, `OriginAmount`, `DestinationAmount`, `FormsOfPaymentsTypeId`, `IsEmployee`; para cotización por pagador avanzada también `DestinationState`, `DestinationCity`.

Este repositorio **no incluye middleware de autenticación**; en despliegues reales suele colocarse detrás de un API gateway o reverse proxy que aplique la política de identidad.

---

## Configuración (variables de entorno)

| Variable | Uso |
|----------|-----|
| `SQLSERVER_CONNECTION_STRING` | Obligatoria. Cadena del driver `go-mssqldb`. |
| `HTTP_ADDR` | Dirección de escucha (`:8080` por defecto). |
| `PRICING_USE_FEES_V2` | Si no hay lista de partners V2, activa los SP `_v2` para todos. |
| `PRICING_V2_PARTNER_IDS` | Lista `1,2,3`: solo esos partners usan `_v2`; el resto usa los SP sin sufijo. |
| `BASE_RATE_OPTION` | Forzar opción de tasa base en SP V1 (`1` / `2` / `3`); `0` usa listas o valor por defecto. |
| `PRICING_SIMULATOR_PAIRS` / `PRICING_TREASURY_PAIRS` | Pares `País|MON` separados por `;` para elegir opción de base 1 o 2 según destino. |
| `DELIVERY_METHODS_USE_SP` | Si es `false`, la ruta de métodos de entrega responde **501**. |
| `EXCHANGE_RATE_BY_PAYER_USE_SP` | Si es `false`, tasas por pagador → **501**. |
| `ORIGIN_AMOUNT_USE_SP` | Si es `false`, monto origen → **501**. |

Booleanos: `true`, `false`, `1`, `yes` donde aplique.

---

## Ejecución

```powershell
cd D:\Christian\Repo\Golang_PricingEngine
$env:SQLSERVER_CONNECTION_STRING = "sqlserver://usuario:clave@servidor:1433?database=NombreBD&encrypt=disable"
go run ./cmd/server
```

```bash
go build -o bin/pricingengine ./cmd/server
```

---

## Dependencias

Go **1.22+**. Librerías: **Chi** (HTTP), **go-mssqldb** (SQL Server), **shopspring/decimal** (redondeos en el dominio).

---

## Procedimientos almacenados utilizados

- `[dbo].[DeliveryMethodsAndStyleOptions]`
- `[dbo].[FeesAndExchangesByFormOfPayment]` / `[dbo].[FeesAndExchangesByFormOfPayment_v2]`
- `[dbo].[OriAmountCalculator]` / `[dbo].[OriAmountCalculator_v2]`
- `[dbo].[ExchangesRatesByPayers]` / `[dbo].[ExchangesRatesByPayers_v2]`

Los parámetros y orden de llamada deben coincidir con la definición publicada en la base de datos que aloja estos objetos.

---

## Nota sobre precisión

Los valores numéricos que vienen del motor SQL se interpretan en Go principalmente como **`float64`**. Para auditoría financiera estricta puede valorarse tipos decimales también en la capa de lectura de filas.
