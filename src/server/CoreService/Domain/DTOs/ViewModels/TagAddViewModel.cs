using System.ComponentModel.DataAnnotations;

namespace CoreService.Domain.DTOs.ViewModels
{
    public class TagAddViewModel
    {
        [Required(ErrorMessage = "O campo 'tag' é obrigatorio")]
        public string? tag { get; set; }

        [Required(ErrorMessage = "O campo 'description' é obrigatorio.")]
        public string? description { get; set; }
    }
}
