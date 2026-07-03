# BRE-B PSE Recaudos App 🚀

Aplicación en Go para gestionar recaudos de BRE-B mediante PSE con integración de ePayco, split de pagos y notificaciones multi-canal (Email, SMS, WhatsApp).

## Características ✨

- ✅ Integración con ePayco para pagos PSE
- ✅ Split automático de pagos entre beneficiarios
- ✅ Notificaciones por Email, SMS y WhatsApp
- ✅ Webhooks para confirmar pagos
- ✅ Generación de reportes en Excel
- ✅ Estadísticas y análisis
- ✅ Auditoría de transacciones
- ✅ API REST con autenticación JWT

## Requisitos 📋

- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (opcional)
- Cuentas activas en:
  - ePayco
  - Twilio
  - SendGrid

## Instalación 🔧

### Opción 1: Usando Docker Compose

```bash
# Clonar repositorio
git clone https://github.com/fabianherreracruz/bre-b-pse-app.git
cd bre-b-pse-app

# Crear archivo .env
cp .env.example .env

# Configurar variables de entorno
nano .env

# Levantar servicios
docker-compose up -d
```

### Opción 2: Instalación Local

```bash
# Clonar repositorio
git clone https://github.com/fabianherreracruz/bre-b-pse-app.git
cd bre-b-pse-app

# Instalar dependencias
go mod download

# Crear archivo .env
cp .env.example .env

# Configurar variables de entorno
nano .env

# Ejecutar migraciones (si es necesario)
go run cmd/main.go

# La app estará disponible en http://localhost:8080
```

## Configuración de Variables de Entorno 🔐

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bre_b_pse_db

# ePayco
EPAYCO_CLIENT_ID=tu_cliente_id
EPAYCO_CLIENT_SECRET=tu_cliente_secret
EPAYCO_API_KEY=tu_api_key
EPAYCO_PRIVATE_KEY=tu_private_key

# JWT
JWT_SECRET=tu_clave_super_secreta

# Twilio
TWILIO_ACCOUNT_SID=tu_account_sid
TWILIO_AUTH_TOKEN=tu_auth_token
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890

# SendGrid
SENDGRID_API_KEY=tu_sendgrid_key
SENDGRID_FROM_EMAIL=noreply@recaudos.com
```

## API Endpoints 📡

### Pagos

```bash
# Crear recaudo
POST /api/v1/payments/create
Content-Type: application/json

{
  "amount": 100000,
  "description": "Pago de servicios",
  "customer_email": "cliente@example.com",
  "customer_phone": "+573001234567",
  "splits": [
    {
      "beneficiario_id": "uuid-1",
      "percentage": 60
    },
    {
      "beneficiario_id": "uuid-2",
      "percentage": 40
    }
  ]
}

# Verificar estado de pago
GET /api/v1/payments/verify/:reference

# Obtener estado de recaudo
GET /api/v1/payments/status/:id

# Webhook de ePayco
POST /api/v1/payments/webhook
```

### Reportes (Requieren autenticación)

```bash
# Reporte por rango de fechas
GET /api/v1/reports/by-date?start_date=2024-01-01&end_date=2024-12-31
Authorization: Bearer your_jwt_token

# Exportar reporte a Excel
GET /api/v1/reports/export-excel?start_date=2024-01-01&end_date=2024-12-31

# Obtener estadísticas
GET /api/v1/reports/statistics

# Reporte mensual
GET /api/v1/reports/monthly

# Reporte anual
GET /api/v1/reports/yearly
```

## Estructura del Proyecto 📁

```
bre-b-pse-app/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── app/
│   │   └── app.go           # Configuración de la app
│   ├── config/
│   │   └── config.go        # Variables de entorno
│   ├── db/
│   │   └── database.go      # Conexión a BD
│   ├── handlers/
│   │   ├── payment_handler.go
│   │   └── report_handler.go
│   ├── middleware/
│   │   └── auth.go          # JWT middleware
│   ├── models/
│   │   └── models.go        # Modelos de BD
│   └── services/
│       ├── epayco_service.go
│       ├── notification_service.go
│       ├── split_service.go
│       └── report_service.go
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Flujo de Recaudo 💳

1. **Crear Recaudo** → Se genera un link de pago en ePayco
2. **Cliente paga** → Realiza transacción por PSE
3. **Webhook ePayco** → Confirma el pago
4. **Procesar Splits** → Distribuye dinero a beneficiarios
5. **Notificaciones** → Envía confirmaciones por Email/SMS/WhatsApp
6. **Generar Reporte** → Registra en auditoría y reportes

## Comandos Útiles 🛠️

```bash
# Compilar aplicación
make build

# Ejecutar en desarrollo
make dev

# Levantar servicios Docker
make docker-up

# Detener servicios Docker
make docker-down

# Ver logs
make docker-logs

# Instalar dependencias
make install-deps
```

## Testing 🧪

```bash
# Ejecutar tests
go test ./...

# Con cobertura
go test -cover ./...
```

## Licencia 📄

Este proyecto está bajo la Licencia MIT.

## Soporte 💬

Para soporte, contacta con:
- Email: fabianherreracruz@gmail.com
- GitHub Issues

---

**Hecho con ❤️ por Fabian Herrera**
