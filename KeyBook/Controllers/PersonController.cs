using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Controllers
{
    [Authorize]
    public class PersonController : Controller
    {
        private readonly KeyBookDbContext _context;

        public PersonController(KeyBookDbContext context)
        {
            _context = context;
        }

        public IActionResult Index()
        {
            User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Seed"); //replace this - Authentication
            var personDeviceQuery = from person in _context.Set<Person>()
                                    from personDevice in _context.Set<PersonDevice>().Where(personDevice => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                    from device in _context.Set<Device>().Where(device => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                    where person.UserId == user.Id && !person.IsDeleted
                                    select new { person, personDevice, device };

            List<Person> personsWithDuplicates = personDeviceQuery.ToArray().Select(pdq => pdq.person).ToList();
            // remove duplicates
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
                User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
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
            User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
            Person? person = _context.Persons.Where(p => p.Id == id && p.UserId == user.Id).FirstOrDefault();

            if (person == null)
            {
                return NotFound();
            }
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

        [BindProperties]
        public class PersonBindModel
        {
            public string PersonId { get; set; }
            public string Name { get; set; }
            public bool IsGone { get; set; }
        }
        [HttpPost]
        public IActionResult Save(PersonBindModel personBindModel)
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
                Person? personFromDb = _context.Persons.Where(
                    p => p.Id == Guid.Parse(personBindModel.PersonId) && p.UserId == user.Id
                ).FirstOrDefault();
                if (personFromDb == null) return NotFound("Person not found");
                bool isNameChange;
                if (isNameChange = (personFromDb.Name != personBindModel.Name)) personFromDb.Name = personBindModel.Name;
                bool isIsGoneChange;
                if (isIsGoneChange = (personFromDb.IsGone != personBindModel.IsGone)) personFromDb.IsGone = personBindModel.IsGone;
                if (isNameChange || isIsGoneChange)
                {
                    _context.PersonHistories.Add(new PersonHistory
                    {
                        Name = personFromDb.Name,
                        IsGone = personFromDb.IsGone,
                        Type = personFromDb.Type,
                        IsDeleted = personFromDb.IsDeleted,
                        Description = (personFromDb.IsGone) ? "person left" : "edit person",
                        PersonId = personFromDb.Id
                    });
                }
                _context.Persons.Update(personFromDb);
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

        public ActionResult<Dictionary<Guid, string?>> GetPersonNames()
        {
            User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Administrator"); // replace - add user auth
            return _context.Persons.Where(p => p.UserId == user.Id).ToDictionary(p => p.Id, p => p.Name);
        }
    }
}
