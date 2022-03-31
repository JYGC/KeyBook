using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class DevicePersonDetailsPersonIdViewModel
    {
        public Device Device { get; set; }
        public Guid? FromPersonDetailsPersonId { get; set; }
        public List<DeviceActivityHistory> DeviceActivityHistoryList { get; set; }
    }
}
