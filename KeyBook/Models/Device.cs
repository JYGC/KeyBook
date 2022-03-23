using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class Device
    {
        public enum DeviceType
        {
            Fob,
            Key,
            Remote,
            RoomKey,
            MailKey
        }

        public enum DeviceStatus
        {
            NotUsed,
            WithManager,
            Used,
            Lost,
            Damaged,
            Retired,
            Stolen
        }

        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public string? Name { get; set; }
        [Required]
        public string? Identifier { get; set; }
        [Required]
        public DeviceStatus Status { get; set; } = DeviceStatus.NotUsed;
        [Required]
        public DeviceType Type { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public Guid UserId { get; set; }
        public virtual User? User { get; set; }
        public PersonDevice? PersonDevice { get; set; }
        public virtual ICollection<DeviceHistory> DeviceHistories { get; set; } = new List<DeviceHistory>();
    }
}
