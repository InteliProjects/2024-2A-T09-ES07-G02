using CoreService.Application.Interfaces.Messaging;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json;

namespace CoreService.Application.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class MessagingController : ControllerBase
    {
        private readonly IKafkaProducer _kafkaProducer;

        public MessagingController(IKafkaProducer kafkaProducer)
        {
            _kafkaProducer = kafkaProducer;
        }

        [HttpPost("SendNLU")]
        public async Task<IActionResult> SendNLU([FromBody] object jsonBody)
        {
            try
            {
                var jsonString = JsonConvert.SerializeObject(jsonBody);

                await _kafkaProducer.SendMessageAsync("user-messages", jsonString);

                var successLog = $"Message sent to 'user-messages' topic at {DateTime.UtcNow}: {jsonString}";
                await _kafkaProducer.SendMessageAsync("core-service-logs", successLog);

                return Ok("Message sent successfully");
            }
            catch (Exception ex)
            {
                var errorLog = $"Error sending message: {ex.Message}";
                await _kafkaProducer.SendMessageAsync("core-service-logs", errorLog);

                return StatusCode(500, "An error occurred while sending the message.");
            }
        }
    }
}
