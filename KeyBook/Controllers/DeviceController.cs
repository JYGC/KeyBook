using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Controllers
{
    public class DeviceController : Controller
    {
        private readonly KeyBookDbContext _context;

        public DeviceController(KeyBookDbContext context)
        {
            _context = context;
        }

        public IActionResult Index()
        {
            User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
            var devicePersonAssocRowQuery = (from device in _context.Devices
                                             join personDevice in _context.PersonDevices on device.Id equals personDevice.DeviceId into NullablePersonDevice
                                             from nullablePersonDevice in NullablePersonDevice.DefaultIfEmpty()
                                             join person in _context.Persons on nullablePersonDevice.PersonId equals person.Id into NullablePerson
                                             from nullablePerson in NullablePerson.DefaultIfEmpty()
                                             where device.UserId == user.Id && (device.Status == Device.DeviceStatus.NotUsed || device.Status == Device.DeviceStatus.WithManager || device.Status == Device.DeviceStatus.Used)
                                             select new { device, nullablePersonDevice, nullablePerson });
            List<Device> devices = new List<Device>();
            foreach (var row in devicePersonAssocRowQuery.ToArray())
            {
                if (row.nullablePersonDevice != null)
                {
                    row.device.PersonDevice = row.nullablePersonDevice;
                    row.device.PersonDevice.Person = row.nullablePerson;
                }
                devices.Add(row.device);
            }
            return View(new DeviceListViewModel
            {
                Devices = devices
            });
        }

        public IActionResult Welcome(string name, int Id = 1)
        {
            ViewData["name"] = "Hello" + name;
            ViewData["id"] = Id;
            return View();
        }

        public IActionResult New()
        {
            return View();
        }

        public ActionResult<Dictionary<int, string>> GetDeviceTypes()
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
        //[ValidateAntiForgeryToken] - Add XSRF protection later 
        public IActionResult Add(NewDeviceBindModel newDeviceBindModel)
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
                Device newDevice = new Device
                {
                    Name = newDeviceBindModel.Name,
                    Identifier = newDeviceBindModel.Identifier,
                    Type = (Device.DeviceType)Enum.ToObject(typeof(Device.DeviceType), newDeviceBindModel.Type),
                    User = user,
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
                _context.Devices.Add(newDevice);
                _context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Device");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        public IActionResult Edit(Guid id, Guid? fromPersonDetailsPersonId)
        {
            User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); // replace - add user auth
            Device? device = _context.Devices.Find(id);
            device.PersonDevice = _context.PersonDevices.FirstOrDefault(pd => pd.DeviceId == device.Id);

            if (device == null)
            {
                return NotFound();
            }

            return View(new DevicePersonDetailsPersonIdViewModel
            {
                Device = device,
                FromPersonDetailsPersonId = fromPersonDetailsPersonId,
                DeviceActivityHistoryList = _context.DeviceActivityHistory.FromSqlRaw(
                    $"SELECT * FROM public.sp_GetDeviceActivityHistory('{device.Id}')"
                ).ToList()
            });
        }

        [BindProperties]
        public class DevicePersonBindModel
        {
            public string DeviceId { get; set; }
            public string Name { get; set; }
            public string Identifier { get; set; }
            public string Status { get; set; }
            public string Type { get; set; }
            public string PersonId { get; set; }
            public string? FromPersonDetailsPersonId { get; set; } = null;
        }
        [HttpPost]
        //[ValidateAntiForgeryToken] - Add XSRF protection later 
        public IActionResult Save(DevicePersonBindModel devicePersonViewModel) // continue here - not all properties are passing
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                User? user = _context.Users.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
                // Update device
                Device? deviceFromDb = _context.Devices.Where(
                    d => d.Id == Guid.Parse(devicePersonViewModel.DeviceId) && d.UserId == user.Id
                ).FirstOrDefault();
                if (deviceFromDb == null) return NotFound("Device not found");
                bool isNameChange;
                if (isNameChange = (deviceFromDb.Name != devicePersonViewModel.Name)) deviceFromDb.Name = devicePersonViewModel.Name;
                bool isIdentifierChange;
                if (isIdentifierChange = (deviceFromDb.Identifier != devicePersonViewModel.Identifier)) deviceFromDb.Identifier = devicePersonViewModel.Identifier;
                bool isSatusChange;
                Device.DeviceStatus inboundDeviceStatus = (Device.DeviceStatus)Enum.Parse(typeof(Device.DeviceStatus), devicePersonViewModel.Status);
                if (isSatusChange = (deviceFromDb.Status != inboundDeviceStatus)) deviceFromDb.Status = inboundDeviceStatus;
                if (isNameChange || isIdentifierChange || isSatusChange)
                {
                    _context.DeviceHistories.Add(new DeviceHistory
                    {
                        Name = deviceFromDb.Name,
                        Identifier = deviceFromDb.Identifier,
                        Status = deviceFromDb.Status,
                        Type = deviceFromDb.Type,
                        IsDeleted = deviceFromDb.IsDeleted,
                        Description = "edit device",
                        Device = deviceFromDb
                    });
                    _context.SaveChanges();
                }
                PersonDevice? personDevice = _context.PersonDevices.FirstOrDefault(pd => deviceFromDb.Id == pd.DeviceId);
                // Update PersonDevices
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
                    Person? personFromDb = _context.Persons.Where(
                        p => p.Id == inboundPersonId && p.UserId == user.Id
                    ).FirstOrDefault();
                    personDevice = new PersonDevice
                    {
                        Person = personFromDb,
                        Device = deviceFromDb
                    };
                    _context.PersonDevices.Add(personDevice);
                    _context.SaveChanges();
                    _context.PersonDeviceHistories.Add(new PersonDeviceHistory
                    {
                        PersonDeviceId = personDevice.Id,
                        PersonId = personDevice.Person.Id,
                        DeviceId = personDevice.Device.Id,
                        Description = String.Format("{0} assigned to {1} {2}", deviceFromDb.Name, personFromDb.Type, personFromDb.Name)
                    });
                    _context.Persons.Update(personFromDb);
                }
                _context.Devices.Update(deviceFromDb);
                _context.SaveChanges();
                transaction.Commit();
                return (devicePersonViewModel.FromPersonDetailsPersonId == null)
                    ? RedirectToAction("Index", "Device")
                    : RedirectToAction("Edit", "Person", new
                    {
                        id = devicePersonViewModel.FromPersonDetailsPersonId
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
            deviceFromDb.PersonDevice.IsNoLongerHas = true;
            _context.PersonDeviceHistories.Add(new PersonDeviceHistory
            {
                PersonDeviceId = deviceFromDb.PersonDevice.Id,
                PersonId = deviceFromDb.PersonDevice.PersonId,
                DeviceId = deviceFromDb.PersonDevice.DeviceId,
                Description = "person no longer has device",
                IsNoLongerHas = deviceFromDb.PersonDevice.IsNoLongerHas
            });
            _context.PersonDevices.Remove(deviceFromDb.PersonDevice);
            _context.SaveChanges();
        }
    }
}
