import logging
import os
import requests

from src.services.kafka_listener import KafkaListener
from src.services.kafka_publisher import KafkaPublisher
from src.utils.preprocess import TextProcessor
from src.utils.vectorizing import WordVectorizer
from src.models.lh_model import LHModel
from src.utils.kafka_logging_handler import KafkaLoggingHandler

log_folder = 'logs'
if not os.path.exists(log_folder):
    os.makedirs(log_folder)

log_filename = os.path.join(log_folder, 'service.log')

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler(log_filename),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

kafka_log_handler = KafkaLoggingHandler(
    kafka_servers=['kafka:29092'],
    kafka_topic='NLU-service-logs'
)
kafka_log_handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
logger.addHandler(kafka_log_handler)

def download_model():
    model_url = "https://law-hunter.s3.amazonaws.com/Word2VecTrained/skip_s300.txt"
    model_filename = model_url.split('/')[-1]
    folder = 'src/models/w2v_models'
    file_path = os.path.join(folder, model_filename)

    if not os.path.exists(folder):
        os.makedirs(folder)
    
    if not os.path.exists(file_path):
        logger.info(f"Downloading model from {model_url}")
        response = requests.get(model_url)
        response.raise_for_status()  
        with open(file_path, 'wb') as f:
            f.write(response.content)
        logger.info(f"Model downloaded and saved to {file_path}")
    else:
        logger.info(f"Model already downloaded at {file_path}")

    return file_path

model_path = download_model()
preprocess = TextProcessor()
vectorizer = WordVectorizer(model_path)
model = LHModel()

# Process the incoming message and publish the prediction
def process_message(message: dict, publisher: KafkaPublisher):
    logger.info(f"Processando mensagem: {message}")
    
    data = preprocess.process(message['text'])
    logger.info(f"Texto processado: {data}")

    mean_vector = vectorizer.average_vector(data["lemmatized"])
    
    prediction = model.predict(mean_vector)
    logger.info(f"Previsão: {prediction}")
    
    # Publish the prediction to another topic
    publisher.publish_message('user-predictions', {"prediction": prediction.tolist()})
    
if __name__ == "__main__":
    # Set up the KafkaPublisher
    publisher = KafkaPublisher(bootstrap_servers=['kafka:29092'])

    # Configuração do Listener
    listener = KafkaListener(
        topic='user-messages',
        bootstrap_servers=['kafka:29092'],
        group_id='nlu-consumer'
    )

    # Define the callback function for incoming messages
    def on_message(message):
        process_message(message, publisher)

    listener.start_listening(on_message)
    publisher.close()  # Close the publisher when done
