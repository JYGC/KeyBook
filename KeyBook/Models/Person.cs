using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class Person
    {
        public enum PersonType
        {
            Tenant,
            Owner,
            Manager
        }

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
        public Guid UserId { get; set; }
        public virtual User? User { get; set; }
        public virtual ICollection<PersonDevice> PersonDevices { get; set; } = new List<PersonDevice>();
        public virtual ICollection<PersonHistory> PersonHistories { get; set; } = new List<PersonHistory>();
    }
}
