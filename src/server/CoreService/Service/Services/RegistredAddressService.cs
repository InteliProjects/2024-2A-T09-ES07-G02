using System.Collections.Generic;
using System.Threading.Tasks;
using CoreService.Domain.Entities;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Service.Interfaces;

namespace CoreService.Service.Services
{
    public class RegistredAddressService : IRegistredAddressService
    {
        private readonly IRegistredAddressRepository _registredAddressRepository;

        public RegistredAddressService(IRegistredAddressRepository registredAddressRepository)
        {
            _registredAddressRepository = registredAddressRepository;
        }

        public async Task<IEnumerable<RegistredAddress>> GetAllRegistredAddressAsync()
        {
            return await _registredAddressRepository.GetAllRegistredAddressAsync();
        }

        public async Task AddRegistredAddressAsync(RegistredAddress registredAddress)
        {
            await _registredAddressRepository.AddRegistredAddressAsync(registredAddress);
        }
    }
}