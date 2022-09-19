using KeyBook.Models;

namespace KeyBook.ViewModels
{
    public class PersonDetailsViewModel
    {
        public Person? Person { get; set; }
        public Dictionary<int, string> PersonTypes { get; set; }
        public bool IsNewPerson { get; set; }
    }
}
