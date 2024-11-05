from fastapi import FastAPI, BackgroundTasks, HTTPException
from pymongo.mongo_client import MongoClient
from pymongo.server_api import ServerApi
from kafka import KafkaConsumer
import os
import json
from dotenv import load_dotenv
import logging
import time

load_dotenv()

app = FastAPI()

MONGODB_URI = os.getenv("ATLAS_URI")
KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BROKER")
KAFKA_TOPICS = [
    "tag-service-logs",
    "NLU-service-logs",
    "web-scrapping-logs",
    "SpeechService-logs",
    "core-service-logs"
]

# MongoDB Connection
client = MongoClient(MONGODB_URI, server_api=ServerApi('1'))
db = client.get_database('Cluster0')
log_collection = db['service_logs']

def deserialize_message(m):
    try:
        return json.loads(m.decode('utf-8')) if m else None
    except (json.JSONDecodeError, UnicodeDecodeError) as e:
        logging.error(f"Failed to deserialize message: {e}")
        return None

consumer = None

def create_kafka_consumer():
    """Tenta criar um consumidor Kafka com tratamento de exceção e timeout."""
    global consumer
    try:
        consumer = KafkaConsumer(
            *KAFKA_TOPICS,
            bootstrap_servers=[KAFKA_BOOTSTRAP_SERVERS],
            value_deserializer=lambda m: deserialize_message(m),
            auto_offset_reset='earliest',
            enable_auto_commit=True,
            request_timeout_ms=10000,  # Timeout de 10 segundos
            consumer_timeout_ms=10000  # Timeout de 10 segundos para receber mensagens
        )
        print(f"Conectado ao Kafka nos tópicos {KAFKA_TOPICS}")
    except Exception as e:
        logging.error(f"Erro ao conectar ao Kafka: {e}")
        consumer = None

is_running = False

def save_log_to_mongo(log):
    if log:
        try:
            log_collection.insert_one(log)
            print(f"Log inserido no MongoDB: {log}")
        except Exception as e:
            print(f"Erro ao inserir log no MongoDB: {e}")

def consume_logs():
    global is_running, consumer
    while is_running:
        if consumer is None:
            create_kafka_consumer()
        if consumer:
            try:
                for message in consumer:
                    log = message.value
                    if log is not None:
                        print(f"Log consumido: {log}")
                        save_log_to_mongo(log)
                    else:
                        print("Mensagem Kafka recebida é None ou vazia.")
            except Exception as e:
                logging.error(f"Erro ao consumir mensagens do Kafka: {e}")
                consumer = None
                time.sleep(5)  # Espera antes de tentar reconectar

@app.post("/start-consumer/")
def start_consumer(background_tasks: BackgroundTasks):
    global is_running
    if is_running:
        raise HTTPException(status_code=400, detail="O consumidor já está em execução.")
    is_running = True
    background_tasks.add_task(consume_logs)
    return {"message": "Consumidor iniciado."}

@app.post("/stop-consumer/")
def stop_consumer():
    global is_running, consumer
    if not is_running:
        raise HTTPException(status_code=400, detail="O consumidor não está em execução.")
    is_running = False
    if consumer:
        consumer.close()
    client.close()
    return {"message": "Consumidor interrompido."}

@app.get("/status/")
def get_status():
    status = "em execução" if is_running else "parado"
    return {"status": status}

