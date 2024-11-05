    namespace CoreService.Domain.Entities
    {
        public class LrrTag
        {
            public int id { get; set; }
            public int? lrr_id { get; set; }
            public int? tag_id { get; set; }
            public int? subtag_id { get; set; }


            public Lrr? lrr { get; set; }
            public Tag? tag { get; set; }
            public Subtag? subtag { get; set; }
        }
    }
