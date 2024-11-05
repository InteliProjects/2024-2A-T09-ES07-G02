namespace CoreService.Application.Interfaces.Messaging
{
    public interface IKafkaProducer
    {
        Task SendMessageAsync(string topic, string message);
    }
}
