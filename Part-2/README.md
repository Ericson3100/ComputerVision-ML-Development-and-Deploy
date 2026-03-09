# 🖼️ Clasificador de Imágenes - Multiclass vs Mixture of Experts

Sistema completo de clasificación de imágenes con dos arquitecturas de modelos: un clasificador multiclase estándar y un modelo Mixture of Experts (MoE). Todo dockerizado para despliegue fácil.

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────┐
│                     Docker Compose                       │
├──────────────────────┬──────────────────────────────────┤
│   Frontend (Vue.js)  │     Backend (Go)                 │
│   - Puerto: 3000     │     - Puerto: 8080               │
│   - Nginx            │     - API REST                   │
│                      │     - 3 Modelos ML               │
└──────────────────────┴──────────────────────────────────┘
```

### Componentes

- **Frontend**: Vue.js 3 con interfaz moderna para subir imágenes y visualizar resultados
- **Backend**: Go con Gorilla Mux para API REST de clasificación
- **Modelos**:
  - **Multiclass**: Clasificador estándar con 10 clases (animales y vehículos)
  - **Mixture of Experts**: Sistema MoE con expertos especializados y router inteligente
  - **Binary Relevance**: Clasificadores binarios independientes para multi-label

## 🚀 Inicio Rápido

### Prerrequisitos

- Docker Desktop instalado
- Docker Compose v3.8+
- Puertos 3000 y 8080 disponibles

### Instalación y Ejecución

```bash
# 1. Clonar el repositorio
git clone <tu-repo>
cd Part-2

# 2. Construir y ejecutar con Docker Compose
docker-compose up --build

# O en modo detached
docker-compose up -d --build
```

### Acceso a la Aplicación

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/api/health

## 📋 Uso de la Aplicación

1. **Seleccionar Modelo**: Elige entre "Modelo Multiclass", "Mixture of Experts" o "Binary Relevance"
2. **Subir Imagen**: Arrastra o selecciona una imagen (JPG, PNG, GIF hasta 10MB)
3. **Clasificar**: Haz clic en "Clasificar Imagen"
4. **Ver Resultados**: Observa las predicciones ordenadas por confianza

### Ejemplo de Respuesta

```json
{
  "model": "mixture_of_experts",
  "predictions": [
    {
      "class": "cat",
      "confidence": 0.85,
      "expert_used": "animal_expert"
    },
    {
      "class": "dog",
      "confidence": 0.10,
      "expert_used": "animal_expert"
    }
  ],
  "process_time_ms": 45.2
}
```

## 🔧 Desarrollo

### Modo Desarrollo con Hot-Reload

```bash
# Usar docker-compose de desarrollo
docker-compose -f docker-compose.dev.yml up --build
```

En este modo:
- Frontend: Hot-reload automático en cambios de código
- Backend: Recompilación automática en cambios de Go

### Ejecutar Solo Backend

```bash
cd backend
go mod download
go run .
```

### Ejecutar Solo Frontend

```bash
cd frontend
npm install
npm run serve
```

## 📁 Estructura del Proyecto

```
Part-2/
├── backend/                    # Backend en Go
│   ├── main.go                # Servidor principal y API
│   ├── multiclass_model.go    # Modelo clasificador multiclase
│   ├── mixture_of_experts_model.go  # Modelo MoE
│   ├── go.mod                 # Dependencias Go
│   ├── Dockerfile             # Imagen Docker backend
│   └── .dockerignore
│
├── frontend/                   # Frontend en Vue.js
│   ├── src/
│   │   ├── App.vue           # Componente principal
│   │   └── main.js           # Punto de entrada
│   ├── public/
│   │   └── index.html        # HTML base
│   ├── package.json          # Dependencias npm
│   ├── vue.config.js         # Configuración Vue
│   ├── nginx.conf            # Configuración Nginx
│   ├── Dockerfile            # Imagen Docker frontend (producción)
│   ├── Dockerfile.dev        # Imagen Docker frontend (desarrollo)
│   └── .dockerignore
│
├── docker-compose.yml         # Orquestación producción
├── docker-compose.dev.yml     # Orquestación desarrollo
├── .gitignore
└── README.md
```

## 🎯 API Endpoints

### GET /api/health
Health check del servidor

```bash
curl http://localhost:8080/api/health
```

### GET /api/models
Obtener información de modelos disponibles

```bash
curl http://localhost:8080/api/models
```

### POST /api/classify
Clasificar una imagen

```bash
curl -X POST \
  -F "image=@imagen.jpg" \
  -F "model=multiclass" \
  http://localhost:8080/api/classify
```

**Parámetros**:
- `image`: Archivo de imagen (multipart/form-data)
- `model`: "multiclass", "mixture_of_experts" o "binary_relevance"

## 🧠 Modelos de Machine Learning

### Modelo Multiclass

Arquitectura tradicional donde un solo modelo maneja todas las clases:

**Clases**: cat, dog, bird, fish, horse, car, truck, airplane, boat, bicycle

**Ventajas**:
- Simple y directo
- Menos recursos computacionales
- Bueno para problemas balanceados

### Mixture of Experts (MoE)

Sistema con múltiples expertos especializados y un router que decide qué experto usar:

**Expertos**:
- `animal_expert`: cat, dog, bird, fish, horse
- `vehicle_expert`: car, truck, airplane, boat, bicycle

**Componentes**:
1. **Router**: Analiza la imagen y selecciona el experto más apropiado
2. **Expertos**: Modelos especializados en subconjuntos de clases

**Ventajas**:
- Mayor precisión en clases específicas
- Escalable a más expertos
- Mejor para problemas complejos

### Binary Relevance

Sistema de clasificación multi-label con clasificadores binarios independientes:

**Clasificadores**:
- Un clasificador binario por cada clase (10 en total)
- Cada uno decide independientemente si su clase está presente

**Características**:
1. **Independencia**: Cada clasificador funciona de forma autónoma
2. **Multi-label**: Puede predecir múltiples clases simultáneamente
3. **Umbral ajustable**: Cada clasificador tiene su propio umbral de decisión

**Ventajas**:
- Soporte para multi-label (múltiples clases por imagen)
- Fácil de entrenar (cada clasificador es simple)
- Flexible en ajuste de umbrales por clase
- Bueno para datasets desbalanceados

**Casos de uso**:
- Imágenes con múltiples objetos
- Clasificación de tags/etiquetas
- Detección de atributos múltiples

## 🔄 Comandos Docker Útiles

```bash
# Ver logs
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f backend
docker-compose logs -f frontend

# Parar servicios
docker-compose down

# Parar y eliminar volúmenes
docker-compose down -v

# Reconstruir sin caché
docker-compose build --no-cache

# Ver contenedores corriendo
docker-compose ps

# Ejecutar comando en contenedor
docker-compose exec backend sh
docker-compose exec frontend sh
```

## 🧪 Testing

### Backend

```bash
cd backend
go test ./...
```

### Frontend

```bash
cd frontend
npm run test
```

## 📊 Monitoreo y Logs

Los logs de ambos servicios se pueden ver con:

```bash
# Todos los logs
docker-compose logs -f

# Solo backend
docker-compose logs -f backend

# Solo frontend
docker-compose logs -f frontend
```

## 🚢 Despliegue en Producción

### Build de Producción

```bash
docker-compose build
docker-compose up -d
```

### Consideraciones de Producción

1. **Modelos Reales**: Integrar modelos ML reales (TensorFlow, PyTorch, ONNX)
2. **Persistencia**: Agregar volúmenes para caché de modelos
3. **Escalabilidad**: Usar Kubernetes para auto-scaling
4. **Seguridad**: 
   - Rate limiting
   - Autenticación JWT
   - HTTPS con certificados SSL
5. **Monitoreo**: Prometheus + Grafana para métricas

## 🔐 Variables de Entorno

### Backend
```env
PORT=8080
GO_ENV=production
```

### Frontend
```env
VUE_APP_API_URL=http://localhost:8080
```

## 🐛 Troubleshooting

### El frontend no se conecta al backend

Verifica que el backend esté corriendo:
```bash
curl http://localhost:8080/api/health
```

### Error de puertos ocupados

Cambia los puertos en `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Backend
  - "3001:80"    # Frontend
```

### Reconstruir desde cero

```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

## 🤝 Contribuciones

1. Fork el proyecto
2. Crea una rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Próximas Mejoras

- [ ] Integrar modelos ML reales (TensorFlow/PyTorch)
- [ ] Agregar autenticación de usuarios
- [ ] Implementar caché de predicciones
- [ ] Agregar más tipos de expertos en MoE
- [ ] Dashboard de métricas en tiempo real
- [ ] Soporte para batch processing
- [ ] API GraphQL adicional
- [ ] Tests unitarios y de integración completos

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.

## 👥 Autores

- Tu Nombre - Desarrollador Principal

## 🙏 Agradecimientos

- Vue.js por el excelente framework frontend
- Go community por las librerías de calidad
- Docker por facilitar el despliegue

---

**¿Preguntas o problemas?** Abre un issue en GitHub.
