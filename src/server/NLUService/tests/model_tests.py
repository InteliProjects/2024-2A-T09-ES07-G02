import pytest
from unittest.mock import patch, MagicMock
from src.models.lh_model import LHModel  

# Test the LHModel class PREDICT method
class TestLHModel:
    @pytest.fixture

    # Mock for the model object
    def mock_model(self):
        mock_model = MagicMock()
        mock_model.predict.return_value = ['predicted_output']
        return mock_model

    @patch.object(LHModel, '_load_model', return_value=None)

    # Test the predict method
    def test_predict(self, mock_load_model, mock_model):
        # Arrange: Creates an instance of the LHModel class and sets the mock_model as the model attribute
        lh_model = LHModel()
        lh_model.model = mock_model  

        # Act: Calls the predict method with a sample input
        data = ['sample_input']
        prediction = lh_model.predict(data)
        
        # Assert: Checks if the model's predict method was called with the sample input and if the output is as expected
        mock_model.predict.assert_called_once_with(data)
        assert prediction == ['predicted_output']
