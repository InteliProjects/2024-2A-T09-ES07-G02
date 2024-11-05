using Confluent.Kafka;
using CoreService.Application.Interfaces.Messaging;
using Microsoft.Extensions.Logging;

namespace CoreService.Infra.Messaging
{
    public class KafkaConsumer : IKafkaConsumer
    {
        private readonly IConsumer<Null, string> _consumer;
        private readonly ILogger<KafkaConsumer> _logger;

        public KafkaConsumer(IConfiguration configuration, ILogger<KafkaConsumer> logger)
        {
            var config = new ConsumerConfig
            {
                BootstrapServers = configuration["Kafka:BootstrapServers"],
                GroupId = "nlu-consumer-group", // O group id do consumidor
                AutoOffsetReset = AutoOffsetReset.Earliest
            };

            _consumer = new ConsumerBuilder<Null, string>(config).Build();
            _logger = logger;
        }

        public async Task<string?> ConsumeMessageAsync(string topic)
        {
            _consumer.Subscribe(topic);

            try
            {
                var consumeResult = _consumer.Consume();
                _logger.LogInformation($"Mensagem recebida do tópico {topic}: {consumeResult.Message.Value}");
                return consumeResult.Message.Value;
            }
            catch (ConsumeException ex)
            {
                _logger.LogError($"Erro ao consumir mensagem: {ex.Error.Reason}");
                return null;
            }
            finally
            {
                _consumer.Close();
            }
        }

        public void Dispose()
        {
            _consumer?.Dispose();
        }
    }
}
