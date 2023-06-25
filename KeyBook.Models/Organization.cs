using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class Organization
    {
        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public string Name { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        public virtual ICollection<User> Users { get; set; } = new List<User>();
        public virtual ICollection<Person> Persons { get; set; } = new List<Person>();
        public virtual ICollection<Device> Devices { get; set; } = new List<Device>();
    }
}
