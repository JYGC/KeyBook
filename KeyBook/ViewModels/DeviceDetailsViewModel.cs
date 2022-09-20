using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class DeviceDetailsViewModel
    {
        public Device? Device { get; set; }
        public Guid? FromPersonDetailsPersonId { get; set; }
        public List<DeviceActivityHistory>? DeviceActivityHistoryList { get; set; }
        public Dictionary<int, string> DeviceTypes { get; set; }
        public Dictionary<Guid, string?>? PersonNamesTypes { get; set; }
        public bool IsNewDevice { get; set; }
    }
}
