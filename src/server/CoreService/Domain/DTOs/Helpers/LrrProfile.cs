using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.Entities;

namespace CoreService.Domain.DTOs.Helpers
{
    public class LrrProfile :Profile
    {
        public LrrProfile()
        {
            CreateMap<Lrr, LrrDto>();
            CreateMap<LrrDto, Lrr>();
        }
    }
}
