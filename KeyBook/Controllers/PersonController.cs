using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Controllers
{
    public class PersonController : Controller
    {
        private readonly KeyBookDbContext _context;

        public PersonController(KeyBookDbContext context)
        {
            _context = context;
        }

        public IActionResult Index()
        {
            User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
            var personDeviceQuery = from person in _context.Set<Person>()
                                    from personDevice in _context.Set<PersonDevice>().Where(personDevice => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                    from device in _context.Set<Device>().Where(device => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                    where person.UserId == user.Id && !person.IsDeleted
                                    select new { person, personDevice, device };
            List<Person> personsWithDuplicates = personDeviceQuery.ToList().Select(pdq => pdq.person).ToList();
            Dictionary<Guid, Person> personsDict = new Dictionary<Guid, Person>();
            foreach (Person person in personsWithDuplicates)
            {
                if (!personsDict.ContainsKey(person.Id)) personsDict[person.Id] = person;
            }
            return View(new PersonListViewModel { Persons = personsDict.Values.ToList() });
        }

        public IActionResult New()
        {
            return View();
        }

        public ActionResult<Dictionary<int, string>> GetPersonTypes()
        {
            return Enum.GetValues(typeof(Person.PersonType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }

        public class NewPersonBindModel
        {
            public string Name { get; set; }
            public int Type { get; set; }
        }
        [HttpPost]
        public IActionResult Add(NewPersonBindModel newPersonBindModel)
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
                Person newPerson = new Person
                {
                    Name = newPersonBindModel.Name,
                    Type = (Person.PersonType)Enum.ToObject(typeof(Person.PersonType), newPersonBindModel.Type),
                    User = user,
                };
                newPerson.PersonHistories.Add(new PersonHistory
                {
                    Name = newPerson.Name,
                    IsGone = newPerson.IsGone,
                    Type = newPerson.Type,
                    IsDeleted = newPerson.IsDeleted,
                    Description = "create new person",
                    Person = newPerson
                });
                _context.Persons.Add(newPerson);
                _context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Person");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public IActionResult Edit(Guid id)
        {

            Person? person = _context.Persons.Find(id);

            if (person == null)
            {
                return NotFound();
            }
            User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator");
            person.PersonDevices = (from device in _context.Devices
                                    join personDevice in _context.PersonDevices on device.Id equals personDevice.DeviceId
                                    where device.UserId == user.Id && personDevice.PersonId == person.Id
                                    select new PersonDevice
                                    {
                                        Id = personDevice.Id,
                                        PersonId = personDevice.PersonId,
                                        DeviceId = personDevice.DeviceId,
                                        Device = device,
                                        IsDeleted = personDevice.IsDeleted,
                                    }).ToList();

            return View(person);
        }

        [HttpPost]
        public IActionResult Save()
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                return null;
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }
    }
}
