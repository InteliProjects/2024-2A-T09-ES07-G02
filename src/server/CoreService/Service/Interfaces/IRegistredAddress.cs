using System.Collections.Generic;
using System.Threading.Tasks;
using CoreService.Domain.Entities;

namespace CoreService.Service.Interfaces
{
    public interface IRegistredAddressService
    {
        Task<IEnumerable<RegistredAddress>> GetAllRegistredAddressAsync();
        Task AddRegistredAddressAsync(RegistredAddress registredAddress);
    }
}