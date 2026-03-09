<template>
  <div id="app">
    <div class="container">
      <header>
        <h1>🖼️ Clasificador de Imágenes</h1>
        <p class="subtitle">Sistema de clasificación con Multiclass y Mixture of Experts</p>
      </header>

      <main>
        <!-- Selector de Modelo -->
        <section class="model-selector">
          <h2>Selecciona el Modelo</h2>
          <div class="model-cards">
            <div 
              class="model-card" 
              :class="{ active: selectedModel === 'multiclass' }"
              @click="selectModel('multiclass')"
            >
              <div class="model-icon">🎯</div>
              <h3>Modelo Multiclass</h3>
              <p>Clasificador estándar con todas las clases en un solo modelo</p>
              <div class="model-info">
                <span v-if="models.multiclass">{{ models.multiclass.classes.length }} clases</span>
              </div>
            </div>

            <div 
              class="model-card" 
              :class="{ active: selectedModel === 'mixture_of_experts' }"
              @click="selectModel('mixture_of_experts')"
            >
              <div class="model-icon">🧠</div>
              <h3>Mixture of Experts</h3>
              <p>Modelo MoE con expertos especializados y router inteligente</p>
              <div class="model-info">
                <span v-if="models.moe">{{ models.moe.experts.length }} expertos</span>
              </div>
            </div>

            <div 
              class="model-card" 
              :class="{ active: selectedModel === 'binary_relevance' }"
              @click="selectModel('binary_relevance')"
            >
              <div class="model-icon">🔢</div>
              <h3>Binary Relevance</h3>
              <p>Clasificadores binarios independientes para multi-label</p>
              <div class="model-info">
                <span v-if="models.binary">{{ models.binary.classifiers }} clasificadores</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Upload de Imagen -->
        <section class="upload-section">
          <h2>Sube una Imagen</h2>
          <div 
            class="upload-area" 
            :class="{ 'drag-over': isDragging }"
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="handleDrop"
            @click="$refs.fileInput.click()"
          >
            <input 
              type="file" 
              ref="fileInput"
              accept="image/*"
              @change="handleFileSelect"
              style="display: none"
            >
            <div v-if="!selectedImage" class="upload-placeholder">
              <div class="upload-icon">📤</div>
              <p>Arrastra una imagen aquí o haz clic para seleccionar</p>
              <span class="upload-hint">JPG, PNG, GIF hasta 10MB</span>
            </div>
            <div v-else class="image-preview">
              <img :src="imagePreview" alt="Preview">
              <button class="remove-btn" @click.stop="removeImage">✕</button>
            </div>
          </div>
          
          <button 
            class="classify-btn" 
            :disabled="!selectedImage || isClassifying"
            @click="classifyImage"
          >
            <span v-if="!isClassifying">🔍 Clasificar Imagen</span>
            <span v-else>⏳ Clasificando...</span>
          </button>
        </section>

        <!-- Resultados -->
        <section v-if="results" class="results-section">
          <h2>Resultados de Clasificación</h2>
          <div class="results-header">
            <span class="model-badge">{{ getModelBadge() }}</span>
            <span class="time-badge">⚡ {{ results.process_time_ms }}ms</span>
          </div>

          <div class="predictions">
            <div 
              v-for="(pred, index) in results.predictions" 
              :key="index"
              class="prediction-item"
              :class="'rank-' + (index + 1)"
            >
              <div class="prediction-rank">{{ index + 1 }}</div>
              <div class="prediction-content">
                <div class="prediction-header">
                  <span class="prediction-class">{{ pred.class }}</span>
                  <span v-if="pred.expert_used" class="expert-badge">👤 {{ pred.expert_used }}</span>
                </div>
                <div class="confidence-bar">
                  <div 
                    class="confidence-fill" 
                    :style="{ width: (pred.confidence * 100) + '%' }"
                  ></div>
                </div>
                <span class="confidence-text">{{ (pred.confidence * 100).toFixed(2) }}%</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Error Message -->
        <div v-if="error" class="error-message">
          <strong>❌ Error:</strong> {{ error }}
        </div>
      </main>

      <footer>
        <p>Desarrollado con Vue.js y Go | Dockerizado para despliegue fácil</p>
      </footer>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      selectedModel: 'multiclass',
      selectedImage: null,
      imagePreview: null,
      isDragging: false,
      isClassifying: false,
      results: null,
      error: null,
      models: {
        multiclass: null,
        moe: null,
        binary: null
      },
      apiUrl: process.env.VUE_APP_API_URL || 'http://localhost:8080'
    };
  },
  mounted() {
    this.loadModels();
  },
  methods: {
    async loadModels() {
      try {
        const response = await axios.get(`${this.apiUrl}/api/models`);
        this.models.multiclass = response.data.find(m => m.id === 'multiclass');
        this.models.moe = response.data.find(m => m.id === 'mixture_of_experts');
        this.models.binary = response.data.find(m => m.id === 'binary_relevance');
      } catch (err) {
        console.error('Error cargando modelos:', err);
      }
    },
    selectModel(model) {
      this.selectedModel = model;
      this.results = null;
      this.error = null;
    },
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file) {
        this.processFile(file);
      }
    },
    handleDrop(event) {
      this.isDragging = false;
      const file = event.dataTransfer.files[0];
      if (file && file.type.startsWith('image/')) {
        this.processFile(file);
      }
    },
    processFile(file) {
      if (file.size > 10 * 1024 * 1024) {
        this.error = 'La imagen es demasiado grande. Máximo 10MB.';
        return;
      }

      this.selectedImage = file;
      this.error = null;
      this.results = null;

      // Crear preview
      const reader = new FileReader();
      reader.onload = (e) => {
        this.imagePreview = e.target.result;
      };
      reader.readAsDataURL(file);
    },
    removeImage() {
      this.selectedImage = null;
      this.imagePreview = null;
      this.results = null;
      this.error = null;
    },
    async classifyImage() {
      if (!this.selectedImage) return;

      this.isClassifying = true;
      this.error = null;
      this.results = null;

      const formData = new FormData();
      formData.append('image', this.selectedImage);
      formData.append('model', this.selectedModel);

      try {
        const response = await axios.post(
          `${this.apiUrl}/api/classify`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        );
        this.results = response.data;
      } catch (err) {
        this.error = err.response?.data?.message || 'Error al clasificar la imagen. Verifica que el servidor esté corriendo.';
        console.error('Error:', err);
      } finally {
        this.isClassifying = false;
      }
    },
    getModelBadge() {
      if (!this.results) return '';
      const badges = {
        'multiclass': '🎯 Multiclass',
        'mixture_of_experts': '🧠 MoE',
        'binary_relevance': '🔢 Binary Relevance'
      };
      return badges[this.results.model] || this.results.model;
    }
  }
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

#app {
  padding: 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

header {
  text-align: center;
  color: white;
  margin-bottom: 40px;
}

header h1 {
  font-size: 3rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

main {
  background: white;
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

h2 {
  color: #333;
  margin-bottom: 20px;
  font-size: 1.5rem;
}

/* Model Selector */
.model-selector {
  margin-bottom: 40px;
}

.model-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.model-card {
  background: #f8f9fa;
  border: 3px solid transparent;
  border-radius: 15px;
  padding: 30px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.model-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
}

.model-card.active {
  border-color: #667eea;
  background: linear-gradient(135deg, #667eea15, #764ba215);
  box-shadow: 0 5px 20px rgba(102, 126, 234, 0.3);
}

.model-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.model-card h3 {
  color: #333;
  margin-bottom: 10px;
}

.model-card p {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: 15px;
}

.model-info {
  color: #667eea;
  font-weight: 600;
  font-size: 0.85rem;
}

/* Upload Section */
.upload-section {
  margin-bottom: 40px;
}

.upload-area {
  border: 3px dashed #ccc;
  border-radius: 15px;
  padding: 60px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
  margin-bottom: 20px;
}

.upload-area:hover, .upload-area.drag-over {
  border-color: #667eea;
  background: #f0f4ff;
}

.upload-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.upload-placeholder p {
  font-size: 1.2rem;
  color: #333;
  margin-bottom: 10px;
}

.upload-hint {
  color: #999;
  font-size: 0.9rem;
}

.image-preview {
  position: relative;
  max-width: 500px;
  margin: 0 auto;
}

.image-preview img {
  max-width: 100%;
  max-height: 400px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.remove-btn {
  position: absolute;
  top: -10px;
  right: -10px;
  background: #ff4757;
  color: white;
  border: none;
  border-radius: 50%;
  width: 35px;
  height: 35px;
  font-size: 1.2rem;
  cursor: pointer;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.remove-btn:hover {
  background: #ff3838;
  transform: scale(1.1);
}

.classify-btn {
  width: 100%;
  padding: 18px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1.2rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.classify-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.classify-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Results Section */
.results-section {
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.results-header {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.model-badge, .time-badge {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
}

.model-badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.time-badge {
  background: #f0f4ff;
  color: #667eea;
}

.predictions {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.prediction-item {
  display: flex;
  gap: 15px;
  align-items: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.prediction-item:hover {
  transform: translateX(5px);
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
}

.prediction-item.rank-1 {
  background: linear-gradient(135deg, #ffd70015, #ffa50015);
  border-left: 5px solid #ffd700;
}

.prediction-rank {
  font-size: 1.5rem;
  font-weight: 700;
  color: #667eea;
  min-width: 40px;
  text-align: center;
}

.prediction-content {
  flex: 1;
}

.prediction-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.prediction-class {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  text-transform: capitalize;
}

.expert-badge {
  font-size: 0.8rem;
  padding: 4px 12px;
  background: #667eea;
  color: white;
  border-radius: 12px;
}

.confidence-bar {
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.confidence-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transition: width 0.5s ease;
}

.confidence-text {
  font-size: 0.9rem;
  color: #666;
  font-weight: 600;
}

/* Error Message */
.error-message {
  background: #ff475715;
  border: 2px solid #ff4757;
  border-radius: 10px;
  padding: 15px;
  color: #ff4757;
  margin-top: 20px;
}

footer {
  text-align: center;
  color: white;
  margin-top: 40px;
  opacity: 0.8;
}

/* Responsive */
@media (max-width: 768px) {
  header h1 {
    font-size: 2rem;
  }

  main {
    padding: 20px;
  }

  .model-cards {
    grid-template-columns: 1fr;
  }
}
</style>
