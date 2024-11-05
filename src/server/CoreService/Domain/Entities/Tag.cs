namespace CoreService.Domain.Entities
{
    public class Tag
    {
        public int id { get; set; }
        public string? tag { get; set; }
        public string? description { get; set; }
        public DateTime created_at { get; set; }
        public DateTime updated_at { get; set; }

        public ICollection<Subtag>? subtags { get; set; }
        public ICollection<Synonym>? synonyms { get; set; }
    }
}
