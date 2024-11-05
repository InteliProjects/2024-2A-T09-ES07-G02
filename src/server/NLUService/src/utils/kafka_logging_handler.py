import logging
from kafka import KafkaProducer
import json

class KafkaLoggingHandler(logging.Handler):
    def __init__(self, kafka_servers, kafka_topic):
        logging.Handler.__init__(self)
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_servers,
            value_serializer=lambda v: json.dumps(v).encode('utf-8')
        )
        self.kafka_topic = kafka_topic

    def emit(self, record):
        try:
            log_entry = self.format(record)
            self.producer.send(self.kafka_topic, {'log': log_entry})
        except Exception as e:
            print(f"Error sending log to Kafka: {e}")

