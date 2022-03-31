using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class DeviceHistory
    {
        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public string? Name { get; set; }
        [Required]
        public string? Identifier { get; set; }
        [Required]
        public Device.DeviceStatus Status { get; set; }
        [Required]
        public Device.DeviceType Type { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public string? Description { get; set; }
        [Required]
        public DateTime RecordDateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public DateTime DateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public Guid DeviceId { get; set; }
        public virtual Device? Device { get; set; }
    }
}
