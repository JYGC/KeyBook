using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class DeviceHistory
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public string? Name { get; set; }
        [Required]
        public string? Identifier { get; set; }
        [Required]
        public Device.DeviceStatus Status { get; set; } = Device.DeviceStatus.NotUsed;
        [Required]
        public Device.DeviceType Type { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public string? Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public Guid DeviceId { get; set; }
        public virtual Device? Device { get; set; }
    }
}
