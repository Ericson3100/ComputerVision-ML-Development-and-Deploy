package main

import (
	"fmt"
	"log"
	"os"
)

// MultilabelModel represents a multilabel classification model
type MultilabelModel struct {
	modelPath string
	classes   []string
	loaded    bool
}

func NewMultilabelModel() *MultilabelModel {
	modelPath := "./models/simple_model_multilabel/unique_model_multilabel"
	classes := []string{
		"dew", "fogsmog", "frost", "glaze", "hail", "lightning",
		"rain", "rainbow", "rime", "sandstorm", "snow", "volcanos",
	}

	m := &MultilabelModel{
		modelPath: modelPath,
		classes:   classes,
		loaded:    false,
	}

	// Check if the model exists
	if _, err := os.Stat(modelPath); err == nil {
		m.loaded = true
		log.Println("Multilabel model loaded successfully")
		log.Printf("Model path: %s", m.modelPath)
	} else {
		log.Printf("ERROR: Required multilabel model not found at %s", modelPath)
	}

	log.Printf("Classes: %d", len(m.classes))

	return m
}

func (m *MultilabelModel) GetClasses() []string {
	return m.classes
}

func (m *MultilabelModel) Predict(imagePath string) ([]PredictionResult, error) {
	if !m.loaded {
		return nil, fmt.Errorf("model not loaded: model file not found at %s", m.modelPath)
	}

	log.Printf("Classifying with Multilabel model: %s", imagePath)

	probabilities, err := RunPythonInference(m.modelPath, imagePath)
	if err != nil {
		return nil, fmt.Errorf("error in multilabel model inference: %v", err)
	}

	return GetTopKPredictions(probabilities, m.classes, len(m.classes)), nil
}
