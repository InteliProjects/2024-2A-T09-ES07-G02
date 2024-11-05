using CoreService.Domain.Entities;
using CoreService.Infra.Data.Context;
using CoreService.Infra.Data.Repository.Interfaces;
using Dapper;

namespace CoreService.Infra.Data.Repository.Repositories
{
    public class LrrTagRepository : ILrrTagRepository
    {
        private readonly DatabaseContext _context;

        public LrrTagRepository(DatabaseContext context)
        {
            _context = context;
        }

        public async Task<IEnumerable<LrrTag>> GetAllNewLrrTagAsync()
        {
            using var db = _context.Connect();
            var query = @"
                SELECT * 
                FROM lrr_tag 
            ";
            var lrr_tags = await db.QueryAsync<LrrTag>(query);
            return lrr_tags;
        }
    }
}

