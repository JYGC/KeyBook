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
        public string Name { get; set; }
        public string Identifier { get; set; }
        public DeviceStatus Status { get; set; } = DeviceStatus.NotUsed;
        public DeviceType Type { get; set; }
    }
}
