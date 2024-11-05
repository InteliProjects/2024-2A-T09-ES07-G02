using AutoMapper;
using CoreService.Domain.DTOs.Responses;
using CoreService.Domain.DTOs.ViewModels;
using CoreService.Domain.Entities;

namespace CoreService.Domain.DTOs.Helpers
{
    public class TagProfile : Profile
    {
        public TagProfile()
        {
            CreateMap<Tag, TagDto>();
            CreateMap<TagDto, Tag>();
            CreateMap<Tag, TagAddViewModel>();
            CreateMap<TagAddViewModel, Tag>();
            CreateMap<Tag, TagUpdateViewModel>();
            CreateMap<TagUpdateViewModel, Tag>();
        }
    }
}
