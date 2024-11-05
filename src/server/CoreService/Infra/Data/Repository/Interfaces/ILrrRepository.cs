using CoreService.Domain.Entities;

namespace CoreService.Infra.Data.Repository.Interfaces
{
    public interface ILrrRepository
    {
        Task<IEnumerable<Lrr>> GetLatestLrrAsync(DateTime yesterday);
        Task<IEnumerable<Lrr>> GetLrrByIntentAsync(string intent);
        Task<IEnumerable<Lrr>> GetLrrByIntentAndEntityAsync(string intent, string[] entities);
    }
}
