using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;

namespace CoreService.Domain.DTOs.Helpers
{
    public class LrrTagProfile : Profile
    {
        public LrrTagProfile()
        {
            CreateMap<LrrTag, LrrTagDto>();
            CreateMap<LrrTagDto, LrrTag>();
        }
    }
}
