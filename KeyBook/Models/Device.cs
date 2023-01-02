using System.ComponentModel;
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
            [Description("Room key")]
            RoomKey,
            [Description("Mailbox key")]
            MailboxKey
        }

        public enum DeviceDefunctReason
        {
            None,
            Lost,
            Damaged,
            Retired,
            Stolen
        }

        [Key]
        public Guid Id { get; set; }
        [Required]
        public string? Name { get; set; }
        [Required]
        public string? Identifier { get; set; }
        [Required]
        public DeviceDefunctReason DefunctReason { get; set; } = DeviceDefunctReason.None;
        [Required]
        public DeviceType Type { get; set; }
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public Guid OrganizationId { get; set; }
        public virtual Organization? Organization { get; set; }
        public PersonDevice? PersonDevice { get; set; }
        public virtual ICollection<DeviceHistory> DeviceHistories { get; set; } = new List<DeviceHistory>();
    }
}
