using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Storage;
using System;

namespace KeyBook.Controllers
{
    [Authorize]
    public class PersonController : Controller
    {
        private readonly UserManager<User> __userManager;
        private readonly KeyBookDbContext __context;

        public PersonController(UserManager<User> userManager, KeyBookDbContext context)
        {
            __userManager = userManager;
            __context = context;
        }

        public async Task<IActionResult> Index()
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            var personDeviceQuery = from person in __context.Set<Person>()
                                    from personDevice in __context.Set<PersonDevice>().Where(personDevice => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                    from device in __context.Set<Device>().Where(device => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                    where person.OrganizationId == user.OrganizationId && !person.IsDeleted
                                    select new { person, personDevice, device };
            List<Person> personsWithDuplicates = personDeviceQuery.ToArray().Select(pdq => pdq.person).ToList();
            // group device possession by person
            Dictionary<Guid, Person> personsDict = new Dictionary<Guid, Person>();
            foreach (Person person in personsWithDuplicates)
            {
                if (!personsDict.ContainsKey(person.Id)) personsDict[person.Id] = person;
            }
            return View(new PersonListViewModel
            {
                Persons = personsDict.Values.ToList(),
                PersonTypes = __GetPersonTypes()
            });
        }

        public IActionResult New()
        {
            return View(new PersonDetailsViewModel
            {
                PersonTypes = __GetPersonTypes(),
                IsNewPerson = true
            });
        }

        private Dictionary<int, string> __GetPersonTypes()
        {
            return Enum.GetValues(typeof(Person.PersonType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }

        public class NewPersonBindModel
        {
            public string Name { get; set; }
            public int Type { get; set; }
        }
        [HttpPost]
        public async Task<IActionResult> Add(NewPersonBindModel newPersonBindModel)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                Person newPerson = new Person
                {
                    Name = newPersonBindModel.Name,
                    Type = (Person.PersonType)Enum.ToObject(typeof(Person.PersonType), newPersonBindModel.Type),
                    OrganizationId = user.OrganizationId,
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
                __context.Persons.Add(newPerson);
                __context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Person");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public async Task<IActionResult> Edit(Guid personId)
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            Person? person = __context.Persons.Where(p => p.Id == personId && p.OrganizationId == user.OrganizationId).FirstOrDefault();

            if (person == null)
            {
                return NotFound();
            }
            person.PersonDevices = (from device in __context.Devices
                                    join personDevice in __context.PersonDevices on device.Id equals personDevice.DeviceId
                                    where device.OrganizationId == user.OrganizationId && personDevice.PersonId == person.Id
                                    select new PersonDevice
                                    {
                                        Id = personDevice.Id,
                                        PersonId = personDevice.PersonId,
                                        DeviceId = personDevice.DeviceId,
                                        Device = device,
                                        IsDeleted = personDevice.IsDeleted,
                                    }).ToList();

            return View(new PersonDetailsViewModel
            {
                Person = person,
                PersonTypes = __GetPersonTypes()
            });
        }

        [BindProperties]
        public class PersonBindModel
        {
            public string PersonId { get; set; }
            public string Name { get; set; }
            public string IsGoneString { get; set; } // Checkbox returns string "on" if checked

            public bool IsGone
            {
                get
                {
                    return IsGoneString == "on";
                }
            }
        }
        [HttpPost]
        public async Task<IActionResult> Save(PersonBindModel personBindModel)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                Person? personFromDb = __context.Persons.Where(
                    p => p.Id == Guid.Parse(personBindModel.PersonId) && p.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (personFromDb == null) return NotFound("Person not found");
                bool isNameChange;
                if (isNameChange = (personFromDb.Name != personBindModel.Name)) personFromDb.Name = personBindModel.Name;
                bool isIsGoneChange;
                if (isIsGoneChange = (personFromDb.IsGone != personBindModel.IsGone)) personFromDb.IsGone = personBindModel.IsGone;
                if (isNameChange || isIsGoneChange)
                {
                    __context.PersonHistories.Add(new PersonHistory
                    {
                        Name = personFromDb.Name,
                        IsGone = personFromDb.IsGone,
                        Type = personFromDb.Type,
                        IsDeleted = personFromDb.IsDeleted,
                        Description = (personFromDb.IsGone) ? "person left" : "edit person",
                        PersonId = personFromDb.Id
                    });
                }
                __context.Persons.Update(personFromDb);
                __context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Person");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public async Task<Dictionary<Guid, string?>> GetPersonNames()
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            return __context.Persons.Where(p => p.OrganizationId == user.OrganizationId).ToDictionary(p => p.Id, p => p.Name);
        }
    }
}
