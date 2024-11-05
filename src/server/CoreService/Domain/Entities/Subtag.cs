namespace CoreService.Domain.Entities
{
    public class Subtag
    {
        public int id { get; set; }
        public string? name { get; set; }
        public int? parent_tag_id { get; set; }
        public int? parent_subtag_id { get; set; }
        public int? child_tag_id { get; set; }
        public int? child_subtag_id { get; set; }
        public DateTime created_at { get; set; }
        public DateTime updated_at { get; set; }

        public Tag? parent_tag { get; set; }
        public Subtag? parent_subtag { get; set; }

        public Tag? child_tag { get; set; }
        public Subtag? child_subtag { get; set; }

        public ICollection<Synonym>? synonyms { get; set; }
    }
}
