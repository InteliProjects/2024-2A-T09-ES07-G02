using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;
using CoreService.Service.Interfaces;
using Microsoft.AspNetCore.Mvc;

namespace CoreService.Application.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class RegistredAddressController : ControllerBase
    {
        private readonly IRegistredAddressService _registredAddressService;
        private readonly ILrrTagService _lrrTagService;
        private readonly IMapper _mapper;

        public RegistredAddressController(IRegistredAddressService registredAddressService, ILrrTagService lrrTagService, IMapper mapper)
        {
            _registredAddressService = registredAddressService;
            _lrrTagService = lrrTagService;
            _mapper = mapper;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<RegistredAddressDto>>> GetAllRegisters()
        {
            IEnumerable<RegistredAddress> registredAddress = await _registredAddressService.GetAllRegistredAddressAsync();
            var registredAddressDto = _mapper.Map<IEnumerable<RegistredAddressDto>>(registredAddress);
            return Ok(registredAddressDto);
        }

        [HttpPost]
        public async Task<ActionResult> AddRegistredAddress([FromBody] RegistredAddressAddViewModel registredAddressAddViewModel)
        {
            var registredAddress = _mapper.Map<RegistredAddress>(registredAddressAddViewModel);
            await _registredAddressService.AddRegistredAddressAsync(registredAddress);
            var registredAddressDto = _mapper.Map<RegistredAddressDto>(registredAddress);
            return CreatedAtAction(nameof(GetAllRegisters), new { address = registredAddress.address }, registredAddressDto);
        }
        [HttpGet("serve-webhook")]
        public async Task<ActionResult> GetWebhookData()
        {
            var registredAddresses = await _registredAddressService.GetAllRegistredAddressAsync();
            var result = new List<object>();

            foreach (var address in registredAddresses)
            {
                var lrrTags = await _lrrTagService.GetAllNewLrrTagAsync();
                result.Add(new { Address = address, LrrTags = lrrTags });
            }

            return Ok(result);
        }
    }
}