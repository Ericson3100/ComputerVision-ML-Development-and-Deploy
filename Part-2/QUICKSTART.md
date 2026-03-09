# Clasificador de Imágenes - Instrucciones Rápidas

## 🎯 3 Modelos de Clasificación

1. **🎯 Multiclass**: Clasificador estándar (una clase por imagen)
2. **🧠 Mixture of Experts**: Expertos especializados con router
3. **🔢 Binary Relevance**: Multi-label (múltiples clases por imagen)

## 🚀 Inicio Rápido con Docker (Recomendado)

### Windows
```cmd
start.bat
```

### Linux/Mac
```bash
chmod +x start.sh
./start.sh
```

O usando Docker Compose directamente:
```bash
docker-compose up --build
```

Accede a:
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080

## 💻 Modo Desarrollo

### Windows
```cmd
dev.bat
```

### Linux/Mac
```bash
chmod +x dev.sh
./dev.sh
```

## 📝 Con Makefile

```bash
make help          # Ver todos los comandos
make up            # Iniciar aplicación
make logs          # Ver logs
make down          # Detener aplicación
make dev           # Modo desarrollo
```

## 🛠️ Sin Docker

### Backend
```bash
cd backend
go mod download
go run .
```

### Frontend
```bash
cd frontend
npm install
npm run serve
```

## 📚 Documentación Completa

Ver [README.md](README.md) para documentación detallada.
