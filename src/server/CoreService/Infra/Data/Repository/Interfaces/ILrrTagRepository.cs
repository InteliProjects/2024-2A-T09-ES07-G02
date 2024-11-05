using CoreService.Domain.Entities;

namespace CoreService.Infra.Data.Repository.Interfaces
{
    public interface ILrrTagRepository
    {
        Task<IEnumerable<LrrTag>> GetAllNewLrrTagAsync();
    }
}
