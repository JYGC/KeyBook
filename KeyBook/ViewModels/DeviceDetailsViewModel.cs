using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class DeviceDetailsViewModel
    {
        public Device? Device { get; set; }
        public Guid? FromPersonDetailsPersonId { get; set; }
        public bool IsNewDevice { get; set; }
    }
}
