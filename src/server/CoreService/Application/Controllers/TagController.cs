using AutoMapper;
using CoreService.Application.Interfaces.Messaging;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;
using CoreService.Infra.Messaging;
using CoreService.Service.Interfaces;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace CoreService.Application.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class TagController : Controller
    {
        private readonly ITagService _tagService;
        private readonly IMapper _mapper;
        private readonly IKafkaProducer _kafkaProducer;
        private readonly ILogger<TagController> _logger;

        public TagController(ITagService tagService, IMapper mapper, ILogger<TagController> logger, IKafkaProducer kafkaProducer)
        {
            _tagService = tagService;
            _mapper = mapper;
            _logger = logger;
            _kafkaProducer = kafkaProducer;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<TagDto>>> GetAllProject()
        {
            IEnumerable<Tag> tags = await _tagService.GetAllTagAsync();
            IEnumerable<TagDto> response = _mapper.Map<IEnumerable<TagDto>>(tags);

            return Ok(response);
        }

        [HttpGet("{id:int}")]
        public async Task<ActionResult<TagDto>> GetTagById(int id)
        {
            Tag tag = await _tagService.GetTagByIdAsync(id);
            TagDto response = _mapper.Map<TagDto>(tag);

            return Ok(response);
        }

        [HttpPost]
        public async Task<ActionResult<TagDto>> CreateTag([FromBody] TagAddViewModel tag)
        {
            if (!ModelState.IsValid)
            {
                _logger.LogWarning("Invalid model state for CreateTag request.");
                await _kafkaProducer.SendMessageAsync("core-service-logs", "Invalid model state for CreateTag request.");
                return BadRequest(ModelState);
            }

            try
            {
                Tag tagEntity = _mapper.Map<Tag>(tag);
                await _tagService.CreateTagAsync(tagEntity);

                TagDto response = _mapper.Map<TagDto>(tagEntity);

                _logger.LogInformation($"Tag created successfully with ID: {tagEntity.id}");
                await _kafkaProducer.SendMessageAsync("core-service-logs", $"Tag created successfully with ID: {tagEntity.id}");

                return CreatedAtAction(nameof(CreateTag), response);
            }
            catch (Exception ex)
            {
                _logger.LogError($"Error occurred while creating tag: {ex.Message}");
                await _kafkaProducer.SendMessageAsync("core-service-logs", $"Error occurred while creating tag: {ex.Message}");
                return StatusCode(500, "An error occurred while creating the tag.");
            }
        }

        [HttpPut("{id:int}")]
        public async Task<ActionResult<TagDto>> UpdateTag(int id, [FromBody] TagUpdateViewModel tag)
        {
            if (!ModelState.IsValid)
            {
                return BadRequest(ModelState);
            }

            if (id != tag.id)
            {
                return BadRequest("Id da requisição difere do Id do objeto.");
            }

            if (_tagService.GetTagByIdAsync(tag.id) == null)
            {
                return NotFound();
            }

            Tag tagEntity = _mapper.Map<Tag>(tag);

            await _tagService.UpdateTagAsync(tagEntity);
            return NoContent();
        }

        [HttpDelete("{id:int}")]
        public async Task<ActionResult> DeleteTag(int id)
        {
            if (_tagService.GetTagByIdAsync(id) == null)
            {
                return NotFound();
            }

            await _tagService.DeleteTagAsync(id);
            return NoContent();
        }
    }
}
