from kafka import KafkaProducer
from kafka.errors import KafkaError
from typing import List
import json
import logging

# Set up logging configuration
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Class definition for KafkaPublisher
class KafkaPublisher:
    # Constructor method
    def __init__(self, bootstrap_servers: List[str]):
        # Initialize Kafka producer configuration
        self.bootstrap_servers = bootstrap_servers
        # Create a KafkaProducer instance with the specified configurations
        self.producer = KafkaProducer(
            bootstrap_servers=self.bootstrap_servers,
            # Serialize Python dict to JSON format before sending to Kafka
            value_serializer=lambda v: json.dumps(v).encode('utf-8')
        )

    # Method to publish a message to the specified topic
    def publish_message(self, topic: str, message: dict):
        """
        Publishes a message to the specified topic.

        :param topic: Topic where the message will be published.
        :param message: Message to be published (in dict format).
        """
        try:
            logger.info(f"Sending message to topic {topic}: {message}")
            # Send the message to the specified Kafka topic
            future = self.producer.send(topic, value=message)
            # Block until a single message is sent (or timeout is reached)
            result = future.get(timeout=10)
            logger.info(f"Message successfully published: {result}")
        except KafkaError as e:
            # Handle any Kafka-related errors during message publishing
            logger.error(f"Error while publishing message: {str(e)}")

    # Method to safely close the Kafka producer
    def close(self):
        """Safely closes the Kafka producer."""
        # Flush any pending messages and close the producer
        self.producer.flush()
        self.producer.close()
