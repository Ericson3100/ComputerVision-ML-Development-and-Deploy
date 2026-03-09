package main

import (
	"fmt"
	"log"
	"math/rand"
)

// MixtureOfExpertsModel representa un modelo MoE con múltiples expertos especializados
type MixtureOfExpertsModel struct {
	experts map[string][]string // expert name -> classes
	router  *Router
	loaded  bool
}

type Router struct {
	// El router decide qué experto usar basado en la entrada
}

type Expert struct {
	name    string
	classes []string
}

func NewMixtureOfExpertsModel() *MixtureOfExpertsModel {
	model := &MixtureOfExpertsModel{
		experts: map[string][]string{
			"animal_expert": {"cat", "dog", "bird", "fish", "horse"},
			"vehicle_expert": {"car", "truck", "airplane", "boat", "bicycle"},
		},
		router: &Router{},
		loaded: true,
	}
	log.Println("Modelo Mixture of Experts cargado exitosamente")
	log.Printf("Expertos disponibles: %d", len(model.experts))
	for expertName, classes := range model.experts {
		log.Printf("  - %s: %d clases", expertName, len(classes))
	}
	return model
}

func (m *MixtureOfExpertsModel) GetClasses() []string {
	allClasses := make([]string, 0)
	for _, classes := range m.experts {
		allClasses = append(allClasses, classes...)
	}
	return allClasses
}

func (m *MixtureOfExpertsModel) GetExperts() []string {
	experts := make([]string, 0, len(m.experts))
	for expertName := range m.experts {
		experts = append(experts, expertName)
	}
	return experts
}

func (m *MixtureOfExpertsModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	log.Printf("Clasificando con modelo Mixture of Experts: %s", imagePath)

	// Paso 1: El router selecciona el experto más apropiado
	// En un modelo real, el router analizaría la imagen y elegiría el experto
	selectedExpert, expertClasses := m.selectExpert(imagePath)
	log.Printf("Router seleccionó experto: %s", selectedExpert)

	// Paso 2: El experto seleccionado hace la predicción
	results := make([]PredictionResult, 0)
	
	// Generar probabilidades para las clases del experto
	confidences := make([]float64, len(expertClasses))
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

	// Crear resultados con información del experto usado
	for i := 0; i < len(expertClasses); i++ {
		results = append(results, PredictionResult{
			Class:      expertClasses[i],
			Confidence: confidences[i],
			ExpertUsed: selectedExpert,
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

	// Retornar top 5
	if len(results) > 5 {
		results = results[:5]
	}

	return results, nil
}

func (m *MixtureOfExpertsModel) selectExpert(imagePath string) (string, []string) {
	// El router analiza la imagen y decide qué experto usar
	// En este ejemplo simplificado, elegimos aleatoriamente
	// En producción, el router sería una red neuronal que aprende a elegir
	
	expertNames := make([]string, 0, len(m.experts))
	for name := range m.experts {
		expertNames = append(expertNames, name)
	}
	
	selectedExpert := expertNames[rand.Intn(len(expertNames))]
	return selectedExpert, m.experts[selectedExpert]
}

// Método adicional para obtener predicciones de todos los expertos con sus pesos
func (m *MixtureOfExpertsModel) PredictWithAllExperts(imagePath string) (map[string][]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("modelo no cargado")
	}

	allPredictions := make(map[string][]PredictionResult)

	// Obtener predicciones de cada experto
	for expertName, classes := range m.experts {
		results := make([]PredictionResult, 0)
		
		confidences := make([]float64, len(classes))
		total := 0.0
		for i := range confidences {
			val := rand.Float64()
			confidences[i] = val
			total += val
		}
		
		for i := range confidences {
			confidences[i] /= total
		}

		for i := 0; i < len(classes); i++ {
			results = append(results, PredictionResult{
				Class:      classes[i],
				Confidence: confidences[i],
				ExpertUsed: expertName,
			})
		}

		allPredictions[expertName] = results
	}

	return allPredictions, nil
}
