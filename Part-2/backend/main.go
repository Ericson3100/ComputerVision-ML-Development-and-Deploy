package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ClassificationRequest struct {
	Model string `json:"model"` // "multiclass", "mixture_of_experts" o "binary_relevance"
}

type ClassificationResponse struct {
	Model       string             `json:"model"`
	Predictions []PredictionResult `json:"predictions"`
	ProcessTime float64            `json:"process_time_ms"`
}

type PredictionResult struct {
	Class       string  `json:"class"`
	Confidence  float64 `json:"confidence"`
	ExpertUsed  string  `json:"expert_used,omitempty"`
}

type Server struct {
	multiclassModel       *MulticlassModel
	mixtureOfExpertsModel *MixtureOfExpertsModel
	binaryRelevanceModel  *BinaryRelevanceModel
}

func main() {
	log.Println("Inicializando servidor de clasificación de imágenes...")

	// Inicializar modelos
	multiclassModel := NewMulticlassModel()
	mixtureModel := NewMixtureOfExpertsModel()
	binaryRelevanceModel := NewBinaryRelevanceModel()

	server := &Server{
		multiclassModel:       multiclassModel,
		mixtureOfExpertsModel: mixtureModel,
		binaryRelevanceModel:  binaryRelevanceModel,
	}

	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/api/classify", server.classifyHandler).Methods("POST")
	router.HandleFunc("/api/models", server.getModelsHandler).Methods("GET")

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

	log.Printf("Servidor escuchando en puerto %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) getModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := []map[string]interface{}{
		{
			"id":          "multiclass",
			"name":        "Modelo Multiclass",
			"description": "Clasificador multiclase estándar con todas las clases en un solo modelo",
			"classes":     s.multiclassModel.GetClasses(),
		},
		{
			"id":          "mixture_of_experts",
			"name":        "Mixture of Experts",
			"description": "Modelo MoE con múltiples expertos especializados y un router",
			"classes":     s.mixtureOfExpertsModel.GetClasses(),
			"experts":     s.mixtureOfExpertsModel.GetExperts(),
		},
		{
			"id":          "binary_relevance",
			"name":        "Binary Relevance",
			"description": "Clasificadores binarios independientes para cada clase (multi-label)",
			"classes":     s.binaryRelevanceModel.GetClasses(),
			"classifiers": s.binaryRelevanceModel.GetClassifierCount(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

func (s *Server) classifyHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Parsear multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error al parsear form", http.StatusBadRequest)
		return
	}

	// Obtener modelo seleccionado
	modelType := r.FormValue("model")
	if modelType == "" {
		modelType = "multiclass"
	}

	// Obtener archivo
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No se pudo leer la imagen", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Procesando imagen: %s con modelo: %s", header.Filename, modelType)

	// Guardar temporalmente la imagen
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d%s", time.Now().Unix(), filepath.Ext(header.Filename)))
	dst, err := os.Create(tempFile)
	if err != nil {
		http.Error(w, "Error al guardar imagen", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	defer os.Remove(tempFile)

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error al copiar imagen", http.StatusInternalServerError)
		return
	}

	// Clasificar con el modelo seleccionado
	var predictions []PredictionResult
	if modelType == "mixture_of_experts" {
		predictions, err = s.mixtureOfExpertsModel.Predict(tempFile)
	} else if modelType == "binary_relevance" {
		predictions, err = s.binaryRelevanceModel.Predict(tempFile)
	} else {
		predictions, err = s.multiclassModel.Predict(tempFile)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error en clasificación: %v", err), http.StatusInternalServerError)
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
