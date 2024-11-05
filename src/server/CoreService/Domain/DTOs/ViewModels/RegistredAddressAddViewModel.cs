using System.ComponentModel.DataAnnotations;

namespace CoreService.Domain.DTOs.ViewModels
{
    public class RegistredAddressAddViewModel
    {
        [Required(ErrorMessage = "O campo 'address' é obrigatorio")]
        public string? address { get; set; }
    }
}
