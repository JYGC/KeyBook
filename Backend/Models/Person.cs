using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public enum PersonType
    {
        Tenant,
        Owner,
        Manager
    }

    public class Person
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
        public virtual ICollection<PersonDevice> PersonDevice { get; set; }
        public virtual ICollection<PersonHistory> PersonHistory { get; set; }
    }
}
