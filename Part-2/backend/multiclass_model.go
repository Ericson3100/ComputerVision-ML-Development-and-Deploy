package main

import (
	"fmt"
	"log"
	"os"
)

// MulticlassModel represents a standard multiclass classification model
type MulticlassModel struct {
	modelPath string
	classes   []string
	loaded    bool
}

func NewMulticlassModel() *MulticlassModel {
	modelPath := "./models/simple_model_multiclass/unique_model_multiclass"
	classes := []string{
		"dew", "fogsmog", "frost", "glaze", "hail", "lightning",
		"rain", "rainbow", "rime", "sandstorm", "snow", "volcanos",
	}

	m := &MulticlassModel{
		modelPath: modelPath,
		classes:   classes,
		loaded:    false,
	}

	// Check if the model exists
	if _, err := os.Stat(modelPath); err == nil {
		m.loaded = true
		log.Println("Multiclass model loaded successfully")
		log.Printf("Model path: %s", m.modelPath)
	} else {
		log.Printf("ERROR: Required multiclass model not found at %s", modelPath)
	}

	log.Printf("Classes: %d", len(m.classes))

	return m
}

func (m *MulticlassModel) GetClasses() []string {
	return m.classes
}

func (m *MulticlassModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("model not loaded: model file not found at %s", m.modelPath)
	}

	log.Printf("Classifying with Multiclass model: %s", imagePath)

	// Execute inference via Python
	probabilities, err := RunPythonInference(m.modelPath, imagePath)
	if err != nil {
		return nil, fmt.Errorf("error in multiclass model inference: %v", err)
	}

	// Return all class probabilities sorted by confidence
	results := GetTopKPredictions(probabilities, m.classes, len(m.classes))

	return results, nil
}
