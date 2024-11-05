namespace CoreService.Domain.Entities
{
    public class Lrr
    {
        public int id { get; set; }
        public string? agency { get; set; }
        public string? rule { get; set; }
        public string? segment { get; set; }
        public string? s3_url { get; set; }
        public DateTime created_at { get; set; }
        public DateTime updated_at { get; set; }

        public string? tags { get; set; }
    }
}
