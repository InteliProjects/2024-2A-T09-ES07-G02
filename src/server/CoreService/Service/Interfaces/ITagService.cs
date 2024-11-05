using CoreService.Domain.Entities;

namespace CoreService.Service.Interfaces
{
    public interface ITagService
    {
        Task<IEnumerable<Tag>> GetAllTagAsync();
        Task<Tag> GetTagByIdAsync(int id);
        Task CreateTagAsync(Tag tag);
        Task UpdateTagAsync(Tag tag);
        Task DeleteTagAsync(int id);
    }
}
