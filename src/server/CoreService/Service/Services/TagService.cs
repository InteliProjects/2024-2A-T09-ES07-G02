using CoreService.Domain.Entities;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Service.Interfaces;

namespace CoreService.Service.Services
{
    public class TagService : ITagService
    {
        private readonly ITagRepository _tagRepository;

        public TagService(ITagRepository tagRepository)
        {
            _tagRepository = tagRepository;
        }

        public async Task<IEnumerable<Tag>> GetAllTagAsync()
        {
            return await _tagRepository.GetAllTagAsync();
        }

        public async Task<Tag> GetTagByIdAsync(int id)
        {
            return await _tagRepository.GetTagByIdAsync(id);
        }

        public async Task CreateTagAsync(Tag tag)
        {
            await _tagRepository.CreateTagAsync(tag);
        }

        public async Task UpdateTagAsync(Tag tag)
        {
            await _tagRepository.UpdateTagAsync(tag);
        }

        public async Task DeleteTagAsync(int id)
        {
            await _tagRepository.DeleteTagAsync(id);
        }
    }
}
