package main

import (
	"fmt"
	"log"
	"math/rand"
)

// BinaryRelevanceModel representa un modelo de Binary Relevance
// Cada clase tiene su propio clasificador binario independiente
type BinaryRelevanceModel struct {
	classifiers map[string]*BinaryClassifier
	classes     []string
	loaded      bool
}

// BinaryClassifier es un clasificador binario para una clase específica
type BinaryClassifier struct {
	className string
	threshold float64 // Umbral para considerar la clase como positiva
}

func NewBinaryRelevanceModel() *BinaryRelevanceModel {
	classes := []string{
		"cat", "dog", "bird", "fish", "horse",
		"car", "truck", "airplane", "boat", "bicycle",
	}

	classifiers := make(map[string]*BinaryClassifier)
	for _, class := range classes {
		classifiers[class] = &BinaryClassifier{
			className: class,
			threshold: 0.5, // Umbral por defecto
		}
	}

	model := &BinaryRelevanceModel{
		classifiers: classifiers,
		classes:     classes,
		loaded:      true,
	}

	log.Println("Modelo Binary Relevance cargado exitosamente")
	log.Printf("Clasificadores binarios creados: %d", len(classifiers))
	return model
}

func (m *BinaryRelevanceModel) GetClasses() []string {
	return m.classes
}

func (m *BinaryRelevanceModel) GetClassifierCount() int {
	return len(m.classifiers)
}

func (m *BinaryRelevanceModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	log.Printf("Clasificando con modelo Binary Relevance: %s", imagePath)

	results := make([]PredictionResult, 0)

	// Cada clasificador binario decide independientemente
	// si su clase está presente en la imagen
	for _, className := range m.classes {
		classifier := m.classifiers[className]
		
		// Simular la salida del clasificador binario (en producción, aquí iría la inferencia real)
		// Cada clasificador devuelve una probabilidad de que la clase esté presente
		probability := m.predictBinary(classifier, imagePath)
		
		// Solo incluir clases que superen el umbral
		if probability >= classifier.threshold {
			results = append(results, PredictionResult{
				Class:      className,
				Confidence: probability,
				ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
			})
		} else {
			// También podemos incluir todas las clases con sus probabilidades
			// para dar contexto completo (aunque estén por debajo del umbral)
			results = append(results, PredictionResult{
				Class:      className,
				Confidence: probability,
				ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
			})
		}
	}

	// Ordenar por confianza descendente
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Retornar solo top 5 para consistencia con otros modelos
	if len(results) > 5 {
		results = results[:5]
	}

	return results, nil
}

// predictBinary simula la inferencia de un clasificador binario
// En producción, esto llamaría al modelo real
func (m *BinaryRelevanceModel) predictBinary(classifier *BinaryClassifier, imagePath string) float64 {
	// Simulación: generar probabilidad aleatoria
	// En un modelo real, aquí se procesaría la imagen con el clasificador específico
	probability := rand.Float64()
	
	// Ajustar probabilidades para hacer algunas clases más probables que otras
	// Esto simula un comportamiento más realista
	if rand.Float64() > 0.7 {
		// 30% de las veces, dar una probabilidad alta
		probability = 0.5 + (probability * 0.5) // Entre 0.5 y 1.0
	}
	
	return probability
}

// PredictWithThreshold permite especificar un umbral personalizado
func (m *BinaryRelevanceModel) PredictWithThreshold(imagePath string, threshold float64) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	log.Printf("Clasificando con Binary Relevance (umbral: %.2f): %s", threshold, imagePath)

	results := make([]PredictionResult, 0)
	positiveClasses := make([]string, 0)

	// Ejecutar clasificadores binarios
	for _, className := range m.classes {
		classifier := m.classifiers[className]
		probability := m.predictBinary(classifier, imagePath)
		
		result := PredictionResult{
			Class:      className,
			Confidence: probability,
			ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
		}
		
		// Marcar clases positivas (por encima del umbral)
		if probability >= threshold {
			positiveClasses = append(positiveClasses, className)
		}
		
		results = append(results, result)
	}

	log.Printf("Clases detectadas (umbral %.2f): %d/%d - %v", 
		threshold, len(positiveClasses), len(m.classes), positiveClasses)

	// Ordenar por confianza
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

// GetAllPredictions devuelve predicciones para todas las clases sin filtrar
func (m *BinaryRelevanceModel) GetAllPredictions(imagePath string) (map[string]float64, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	predictions := make(map[string]float64)
	
	for _, className := range m.classes {
		classifier := m.classifiers[className]
		probability := m.predictBinary(classifier, imagePath)
		predictions[className] = probability
	}

	return predictions, nil
}

// SetThreshold permite ajustar el umbral de un clasificador específico
func (m *BinaryRelevanceModel) SetThreshold(className string, threshold float64) error {
	classifier, exists := m.classifiers[className]
	if !exists {
		return fmt.Errorf("clase '%s' no encontrada", className)
	}
	
	if threshold < 0.0 || threshold > 1.0 {
		return fmt.Errorf("umbral debe estar entre 0.0 y 1.0")
	}
	
	classifier.threshold = threshold
	log.Printf("Umbral actualizado para '%s': %.2f", className, threshold)
	return nil
}

// GetPositiveClasses devuelve solo las clases que superan el umbral
func (m *BinaryRelevanceModel) GetPositiveClasses(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	results := make([]PredictionResult, 0)

	for _, className := range m.classes {
		classifier := m.classifiers[className]
		probability := m.predictBinary(classifier, imagePath)
		
		// Solo incluir si supera el umbral
		if probability >= classifier.threshold {
			results = append(results, PredictionResult{
				Class:      className,
				Confidence: probability,
				ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
			})
		}
	}

	// Ordenar por confianza
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}
