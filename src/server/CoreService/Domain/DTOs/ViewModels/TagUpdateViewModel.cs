using System.ComponentModel.DataAnnotations;

namespace CoreService.Domain.DTOs.ViewModels
{
    public class TagUpdateViewModel
    {
        [Required(ErrorMessage = "O campo 'id' é obrigatorio.")]
        public int id { get; set; }
        public string? tag { get; set; }
        public string? description { get; set; }
    }
}
