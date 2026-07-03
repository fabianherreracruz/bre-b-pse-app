# Documentación de API

## Base URL
```
http://localhost:8080/api/v1
```

## Autenticación
Todos los endpoints de perfil y reportes requieren JWT:
```
Authorization: Bearer <token>
```

---

## Auth Endpoints

### POST /auth/register
Registra un nuevo usuario

**Request:**
```json
{
  "first_name": "Juan",
  "last_name": "Pérez",
  "email": "juan@example.com",
  "phone": "+573001234567",
  "password": "password123",
  "role": "recaudador"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "juan@example.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": "recaudador",
  "token": "jwt_token",
  "expires_at": 1234567890
}
```

---

### POST /auth/login
Inicia sesión

**Request:**
```json
{
  "email": "juan@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "juan@example.com",
  "token": "jwt_token",
  "expires_at": 1234567890
}
```

---

## Profile Endpoints

### GET /profile
Obtiene el perfil del usuario actual

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": "uuid",
  "email": "juan@example.com",
  "first_name": "Juan",
  "phone": "+573001234567",
  "role": "recaudador",
  "status": "active"
}
```

---

### PUT /profile
Actualiza el perfil

**Request:**
```json
{
  "first_name": "Juan",
  "last_name": "Pérez",
  "phone": "+573001234567"
}
```

---

### POST /profile/change-password
Cambia la contraseña

**Request:**
```json
{
  "old_password": "password123",
  "new_password": "newpassword123"
}
```

---

## Payment Endpoints

### POST /payments/create
Crea un recaudo

**Request:**
```json
{
  "amount": 100000,
  "description": "Pago de servicios",
  "customer_email": "cliente@example.com",
  "customer_phone": "+573001234567",
  "splits": [
    {
      "beneficiario_id": "uuid",
      "percentage": 60
    }
  ]
}
```

---

### GET /payments/status/:id
Obtiene el estado de un recaudo

---

### GET /payments/verify/:reference
Verifica un pago

---

## Report Endpoints

### GET /reports/by-date
Reporte por rango de fechas

**Query Params:**
- `start_date` (YYYY-MM-DD)
- `end_date` (YYYY-MM-DD)

---

### GET /reports/statistics
Estadísticas del usuario

---

### GET /reports/monthly
Reporte mensual

---

### GET /reports/yearly
Reporte anual

---

### GET /reports/export-excel
Exporta reporte a Excel

**Query Params:**
- `start_date`
- `end_date`

**Response:** Archivo .xlsx
