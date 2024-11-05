using CoreService.Domain.Entities;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Service.Interfaces;

namespace CoreService.Service.Services
{
    public class LrrTagService : ILrrTagService
    {
        private readonly ILrrTagRepository _lrrTagRepository;

        public LrrTagService(ILrrTagRepository lrrTagRepository)
        {
            _lrrTagRepository = lrrTagRepository;
        }

        public async Task<IEnumerable<LrrTag>> GetAllNewLrrTagAsync()
        {
            return await _lrrTagRepository.GetAllNewLrrTagAsync();
        }
    }
}
