using Confluent.Kafka;
using CoreService.Application.Interfaces.Messaging;

namespace CoreService.Infra.Messaging
{
    public class KafkaProducer : IKafkaProducer
    {
        private readonly IProducer<Null, string> _producer;
        private readonly ILogger<KafkaProducer> _logger;

        public KafkaProducer(IConfiguration configuration, ILogger<KafkaProducer> logger)
        {
            var config = new ProducerConfig
            {
                BootstrapServers = configuration.GetValue<string>("Kafka:BootstrapServers"),
                Acks = Acks.All
            };

            _producer = new ProducerBuilder<Null, string>(config).Build();
            _logger = logger;
        }

        public async Task SendMessageAsync(string topic, string message)
        {
            try
            {
                var result = await _producer.ProduceAsync(topic, new Message<Null, string> { Value = message });
                _logger.LogInformation($"Message '{result.Value}' sent to '{result.TopicPartitionOffset}'");
            }
            catch (ProduceException<Null, string> ex)
            {
                _logger.LogError($"Delivery failed: {ex.Error.Reason}");
            }
        }

        public void Dispose()
        {
            _producer?.Dispose();
        }
    }
}
