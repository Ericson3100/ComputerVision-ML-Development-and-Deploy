package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ClassificationRequest struct {
	Model string `json:"model"` // "multiclass", "mixture_of_experts", "binary_relevance", "multilabel"
}

type ClassificationResponse struct {
	Model       string             `json:"model"`
	Predictions []PredictionResult `json:"predictions"`
	ProcessTime float64            `json:"process_time_ms"`
}

type MultiModelResponse struct {
	Images  []ImageResults `json:"images"`
	Overall OverallStats   `json:"overall"`
}

type ImageResults struct {
	ImageName string        `json:"image_name"`
	Results   []ModelResult `json:"results"`
}

type ModelResult struct {
	ModelID     string             `json:"model_id"`
	ModelName   string             `json:"model_name"`
	Predictions []PredictionResult `json:"predictions"`
	ProcessTime float64            `json:"process_time_ms"`
	TopClass    string             `json:"top_class"`
	TopConfidence float64          `json:"top_confidence"`
}

type PredictionResult struct {
	Class      string  `json:"class"`
	Confidence float64 `json:"confidence"`
	ExpertUsed string  `json:"expert_used,omitempty"`
}

type ProgressiveModelResult struct {
	ImageName     string             `json:"image_name"`
	ModelID       string             `json:"model_id"`
	ModelName     string             `json:"model_name"`
	Predictions   []PredictionResult `json:"predictions"`
	ProcessTime   float64            `json:"process_time_ms"`
	TopClass      string             `json:"top_class"`
	TopConfidence float64            `json:"top_confidence"`
	Complete      bool               `json:"complete"`
}

type OverallStats struct {
	TotalImages      int     `json:"total_images"`
	TotalProcessTime float64 `json:"total_process_time_ms"`
	AveragePerImage  float64 `json:"average_per_image_ms"`
}

type Server struct {
	multiclassModel       *MulticlassModel
	mixtureOfExpertsModel *MixtureOfExpertsModel
	binaryRelevanceModel  *BinaryRelevanceModel
	multilabelModel       *MultilabelModel
}

func main() {
	log.Println("Initializing image classification server...")

	// Initialize models
	multiclassModel := NewMulticlassModel()
	mixtureModel := NewMixtureOfExpertsModel()
	binaryRelevanceModel := NewBinaryRelevanceModel()
	multilabelModel := NewMultilabelModel()

	if err := validateRequiredModels(multiclassModel, mixtureModel, binaryRelevanceModel, multilabelModel); err != nil {
		log.Fatalf("Startup failed: %v", err)
	}

	server := &Server{
		multiclassModel:       multiclassModel,
		mixtureOfExpertsModel: mixtureModel,
		binaryRelevanceModel:  binaryRelevanceModel,
		multilabelModel:       multilabelModel,
	}

	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/api/classify", server.classifyHandler).Methods("POST")
	router.HandleFunc("/api/classify-stream", server.classifyStreamHandler).Methods("POST")
	router.HandleFunc("/api/models", server.getModelsHandler).Methods("GET")
	router.HandleFunc("/swagger/doc.yaml", swaggerDocHandler).Methods("GET")
	router.HandleFunc("/swagger", swaggerUIHandler).Methods("GET")
	router.HandleFunc("/swagger/", swaggerUIHandler).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func validateRequiredModels(
	multiclass *MulticlassModel,
	mixture *MixtureOfExpertsModel,
	binary *BinaryRelevanceModel,
	multilabel *MultilabelModel,
) error {
	if !multiclass.loaded {
		return fmt.Errorf("required multiclass model not found at %s", multiclass.modelPath)
	}

	if !multilabel.loaded {
		return fmt.Errorf("required multilabel model not found at %s", multilabel.modelPath)
	}

	if !mixture.loaded {
		return fmt.Errorf("required mixture_of_experts model not found at %s", mixture.modelPath)
	}

	if !binary.loaded {
		return fmt.Errorf("required binary_relevance models not found under %s", binary.modelBasePath)
	}

	for _, className := range binary.classes {
		classifier := binary.classifiers[className]
		if !classifier.loaded {
			return fmt.Errorf("required binary_relevance classifier for class '%s' not found at %s", className, classifier.modelPath)
		}
	}

	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// swaggerDocHandler serves the OpenAPI YAML document.
func swaggerDocHandler(w http.ResponseWriter, r *http.Request) {
	docPath := filepath.Join("swagger", "openapi.yaml")

	content, err := os.ReadFile(docPath)
	if err != nil {
		http.Error(w, "Swagger document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Write(content)
}

// swaggerUIHandler serves a minimal Swagger UI page.
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	swaggerURL := "/swagger/doc.yaml"
	if strings.HasSuffix(r.URL.Path, "/") {
		swaggerURL = "doc.yaml"
	}

	html := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Image Classifier API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '%s',
      dom_id: '#swagger-ui'
    });
  </script>
</body>
</html>`, swaggerURL)

	w.Write([]byte(html))
}

// getModelsHandler returns all available models with metadata.
func (s *Server) getModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := []map[string]interface{}{
		{
			"id":          "multiclass",
			"name":        "Multiclass Model",
			"description": "Standard multiclass classifier with all classes in a single model",
			"classes":     s.multiclassModel.GetClasses(),
		},
		{
			"id":          "multilabel",
			"name":        "Multilabel Model",
			"description": "Multilabel classifier that can assign multiple labels simultaneously",
			"classes":     s.multilabelModel.GetClasses(),
		},
		{
			"id":          "mixture_of_experts",
			"name":        "Mixture of Experts",
			"description": "MoE model with multiple specialized experts and a router",
			"classes":     s.mixtureOfExpertsModel.GetClasses(),
		},
		{
			"id":          "binary_relevance",
			"name":        "Binary Relevance",
			"description": "Independent binary classifiers for each class (multi-label)",
			"classes":     s.binaryRelevanceModel.GetClasses(),
			"classifiers": s.binaryRelevanceModel.GetClassifierCount(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

// classifyHandler predicts classes for one image using one selected model.
func (s *Server) classifyHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20) // 100 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	// Get selected model
	modelType := r.FormValue("model")
	if modelType == "" {
		modelType = "multiclass"
	}

	// Get file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing image: %s with model: %s", header.Filename, modelType)

	// Save image temporarily
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d%s", time.Now().Unix(), filepath.Ext(header.Filename)))
	dst, err := os.Create(tempFile)
	if err != nil{
		http.Error(w, "Error saving image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	defer os.Remove(tempFile)

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error copying image", http.StatusInternalServerError)
		return
	}

	// Classify with selected model
	var predictions []PredictionResult
	if modelType == "mixture_of_experts" {
		predictions, err = s.mixtureOfExpertsModel.Predict(tempFile)
	} else if modelType == "binary_relevance" {
		predictions, err = s.binaryRelevanceModel.Predict(tempFile)
	} else if modelType == "multilabel" {
		predictions, err = s.multilabelModel.Predict(tempFile)
	} else {
		predictions, err = s.multiclassModel.Predict(tempFile)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Classification error: %v", err), http.StatusInternalServerError)
		return
	}

	processTime := time.Since(startTime).Milliseconds()

	response := ClassificationResponse{
		Model:       modelType,
		Predictions: predictions,
		ProcessTime: float64(processTime),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// classifyStreamHandler executes selected models over uploaded images and streams progressive SSE results.
func (s *Server) classifyStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Configure headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"Error parsing form\"}\n\n")
		flusher.Flush()
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	log.Printf("Starting stream classification")

	// Process all images
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"No images uploaded\"}\n\n")
		flusher.Flush()
		return
	}

	// Prepare temporary images
	type imageFile struct {
		name     string
		tempPath string
	}

	imageFiles := make([]imageFile, 0)
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Save image temporarily
		tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d_%s", time.Now().UnixNano(), fileHeader.Filename))
		dst, err := os.Create(tempFile)
		if err != nil {
			file.Close()
			log.Printf("Error creating temporary file: %v", err)
			continue
		}

		_, err = io.Copy(dst, file)
		dst.Close()
		file.Close()
		
		if err != nil {
			os.Remove(tempFile)
			log.Printf("Error copying image: %v", err)
			continue
		}

		imageFiles = append(imageFiles, imageFile{
			name:     fileHeader.Filename,
			tempPath: tempFile,
		})
	}

	if len(imageFiles) == 0 {
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"Could not process images\"}\n\n")
		flusher.Flush()
		return
	}

	cleanupTempFiles := func() {
		for _, img := range imageFiles {
			if err := os.Remove(img.tempPath); err != nil && !os.IsNotExist(err) {
				log.Printf("Error removing temporary file %s: %v", img.tempPath, err)
			}
		}
	}

	// Determine selected models (optional); defaults to all 4
	modelCatalog := map[string]struct {
		name string
		fn   func(string) ([]PredictionResult, error)
	}{
		"multiclass": {
			name: "Multiclass Model",
			fn:   func(path string) ([]PredictionResult, error) { return s.multiclassModel.Predict(path) },
		},
		"multilabel": {
			name: "Multilabel Model",
			fn:   func(path string) ([]PredictionResult, error) { return s.multilabelModel.Predict(path) },
		},
		"mixture_of_experts": {
			name: "Mixture of Experts",
			fn:   func(path string) ([]PredictionResult, error) { return s.mixtureOfExpertsModel.Predict(path) },
		},
		"binary_relevance": {
			name: "Binary Relevance",
			fn:   func(path string) ([]PredictionResult, error) { return s.binaryRelevanceModel.Predict(path) },
		},
	}

	selectedModelIDs := r.MultipartForm.Value["models"]
	if len(selectedModelIDs) == 0 {
		selectedModelIDs = []string{"multiclass", "multilabel", "mixture_of_experts", "binary_relevance"}
	}

	type modelSpec struct {
		id   string
		name string
		fn   func(string) ([]PredictionResult, error)
	}

	models := make([]modelSpec, 0, len(selectedModelIDs))
	for _, modelID := range selectedModelIDs {
		if model, ok := modelCatalog[modelID]; ok {
			models = append(models, modelSpec{id: modelID, name: model.name, fn: model.fn})
		}
	}

	if len(models) == 0 {
		cleanupTempFiles()
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"No valid models selected\"}\n\n")
		flusher.Flush()
		return
	}

	// Channel to receive results asynchronously
	type modelResult struct {
		imageName     string
		modelID       string
		modelName     string
		predictions   []PredictionResult
		processTime   float64
		topClass      string
		topConfidence float64
	}

	results := make(chan modelResult, len(imageFiles)*len(models))
	var wg sync.WaitGroup

	// Launch goroutine to process each image x model combination
	totalTasks := len(imageFiles) * len(models)
	wg.Add(totalTasks)

	for _, img := range imageFiles {
		for _, model := range models {
			go func(imageName, imagePath, modelID, modelName string, modelFn func(string) ([]PredictionResult, error)) {
				defer wg.Done()

				start := time.Now()
				log.Printf("Starting %s for %s", modelName, imageName)

				preds, err := modelFn(imagePath)

				processTime := float64(time.Since(start).Milliseconds())

				if err != nil {
					log.Printf("Error in %s for %s: %v", modelName, imageName, err)
					return
				}

				log.Printf("Completed %s for %s in %.0fms", modelName, imageName, processTime)

				// Determine top class
				topClass := ""
				topConfidence := 0.0
				if len(preds) > 0 {
					topClass = preds[0].Class
					topConfidence = preds[0].Confidence
				}


				results <- modelResult{
					imageName:     imageName,
					modelID:       modelID,
					modelName:     modelName,
					predictions:   preds,
					processTime:   processTime,
					topClass:      topClass,
					topConfidence: topConfidence,
				}
			}(img.name, img.tempPath, model.id, model.name, model.fn)
		}
	}

	// Goroutine to close channel when all tasks finish or are cancelled
	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		wg.Wait()
		cleanupTempFiles()
	}()

	// Send results as they arrive
	resultCount := 0
	for {
		select {
		case <-r.Context().Done():
			log.Printf("Client disconnected, stopping result sending")
			return
		case result, ok := <-results:
			if !ok {
				// Channel closed, all tasks completed
				goto complete
			}

			resultCount++

			progressiveResult := ProgressiveModelResult{
				ImageName:     result.imageName,
				ModelID:       result.modelID,
				ModelName:     result.modelName,
				Predictions:   result.predictions,
				ProcessTime:   result.processTime,
				TopClass:      result.topClass,
				TopConfidence: result.topConfidence,
				Complete:      false,
			}

			data, _ := json.Marshal(progressiveResult)
			_, err := fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
			if err != nil {
				log.Printf("Error sending result, client disconnected")
				return
			}
			flusher.Flush()

			log.Printf("Sent result %d/%d: %s - %s", resultCount, totalTasks, result.modelName, result.imageName)
			predictionsJSON, _ := json.Marshal(result.predictions)
			log.Printf("Sent probabilities for %s - %s: %s", result.modelName, result.imageName, string(predictionsJSON))
		}
	}

	complete:
	// Send completion event
	log.Printf("Stream classification completed: %d results sent", resultCount)
	fmt.Fprintf(w, "event: complete\ndata: {\"complete\": true}\n\n")
	flusher.Flush()
}
