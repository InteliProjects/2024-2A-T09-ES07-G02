using CoreService.Domain.Entities;

namespace CoreService.Infra.Data.Repository.Interfaces
{
    public interface ITagRepository
    {
        Task<IEnumerable<Tag>> GetAllTagAsync();
        Task<Tag> GetTagByIdAsync(int id);
        Task CreateTagAsync(Tag tag);
        Task UpdateTagAsync(Tag tag);
        Task DeleteTagAsync(int id);
    }
}
