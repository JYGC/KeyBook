using KeyBook.Models;
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
                                            select new { device, personDevice, person };
            List <Device> devices = new List<Device>();
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
            return View(new DeviceDetailsViewModel
            {
                DeviceTypes = __GetDeviceTypes(),
                IsNewDevice = true
            });
        }

        private Dictionary<int, string> __GetDeviceTypes()
        {
            return Enum.GetValues(typeof(Device.DeviceType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }

        public class NewDeviceBindModel
        {
            public string Name { get; set; }
            public string Identifier { get; set; }
            public int Type { get; set; }
        }
        [HttpPost] 
        public async Task<IActionResult> Add(NewDeviceBindModel newDeviceBindModel)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                Device newDevice = new Device
                {
                    Name = newDeviceBindModel.Name,
                    Identifier = newDeviceBindModel.Identifier,
                    Type = (Device.DeviceType)Enum.ToObject(typeof(Device.DeviceType), newDeviceBindModel.Type),
                    OrganizationId = user.OrganizationId,
                };
                newDevice.DeviceHistories.Add(new DeviceHistory
                {
                    Name = newDevice.Name,
                    Identifier = newDevice.Identifier,
                    Status = newDevice.Status,
                    Type = newDevice.Type,
                    IsDeleted = newDevice.IsDeleted,
                    Description = "create new device",
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
            device.PersonDevice = __context.PersonDevices.FirstOrDefault(pd => pd.DeviceId == device.Id);

            if (device == null)
            {
                return NotFound();
            }

            return View(new DeviceDetailsViewModel
            {
                Device = device,
                FromPersonDetailsPersonId = fromPersonDetailsPersonId,
                DeviceActivityHistoryList = __context.DeviceActivityHistory.FromSqlRaw(
                    $"SELECT * FROM \"KeyBook\".sp_GetDeviceActivityHistory('{device.Id}')"
                ).ToList(),
                DeviceTypes = __GetDeviceTypes(),
                PersonNames = __context.Persons.Where(p => p.OrganizationId == user.OrganizationId).ToDictionary(p => p.Id, p => p.Name)
            });
        }

        [BindProperties]
        public class DevicePersonBindModel
        {
            public string DeviceId { get; set; }
            public string Name { get; set; }
            public string Identifier { get; set; }
            public string Status { get; set; }
            public string PersonId { get; set; }
            public string FromPersonDetailsPersonId { get; set; } = null;
        }
        [HttpPost]
        public async Task<IActionResult> Save(DevicePersonBindModel devicePersonViewModel) // continue here - not all properties are passing
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                User? user = await __userManager.GetUserAsync(HttpContext.User);
                // Update device
                Device? deviceFromDb = __context.Devices.Where(
                    d => d.Id == Guid.Parse(devicePersonViewModel.DeviceId) && d.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (deviceFromDb == null) return NotFound("Device not found");

                Device.DeviceStatus inboundDeviceStatus = (Device.DeviceStatus)Enum.Parse(typeof(Device.DeviceStatus), devicePersonViewModel.Status);
                bool detailsOrStatusChanged = (deviceFromDb.Name != devicePersonViewModel.Name || deviceFromDb.Identifier != devicePersonViewModel.Identifier || deviceFromDb.Status != inboundDeviceStatus);
                deviceFromDb.Name = devicePersonViewModel.Name;
                deviceFromDb.Identifier = devicePersonViewModel.Identifier;
                deviceFromDb.Status = inboundDeviceStatus;
                if (detailsOrStatusChanged)
                {
                    __context.DeviceHistories.Add(new DeviceHistory
                    {
                        Name = deviceFromDb.Name,
                        Identifier = deviceFromDb.Identifier,
                        Status = deviceFromDb.Status,
                        Type = deviceFromDb.Type,
                        IsDeleted = deviceFromDb.IsDeleted,
                        Description = "device details and status edited",
                        Device = deviceFromDb
                    });
                    __context.SaveChanges();
                }
                // Update PersonDevices
                PersonDevice? personDevice = __context.PersonDevices.FirstOrDefault(pd => deviceFromDb.Id == pd.DeviceId);
                Guid inboundPersonId;
                Guid.TryParse(devicePersonViewModel.PersonId, out inboundPersonId);
                if (inboundPersonId == Guid.Empty && personDevice != null)
                {
                    __RemovePersonDeviceAndEditAssociateHistory(deviceFromDb);
                    personDevice = null;
                }
                else if (inboundPersonId != Guid.Empty && (personDevice == null || personDevice.PersonId != inboundPersonId))
                {
                    if (personDevice != null) __RemovePersonDeviceAndEditAssociateHistory(deviceFromDb);
                    Person? personFromDb = __context.Persons.Where(
                        p => p.Id == inboundPersonId && p.OrganizationId == user.OrganizationId
                    ).FirstOrDefault();
                    personDevice = new PersonDevice
                    {
                        Person = personFromDb,
                        Device = deviceFromDb
                    };
                    __context.PersonDevices.Add(personDevice);
                    __context.SaveChanges();
                    __context.PersonDeviceHistories.Add(new PersonDeviceHistory
                    {
                        PersonDeviceId = personDevice.Id,
                        PersonId = personDevice.Person.Id,
                        DeviceId = personDevice.Device.Id,
                        Description = String.Format("{0} assigned to {1} {2}", deviceFromDb.Name, personFromDb.Type, personFromDb.Name)
                    });
                    __context.Persons.Update(personFromDb);
                }
                __context.Devices.Update(deviceFromDb);
                __context.SaveChanges();
                transaction.Commit();
                return (devicePersonViewModel.FromPersonDetailsPersonId == null)
                    ? RedirectToAction("Index", "Device")
                    : RedirectToAction("Edit", "Person", new
                    {
                        personId = devicePersonViewModel.FromPersonDetailsPersonId
                    });
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        private void __RemovePersonDeviceAndEditAssociateHistory(Device deviceFromDb)
        {
            __context.PersonDeviceHistories.Add(new PersonDeviceHistory
            {
                PersonDeviceId = deviceFromDb.PersonDevice.Id,
                PersonId = deviceFromDb.PersonDevice.PersonId,
                DeviceId = deviceFromDb.PersonDevice.DeviceId,
                Description = "person no longer has device",
                IsNoLongerHas = true
            });
            __context.PersonDevices.Remove(deviceFromDb.PersonDevice);
            __context.SaveChanges();
        }
    }
}
