using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class PersonHistory
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public string Name { get; set; }
        [Required]
        public bool IsGone { get; set; } = false;
        [Required]
        public PersonType Type { get; set; } = PersonType.Tenant;
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public string Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.Now;
        [Required]
        public Person Person { get; set; }
    }
}
