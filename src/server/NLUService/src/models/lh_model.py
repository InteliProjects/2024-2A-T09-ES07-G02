import pickle
import os
import logging

# Class definition for the LHModel
class LHModel:
    # Constructor method
    def __init__(self):
        model_path = os.path.join(os.path.dirname(__file__), 'model.pkl')
        self.model = self._load_model(model_path)
    
    # Method to load the model from a file
    def _load_model(self, model_path):
        # Log the model loading process
        logging.info(f"Loading model from {model_path}")
        # Open the model file in binary read mode and load it using pickle
        with open(model_path, 'rb') as model_file:
            model = pickle.load(model_file)
        return model
    
    # Method to make predictions using the loaded model
    def predict(self, data):
        # Log the prediction process
        logging.info(f"Making prediction for: {data}")
        # Perform the prediction using the loaded model
        prediction = self.model.predict(data)
        return prediction
