from kafka import KafkaConsumer
from kafka.errors import KafkaError
from typing import List, Callable
import json
import logging

# Set up logging configuration
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Class definition for KafkaListener
class KafkaListener:
    # Constructor method
    def __init__(self, topic: str, bootstrap_servers: List[str], group_id: str, auto_offset_reset: str = 'earliest'):
        # Initialize Kafka consumer configuration
        self.topic = topic
        self.bootstrap_servers = bootstrap_servers
        self.group_id = group_id
        self.auto_offset_reset = auto_offset_reset
        # Create a KafkaConsumer instance with the specified configurations
        self.consumer = KafkaConsumer(
            self.topic,
            bootstrap_servers=self.bootstrap_servers,
            group_id=self.group_id,
            auto_offset_reset=self.auto_offset_reset,
            enable_auto_commit=True,
            # Deserialize JSON messages from Kafka
            value_deserializer=lambda x: json.loads(x.decode('utf-8'))
        )

    # Method to start listening for messages on the specified topic
    def start_listening(self, callback: Callable[[dict], None]):
        """
        Starts consuming messages and applies the provided callback to each received message.
        
        :param callback: Function to be executed for each consumed message.
        """
        logger.info(f"Listener started for topic: {self.topic}")
        try:
            # Iterate over incoming messages and apply the callback function
            for message in self.consumer:
                logger.info(f"Message received: {message.value}")
                callback(message.value)
        except KafkaError as e:
            # Handle any Kafka-related errors
            logger.error(f"Error in Kafka Listener: {str(e)}")
        finally:
            # Ensure the consumer is properly closed
            self.consumer.close()
