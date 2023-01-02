using KeyBook.Database;
using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Services
{
    public class PersonService
    {
        private readonly UserManager<User> __userManager;
        private readonly KeyBookDbContext __context;
        private readonly IHttpContextAccessor __httpContextAccessor;

        public PersonService(UserManager<User> userManager, KeyBookDbContext context, IHttpContextAccessor httpContextAccessor)
        {
            __userManager = userManager;
            __context = context;
            __httpContextAccessor = httpContextAccessor;
        }

        public async Task<(bool, string?)> AddPerson(Person newPerson)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No Http Context");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                newPerson.OrganizationId = user.OrganizationId;
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
                return (true, null);
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return (false, ex.Message);
            }
        }

        public async Task<List<Person>?> GetPersonForUser(bool showPersonsHowLeft)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (user == null) return null;
            var personDeviceQuery = from person in __context.Set<Person>()
                                    from personDevice in __context.Set<PersonDevice>().Where(personDevice => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                    from device in __context.Set<Device>().Where(device => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                    where person.OrganizationId == user.OrganizationId && (!person.IsGone || showPersonsHowLeft) && !person.IsDeleted
                                    orderby person.Name ascending
                                    select new { person, personDevice, device };
            List<Person> personsWithDuplicates = personDeviceQuery.ToArray().Select(pdq => pdq.person).ToList();
            // get rid of deplicated person due to joining with personDevice
            Dictionary<Guid, Person> personsDict = new Dictionary<Guid, Person>();
            foreach (Person person in personsWithDuplicates)
            {
                if (!personsDict.ContainsKey(person.Id)) personsDict[person.Id] = person;
            }
            return personsDict.Values.ToList();
        }

        public async Task<Dictionary<Guid, string>> GetPersonNamesTypesForUser()
        {
            if (__httpContextAccessor.HttpContext == null) throw new Exception("No HttpContext");
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            return __context.Persons.Where(p => p.OrganizationId == user.OrganizationId && !p.IsGone).OrderBy(p => p.Name).ToDictionary(p => p.Id, p => string.Format("{0} - {1}", p.Name, p.Type.ToString()));
        }

        public Dictionary<Enum, string> GetPersonTypes()
        {
            return Enum.GetValues(typeof(Person.PersonType)).Cast<Enum>().ToDictionary(t => t, t => t.ToString());
        }

        public async Task<List<PersonDevice>> GetPersonDevices(Guid personId)
        {
            if (__httpContextAccessor.HttpContext == null) throw new Exception("No HttpContext");
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            return (from device in __context.Devices
                    join personDevice in __context.PersonDevices on device.Id equals personDevice.DeviceId
                    where device.OrganizationId == user.OrganizationId && personDevice.PersonId == personId
                    select new PersonDevice
                    {
                        Id = personDevice.Id,
                        PersonId = personDevice.PersonId,
                        DeviceId = personDevice.DeviceId,
                        Device = device,
                        IsDeleted = personDevice.IsDeleted,
                    }).ToList();
        }

        public async Task<Person?> GetPersonById(Guid personId)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            Person? person = __context.Persons.Where(p => p.Id == personId && p.OrganizationId == user.OrganizationId).FirstOrDefault();
            return person;
        }

        public async Task<(bool, string?)> SavePerson(Person personFromView)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No HttpContext");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                Person? personFromDb = __context.Persons.Where(
                    p => p.Id == personFromView.Id && p.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (personFromDb == null) throw new Exception("Person not found");
                bool isNameChange;
                if (isNameChange = (personFromDb.Name != personFromView.Name)) personFromDb.Name = personFromView.Name;
                bool isIsGoneChange;
                if (isIsGoneChange = (personFromDb.IsGone != personFromView.IsGone)) personFromDb.IsGone = personFromView.IsGone;
                if (isNameChange || isIsGoneChange)
                {
                    __context.PersonHistories.Add(new PersonHistory
                    {
                        Name = personFromDb.Name,
                        IsGone = personFromDb.IsGone,
                        Type = personFromDb.Type,
                        IsDeleted = personFromDb.IsDeleted,
                        Description = (personFromDb.IsGone) ? "Person mark as left" : "Person's details changed",
                        PersonId = personFromDb.Id
                    });
                }
                __context.Persons.Update(personFromDb);
                __context.SaveChanges();
                transaction.Commit();
                return (true, null);
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return (false, ex.Message);
            }
        }
    }
}
