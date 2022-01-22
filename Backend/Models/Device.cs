using System.ComponentModel.DataAnnotations;

namespace Backend.Models
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
            Retired
        }

        [Key]
        public Guid Id { get; set; }
        [Required]
        public string Name { get; set; }
        [Required]
        public string Identifier { get; set; }
        [Required]
        public DeviceStatus Status { get; set; } = DeviceStatus.NotUsed;
        [Required]
        public DeviceType Type { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        public ICollection<PersonDevice> PersonDevice { get; set; }
    }
}
