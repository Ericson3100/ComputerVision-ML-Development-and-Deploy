"""
Script to inspect exported Keras 3 models and get operation names

Useful for configuring input/output names in Go
"""

import sys
import subprocess

def inspect_savedmodel(model_path):
    """
    Inspects a SavedModel and shows relevant information
    for TensorFlow Go integration
    """
    print("=" * 80)
    print(f"Inspecting model: {model_path}")
    print("=" * 80)
    
    # Execute saved_model_cli
    try:
        result = subprocess.run(
            ["saved_model_cli", "show", "--dir", model_path, "--all"],
            capture_output=True,
            text=True,
            check=True
        )
        
        output = result.stdout
        print(output)
        
        # Extract and highlight key information
        print("\n" + "=" * 80)
        print("INFORMATION FOR GO:")
        print("=" * 80)
        
        # Search for input names
        for line in output.split('\n'):
            if 'name:' in line and 'input' in line.lower():
                # Extract operation name (without :0)
                parts = line.split('name:')
                if len(parts) > 1:
                    op_name = parts[1].strip().split(':')[0]
                    print(f"[OK] Input Operation: \"{op_name}\"")
            
            if 'name:' in line and 'output' in line.lower():
                parts = line.split('name:')
                if len(parts) > 1:
                    op_name = parts[1].strip().split(':')[0]
                    print(f"Output Operation: \"{op_name}\"")
        
        print("\nUse in mixture_of_experts_model.go:")
        print("expert := &ExpertModel{")
        print(f"    inputOp:  \"<copy_input_operation>\",")
        print(f"    outputOp: \"<copy_output_operation>\",")
        print("}")
        
    except subprocess.CalledProcessError as e:
        print(f"Error executing saved_model_cli: {e}")
        print("Make sure TensorFlow is installed: pip install tensorflow")
    except FileNotFoundError:
        print("saved_model_cli not found.")
        print("Install: pip install tensorflow")

def inspect_all_models():
    """Inspects all models in the models/ directory"""
    import os
    
    models_dir = "./models"
    
    if not os.path.exists(models_dir):
        print(f"Directory {models_dir} not found")
        print("Run first: python train_simple_moe.py --export")
        return
    
    for model_name in os.listdir(models_dir):
        model_path = os.path.join(models_dir, model_name)
        
        # Check if it's a directory and contains saved_model.pb
        if os.path.isdir(model_path):
            pb_file = os.path.join(model_path, "saved_model.pb")
            if os.path.exists(pb_file):
                inspect_savedmodel(model_path)
                print("\n")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        # Inspect specific model
        model_path = sys.argv[1]
        inspect_savedmodel(model_path)
    else:
        # Inspect all models
        print("Inspecting all models in ./models/\n")
        inspect_all_models()
        
        print("\nTo inspect a specific model:")
        print("python inspect_model.py ./models/expert_1")
