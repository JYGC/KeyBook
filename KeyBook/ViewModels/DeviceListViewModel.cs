using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class DeviceListViewModel
    {
        public List<Device> Devices { get; set; }
        public Dictionary<int, string> DeviceTypes { get; set; }
    }
}
