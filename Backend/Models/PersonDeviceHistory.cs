using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class PersonDeviceHistory
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public Guid PersonId { get; set; }
        [Required]
        public Guid DeviceId { get; set; }
        [Required]
        public bool IsNotHave { get; set; } = false;
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public string? Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public Guid PersonDeviceId { get; set; }
        public virtual PersonDevice? PersonDevice { get; set; }
    }
}
