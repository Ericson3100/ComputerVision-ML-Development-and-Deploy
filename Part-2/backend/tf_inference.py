"""
TensorFlow Model Inference Server
Loads and runs inference on TensorFlow SavedModel format models.
"""
import sys
import json
import os
import gc
from pathlib import Path

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
os.environ['TF_ENABLE_ONEDNN_OPTS'] = '0'

TF_IMPORT_ERROR = None
try:
    import tensorflow as tf
    tf.get_logger().setLevel('ERROR')
    TF_AVAILABLE = True
except ImportError as e:
    TF_AVAILABLE = False
    TF_IMPORT_ERROR = str(e)

efficientnet_preprocess = None
if TF_AVAILABLE:
    try:
        efficientnet_preprocess = tf.keras.applications.efficientnet.preprocess_input
    except AttributeError:
        efficientnet_preprocess = None


class ModelInference:
    """Handles loading and prediction of TensorFlow models."""
    
    def __init__(self, model_path: str):
        self.model_path = model_path
        self.model = None
        self.error_message = None
        
        if not TF_AVAILABLE:
            details = f": {TF_IMPORT_ERROR}" if TF_IMPORT_ERROR else ""
            self.error_message = f"TensorFlow is not installed or could not be imported{details}"
            return
        
        if not Path(model_path).exists():
            self.error_message = f"Model does not exist at path: {model_path}"
            return
            
        try:
            self.model = tf.saved_model.load(model_path)
            self.infer = self.model.signatures["serving_default"]
        except Exception as e:
            self.error_message = f"Error loading model: {str(e)}"
            self.model = None
    
    def preprocess_image(self, image_path: str, target_size=(224, 224)):
        """Loads and preprocesses an image using official EfficientNetB0 preprocessing."""
        if not TF_AVAILABLE or self.model is None or efficientnet_preprocess is None:
            return None
        
        try:
            img = tf.io.read_file(image_path)
            img = tf.image.decode_image(img, channels=3, expand_animations=False)
            
            img = tf.image.resize(img, target_size, method='bilinear')
            
            if efficientnet_preprocess is not None:
                img = efficientnet_preprocess(img)
            else:
                img = tf.cast(img, tf.float32)
            
            img = tf.expand_dims(img, 0)
            
            return img
        except Exception as e:
            print(f"ERROR preprocessing image {image_path}: {e}", file=sys.stderr)
            return None
    
    def predict(self, image_path: str):
        """Executes prediction on an image."""
        if self.error_message:
            raise Exception(self.error_message)
        
        if not TF_AVAILABLE or self.model is None:
            raise Exception("Model not available for inference")
        
        img_tensor = None
        try:
            img_tensor = self.preprocess_image(image_path)
            if img_tensor is None:
                raise Exception("Error preprocessing image")
            
            predictions = self.infer(img_tensor)
            
            output_key = None
            for key in predictions.keys():
                output_key = key
                break
            
            if output_key is None:
                raise Exception("Model output not found")
            
            output = predictions[output_key].numpy()[0]

            result = output.tolist()
            
            del predictions
            del output
            if img_tensor is not None:
                del img_tensor
            
            return result
            
        except Exception as e:
            if img_tensor is not None:
                del img_tensor
            raise Exception(f"Error during inference: {str(e)}")


def main():
    """
    Main script for inference.
    
    Usage:
        python tf_inference.py <model_path> <image_path>
    
    Output (JSON):
        {"predictions": [0.1, 0.2, ...], "error": null}
    or
        {"predictions": null, "error": "error message"}
    """
    if len(sys.argv) != 3:
        response = {
            "predictions": None,
            "error": "Usage: python tf_inference.py <model_path> <image_path>"
        }
        print(json.dumps(response))
        sys.exit(1)
    
    model_path = sys.argv[1]
    image_path = sys.argv[2]
    
    if not Path(image_path).exists():
        response = {
            "predictions": None,
            "error": f"Image does not exist: {image_path}"
        }
        print(json.dumps(response))
        sys.exit(1)
    
    inference = ModelInference(model_path)
    
    if inference.error_message:
        response = {
            "predictions": None,
            "error": inference.error_message
        }
        print(json.dumps(response))
        sys.exit(1)
    
    try:
        predictions = inference.predict(image_path)
        response = {
            "predictions": predictions,
            "error": None
        }
        print(json.dumps(response))
        
        del inference
        del predictions
        gc.collect()
        
    except Exception as e:
        response = {
            "predictions": None,
            "error": f"Prediction error: {str(e)}"
        }
        print(json.dumps(response))
        
        del inference
        gc.collect()
        
        sys.exit(1)


if __name__ == "__main__":
    main()
