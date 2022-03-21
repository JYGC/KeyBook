#nullable disable
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Backend.Models;
using Microsoft.EntityFrameworkCore.Storage;

namespace Backend.Controllers
{
    [Route("[controller]")]
    [ApiController]
    public class DeviceController : ControllerBase
    {
        private readonly KeyBookDbContext _context;

        public DeviceController(KeyBookDbContext context)
        {
            _context = context;
        }

        // GET: Device
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Device>>> GetDevices()
        {
            return await _context.Devices.Where(
                device => device.Status == Device.DeviceStatus.NotUsed || device.Status == Device.DeviceStatus.WithManager || device.Status == Device.DeviceStatus.Used
            ).ToListAsync();
        }

        // GET: Device/view/id/149BE541-9271-4E3B-8766-08D9D36C9255
        [HttpGet("view/id/{id}")]
        public async Task<ActionResult<Device>> DeviceView(Guid id)
        {
            Device device = await _context.Devices.FindAsync(id);
            device.PersonDevice = await _context.PersonDevices.FirstOrDefaultAsync(pd => pd.DeviceId == device.Id);

            if (device == null)
            {
                return NotFound();
            }

            return device;
        }

        // POST: Device/add
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost("add")]
        public async Task<ActionResult<Device>> DeviceAdd(Device device)
        {
            User user = await _context.Users.FirstOrDefaultAsync(u => u.Name == "Administrator"); //replace this - Authentication
            device.User = user;
            device.DeviceHistories.Add(new DeviceHistory
            {
                Name = device.Name,
                Identifier = device.Identifier,
                Status = device.Status,
                Type = device.Type,
                IsDeleted = device.IsDeleted,
                Description = "create new device",
                Device = device
            });
            _context.Devices.Add(device);
            await _context.SaveChangesAsync();

            return CreatedAtAction("DeviceView", new { id = device.Id }, device);
        }

        [BindProperties]
        public class DeviceWithPersonRequest
        {
            public Device Device { get; set; }
            public Person Person { get; set; }
        }
        // POST: Device/save
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost("save")]
        public async Task<ActionResult<Device>> DeviceSave(DeviceWithPersonRequest devicePerson)
        {
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                User user = await _context.Users.FirstOrDefaultAsync(u => u.Name == "Administrator"); //replace this - Authentication
                // Update device
                Device deviceFromDb = await _context.Devices.Where(
                    d => d.Id == devicePerson.Device.Id && d.UserId == user.Id
                ).FirstOrDefaultAsync();
                if (deviceFromDb == null)
                {
                    return this.StatusCode(StatusCodes.Status404NotFound, "Device not found");
                }
                bool isNameChange;
                if (isNameChange = (deviceFromDb.Name != devicePerson.Device.Name)) deviceFromDb.Name = devicePerson.Device.Name;
                bool isIdentifierChange;
                if (isIdentifierChange = (deviceFromDb.Identifier != devicePerson.Device.Identifier)) deviceFromDb.Identifier = devicePerson.Device.Identifier;
                bool isSatusChange;
                if (isSatusChange = (deviceFromDb.Status != devicePerson.Device.Status)) deviceFromDb.Status = devicePerson.Device.Status;
                if (isNameChange || isIdentifierChange || isSatusChange)
                {
                    await _context.DeviceHistories.AddAsync(new DeviceHistory
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
                PersonDevice personDevice = await _context.PersonDevices.FirstOrDefaultAsync(pd => deviceFromDb.Id == pd.DeviceId);
                // Update PersonDevices
                if (devicePerson.Person == null && personDevice != null)
                {
                    __RemovePersonDeviceAndEditAssociateHistory(deviceFromDb);
                    personDevice = null;
                }
                else if (devicePerson.Person != null && (personDevice == null || personDevice.PersonId != devicePerson.Person.Id))
                {
                    if (personDevice != null) __RemovePersonDeviceAndEditAssociateHistory(deviceFromDb);
                    Person personFromDb = await _context.Persons.Where(
                        p => p.Id == devicePerson.Person.Id && p.UserId == user.Id
                    ).FirstOrDefaultAsync();
                    personDevice = new PersonDevice
                    {
                        Person = personFromDb,
                        Device = deviceFromDb
                    };
                    await _context.PersonDevices.AddAsync(personDevice);
                    await _context.SaveChangesAsync();
                    await _context.PersonDeviceHistories.AddAsync(new PersonDeviceHistory
                    {
                        PersonDeviceId = personDevice.Id,
                        PersonId = personDevice.Person.Id,
                        DeviceId = personDevice.Device.Id,
                        Description = String.Format("{0} assigned to {1} {2}", deviceFromDb.Name, personFromDb.Type, personFromDb.Name)
                    });
                    _context.Persons.Update(personFromDb);
                }
                _context.Devices.Update(deviceFromDb);
                await _context.SaveChangesAsync();
                transaction.Commit();
                return deviceFromDb;
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        private async void __RemovePersonDeviceAndEditAssociateHistory(Device deviceFromDb)
        {
            deviceFromDb.PersonDevice.IsNoLongerHas = true;
            await _context.PersonDeviceHistories.AddAsync(new PersonDeviceHistory
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

        // DELETE: api/Device/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteDevice(Guid id)
        {
            var device = await _context.Devices.FindAsync(id);
            if (device == null)
            {
                return NotFound();
            }

            _context.Devices.Remove(device);
            await _context.SaveChangesAsync();

            return NoContent();
        }

        private bool DeviceExists(Guid id)
        {
            return _context.Devices.Any(e => e.Id == id);
        }
    }
}
