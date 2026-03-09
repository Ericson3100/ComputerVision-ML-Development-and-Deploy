package main

import (
	"fmt"
	"log"
	"math/rand"
)

// MulticlassModel representa un modelo de clasificación multiclase estándar
type MulticlassModel struct {
	classes []string
	loaded  bool
}

func NewMulticlassModel() *MulticlassModel {
	model := &MulticlassModel{
		classes: []string{
			"cat", "dog", "bird", "fish", "horse",
			"car", "truck", "airplane", "boat", "bicycle",
		},
		loaded: true,
	}
	log.Println("Modelo Multiclass cargado exitosamente")
	return model
}

func (m *MulticlassModel) GetClasses() []string {
	return m.classes
}

func (m *MulticlassModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	log.Printf("Clasificando con modelo Multiclass: %s", imagePath)

	// Simulación de inferencia del modelo
	// En producción, aquí cargarías y procesarías la imagen con un modelo real
	// como TensorFlow, PyTorch (vía ONNX), o usando un servicio externo

	results := make([]PredictionResult, 0)
	
	// Generar probabilidades aleatorias que sumen 1.0
	confidences := make([]float64, len(m.classes))
	total := 0.0
	for i := range confidences {
		val := rand.Float64()
		confidences[i] = val
		total += val
	}
	
	// Normalizar
	for i := range confidences {
		confidences[i] /= total
	}

	// Ordenar por confianza y retornar top 5
	for i := 0; i < len(m.classes); i++ {
		results = append(results, PredictionResult{
			Class:      m.classes[i],
			Confidence: confidences[i],
		})
	}

	// Ordenar por confianza descendente
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Retornar solo top 5
	if len(results) > 5 {
		results = results[:5]
	}

	return results, nil
}
