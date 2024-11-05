using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;
using CoreService.Service.Interfaces;
using Microsoft.AspNetCore.Mvc;

namespace CoreService.Application.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class LrrTagController : Controller
    {
        private readonly ILrrTagService _lrrTagService;
        private readonly IMapper _mapper;

        public LrrTagController(ILrrTagService lrrTagService, IMapper mapper)
        {
            _lrrTagService = lrrTagService;
            _mapper = mapper;
        }

        [HttpGet("/webhook/response")]
        public async Task<ActionResult<IEnumerable<LrrTagDto>>> GetAllNewLrrTag()
        {
            IEnumerable<LrrTag> lrr_tags = await _lrrTagService.GetAllNewLrrTagAsync();
            IEnumerable<LrrTagDto> response = _mapper.Map<IEnumerable<LrrTagDto>>(lrr_tags);

            return Ok(response);
        }
    }
}

