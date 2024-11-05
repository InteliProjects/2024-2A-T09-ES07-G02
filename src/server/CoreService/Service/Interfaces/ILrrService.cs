using CoreService.Domain.Entities;

namespace CoreService.Service.Interfaces
{
    public interface ILrrService
    {
        Task<IEnumerable<Lrr>> GetLatestLrrAsync();
        Task<IEnumerable<Lrr>> GetLrrByIntentAsync(string intent, string[] entity);
    }
}
