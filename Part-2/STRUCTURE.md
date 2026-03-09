# 📁 Estructura del Proyecto

## Vista de Árbol

```
Part-2/
├── 📋 README.md                    # Documentación principal completa
├── 📋 QUICKSTART.md                # Guía de inicio rápido
├── 📋 MODELS.md                    # Arquitectura técnica de los modelos
├── 🐳 docker-compose.yml           # Orquestación Docker (producción)
├── 🐳 docker-compose.dev.yml       # Orquestación Docker (desarrollo)
├── ⚙️  Makefile                     # Comandos útiles automatizados
├── 🚀 start.sh                     # Script de inicio (Linux/Mac)
├── 🚀 start.bat                    # Script de inicio (Windows)
├── 🔧 dev.sh                       # Script desarrollo (Linux/Mac)
├── 🔧 dev.bat                      # Script desarrollo (Windows)
├── 📝 .env.example                 # Plantilla de variables de entorno
├── 🚫 .gitignore                   # Archivos a ignorar por Git
├── 🚫 .dockerignore                # Archivos a ignorar por Docker
│
├── 📁 backend/                     # Backend en Go
│   ├── 🔷 main.go                  # Servidor HTTP y API REST
│   ├── 🔷 multiclass_model.go      # Modelo clasificador multiclase
│   ├── 🔷 mixture_of_experts_model.go  # Modelo MoE
│   ├── � binary_relevance_model.go    # Modelo Binary Relevance
│   ├── �📦 go.mod                   # Dependencias Go
│   ├── 📦 go.sum                   # Checksums de dependencias
│   ├── 🐳 Dockerfile               # Imagen Docker del backend
│   └── 🚫 .dockerignore            
│
└── 📁 frontend/                    # Frontend en Vue.js
    ├── 📁 src/
    │   ├── 🎨 App.vue              # Componente principal de Vue
    │   └── 🔷 main.js              # Punto de entrada de Vue
    ├── 📁 public/
    │   └── 📄 index.html           # HTML base
    ├── 📦 package.json             # Dependencias npm
    ├── ⚙️  vue.config.js            # Configuración de Vue CLI
    ├── 🌐 nginx.conf               # Configuración Nginx
    ├── 🐳 Dockerfile               # Imagen Docker (producción)
    ├── 🐳 Dockerfile.dev           # Imagen Docker (desarrollo)
    ├── 📝 .env                     # Variables de entorno
    └── 🚫 .dockerignore
```

## 🗂️ Descripción de Directorios

### Raíz del Proyecto

| Archivo | Descripción |
|---------|-------------|
| `README.md` | Documentación completa con instrucciones de uso |
| `QUICKSTART.md` | Guía rápida para iniciar en minutos |
| `MODELS.md` | Documentación técnica de arquitecturas ML |
| `docker-compose.yml` | Configuración Docker Compose para producción |
| `docker-compose.dev.yml` | Configuración Docker Compose para desarrollo |
| `Makefile` | Comandos automatizados (make up, make logs, etc.) |
| `start.sh/bat` | Scripts para iniciar aplicación fácilmente |
| `dev.sh/bat` | Scripts para modo desarrollo con hot-reload |

### Backend (`/backend`)

```
backend/
├── main.go                      # Servidor HTTP principal
│   ├── Routers y middleware
│   ├── Handlers de API
│   └── Inicialización de modelos
│
├── multiclass_model.go          # Modelo Multiclass
│   ├── Estructura MulticlassModel
│   ├── Método Predict()
│   └── Gestión de clases
│
├── mixture_of_experts_model.go  # Modelo MoE
│   ├── Estructura MixtureOfExpertsModel
│   ├── Sistema de Router
│   ├── Gestión de Expertos
│   └── Método Predict()
│
├── binary_relevance_model.go   # Modelo Binary Relevance
│   ├── Estructura BinaryRelevanceModel
│   ├── Clasificadores binarios independientes
│   ├── Gestión de umbrales
│   └── Métodos Predict() y PredictWithThreshold()
│
└── go.mod/go.sum                # Dependencias Go
```

**Dependencias**:
- `gorilla/mux`: Router HTTP
- `rs/cors`: Manejo de CORS

### Frontend (`/frontend`)

```
frontend/
├── src/
│   ├── App.vue                  # Componente principal
│   │   ├── UI de selección de modelo
│   │   ├── Drag & drop de imágenes
│   │   ├── Visualización de resultados
│   │   └── Gestión de estado
│   │
│   └── main.js                  # Inicialización de Vue
│
├── public/
│   └── index.html              # HTML template
│
├── package.json                # Dependencias npm
├── vue.config.js               # Configuración Vue CLI
└── nginx.conf                  # Configuración servidor web
```

**Dependencias**:
- `vue`: Framework frontend
- `axios`: Cliente HTTP
- `nginx`: Servidor web (producción)

## 🐳 Configuración Docker

### docker-compose.yml (Producción)

```yaml
services:
  backend:
    - Build optimizado multi-stage
    - Puerto 8080
    - Health checks
    
  frontend:
    - Build optimizado con Nginx
    - Puerto 3000
    - Compresión gzip
```

### docker-compose.dev.yml (Desarrollo)

```yaml
services:
  backend:
    - Volúmenes montados
    - Hot-reload automático
    
  frontend:
    - Volúmenes montados
    - Hot-reload automático
    - Dev server de Vue
```

## 📊 Flujo de Datos

```
┌──────────────────────────────────────────────────────────┐
│                     Usuario                              │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│               Frontend (Vue.js - Puerto 3000)            │
│  ┌────────────────────────────────────────────────────┐  │
│  │ 1. Seleccionar Modelo                              │  │
│  │ 2. Upload Imagen                                   │  │
│  │ 3. Enviar a Backend                                │  │
│  │ 4. Mostrar Resultados                              │  │
│  └────────────────────────────────────────────────────┘  │
└────────────────────────┬─────────────────────────────────┘
                         │ HTTP POST /api/classify
                         ▼
┌──────────────────────────────────────────────────────────┐
│               Backend (Go - Puerto 8080)                 │
│  ┌────────────────────────────────────────────────────┐  │
│  │ API REST                                           │  │
│  │  ├─ GET  /api/health                               │  │
│  │  ├─ GET  /api/models                               │  │
│  │  └─ POST /api/classify                             │  │
│  └────────────────────┬───────────────────────────────┘  │
│                       │                                  │
│  ┌────────────────────┴───────────────────────────────┐  │
│  │ Procesamiento                                      │  │
│  │  1. Recibir imagen                                 │  │
│  │  2. Seleccionar modelo                             │  │
│  │  3. Ejecutar inferencia                            │  │
│  │  4. Formatear respuesta                            │  │
│  └────────────────────┬───────────────────────────────┘  │
│                       │                                  │
│  ┌────────────────────┴───────────────────────────────┐  │
│  │ Modelos ML                                         │  │
│  │  ├─ MulticlassModel                                │  │
│  │  └─ MixtureOfExpertsModel                          │  │
│  │      ├─ Router                                     │  │
│  │      ├─ Animal Expert                              │  │
│  │      └─ Vehicle Expert                             │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

## 🔄 Ciclo de Vida del Request

1. **Usuario sube imagen** en frontend
2. **Frontend envía** POST a `/api/classify` con:
   - Archivo de imagen
   - Modelo seleccionado
3. **Backend procesa**:
   - Valida imagen
   - Guarda temporalmente
   - Selecciona modelo (multiclass o MoE)
   - Ejecuta inferencia
4. **Modelo devuelve** predicciones
5. **Backend responde** JSON con:
   - Modelo usado
   - Top-5 predicciones
   - Confianzas
   - Tiempo de procesamiento
6. **Frontend muestra** resultados visuales

## 🎨 Componentes de UI

### App.vue

```
┌─────────────────────────────────────────┐
│           Header                         │
│  Clasificador de Imágenes                │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│       Selector de Modelo                 │
│  ┌──────────────┐  ┌─────────────────┐  │
│  │ 🎯 Multiclass│  │ 🧠 MoE          │  │
│  └──────────────┘  └─────────────────┘  │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│       Área de Upload                     │
│  ┌─────────────────────────────────────┐│
│  │  📤 Drag & Drop                     ││
│  │  o Click para Seleccionar           ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│    🔍 Botón Clasificar                   │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│       Resultados                         │
│  1. cat         ████████████ 85%        │
│  2. dog         ███          10%        │
│  3. bird        ██           5%         │
└─────────────────────────────────────────┘
```

## 🏃 Próximos Pasos para Extender

1. **Agregar nuevos modelos**: Crear nuevo archivo `new_model.go`
2. **Agregar clases**: Modificar arrays de clases en modelos
3. **Integrar ML real**: Ver `MODELS.md` para opciones
4. **Agregar autenticación**: Middleware en `main.go`
5. **Agregar base de datos**: Para almacenar historial
6. **Agregar tests**: Crear `*_test.go` files

## 📖 Recursos Útiles

- [Documentación Vue.js](https://vuejs.org/)
- [Documentación Go](https://golang.org/doc/)
- [Documentación Docker](https://docs.docker.com/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Axios](https://axios-http.com/)
