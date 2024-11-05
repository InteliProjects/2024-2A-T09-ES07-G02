using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;

namespace CoreService.Domain.DTOs.Helpers
{
    public class RegistredAddressProfile : Profile
    {
        public RegistredAddressProfile()
        {
            CreateMap<RegistredAddress, RegistredAddressDto>();
            CreateMap<RegistredAddressDto, RegistredAddress>();
            CreateMap<RegistredAddress, RegistredAddressAddViewModel>();
            CreateMap<RegistredAddressAddViewModel, RegistredAddress>();
        }
    }
}
