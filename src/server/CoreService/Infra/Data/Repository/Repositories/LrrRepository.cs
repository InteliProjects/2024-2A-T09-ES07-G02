using CoreService.Domain.Entities;
using CoreService.Infra.Data.Context;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Service.Interfaces;
using Dapper;

namespace CoreService.Infra.Data.Repository.Repositories
{
    public class LrrRepository : ILrrRepository
    {
        private readonly DatabaseContext _context;
        private readonly IIntentTagService _intentTagService;

        public LrrRepository(DatabaseContext context, IIntentTagService intentTagService)
        {
            _context = context;
            _intentTagService = intentTagService;
        }

        public async Task<IEnumerable<Lrr>> GetLatestLrrAsync(DateTime yesterday)
        {
            using var db = _context.Connect();

            var parameters = new
            {
                Yesterday = yesterday
            };

            IEnumerable<Lrr> lrrs = await db.QueryAsync<Lrr>(
                @"SELECT * FROM lrr
                  WHERE created_at >= @Yesterday",
                parameters
            );

            if (!lrrs.Any())
            {
                throw new KeyNotFoundException("Nenhum LRR encontrado.");
            }

            return lrrs;
        }

        public async Task<IEnumerable<Lrr>> GetLrrByIntentAsync(string intent)
        {
            using var db = _context.Connect();

            var tags = _intentTagService.GetTagsByIntent(intent);

            if (tags.Length == 0)
            {
                return Enumerable.Empty<Lrr>();
            }

            var parameters = new
            {
                Tags = tags
            };

            IEnumerable<Lrr> lrrs = await db.QueryAsync<Lrr>(
                @"SELECT lrr.*,
                        tag.tag AS tags
                FROM lrr
                JOIN lrr_tag ON lrr.id = lrr_tag.lrr_id
                JOIN tag ON lrr_tag.tag_id = tag.id
                WHERE tag.tag = ANY(@Tags)",
                parameters
            );


            if (!lrrs.Any())
            {
                throw new KeyNotFoundException("Nenhum LRR encontrado para a intenção informada.");
            }

            return lrrs;
        }

        public async Task<IEnumerable<Lrr>> GetLrrByIntentAndEntityAsync(string intent, string[] entities)
        {
            using var db = _context.Connect();

            var tags = _intentTagService.GetTagsByIntent(intent);

            if (tags == null || !tags.Any())
            {
                throw new KeyNotFoundException("Nenhuma tag encontrada para a intenção informada.");
            }

            if (entities == null || entities.Length == 0)
            {
                return Enumerable.Empty<Lrr>();
            }

            var parameters = new
            {
                Tags = tags,
                Entities = entities
            };

            IEnumerable<Lrr> lrrs = await db.QueryAsync<Lrr>(
                @"SELECT DISTINCT lrr.* 
                  FROM lrr
                  JOIN lrr_tag ON lrr.id = lrr_tag.lrr_id
                  JOIN tag ON lrr_tag.tag_id = tag.id
                  JOIN subtag ON lrr_tag.subtag_id = subtag.id
                  WHERE tag.tag = ANY(@Tags)
                  AND subtag.name = ANY(@Entities)",
                parameters);

            if (!lrrs.Any())
            {
                throw new KeyNotFoundException("Nenhum LRR encontrado para a intenção e as entidades informadas.");
            }

            return lrrs;
        }
    }
}