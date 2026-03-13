package main

import (
	"fmt"
	"log"
	"os"
)

// MixtureOfExpertsModel represents an integrated MoE model
type MixtureOfExpertsModel struct {
	modelPath  string
	allClasses []string
	loaded     bool
}

func NewMixtureOfExpertsModel() *MixtureOfExpertsModel {
	allClasses := []string{
		"dew", "fogsmog", "frost", "glaze", "hail", "lightning",
		"rain", "rainbow", "rime", "sandstorm", "snow", "volcanos",
	}

	modelPath := "./models/MoE/moe_model"

	m := &MixtureOfExpertsModel{
		modelPath:  modelPath,
		allClasses: allClasses,
		loaded:     false,
	}

	// Check if the model exists
	if _, err := os.Stat(modelPath); err == nil {
		m.loaded = true
		log.Println("Simple Mixture of Experts model loaded successfully")
		log.Printf("Type: Integrated MoE (experts + gating in a single model)")
		log.Printf("Model path: %s", m.modelPath)
	} else {
		log.Printf("ERROR: Required Mixture of Experts model not found at %s", modelPath)
	}

	log.Printf("Total classes: %d", len(m.allClasses))

	return m
}

func (m *MixtureOfExpertsModel) GetClasses() []string {
	return m.allClasses
}

func (m *MixtureOfExpertsModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("model not loaded: model file not found in %s", m.modelPath)
	}

	log.Printf("Classifying with integrated MoE: %s", imagePath)

	// Execute inference via Python
	probabilities, err := RunPythonInference(m.modelPath, imagePath)
	if err != nil {
		return nil, fmt.Errorf("error in MoE model inference: %v", err)
	}

	// Return all class probabilities sorted by confidence
	results := GetTopKPredictions(probabilities, m.allClasses, len(m.allClasses))

	// Add MoE metadata
	for i := range results {
		results[i].ExpertUsed = "moe_integrated"
	}

	return results, nil
}
