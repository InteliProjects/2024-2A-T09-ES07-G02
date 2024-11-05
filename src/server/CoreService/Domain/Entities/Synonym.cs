namespace CoreService.Domain.Entities
{
    public class Synonym
    {
        public int id { get; set; }
        public string? synonym { get; set; }
        public int? parent_tag_id { get; set; }
        public int? parent_subtag_id { get; set; }
        public DateTime created_at { get; set; }
        public DateTime updated_at { get; set; }

        public Tag? parent_tag { get; set; }
        public Subtag? parent_subtag { get; set; }
    }
}
