using CoreService.Domain.Entities;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Service.Interfaces;

namespace CoreService.Service.Services
{
    public class LrrService : ILrrService
    {
        private readonly ILrrRepository _lrrRepository;

        public LrrService(ILrrRepository lrrRepository)
        {
            _lrrRepository = lrrRepository;
        }

        public async Task<IEnumerable<Lrr>> GetLatestLrrAsync()
        {
            DateTime yesterday = DateTime.Now.AddDays(-1);

            return await _lrrRepository.GetLatestLrrAsync(yesterday);
        }

        public async Task<IEnumerable<Lrr>> GetLrrByIntentAsync(string intent, string[] entities)
        {
            if (entities != null && entities.Length > 0)
            {
                return await _lrrRepository.GetLrrByIntentAndEntityAsync(intent, entities);
            }

            return await _lrrRepository.GetLrrByIntentAsync(intent);
        }
    }
}
