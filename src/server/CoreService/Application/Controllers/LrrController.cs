using AutoMapper;
using CoreService.Application.Interfaces.Messaging;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.Entities;
using CoreService.Service.Interfaces;
using Microsoft.AspNetCore.Mvc;

namespace CoreService.Application.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class LrrController : ControllerBase
    {
        private readonly ILrrService _lrrService;
        private readonly IMapper _mapper;
        private readonly IKafkaProducer _kafkaProducer;

        public LrrController(ILrrService lrrService, IMapper mapper, IKafkaProducer kafkaProducer)
        {
            _lrrService = lrrService;
            _mapper = mapper;
            _kafkaProducer = kafkaProducer;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<LrrDto>>> GetAllLastLrr()
        {
            try
            {
                
                IEnumerable<Lrr> lrrs = await _lrrService.GetLatestLrrAsync();
                IEnumerable<LrrDto> response = _mapper.Map<IEnumerable<LrrDto>>(lrrs);

                string logMessage = $"Successfully retrieved {lrrs.Count()} LRR records at {DateTime.UtcNow}.";
                await _kafkaProducer.SendMessageAsync("core-service-logs", logMessage);

                return Ok(response);
            }
            catch (Exception ex)
            {
                string errorMessage = $"Error retrieving LRR records: {ex.Message}";
                await _kafkaProducer.SendMessageAsync("core-service-logs", errorMessage);

                return StatusCode(500, "An error occurred while retrieving LRR records.");
            }
        }
    }
}
