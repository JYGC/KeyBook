using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class PersonDevice
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public virtual Person Person { get; set; }
        [Required]
        public virtual Device Device { get; set; }
        [Required]
        public bool IsNotHave { get; set; } = false;
        [Required]
        public bool IsDeleted { get; set; } = false;
        public virtual ICollection<PersonDeviceHistory> PersonDeviceHistory { get; set; }
    }
}
