# 🧠 Arquitectura de Modelos ML

## Visión General

Este proyecto implementa dos enfoques de clasificación de imágenes:

1. **Modelo Multiclass**: Clasificador tradicional
2. **Mixture of Experts (MoE)**: Sistema con múltiples expertos especializados

## 🎯 Modelo Multiclass

### Concepto

Un único modelo neuronal que aprende a clasificar todas las clases simultáneamente.

```
┌─────────────────────────────────────────┐
│          Entrada (Imagen)               │
└──────────────┬──────────────────────────┘
               │
               ▼
      ┌────────────────┐
      │  Feature       │
      │  Extractor     │
      │  (CNN)         │
      └────────┬───────┘
               │
               ▼
      ┌────────────────┐
      │  Classifier    │
      │  (Dense)       │
      └────────┬───────┘
               │
               ▼
      ┌────────────────┐
      │  Softmax       │
      │  10 clases     │
      └────────────────┘
```

### Clases

```
Animales:    cat, dog, bird, fish, horse
Vehículos:   car, truck, airplane, boat, bicycle
```

### Características

- **Ventajas**:
  - Arquitectura simple
  - Menos parámetros
  - Rápido en inferencia
  - Fácil de entrenar

- **Desventajas**:
  - Puede tener dificultades con clases similares
  - No aprovecha especialización
  - Menos flexible para agregar clases

### Implementación Actual

```go
type MulticlassModel struct {
    classes []string  // Todas las clases
    loaded  bool      // Estado del modelo
}

func (m *MulticlassModel) Predict(imagePath string) ([]PredictionResult, error) {
    // 1. Cargar y preprocesar imagen
    // 2. Pasar por red neuronal
    // 3. Aplicar softmax
    // 4. Retornar top-k predicciones
}
```

## 🧠 Mixture of Experts (MoE)

### Concepto

Sistema modular con múltiples "expertos" especializados y un "router" que decide qué experto usar.

```
┌─────────────────────────────────────────┐
│          Entrada (Imagen)               │
└──────────────┬──────────────────────────┘
               │
               ├──────────────┬──────────────┐
               │              │              │
               ▼              ▼              ▼
      ┌────────────┐  ┌────────────┐  ┌────────────┐
      │   Router   │  │  Expert 1  │  │  Expert 2  │
      │  Network   │  │ (Animals)  │  │ (Vehicles) │
      └──────┬─────┘  └──────┬─────┘  └──────┬─────┘
             │               │                │
             │  Selecciona   │                │
             └───────────────┼────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Predicción    │
                    │  Especializada │
                    └────────────────┘
```

### Arquitectura

#### 1. Router Network

El router analiza la imagen y decide qué experto debe procesarla.

```python
# Ejemplo conceptual
def router(image):
    features = extract_features(image)
    expert_scores = router_network(features)  # [0.8, 0.2] 
    selected_expert = argmax(expert_scores)
    return selected_expert
```

#### 2. Expert Networks

Cada experto se especializa en un subconjunto de clases.

```
Expert 1 (animal_expert):
  - Especializado en: cat, dog, bird, fish, horse
  - Arquitectura optimizada para texturas animales
  - Mejor rendimiento en formas orgánicas

Expert 2 (vehicle_expert):
  - Especializado en: car, truck, airplane, boat, bicycle
  - Arquitectura optimizada para geometría
  - Mejor rendimiento en objetos mecánicos
```

### Ventajas de MoE

1. **Especialización**: Cada experto aprende mejor su dominio
2. **Escalabilidad**: Fácil agregar nuevos expertos
3. **Eficiencia**: Solo se ejecuta un experto por inferencia
4. **Mejor precisión**: Modelos especializados superan a generalistas
5. **Modularidad**: Puedes actualizar expertos independientemente

### Desventajas

1. **Complejidad**: Mayor complejidad arquitectónica
2. **Entrenamiento**: Requiere estrategia de entrenamiento cuidadosa
3. **Router crítico**: El router debe ser muy preciso
4. **Recursos**: Más modelos que almacenar

### Implementación Actual

```go
type MixtureOfExpertsModel struct {
    experts map[string][]string  // expert_name -> classes
    router  *Router               // Red del router
    loaded  bool
}

func (m *MixtureOfExpertsModel) Predict(imagePath string) ([]PredictionResult, error) {
    // 1. Router selecciona experto
    selectedExpert := m.router.SelectExpert(imagePath)
    
    // 2. Experto hace predicción
    predictions := m.experts[selectedExpert].Predict(imagePath)
    
    // 3. Retornar con metadata del experto
    return predictions
}
```

## � Binary Relevance

### Concepto

Enfoque de clasificación multi-label donde se entrena un clasificador binario independiente para cada clase. Cada clasificador decide si su clase está presente o no.

```
┌─────────────────────────────────────────┐
│          Entrada (Imagen)               │
└──────────────┬──────────────────────────┘
               │
               ├─────┬─────┬─────┬─────┬────┐
               │     │     │     │     │    │
               ▼     ▼     ▼     ▼     ▼    ▼
          ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐
          │BC1 │ │BC2 │ │BC3 │ │... │ │BC10│
          │cat │ │dog │ │bird│ │    │ │bike│
          └──┬─┘ └──┬─┘ └──┬─┘ └──┬─┘ └──┬─┘
             │     │     │     │     │    │
             ▼     ▼     ▼     ▼     ▼    ▼
          ┌──────────────────────────────────┐
          │  Predicciones Multi-Label        │
          │  (múltiples clases posibles)     │
          └──────────────────────────────────┘

BC = Binary Classifier (Clasificador Binario)
```

### Arquitectura

Cada clase tiene su propio clasificador binario:

```go
type BinaryRelevanceModel struct {
    classifiers map[string]*BinaryClassifier
    classes     []string
    loaded      bool
}

type BinaryClassifier struct {
    className string
    threshold float64  // Umbral de decisión
}
```

### Funcionamiento

**Clasificadores Binarios**:
- `binary_classifier_cat`: ¿Es un gato? (Sí/No)
- `binary_classifier_dog`: ¿Es un perro? (Sí/No)
- `binary_classifier_bird`: ¿Es un pájaro? (Sí/No)
- ... (uno por cada clase)

**Proceso de Inferencia**:
1. Cada clasificador binario evalúa la imagen independientemente
2. Produce una probabilidad (0-1) para su clase
3. Si supera el umbral (default: 0.5), la clase se considera presente
4. Resultado final: todas las clases que superaron el umbral

### Ventajas de Binary Relevance

1. **Multi-label**: Permite múltiples clases simultáneas
2. **Independencia**: Cada clasificador es autónomo
3. **Simplicidad**: Cada clasificador binario es más simple que un multiclase
4. **Flexibilidad**: Umbrales ajustables por clase
5. **Balanceo**: Maneja bien clases desbalanceadas
6. **Escalabilidad**: Fácil agregar nuevas clases

### Desventajas

1. **Correlación ignorada**: No captura relaciones entre clases
2. **Más modelos**: Necesita tantos modelos como clases
3. **Ineficiencia**: Potencialmente más lento que otros enfoques
4. **Ajuste de umbrales**: Requiere tuning de múltiples umbrales

### Casos de Uso Ideales

- **Imagenes con múltiples objetos**: Foto con perro Y gato
- **Tags/etiquetas**: Una imagen puede tener múltiples tags
- **Atributos múltiples**: Color, forma, textura simultáneos
- **Categorías no excluyentes**: Clases que pueden coexistir

### Implementación Actual

```go
func (m *BinaryRelevanceModel) Predict(imagePath string) ([]PredictionResult, error) {
    results := make([]PredictionResult, 0)

    // Cada clasificador evalúa independientemente
    for _, className := range m.classes {
        classifier := m.classifiers[className]
        probability := m.predictBinary(classifier, imagePath)
        
        results = append(results, PredictionResult{
            Class:      className,
            Confidence: probability,
            ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
        })
    }
    
    return results, nil
}

// Versión con umbral personalizado
func (m *BinaryRelevanceModel) PredictWithThreshold(imagePath string, threshold float64) ([]PredictionResult, error) {
    results := make([]PredictionResult, 0)
    
    for _, className := range m.classes {
        classifier := m.classifiers[className]
        probability := m.predictBinary(classifier, imagePath)
        
        // Solo incluir si supera el umbral
        if probability >= threshold {
            results = append(results, PredictionResult{
                Class:      className,
                Confidence: probability,
            })
        }
    }
    
    return results, nil
}
```

### Ejemplo de Uso

**Entrada**: Imagen con un perro y un coche

**Salida con Binary Relevance**:
```json
{
  "predictions": [
    {"class": "dog", "confidence": 0.92, "expert_used": "binary_classifier_dog"},
    {"class": "car", "confidence": 0.87, "expert_used": "binary_classifier_car"},
    {"class": "cat", "confidence": 0.12, "expert_used": "binary_classifier_cat"}
  ]
}
```

Con umbral 0.5, se detectan **ambas** clases: dog Y car (multi-label).

### Ajuste de Umbrales

```go
// Ajustar umbral para una clase específica
model.SetThreshold("cat", 0.7)  // Más estricto para gatos
model.SetThreshold("dog", 0.3)  // Más permisivo para perros

// Obtener solo clases positivas (por encima del umbral)
positives, err := model.GetPositiveClasses(imagePath)
```

## �🔄 Comparación

| Característica | Multiclass | MoE | Binary Relevance |
|----------------|------------|-----|------------------|
| Complejidad | Baja | Alta | Media |
| Precisión | Media | Alta | Media-Alta |
| Velocidad | Rápida | Media | Lenta |
| Escalabilidad | Baja | Alta | Media |
| Memoria | Baja | Media-Alta | Alta |
| Entrenamiento | Simple | Complejo | Medio |
| Multi-label | ❌ No | ❌ No | ✅ Sí |
| Independencia clases | ❌ | Parcial | ✅ |
| Ajuste por clase | ❌ | Parcial | ✅ |

### Cuándo usar cada modelo

**Multiclass**:
- Una sola clase por imagen
- Recursos limitados
- Prototipado rápido
- Clases mutuamente excluyentes

**Mixture of Experts**:
- Grupos naturales de clases
- Alta precisión requerida
- Recursos disponibles
- Clases con características distintas

**Binary Relevance**:
- Múltiples clases por imagen
- Necesidad de multi-label
- Umbrales específicos por clase
- Clases independientes entre sí

## 🚀 Integración de Modelos Reales

### Opción 1: TensorFlow Serving

Desplegar modelos TensorFlow como microservicio:

```yaml
# docker-compose.yml
services:
  tensorflow-serving:
    image: tensorflow/serving
    ports:
      - "8501:8501"
    volumes:
      - ./models:/models
    environment:
      - MODEL_NAME=classifier
```

```go
// En Go, hacer peticiones HTTP al servidor TF
func (m *Model) Predict(imagePath string) error {
    imageData := preprocessImage(imagePath)
    
    resp, err := http.Post(
        "http://tensorflow-serving:8501/v1/models/classifier:predict",
        "application/json",
        bytes.NewBuffer(imageData),
    )
    // Procesar respuesta
}
```

### Opción 2: ONNX Runtime

Usar modelos ONNX directamente en Go:

```go
import "github.com/yalue/onnxruntime_go"

type ONNXModel struct {
    session *onnxruntime.Session
}

func (m *ONNXModel) Predict(imagePath string) ([]float32, error) {
    // 1. Cargar imagen
    tensor := preprocessImageToTensor(imagePath)
    
    // 2. Ejecutar inferencia
    outputs, err := m.session.Run([]onnxruntime.Value{tensor})
    
    // 3. Procesar resultados
    return outputs[0].GetData().([]float32), nil
}
```

### Opción 3: Python Microservicio

Crear un microservicio Python separado:

```yaml
# docker-compose.yml
services:
  ml-service:
    build: ./ml-service
    ports:
      - "5000:5000"
    volumes:
      - ./models:/models
```

```python
# ml-service/app.py
from flask import Flask, request
import torch
import torchvision

app = Flask(__name__)
model = torch.load('model.pth')

@app.route('/predict', methods=['POST'])
def predict():
    image = preprocess(request.files['image'])
    prediction = model(image)
    return jsonify(prediction)
```

## 📊 Entrenamiento de Modelos

### Modelo Multiclass

```python
import torch
import torch.nn as nn

class MulticlassClassifier(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        self.backbone = torchvision.models.resnet50(pretrained=True)
        self.classifier = nn.Linear(2048, num_classes)
    
    def forward(self, x):
        features = self.backbone(x)
        return self.classifier(features)

# Entrenar
model = MulticlassClassifier(num_classes=10)
criterion = nn.CrossEntropyLoss()
optimizer = torch.optim.Adam(model.parameters())

for epoch in range(epochs):
    for images, labels in dataloader:
        outputs = model(images)
        loss = criterion(outputs, labels)
        
        optimizer.zero_grad()
        loss.backward()
        optimizer.step()
```

### Modelo MoE

```python
class Router(nn.Module):
    def __init__(self, num_experts=2):
        super().__init__()
        self.backbone = torchvision.models.resnet50(pretrained=True)
        self.router = nn.Linear(2048, num_experts)
    
    def forward(self, x):
        features = self.backbone(x)
        return F.softmax(self.router(features), dim=1)

class Expert(nn.Module):
    def __init__(self, num_classes):
        super().__init__()
        self.backbone = torchvision.models.resnet50(pretrained=True)
        self.classifier = nn.Linear(2048, num_classes)
    
    def forward(self, x):
        features = self.backbone(x)
        return self.classifier(features)

class MoEModel(nn.Module):
    def __init__(self):
        super().__init__()
        self.router = Router(num_experts=2)
        self.animal_expert = Expert(num_classes=5)
        self.vehicle_expert = Expert(num_classes=5)
    
    def forward(self, x):
        # Router decide
        router_probs = self.router(x)
        
        # Ejecutar expertos
        animal_pred = self.animal_expert(x)
        vehicle_pred = self.vehicle_expert(x)
        
        # Combinar con pesos del router
        output = (router_probs[:, 0:1] * animal_pred + 
                  router_probs[:, 1:2] * vehicle_pred)
        
        return output
```

## 🎓 Mejores Prácticas

### 1. Preprocesamiento

```python
def preprocess_image(image_path):
    transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize(
            mean=[0.485, 0.456, 0.406],
            std=[0.229, 0.224, 0.225]
        )
    ])
    return transform(Image.open(image_path))
```

### 2. Data Augmentation

```python
train_transform = transforms.Compose([
    transforms.RandomResizedCrop(224),
    transforms.RandomHorizontalFlip(),
    transforms.ColorJitter(0.4, 0.4, 0.4),
    transforms.RandomRotation(15),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                        std=[0.229, 0.224, 0.225])
])
```

### 3. Validación

```python
def evaluate(model, dataloader):
    model.eval()
    correct = 0
    total = 0
    
    with torch.no_grad():
        for images, labels in dataloader:
            outputs = model(images)
            _, predicted = torch.max(outputs.data, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()
    
    return 100 * correct / total
```

## 📈 Métricas de Evaluación

```python
from sklearn.metrics import classification_report, confusion_matrix

def evaluate_detailed(model, dataloader, class_names):
    all_preds = []
    all_labels = []
    
    model.eval()
    with torch.no_grad():
        for images, labels in dataloader:
            outputs = model(images)
            _, preds = torch.max(outputs, 1)
            all_preds.extend(preds.cpu().numpy())
            all_labels.extend(labels.cpu().numpy())
    
    # Reporte detallado
    print(classification_report(all_labels, all_preds, 
                                target_names=class_names))
    
    # Matriz de confusión
    cm = confusion_matrix(all_labels, all_preds)
    plot_confusion_matrix(cm, class_names)
```

## 🔮 Próximos Pasos

1. **Integrar modelos reales**: Usar TensorFlow/PyTorch preentrenados
2. **Implementar caché**: Guardar predicciones frecuentes
3. **Batch processing**: Procesar múltiples imágenes simultáneamente
4. **A/B Testing**: Comparar rendimiento de ambos modelos
5. **Monitoring**: Métricas de latencia y precisión en producción
6. **Fine-tuning**: Ajustar modelos con datos específicos
7. **Ensemble**: Combinar predicciones de ambos modelos

## 📚 Referencias

- [Mixture of Experts Paper](https://arxiv.org/abs/1701.06538)
- [PyTorch MoE Tutorial](https://pytorch.org/tutorials/)
- [TensorFlow Serving](https://www.tensorflow.org/tfx/guide/serving)
- [ONNX Runtime](https://onnxruntime.ai/)
