using CoreService.Application.Interfaces.Messaging;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json.Linq;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

namespace CoreService.Application.Controllers
{
    [Route("api/webhook/[controller]")]
    [ApiController]
    public class ErrorController : ControllerBase
    {
        private readonly IKafkaConsumer _kafkaConsumer;
        private readonly HttpClient _httpClient;

        private const string WebhookUrl = "https://webhook.site/e03e6b60-46ea-4796-9fe5-507310cfc529";

        public ErrorController(IKafkaConsumer kafkaConsumer)
        {
            _kafkaConsumer = kafkaConsumer;
            _httpClient = new HttpClient(); // Ensure HttpClient is properly instantiated
        }

        [HttpPost]
        public async Task<IActionResult> ProcessErrors()
        {
            // Consume messages from the Kafka topic "web-scrapping-logs"
            var logMessage = await _kafkaConsumer.ConsumeMessageAsync("web-scrapping-logs");

            // Parse the JSON received from Kafka
            var logJson = JObject.Parse(logMessage);
            var errorMessage = logJson["error"]?.ToString();

            // If no error message is found, return a NoContent response
            if (string.IsNullOrEmpty(errorMessage))
            {
                return NoContent();
            }

            // Send the error message to the webhook
            var content = new StringContent(errorMessage, Encoding.UTF8, "application/json");
            var response = await _httpClient.PostAsync(WebhookUrl, content);

            // Check the response from the webhook
            if (response.IsSuccessStatusCode)
            {
                return Ok("Error message successfully sent to webhook.");
            }
            else
            {
                return StatusCode((int)response.StatusCode, "Failed to send error message to webhook.");
            }
        }
    }
}
