using System.Collections.Generic;
using System.Threading.Tasks;
using CoreService.Domain.Entities;

namespace CoreService.Infra.Data.Repository.Interfaces
{
    public interface IRegistredAddressRepository
    {
        Task<IEnumerable<RegistredAddress>> GetAllRegistredAddressAsync();
        Task AddRegistredAddressAsync(RegistredAddress registredAddress);
    }
}