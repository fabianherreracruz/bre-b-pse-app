# Guía de Deployment

## Desplegar en Railway

### 1. Preparar el repositorio
```bash
git push origin main
```

### 2. Conectar Railway
- Ve a railway.app
- Click en "New Project"
- Selecciona "Deploy from GitHub"
- Autoriza y selecciona tu repositorio

### 3. Configurar variables de entorno
En el panel de Railway, agrega:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `EPAYCO_CLIENT_ID`, `EPAYCO_CLIENT_SECRET`, etc.
- `JWT_SECRET`
- `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`
- `SENDGRID_API_KEY`

### 4. Deploy automático
Railway deployará automáticamente con cada push a main

---

## Desplegar en Docker

```bash
# Build
docker build -t bre-b-pse-app:latest .

# Tag para Docker Hub
docker tag bre-b-pse-app:latest username/bre-b-pse-app:latest

# Push
docker push username/bre-b-pse-app:latest

# Run en producción
docker run -d \
  -e DB_HOST=postgres-server \
  -e EPAYCO_CLIENT_ID=xxx \
  -p 8080:8080 \
  username/bre-b-pse-app:latest
```

---

## Desplegar Frontend en Vercel

```bash
# En la carpeta web/
cd web
npm run build
```

- Ve a vercel.com
- Importa el proyecto desde GitHub
- Selecciona la carpeta `web` como root
- Deploy automático

---

## Monitoreo

### Health Check
```bash
curl https://tu-app.com/health
```

### Logs
```bash
# Railway
railway logs

# Docker
docker logs container_id
```

### Métricas
- Usa DataDog, New Relic o Similar
- Monitorea: uptime, latencia, errores
