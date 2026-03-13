package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// TensorFlowInferenceResponse represents the Python script response
type TensorFlowInferenceResponse struct {
	Predictions []float32 `json:"predictions"`
	Error       *string   `json:"error"`
}

// RunPythonInference executes inference using the Python script with TensorFlow
func RunPythonInference(modelPath, imagePath string) ([]float32, error) {
	// Execute Python script
	cmd := exec.Command("python3", "tf_inference.py", modelPath, imagePath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python inference: %v\nOutput: %s", err, string(output))
		return nil, fmt.Errorf("error executing inference: %v", err)
	}

	// Parse JSON response
	var response TensorFlowInferenceResponse
	if err := json.Unmarshal(output, &response); err != nil {
		log.Printf("Error parsing JSON response: %v\nOutput: %s", err, string(output))
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Check if there was an error
	if response.Error != nil && *response.Error != "" {
		return nil, fmt.Errorf("error in Python inference: %s", *response.Error)
	}

	if response.Predictions == nil {
		return nil, fmt.Errorf("no predictions received")
	}

	return response.Predictions, nil
}

// GetTopKPredictions returns the top-k predictions sorted by confidence
func GetTopKPredictions(probabilities []float32, classes []string, k int) []PredictionResult {
	if len(probabilities) != len(classes) {
		return nil
	}

	results := make([]PredictionResult, len(classes))
	for i := range classes {
		results[i] = PredictionResult{
			Class:      classes[i],
			Confidence: float64(probabilities[i]),
		}
	}

	// Sort by descending confidence
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Confidence > results[i].Confidence {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Return only top-k
	if k > 0 && k < len(results) {
		results = results[:k]
	}

	return results
}
