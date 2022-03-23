using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class PersonHistory
    {
        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public string? Name { get; set; }
        [Required]
        public bool IsGone { get; set; } = false;
        [Required]
        public PersonType Type { get; set; } = PersonType.Tenant;
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public string? Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public Guid PersonId { get; set; }
        public virtual Person? Person { get; set; }
    }
}
