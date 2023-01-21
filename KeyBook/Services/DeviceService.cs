using KeyBook.Database;
using KeyBook.DataHandling;
using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage;

namespace KeyBook.Services
{
    public class DeviceService
    {
        private readonly UserManager<User> __userManager;
        private readonly KeyBookDbContext __context;
        private readonly IHttpContextAccessor __httpContextAccessor;

        public DeviceService(UserManager<User> userManager, KeyBookDbContext context, IHttpContextAccessor httpContextAccessor)
        {
            __userManager = userManager;
            __context = context;
            __httpContextAccessor = httpContextAccessor;
        }

        public async Task<(bool, string?)> AddDevice(Device newDevice)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No Http Context");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                newDevice.OrganizationId = user.OrganizationId; ;
                newDevice.DeviceHistories.Add(new DeviceHistory
                {
                    Name = newDevice.Name,
                    Identifier = newDevice.Identifier,
                    DefunctReason = newDevice.DefunctReason,
                    Type = newDevice.Type,
                    IsDeleted = newDevice.IsDeleted,
                    Description = "New device registered",
                    Device = newDevice
                });
                __context.Devices.Add(newDevice);
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

        public async Task<List<Device>?> GetDevicesForUser(bool showDefunctedDevices, bool switchToDeletedDevices)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (user == null) return null;
            var devicePersonAssocRowQuery = from device in __context.Devices
                                            from personDevice in __context.PersonDevices.Where(personDevice => device.Id == personDevice.DeviceId).DefaultIfEmpty()
                                            from person in __context.Persons.Where(person => personDevice.PersonId == person.Id).DefaultIfEmpty()
                                            where device.OrganizationId == user.OrganizationId && ((device.DefunctReason == Device.DeviceDefunctReason.None) || showDefunctedDevices)
                                                                                               && ((!switchToDeletedDevices && !device.IsDeleted) || (switchToDeletedDevices && device.IsDeleted))
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

            return devices;
        }

        public async Task<Device?> GetDeviceById(Guid deviceId)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            Device? device = __context.Devices.Where(d => d.Id == deviceId && d.OrganizationId == user.OrganizationId).FirstOrDefault();
            return device;
        }

        public Dictionary<Enum, string> GetDeviceTypes()
        {
            return Enum.GetValues(typeof(Device.DeviceType)).Cast<Enum>().ToDictionary(t => t, t => t.GetDescription());
        }

        public Dictionary<Enum, string> GetDeviceDefunctReason()
        {
            return Enum.GetValues(typeof(Device.DeviceDefunctReason)).Cast<Enum>().ToDictionary(t => t, t => t.ToString());
        }

        public async Task<List<DeviceActivityHistory>> GetDeviceActivityHistoryList(Guid deviceId)
        {
            if (__httpContextAccessor.HttpContext == null) throw new Exception("No HTTP context");
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            List<DeviceActivityHistory> activityHistoryList = __context.DeviceActivityHistory.FromSqlRaw(
                $"SELECT * FROM \"KeyBook\".sp_GetDeviceActivityHistory('{deviceId}', '{user.OrganizationId}')").ToList();
            return activityHistoryList;
        }

        public async Task<(bool, string?)> SaveDevice(Device deviceFromView)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No HttpContext");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);

                // Update device
                Device? deviceFromDb = __context.Devices.Where(
                    d => d.Id == deviceFromView.Id && d.OrganizationId == user.OrganizationId
                ).FirstOrDefault();
                if (deviceFromDb == null) throw new Exception("No device in database");
                bool detailsOrStatusChanged = (deviceFromDb.Name != deviceFromView.Name || deviceFromDb.Identifier != deviceFromView.Identifier || deviceFromDb.DefunctReason != deviceFromView.DefunctReason);
                deviceFromDb.Name = deviceFromView.Name;
                deviceFromDb.Identifier = deviceFromView.Identifier;
                deviceFromDb.DefunctReason = deviceFromView.DefunctReason;

                __context.DeviceHistories.Add(new DeviceHistory
                {
                    Name = deviceFromDb.Name,
                    Identifier = deviceFromDb.Identifier,
                    DefunctReason = deviceFromDb.DefunctReason,
                    Type = deviceFromDb.Type,
                    IsDeleted = deviceFromDb.IsDeleted,
                    Description = "Device details and status edited",
                    DeviceId = deviceFromDb.Id
                });
                __context.SaveChanges();

                __context.Devices.Update(deviceFromView);
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

        public async Task<(bool, string?)> SavePersonDevice(Guid deviceId, Guid personId)
        {
            using IDbContextTransaction transaction = __context.Database.BeginTransaction();
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No HTTP Context");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                Device? device = __context.Devices.Where(d => d.Id == deviceId && d.OrganizationId == user.OrganizationId).FirstOrDefault();
                Person? person = __context.Persons.Where(p => p.Id == personId && p.OrganizationId == user.OrganizationId).FirstOrDefault();
                if (device == null || (personId != Guid.Empty && person == null)) return (false, "PersonId not belonging to any person");
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
                return (true, null);
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return (false, ex.Message);
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

        public async Task<(bool, string?)> ToggleDeleteDevice(Guid deviceId)
        {
            try
            {
                if (__httpContextAccessor.HttpContext == null) throw new Exception("No HTTP Context");
                User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                Device? deviceFromDb = __context.Devices.Where(d => d.Id == deviceId && d.OrganizationId == user.OrganizationId).First();
                deviceFromDb.IsDeleted = !deviceFromDb.IsDeleted;
                __context.Devices.Update(deviceFromDb);

                __context.DeviceHistories.Add(new DeviceHistory
                {
                    Name = deviceFromDb.Name,
                    Identifier = deviceFromDb.Identifier,
                    DefunctReason = deviceFromDb.DefunctReason,
                    Type = deviceFromDb.Type,
                    IsDeleted = deviceFromDb.IsDeleted,
                    Description = deviceFromDb.IsDeleted ? "Device deleted" : "Device undeleted",
                    DeviceId = deviceFromDb.Id
                });

                __context.SaveChanges();
                return (true, null);
            }
            catch(Exception ex)
            {
                return (false, ex.Message);
            }
        }
    }
}
