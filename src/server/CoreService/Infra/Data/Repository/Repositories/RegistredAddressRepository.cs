using System.Collections.Generic;
using System.Threading.Tasks;
using CoreService.Domain.Entities;
using CoreService.Infra.Data.Context;
using CoreService.Infra.Data.Repository.Interfaces;
using Dapper;

namespace CoreService.Infra.Data.Repository.Repositories
{
    public class RegistredAddressRepository : IRegistredAddressRepository
    {
        private readonly DatabaseContext _context;

        public RegistredAddressRepository(DatabaseContext context)
        {
            _context = context;
        }

        public async Task<IEnumerable<RegistredAddress>> GetAllRegistredAddressAsync()
        {
            using var db = _context.Connect();
            var query = @"
                SELECT * 
                FROM registred_address 
            ";
            var registred_addresses = await db.QueryAsync<RegistredAddress>(query);
            return registred_addresses;
        }

        public async Task AddRegistredAddressAsync(RegistredAddress registredAddress)
        {
            using var db = _context.Connect();
            var query = @"
                INSERT INTO registred_address (Address) 
                VALUES (@Address)
            ";
            await db.ExecuteAsync(query, registredAddress);
        }
    }
}