using CoreService.Application.Interfaces.Messaging;
using CoreService.Service.Interfaces;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace CoreService.Application.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class SearchController : ControllerBase
    {
        private readonly IKafkaProducer _kafkaProducer;
        private readonly IKafkaConsumer _kafkaConsumer;
        private readonly ILrrService _lrrService;

        public SearchController(IKafkaProducer kafkaProducer, IKafkaConsumer kafkaConsumer, ILrrService lrrService)
        {
            _kafkaProducer = kafkaProducer;
            _kafkaConsumer = kafkaConsumer;
            _lrrService = lrrService;
        }

        [HttpPost]
        public async Task<IActionResult> Search([FromBody] JObject input)
        {
            if (input == null)
            {
                return BadRequest("Request body is required.");
            }

            // Serialize the request body and send it to the Kafka topic "user-messages"
            var jsonString = JsonConvert.SerializeObject(input);
            await _kafkaProducer.SendMessageAsync("user-messages", jsonString);

            // Consume the predicted intent from the Kafka topic "user-predictions"
            var predictionMessage = await _kafkaConsumer.ConsumeMessageAsync("user-predictions");

            // Check if predictionMessage is null or empty
            if (string.IsNullOrEmpty(predictionMessage))
            {
                return BadRequest("No message was received from Kafka.");
            }

            // Parse the JSON received from Kafka to extract the predicted intent and entities
            JObject predictionJson;
            try
            {
                predictionJson = JObject.Parse(predictionMessage);
            }
            catch (JsonReaderException ex)
            {
                return BadRequest($"Invalid JSON format: {ex.Message}");
            }

            var intent = predictionJson["prediction"]?.FirstOrDefault()?.ToString();

            // Parse "entities" as an array of strings
            var entitiesArray = predictionJson["entities"]?.ToObject<string[]>();

            // If no intent is found, return a BadRequest response
            if (string.IsNullOrEmpty(intent))
            {
                return BadRequest("No intent was predicted.");
            }

            // If no entities were provided, create an empty array
            entitiesArray ??= Array.Empty<string>();

            // Fetch relevant data from the database based on the predicted intent and entities
            var lrrResults = await _lrrService.GetLrrByIntentAsync(intent, entitiesArray);

            // Return the fetched data in the response
            return Ok(lrrResults);
        }

        [HttpPost("search-direct")]
        public async Task<IActionResult> SearchDirect([FromBody] JObject input)
        {
            if (input == null)
            {
                return BadRequest("Request body is required.");
            }

            // Extract the "prediction" and "entities" fields directly from the input JSON
            var intent = input["prediction"]?.ToString();

            // Parse "entities" as an array of strings
            var entitiesArray = input["entities"]?.ToObject<string[]>();

            // If no intent is provided, return an error
            if (string.IsNullOrEmpty(intent))
            {
                return BadRequest("Prediction (intent) is required.");
            }

            // If no entities are provided, initialize an empty array
            entitiesArray ??= Array.Empty<string>();

            // Fetch relevant data from the database based on the provided intent and entities
            var lrrResults = await _lrrService.GetLrrByIntentAsync(intent, entitiesArray);

            // Return the fetched data
            return Ok(lrrResults);
        }

    }
}
