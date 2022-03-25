using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class DevicePersonListViewModel
    {
        public Device Device { get; set; }
        public List<Person> PersonList { get; set; }
        public Guid? FromPersonDetailsPersonId { get; set; }
    }
}
