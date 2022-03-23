using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class PersonDevice
    {
        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public Guid PersonId { get; set; }
        public virtual Person? Person { get; set; }
        [Required]
        public Guid DeviceId { get; set; }
        public virtual Device? Device { get; set; }
        [Required]
        public bool IsNoLongerHas { get; set; } = false;
        [Required]
        public bool IsDeleted { get; set; } = false;
    }
}
