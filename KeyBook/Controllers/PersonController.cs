﻿using KeyBook.Models;
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
                                    orderby person.Name ascending
                                    select new { person, personDevice, device };
            List<Person> personsWithDuplicates = personDeviceQuery.ToArray().Select(pdq => pdq.person).ToList();
            // get rid of deplicated person due to joining with personDevice
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
            return View();
        }

        private Dictionary<int, string> __GetPersonTypes()
        {
            return Enum.GetValues(typeof(Person.PersonType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }

        [HttpPost]
        public async Task<IActionResult> New(Person newPerson)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                newPerson.OrganizationId = user.OrganizationId;
                if (!ModelState.IsValid) return View(newPerson);
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
            return View(person);
        }

        [HttpPost]
        public async Task<IActionResult> Edit(Person personFromView)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                if (!ModelState.IsValid) return View(personFromView);
                Person? personFromDb = __context.Persons.Where(
                    p => p.Id == personFromView.Id && p.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (personFromDb == null) return NotFound("Person not found");
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
                return RedirectToAction("Index", "Person");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public async Task<ActionResult<Dictionary<Guid, string?>>> GetPersonNamesTypesAPI()
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            return __context.Persons.Where(p => p.OrganizationId == user.OrganizationId).OrderBy(p => p.Name).ToDictionary(p => p.Id, p => string.Format("{0} - {1}", p.Name, p.Type.ToString()));
        }

        public ActionResult<Dictionary<int, string>> GetPersonTypesAPI()
        {
            return __GetPersonTypes();
        }

        public async Task<List<PersonDevice>> GetPersonDevicesAPI(Guid personId)
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
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
    }
}
