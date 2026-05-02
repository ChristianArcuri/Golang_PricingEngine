# Golang_PricingEngine

Servicio HTTP en Go que replica las rutas públicas del **`PricingEngineController`** del API .NET (ASP.NET Core / EF Core). La lógica operativa se apoya en los **mismos stored procedures** `[dbo]` que consumía el contexto EF (`DeliveryMethodsAndStyleOptions`, `FeesAndExchangesByFormOfPayment`, versiones `_v2`, `OriAmountCalculator`, `ExchangesRatesByPayers`, etc.).

No incluye: **Entity Framework**, **LaunchDarkly**, **JWT / `[Authorize]`**, ni consultas a bases auxiliares (**WireTransac**, **WirePricing**) para tasas base o definiciones de precio.

---

## Estructura del repositorio

| Ruta | Rol |
|------|-----|
| `cmd/server/main.go` | Punto de entrada: carga configuración, abre SQL Server, ping opcional, servidor HTTP y apagado con señales. |
| `internal/config` | Variables de entorno (equivalentes simplificados a flags y listas del appsettings .NET). |
| `internal/repo` | `database/sql` + `go-mssqldb`: ejecución de SP y lectura de resultados en `map[string]interface{}`. |
| `internal/pricing` | Reglas de negocio alineadas al servicio .NET (filtros por `OriginStateCode` / `ALL`, exclusión forma de pago 4, métodos de pago ACH/Debit/Credit, etc.). |
| `internal/validate` | Validaciones por endpoint (equivalentes reducidas a las RuleSets de FluentValidation del modelo original). |
| `internal/models` | DTOs de entrada/salida y etiquetas JSON. |
| `internal/httpserver` | Chi: rutas bajo `/PricingEngine`, parsing de query string. |

---

## Dependencias (Go 1.22+)

- [`github.com/go-chi/chi/v5`](https://github.com/go-chi/chi) — enrutador HTTP y middleware mínimo.
- [`github.com/microsoft/go-mssqldb`](https://github.com/microsoft/go-mssqldb) — driver SQL Server.
- [`github.com/shopspring/decimal`](https://github.com/shopspring/decimal) — redondeos de importes.

---

## API HTTP

Base: `HTTP_ADDR` (por defecto `:8080`).

| Método | Ruta | Comportamiento |
|--------|------|----------------|
| GET | `/health` | `{"status":"ok"}` |
| GET | `/PricingEngine/GetDeliveryMethodAndStyleOptions` | SP `DeliveryMethodsAndStyleOptions` si `DELIVERY_METHODS_USE_SP=true`; si `false`, **501**. |
| GET | `/PricingEngine/GetFees` | Fees desde SP V1/V2 según configuración del partner. |
| GET | `/PricingEngine/GetExchangeRates` | Misma fuente de datos que fees; selección de mejor tasa disponible. |
| GET | `/PricingEngine/GetFeesExchangeRates` | Combina fee + tipo de cambio (parámetro `FormsOfPaymentsTypeId` al SP según V1/V2, coherente con el flujo .NET). |
| GET | `/PricingEngine/GetOriginAmount` | `OriAmountCalculator` / `_v2`; si `ORIGIN_AMOUNT_USE_SP=false`, **501**. |
| GET | `/PricingEngine/GetExchangeRatesByPayer` | `ExchangesRatesByPayers` / `_v2`; si `EXCHANGE_RATE_BY_PAYER_USE_SP=false`, **501**. |
| GET | `/PricingEngine/GetBestRatesByCountry` | Validación del modelo; respuesta **501** (requiere EF / tablas de precios del .NET). |
| GET | `/PricingEngine/GetCountriesAndCurrencies` | **501** (requiere `PriceDefinitionService` / EF). |

### Códigos de respuesta

- **200** — Cuerpo JSON con el resultado.
- **204** — Sin contenido (equivalente a `NoContent` del controlador .NET cuando el servicio devuelve `null`).
- **400** — Error de validación u otro error de dominio (mensaje en JSON `{ "error": "..." }`).
- **501** — Endpoint o rama no migrada (`ErrNotImplemented`).

### Query string

Se aceptan nombres en **PascalCase** (como el binding .NET) o **camelCase**, por ejemplo `PartnerId` / `partnerId`, `OriginAmount` / `originAmount`.

Parámetros habituales en `PriceEngineModel`: `OriginCountryName`, `OriginCountryCode`, `OriginStateCode`, `OriginCurrencyCode`, `DestinationCountryName`, `DestinationCountryCode`, `DestinationCurrencyCode`, `PartnerId`, `StyleId`, `ReceivingOptionId`, `PayerId`, `OriginAmount`, `DestinationAmount`, `FormsOfPaymentsTypeId`, `IsEmployee`, y para `_v2` por pagador: `DestinationState`, `DestinationCity`.

### Formato JSON

Las respuestas usan **camelCase** (`encoding/json` estándar), que no coincide necesariamente con Newtonsoft usando propiedades en PascalCase en el API original.

### Autenticación

El API .NET marcaba otros controladores con `[Authorize]`; este servicio **no implementa JWT ni middleware de auth** salvo que lo añadas delante del router.

---

## Configuración por entorno

| Variable | Obligatoria | Default | Descripción |
|----------|-------------|---------|-------------|
| `SQLSERVER_CONNECTION_STRING` | Sí | — | Cadena para `go-mssqldb`. Ejemplo: `sqlserver://user:pass@host:1433?database=PricingEngine&encrypt=disable` |
| `HTTP_ADDR` | No | `:8080` | `ListenAndServe` |
| `PRICING_USE_FEES_V2` | No | `false` | Si **no** defines `PRICING_V2_PARTNER_IDS`, indica si **todos** los partners usan los SP `_v2` para fees/origen/por pagador. |
| `PRICING_V2_PARTNER_IDS` | No | vacío | Lista `1,2,3`. Si está definida, **solo** esos partners usan rutas `_v2`; el resto usa siempre V1 (`FeesAndExchangesByFormOfPayment`, etc.). |
| `BASE_RATE_OPTION` | No | `0` | Forzar opción de tasa base para SP V1: `1` simulador, `2` tesorería, `3` diario. `0` = derivar de listas o `3`. |
| `PRICING_SIMULATOR_PAIRS` | No | vacío | Pares `País|MON` separados por `;` → opción base **1** si coincide país y moneda destino. |
| `PRICING_TREASURY_PAIRS` | No | vacío | Igual → opción base **2**. |
| `DELIVERY_METHODS_USE_SP` | No | `true` | `false` desactiva la rama SP de métodos de entrega → **501**. |
| `EXCHANGE_RATE_BY_PAYER_USE_SP` | No | `true` | `false` → **501** en `GetExchangeRatesByPayer`. |
| `ORIGIN_AMOUNT_USE_SP` | No | `true` | `false` → **501** en `GetOriginAmount`. |

Valores booleanos aceptados: `true` / `false` / `1` / `yes` (solo para las variables `*_USE_SP` y `PRICING_USE_FEES_V2`).

---

## Limitaciones respecto al API .NET

1. **Feature flags (LaunchDarkly)** — Sustituidos por variables de entorno anteriores; no hay evaluación remota de flags.
2. **`GetExchangeRatesByPayer` sin V2** — Si el partner **no** está en modo V2 y falta `OriginAmount` pero hay `DestinationAmount`, el .NET podía resolver origen vía otras rutas EF; aquí se devuelve **501** (`ErrNotImplemented`).
3. **`GetDeliveryMethodAndStyleOptions`** — Con `DELIVERY_METHODS_USE_SP=false` el .NET usaba `StyleService` / `ReceivingOptionsTypeService` (EF); no está portado → **501**.
4. **`GetOriginAmount`** — Con `ORIGIN_AMOUNT_USE_SP=false` el .NET usaba repositorios EF y proyecciones de treasury/simulator/daily; no está portado → **501**.
5. **`GetBestRatesByCountry` / `GetCountriesAndCurrencies`** — Dependen de datos y consultas no incluidas en este repo → **501** tras validación cuando aplique.
6. **Precisión** — Se usa `float64` al leer columnas SQL; el .NET usa `decimal`; para montos críticos conviene revisar tolerancias o migrar a tipos decimales en lectura.

---

## Ejecución y compilación

PowerShell (Windows):

```powershell
cd D:\Christian\Repo\Golang_PricingEngine
$env:SQLSERVER_CONNECTION_STRING = "sqlserver://usuario:clave@localhost:1433?database=TuBase&encrypt=disable"
go mod tidy
go run ./cmd/server
```

Compilar binario:

```powershell
go build -o bin/pricingengine.exe ./cmd/server
```

Bash:

```bash
export SQLSERVER_CONNECTION_STRING='sqlserver://...'
go run ./cmd/server
```

Si el ping inicial a la base falla, el proceso **sigue arrancando** y solo se registra un warning en log.

---

## Referencia de stored procedures usados

- `[dbo].[DeliveryMethodsAndStyleOptions]`
- `[dbo].[FeesAndExchangesByFormOfPayment]` / `[dbo].[FeesAndExchangesByFormOfPayment_v2]`
- `[dbo].[OriAmountCalculator]` / `[dbo].[OriAmountCalculator_v2]`
- `[dbo].[ExchangesRatesByPayers]` / `[dbo].[ExchangesRatesByPayers_v2]`

Los nombres y orden de parámetros siguen el código generado por EF Core Power Tools en el proyecto .NET original (`PricingEngineSPContextProcedures`).
