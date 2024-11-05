namespace CoreService.Application.Interfaces.Messaging
{
    public interface IKafkaConsumer
    {
        Task<string> ConsumeMessageAsync(string topic);
    }
}
