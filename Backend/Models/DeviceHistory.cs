using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class DeviceHistory : Device
    {
        public Device Device { get; set; }
        [Required]
        public string Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.Now;
    }
}
