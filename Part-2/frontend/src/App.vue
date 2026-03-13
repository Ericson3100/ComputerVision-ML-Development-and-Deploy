<template>
  <div class="app-shell">
    <div class="dynamic-background" aria-hidden="true">
      <div class="dynamic-background-stage" :style="backgroundStageStyle">
        <FloatingLines
          :key="floatingLinesKey"
          :enabled-waves="['top', 'middle', 'bottom']"
          :line-count="[5, 6, 7]"
          :line-distance="[1, 1, 1]"
          :bend-radius="10.0"
          :bend-strength="-1"
          :interactive="true"
          :parallax="true"
          :animation-speed="1"
          :parallax-strength="0.1"
          :lines-gradient="['#0ea5e9', '#22d3ee', '#38bdf8', '#818cf8', '#c084fc']"
          mix-blend-mode="screen"
        />
      </div>
    </div>
    <div class="container">
      <main>
        <header class="content-header">
          <h1>🖼️ Image Classifier - Multi-Model</h1>
          <p class="content-subtitle">Upload weather imagery and compare predictions across the available models.</p>
        </header>

        <section class="upload-section">
          <h2>📤 Upload One or More Images</h2>
          <div
            class="upload-area"
            :class="{ 'drag-over': isDragging && !isClassifying, 'upload-disabled': isClassifying }"
            @dragover.prevent="!isClassifying && (isDragging = true)"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="!isClassifying && handleDrop($event)"
            @click="!isClassifying && openFileDialog()"
          >
            <input
              type="file"
              ref="fileInput"
              accept="image/*"
              multiple
              :disabled="isClassifying"
              @change="handleFileSelect"
              style="display: none"
            >
            <div v-if="selectedImages.length === 0" class="upload-placeholder">
              <div class="upload-icon">📁</div>
              <p>📎 Drag images here or click to select</p>
              <span class="upload-hint">🖼️ JPG, PNG, GIF - Multiple images supported</span>
            </div>
            <div v-else class="images-preview">
              <div v-for="(image, index) in selectedImages" :key="image.key" class="preview-item">
                <img :src="image.preview" :alt="image.file.name">
                <div class="preview-overlay">
                  <span class="preview-name">{{ image.file.name }}</span>
                  <button
                    class="remove-btn"
                    :disabled="isClassifying"
                    @click.stop="!isClassifying && removeImage(index)"
                  >X</button>
                </div>
              </div>
              <div v-if="isClassifying" class="classification-overlay">
                <span>⏳ Classification in progress...</span>
              </div>
            </div>
          </div>

          <div class="query-control">
            <label class="query-label">🧭 Models</label>
            <div class="query-actions">
              <label class="model-option select-all-option">
                <input
                  type="checkbox"
                  v-model="allModelsSelected"
                  :disabled="isClassifying || availableModels.length === 0"
                >
                <span>Select all models</span>
              </label>
              <span class="models-count">{{ selectedModelIds.length }}/{{ availableModels.length }} selected</span>
            </div>

            <div class="query-field">
              <label>Available models</label>
              <div class="models-selector" :class="{ disabled: isClassifying || availableModels.length === 0 }">
                <label class="model-option" v-for="model in availableModels" :key="model.id">
                  <input
                    type="checkbox"
                    :value="model.id"
                    v-model="selectedModelIds"
                    :disabled="isClassifying"
                  >
                  <span>{{ model.name }}</span>
                </label>
              </div>
            </div>
          </div>

          <div class="threshold-control">
            <label for="threshold">
              <span class="threshold-label">🎯 Confidence Threshold: <strong>{{ threshold.toFixed(2) }}</strong></span>
              <span class="threshold-hint">📊 Only predictions above this value will be shown</span>
            </label>
            <input
              id="threshold"
              type="range"
              min="0"
              max="1"
              step="0.0001"
              v-model.number="threshold"
              @input="handleThresholdChange"
            >
            <div class="threshold-markers">
              <span>0.0</span>
              <span>0.25</span>
              <span>0.5</span>
              <span>0.75</span>
              <span>1.0</span>
            </div>
          </div>

          <button
            class="classify-btn"
            :disabled="selectedImages.length === 0 || isClassifying || effectiveModelIds.length === 0"
            @click="classifyImages"
          >
            <span v-if="!isClassifying">🚀 Classify with Selected Models</span>
            <span v-else>⏳ Classifying... {{ progress }}</span>
          </button>
        </section>

        <section v-if="results" class="results-section">
          <h2>📋 Classification Results</h2>

          <div class="overall-stats">
            <div class="stat-card">
              <div class="stat-icon">🖼️</div>
              <div class="stat-content">
                <div class="stat-value">{{ results.overall.total_images }}</div>
                <div class="stat-label">Images</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">⏱️</div>
              <div class="stat-content">
                <div class="stat-value">{{ Number(results.overall.total_process_time_ms).toFixed(0) }}ms</div>
                <div class="stat-label">Total Time</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">📈</div>
              <div class="stat-content">
                <div class="stat-value">{{ Number(results.overall.average_per_image_ms).toFixed(0) }}ms</div>
                <div class="stat-label">Average/Image</div>
              </div>
            </div>
          </div>

          <div v-for="(image, imageIndex) in results.images" :key="`${image.image_name}-${imageIndex}`" class="image-result">
            <div class="image-result-header">
              <img v-if="image.preview" :src="image.preview" :alt="image.image_name" class="result-preview-image">
              <h3 class="image-title">{{ image.image_name }}</h3>
            </div>
            <div class="models-grid">
              <div
                v-for="(model, modelIndex) in image.results"
                :key="`${model.model_id}-${modelIndex}`"
                class="model-result-card"
              >
                <div class="model-header">
                  <div class="model-title">
                    <span class="model-icon">{{ getModelIcon(model.model_id) }}</span>
                    <span class="model-name">{{ model.model_name }}</span>
                  </div>
                  <div class="model-time">
                    <span class="time-badge">{{ Number(model.process_time_ms).toFixed(0) }}ms</span>
                  </div>
                </div>

                <div v-if="getTopPrediction(model)" class="winner-prediction">
                  <div class="winner-badge">🏆 Winner</div>
                  <div class="winner-class">{{ getTopPrediction(model)?.class }}</div>
                  <div class="winner-confidence">{{ (Number(getTopPrediction(model)?.confidence) * 100).toFixed(2) }}%</div>
                </div>

                <div v-if="getFilteredPredictions(model).length > 0" class="predictions-list">
                  <div
                    v-for="(pred, predIndex) in getFilteredPredictions(model)"
                    :key="`${pred.class}-${predIndex}`"
                    class="prediction-item"
                    :class="{ 'top-prediction': predIndex === 0 }"
                  >
                    <span class="pred-rank">#{{ predIndex + 1 }}</span>
                    <span class="pred-class">{{ pred.class }}</span>
                    <div class="confidence-bar-mini">
                      <div
                        class="confidence-fill-mini"
                        :style="{ width: (Number(pred.confidence) * 100) + '%' }"
                      ></div>
                    </div>
                    <span class="pred-confidence">{{ (Number(pred.confidence) * 100).toFixed(1) }}%</span>
                  </div>
                </div>

                <div v-else class="no-predictions">
                  <p>😕 No predictions above threshold {{ threshold.toFixed(2) }}</p>
                  <button class="lower-threshold-btn" @click="lowerThreshold">
                    🔽 Lower Threshold
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div v-if="error" class="error-message">
          <strong>❌ Error:</strong> {{ error }}
        </div>
      </main>

      <footer>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import FloatingLines from './FloatingLines.vue';

const fileInput = ref(null);
const selectedImages = ref([]);
const isDragging = ref(false);
const isClassifying = ref(false);
const threshold = ref(0);
const results = ref(null);
const error = ref('');
const progress = ref('');
const availableModels = ref([]);
const selectedModelIds = ref([]);
const backgroundScale = ref(1);
const floatingLinesKey = ref(0);
const zoomRatio = ref(window.devicePixelRatio || 1);
let backgroundRafId = null;

const BACKGROUND_BASE_WIDTH = 1920;
const BACKGROUND_BASE_HEIGHT = 1080;

const updateBackgroundScale = () => {
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const nextScale = Math.max(
    viewportWidth / BACKGROUND_BASE_WIDTH,
    viewportHeight / BACKGROUND_BASE_HEIGHT,
    1,
  );

  if (Math.abs(nextScale - backgroundScale.value) > 0.001) {
    backgroundScale.value = nextScale;
  }
};

const scheduleBackgroundScaleUpdate = () => {
  if (backgroundRafId !== null) {
    return;
  }

  backgroundRafId = window.requestAnimationFrame(() => {
    updateBackgroundScale();

    const nextZoomRatio = window.devicePixelRatio || 1;
    if (Math.abs(nextZoomRatio - zoomRatio.value) > 0.001) {
      zoomRatio.value = nextZoomRatio;
      floatingLinesKey.value += 1;
    }

    backgroundRafId = null;
  });
};

const backgroundStageStyle = computed(() => ({
  width: `${Math.round(BACKGROUND_BASE_WIDTH * backgroundScale.value)}px`,
  height: `${Math.round(BACKGROUND_BASE_HEIGHT * backgroundScale.value)}px`,
  transform: 'translate(-50%, -50%)',
}));

const imageKey = (file) => `${file.name}-${file.size}-${file.lastModified}`;

const allModelsSelected = computed({
  get: () =>
    availableModels.value.length > 0 && selectedModelIds.value.length === availableModels.value.length,
  set: (checked) => {
    selectedModelIds.value = checked ? availableModels.value.map((model) => model.id) : [];
  },
});

const effectiveModelIds = computed(() => {
  return selectedModelIds.value.filter((modelId) =>
    availableModels.value.some((model) => model.id === modelId),
  );
});

const openFileDialog = () => {
  fileInput.value?.click();
};

const addFiles = (fileList) => {
  const files = Array.from(fileList || []);
  if (!files.length) {
    return;
  }

  const existingKeys = new Set(selectedImages.value.map((item) => item.key));

  for (const file of files) {
    if (!file.type.startsWith('image/')) {
      continue;
    }

    const key = imageKey(file);
    if (existingKeys.has(key)) {
      continue;
    }

    selectedImages.value.push({
      key,
      file,
      preview: URL.createObjectURL(file),
    });
    existingKeys.add(key);
  }

  error.value = '';
};

const handleFileSelect = (event) => {
  addFiles(event.target.files);
  event.target.value = '';
};

const handleDrop = (event) => {
  isDragging.value = false;
  addFiles(event.dataTransfer?.files);
};

const removeImage = (index) => {
  const [removed] = selectedImages.value.splice(index, 1);
  if (removed?.preview) {
    URL.revokeObjectURL(removed.preview);
  }
};

const handleThresholdChange = () => {
  error.value = '';
};

const getModelIcon = (modelId) => {
  const icons = {
    multiclass: '🧠',
    multilabel: '🏷️',
    mixture_of_experts: '🧩',
    binary_relevance: '⚖️',
  };
  return icons[modelId] || '🤖';
};

const getModelName = (modelId) => {
  const model = availableModels.value.find((item) => item.id === modelId);
  return model?.name || modelId;
};

const normalizePredictions = (predictions = [], confidenceThreshold = threshold.value) => {
  return (predictions || [])
    .filter((prediction) => Number(prediction.confidence) >= confidenceThreshold)
    .sort((left, right) => Number(right.confidence) - Number(left.confidence));
};

const toRawModelResult = (modelId, predictions, processTimeMs) => {
  const normalized = normalizePredictions(predictions, 0);
  const topPrediction = normalized[0] || null;

  return {
    model_id: modelId,
    model_name: getModelName(modelId),
    predictions: normalized,
    process_time_ms: Number(processTimeMs) || 0,
    top_class: topPrediction?.class || '',
    top_confidence: Number(topPrediction?.confidence || 0),
  };
};

const getFilteredPredictions = (model) => {
  return normalizePredictions(model?.predictions || []);
};

const getTopPrediction = (model) => {
  const filtered = getFilteredPredictions(model);
  return filtered[0] || null;
};

const updateOverallStats = () => {
  if (!results.value) {
    return;
  }

  let totalProcessTime = 0;
  for (const image of results.value.images) {
    for (const model of image.results || []) {
      totalProcessTime += Number(model.process_time_ms || 0);
    }
  }

  results.value.overall.total_process_time_ms = totalProcessTime;
  results.value.overall.average_per_image_ms = results.value.images.length > 0
    ? totalProcessTime / results.value.images.length
    : 0;
};

const classifyAllModels = async () => {
  const totalTasks = selectedImages.value.length * effectiveModelIds.value.length;
  let completedTasks = 0;

  results.value = {
    images: selectedImages.value.map((image) => ({
      image_name: image.file.name,
      preview: image.preview,
      results: [],
    })),
    overall: {
      total_images: selectedImages.value.length,
      total_process_time_ms: 0,
      average_per_image_ms: 0,
    },
  };

  const formData = new FormData();
  for (const image of selectedImages.value) {
    formData.append('images', image.file);
  }
  for (const modelId of effectiveModelIds.value) {
    formData.append('models', modelId);
  }

  const response = await fetch('/api/classify-stream', {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `HTTP ${response.status}`);
  }

  if (!response.body) {
    throw new Error('Streaming response body is not available');
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  const applyStreamEvent = (eventChunk) => {
    const lines = eventChunk
      .split('\n')
      .map((line) => line.trim())
      .filter(Boolean);

    let eventName = 'message';
    let data = '';

    for (const line of lines) {
      if (line.startsWith('event:')) {
        eventName = line.slice(6).trim();
      } else if (line.startsWith('data:')) {
        data += line.slice(5).trim();
      }
    }

    if (!data) {
      return;
    }

    if (eventName === 'error') {
      const payload = JSON.parse(data);
      throw new Error(payload.error || 'Streaming classification error');
    }

    if (eventName === 'result') {
      const payload = JSON.parse(data);
      const imageName = payload.image_name;
      const modelId = payload.model_id;

      const imageEntry = results.value.images.find((image) => image.image_name === imageName);
      if (!imageEntry) {
        return;
      }

      const modelResult = toRawModelResult(
        modelId,
        payload.predictions || [],
        payload.process_time_ms || 0,
      );

      const existingIndex = imageEntry.results.findIndex((item) => item.model_id === modelId);
      if (existingIndex >= 0) {
        imageEntry.results.splice(existingIndex, 1, modelResult);
      } else {
        imageEntry.results.push(modelResult);
      }

      completedTasks += 1;
      progress.value = `${completedTasks}/${totalTasks}`;
      updateOverallStats();
    }
  };

  while (true) {
    const { value, done } = await reader.read();
    if (done) {
      break;
    }

    buffer += decoder.decode(value, { stream: true });

    let separatorIndex = buffer.indexOf('\n\n');
    while (separatorIndex !== -1) {
      const eventChunk = buffer.slice(0, separatorIndex);
      buffer = buffer.slice(separatorIndex + 2);

      if (eventChunk.trim()) {
        applyStreamEvent(eventChunk);
      }

      separatorIndex = buffer.indexOf('\n\n');
    }
  }

  updateOverallStats();

  return results.value;
};

const loadModels = async () => {
  try {
    const response = await fetch('/api/models');
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const models = await response.json();
    availableModels.value = Array.isArray(models) ? models : [];

    if (selectedModelIds.value.length === 0 && availableModels.value.length > 0) {
      selectedModelIds.value = availableModels.value.map((model) => model.id);
    }
  } catch {
    availableModels.value = [
      { id: 'multiclass', name: 'Multiclass Model' },
      { id: 'multilabel', name: 'Multilabel Model' },
      { id: 'mixture_of_experts', name: 'Mixture of Experts' },
      { id: 'binary_relevance', name: 'Binary Relevance' },
    ];

    if (selectedModelIds.value.length === 0) {
      selectedModelIds.value = availableModels.value.map((model) => model.id);
    }
  }
};

const classifyImages = async () => {
  if (!selectedImages.value.length || isClassifying.value) {
    return;
  }

  isClassifying.value = true;
  error.value = '';
  results.value = null;
  progress.value = `0/${selectedImages.value.length * effectiveModelIds.value.length}`;

  try {
    results.value = await classifyAllModels();
  } catch (err) {
    results.value = null;
    error.value = err instanceof Error ? err.message : 'Unexpected error while classifying images';
  } finally {
    isClassifying.value = false;
  }
};

const lowerThreshold = async () => {
  threshold.value = Math.max(0, Number((threshold.value - 0.01).toFixed(4)));
};

onMounted(() => {
  loadModels();
  updateBackgroundScale();
  window.addEventListener('resize', scheduleBackgroundScaleUpdate, { passive: true });
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', scheduleBackgroundScaleUpdate);
  if (backgroundRafId !== null) {
    window.cancelAnimationFrame(backgroundRafId);
    backgroundRafId = null;
  }
  for (const image of selectedImages.value) {
    if (image.preview) {
      URL.revokeObjectURL(image.preview);
    }
  }
});
</script>

<style src="./App.css"></style>
