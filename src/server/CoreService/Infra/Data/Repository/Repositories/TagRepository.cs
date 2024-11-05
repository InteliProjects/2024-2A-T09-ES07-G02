using CoreService.Domain.Entities;
using CoreService.Infra.Data.Context;
using CoreService.Infra.Data.Repository.Interfaces;
using Dapper;

namespace CoreService.Infra.Data.Repository.Repositories
{
    public class TagRepository : ITagRepository
    {
        private readonly DatabaseContext _context;

        public TagRepository(DatabaseContext context)
        {
            _context = context;
        }

        public async Task<IEnumerable<Tag>> GetAllTagAsync()
        {
            using var db = _context.Connect();
            IEnumerable<Tag> tags = await db.QueryAsync<Tag>(
                @"SELECT * FROM tag"
            );

            if (tags == null)
            {
                throw new Exception("Projeto não encontrado.");
            }

            return tags;
        }

        public async Task<Tag> GetTagByIdAsync(int id)
        {
            using var db = _context.Connect();

            var parameters = new
            {
                Id = id
            };


            Tag? tag = await db.QueryFirstOrDefaultAsync<Tag>(
                @"SELECT * FROM tag 
                  WHERE id = @Id",
                parameters
            );

            if (tag == null)
            {
                throw new Exception("Tag não encontrado.");
            }

            return tag;
        }

        public async Task CreateTagAsync(Tag tag)
        {
            using var db = _context.Connect();

            var parameters = new
            {
                Tag = tag.tag,
                Description = tag.description
            };

            await db.ExecuteAsync(
                @"INSERT INTO tag (
                      tag
                    , description
                ) 
                  VALUES (
                      @Tag
                    , @Description      
                )",
                parameters
            );
        }

        public async Task UpdateTagAsync(Tag tag)
        {
            using var db = _context.Connect();

            var parameters = new
            {
                Id = tag.id,
                Tag = tag.tag,
                Description = tag.description
            };

            await db.ExecuteAsync(
                @"UPDATE tag
                  SET
                      tag         = @Tag
                    , description = @Description
                  WHERE id = @Id",
                parameters
            );
        }

        public async Task DeleteTagAsync(int id)
        {
            using var db = _context.Connect();

            var parameters = new
            {
                Id = id
            };

            await db.ExecuteAsync(
                @"DELETE FROM tag
                  WHERE id = @Id",
                parameters
            );
        }
    }
}
