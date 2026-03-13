package main

import (
	"fmt"
	"log"
	"os"
)

// BinaryRelevanceModel represents a Binary Relevance model
type BinaryRelevanceModel struct {
	classifiers   map[string]*BinaryClassifier
	classes       []string
	modelBasePath string
	loaded        bool
}

// BinaryClassifier is a binary classifier for a specific class
type BinaryClassifier struct {
	className string
	modelPath string
	loaded    bool
}

func NewBinaryRelevanceModel() *BinaryRelevanceModel {
	classes := []string{
		"dew", "fogsmog", "frost", "glaze", "hail", "lightning",
		"rain", "rainbow", "rime", "sandstorm", "snow", "volcanos",
	}

	modelBasePath := "./models/binary_relevance"
	classifiers := make(map[string]*BinaryClassifier)

	loadedCount := 0
	for _, class := range classes {
		modelPath := fmt.Sprintf("%s/%s_relevance_expert", modelBasePath, class)
		classifier := &BinaryClassifier{
			className: class,
			modelPath: modelPath,
			loaded:    false,
		}

		// Check if model exists
		if _, err := os.Stat(modelPath); err == nil {
			classifier.loaded = true
			loadedCount++
		}

		classifiers[class] = classifier
	}

	model := &BinaryRelevanceModel{
		classifiers:   classifiers,
		classes:       classes,
		modelBasePath: modelBasePath,
		loaded:        loadedCount > 0,
	}

	log.Println("Binary Relevance model loaded successfully")
	log.Printf("Binary classifiers created: %d", len(classifiers))
	log.Printf("Classifiers loaded successfully: %d", loadedCount)
	log.Printf("Base model path: %s", modelBasePath)
	for _, class := range classes {
		status := "MISSING"
		if classifiers[class].loaded {
			status = "LOADED"
		}
		log.Printf("  - Classifier %s: %s [%s]", class, classifiers[class].modelPath, status)
	}

	return model
}

func (m *BinaryRelevanceModel) GetClasses() []string {
	return m.classes
}

func (m *BinaryRelevanceModel) GetClassifierCount() int {
	return len(m.classifiers)
}

func (m *BinaryRelevanceModel) Predict(imagePath string) ([]PredictionResult, error) {
	log.Printf("Classifying with Binary Relevance: %s", imagePath)

	results := make([]PredictionResult, 0)

	if !m.loaded {
		return nil, fmt.Errorf("no binary classifier is loaded: check models in %s", m.modelBasePath)
	}

	for _, className := range m.classes {
		classifier := m.classifiers[className]

		if !classifier.loaded {
			return nil, fmt.Errorf("classifier for class '%s' not loaded: not found in %s", className, classifier.modelPath)
		}

		probability, err := m.predictBinary(classifier, imagePath)
		if err != nil {
			return nil, fmt.Errorf("error in classifier %s: %v", className, err)
		}

		results = append(results, PredictionResult{
			Class:      className,
			Confidence: probability,
			ExpertUsed: fmt.Sprintf("binary_classifier_%s", className),
		})
	}

	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

// predictBinary executes inference with a specific binary classifier
func (m *BinaryRelevanceModel) predictBinary(classifier *BinaryClassifier, imagePath string) (float64, error) {
	probabilities, err := RunPythonInference(classifier.modelPath, imagePath)
	if err != nil {
		return 0, err
	}

	if len(probabilities) == 0 {
		return 0, fmt.Errorf("no probabilities received")
	}

	// If it returns 2 values, take the second one (positive class probability)
	if len(probabilities) == 2 {
		return float64(probabilities[1]), nil
	}

	// If it returns 1 value, assume it as positive class probability
	return float64(probabilities[0]), nil
}
