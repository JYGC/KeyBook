﻿using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Controllers
{
    [Authorize]
    //[ValidateAntiForgeryToken] - Add XSRF protection later
    public class DeviceController : Controller
    {
        private readonly UserManager<User> __userManager;
        private readonly KeyBookDbContext __context;

        public DeviceController(UserManager<User> userManager, KeyBookDbContext context)
        {
            __context = context;
            __userManager = userManager;
        }

        public async Task<IActionResult> Index()
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            var devicePersonAssocRowQuery = from device in __context.Devices
                                            from personDevice in __context.PersonDevices.Where(personDevice => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                            from person in __context.Persons.Where(person => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                            where device.OrganizationId == user.OrganizationId && (device.Status == Device.DeviceStatus.NotUsed || device.Status == Device.DeviceStatus.WithManager || device.Status == Device.DeviceStatus.Used)
                                            orderby device.Name ascending
                                            select new { device, personDevice, person };
            List<Device> devices = new List<Device>();
            foreach (var row in devicePersonAssocRowQuery.ToArray())
            {
                if (row.personDevice != null)
                {
                    row.device.PersonDevice = row.personDevice;
                    row.device.PersonDevice.Person = row.person;
                }
                devices.Add(row.device);
            }
            return View(new DeviceListViewModel
            {
                Devices = devices,
                DeviceTypes = __GetDeviceTypes()
            });
        }

        public IActionResult New()
        {
            return View();
        }

        [HttpPost] 
        public async Task<IActionResult> New(Device newDevice)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                newDevice.OrganizationId = user.OrganizationId;
                if (!ModelState.IsValid) return View(newDevice);
                newDevice.DeviceHistories.Add(new DeviceHistory
                {
                    Name = newDevice.Name,
                    Identifier = newDevice.Identifier,
                    Status = newDevice.Status,
                    Type = newDevice.Type,
                    IsDeleted = newDevice.IsDeleted,
                    Description = "New device registered",
                    Device = newDevice
                });
                __context.Devices.Add(newDevice);
                __context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Device");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public async Task<IActionResult> Edit(Guid deviceId, Guid? fromPersonDetailsPersonId)
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            Device? device = __context.Devices.Find(deviceId);

            if (device == null)
            {
                return NotFound();
            }

            TempData["fromPersonDetailsPersonId"] = fromPersonDetailsPersonId;
            return View(device);
        }

        [HttpPost]
        public async Task<IActionResult> Edit(Device deviceFromView)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                if (!ModelState.IsValid) return View(deviceFromView);
                // Update device
                Device? deviceFromDb = __context.Devices.Where(
                    d => d.Id == deviceFromView.Id && d.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (deviceFromDb == null) return NotFound("Device not found");
                bool detailsOrStatusChanged = (deviceFromDb.Name != deviceFromView.Name || deviceFromDb.Identifier != deviceFromView.Identifier || deviceFromDb.Status != deviceFromView.Status);
                deviceFromDb.Name = deviceFromView.Name;
                deviceFromDb.Identifier = deviceFromView.Identifier;
                deviceFromDb.Status = deviceFromView.Status;
                if (detailsOrStatusChanged)
                {
                    __context.DeviceHistories.Add(new DeviceHistory
                    {
                        Name = deviceFromDb.Name,
                        Identifier = deviceFromDb.Identifier,
                        Status = deviceFromDb.Status,
                        Type = deviceFromDb.Type,
                        IsDeleted = deviceFromDb.IsDeleted,
                        Description = "Device details and status edited",
                        Device = deviceFromDb
                    });
                    __context.SaveChanges();
                }
                __context.Devices.Update(deviceFromDb);
                __context.SaveChanges();
                transaction.Commit();
                //var fromPersonDetailsPersonId = TempData["fromPersonDetailsPersonId"];
                //TempData["fromPersonDetailsPersonId"] = null; // TODO: TempData["fromPersonDetailsPersonId"] NOT CARRYING OVER
                return (TempData["fromPersonDetailsPersonId"] == null)
                    ? RedirectToAction("Index", "Device")
                    : RedirectToAction("Edit", "Person", new
                    {
                        personId = TempData["fromPersonDetailsPersonId"]
                    });
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public async Task<PersonDevice?> GetPersonDeviceAPI(Guid deviceId)
        {
            User? user = await __userManager.GetUserAsync(HttpContext.User);
            if (__context.Devices.Any(d => d.Id == deviceId && d.OrganizationId == user.OrganizationId))
            {
                return __context.PersonDevices.FirstOrDefault(pd => pd.DeviceId == deviceId);
            }
            return null;
        }

        public class TestViewModel
        {
            public Guid DeviceId { get; set; }
            public Guid PersonId { get; set; }
        }
        [HttpPost]
        public async Task<ActionResult<bool>> SavePersonDeviceAPI([FromBody]TestViewModel body)
        {
            Guid deviceId = body.DeviceId;
            Guid personId = body.PersonId;
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                Device? device = __context.Devices.Where(d => d.Id == deviceId && d.OrganizationId == user.OrganizationId).FirstOrDefault();
                Person? person = __context.Persons.Where(p => p.Id == personId && p.OrganizationId == user.OrganizationId).FirstOrDefault();
                if (device == null || (personId != Guid.Empty && person == null)) return false;
                PersonDevice? personDevice = __context.PersonDevices.FirstOrDefault(pd => pd.DeviceId == deviceId);
                if (personId == Guid.Empty && personDevice != null) // device previous assigned now unassign device to any
                {
                    __RemovePersonDeviceAndEditAssociateHistory(device);
                    personDevice = null;
                }
                if (personId != Guid.Empty && (personDevice == null || personDevice.PersonId != personId)) // assign device to new holding user
                {
                    if (personDevice != null) __RemovePersonDeviceAndEditAssociateHistory(device); // if device previous assigned
                    personDevice = new PersonDevice
                    {
                        Device = device,
                        Person = person
                    };
                    __context.PersonDevices.Add(personDevice);
                    __context.PersonDeviceHistories.Add(new PersonDeviceHistory
                    {
                        PersonDeviceId = personDevice.Id,
                        PersonId = personDevice.Person.Id,
                        DeviceId = personDevice.Device.Id,
                        Description = String.Format("{0} assigned to {1}, {2}", device.Name, person.Type, person.Name)
                    });
                    __context.Persons.Update(person);
                }
                __context.SaveChanges();
                transaction.Commit();
                return true;
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return false;
            }
        }

        private void __RemovePersonDeviceAndEditAssociateHistory(Device device)
        {
            Person? person = __context.Persons.Where(p => p.Id == device.PersonDevice.PersonId).FirstOrDefault();
            __context.PersonDeviceHistories.Add(new PersonDeviceHistory
            {
                PersonDeviceId = device.PersonDevice.Id,
                PersonId = device.PersonDevice.PersonId,
                DeviceId = device.PersonDevice.DeviceId,
                Description = string.Format("{0}, {1}, no longer has {2}", person.Type, person.Name, device.Name),
                IsNoLongerHas = true
            });
            __context.PersonDevices.Remove(device.PersonDevice);
            __context.SaveChanges();
        }

        public ActionResult<Dictionary<int, string>> GetDeviceTypesAPI()
        {
            return __GetDeviceTypes();
        }

        public ActionResult<List<DeviceActivityHistory>> GetDeviceActivityHistoryListAPI(Guid deviceId)
        {
            return __context.DeviceActivityHistory.FromSqlRaw(
                $"SELECT * FROM \"KeyBook\".sp_GetDeviceActivityHistory('{deviceId}')" // sp_GetDeviceActivityHistory Needs to take in OrganizationId
            ).ToList();
        }

        private Dictionary<int, string> __GetDeviceTypes()
        {
            return Enum.GetValues(typeof(Device.DeviceType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }
    }
}
