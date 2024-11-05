using CoreService.Domain.Entities;

namespace CoreService.Service.Interfaces
{
    public interface ILrrTagService
    {
        Task<IEnumerable<LrrTag>> GetAllNewLrrTagAsync();
    }
}

