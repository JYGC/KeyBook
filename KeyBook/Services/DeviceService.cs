using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

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

        public async Task<List<Device>?> GetDevicesForUser()
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? user = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (user == null) return null;
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

            return devices;
        }

        public Dictionary<int, string> GetDeviceTypes()
        {
            return Enum.GetValues(typeof(Device.DeviceType)).Cast<Enum>().ToDictionary(t => (int)(object)t, t => t.ToString());
        }
    }
}
