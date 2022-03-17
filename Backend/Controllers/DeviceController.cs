#nullable disable
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Backend.Models;

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
            User user = await _context.Users.FirstOrDefaultAsync(u => u.Name == "Administrator"); //replace this - Authentication
            // Update device
            Device deviceFromDb = await _context.Devices.Where(
                d => d.Id == devicePerson.Device.Id && d.UserId == user.Id
            ).FirstOrDefaultAsync();
            if (deviceFromDb == null)
            {
                return this.StatusCode(StatusCodes.Status404NotFound, "Device not found");
            }
            deviceFromDb.Name = devicePerson.Device.Name;
            deviceFromDb.Identifier = devicePerson.Device.Identifier;
            deviceFromDb.Status = devicePerson.Device.Status;
            deviceFromDb.DeviceHistories.Add(new DeviceHistory
            {
                Name = deviceFromDb.Name,
                Identifier = deviceFromDb.Identifier,
                Status = deviceFromDb.Status,
                Type = deviceFromDb.Type,
                IsDeleted = deviceFromDb.IsDeleted,
                Description = "edit device",
                Device = deviceFromDb
            });
            // Update PersonDevices
            if (devicePerson.Person == null && deviceFromDb.PersonDevice != null)
            {
                __RemovePersonDevice(deviceFromDb);
                deviceFromDb.PersonDevice = null;
            }
            else if (devicePerson.Person != null && deviceFromDb.PersonDevice.PersonId != devicePerson.Person.Id)
            {
                __RemovePersonDevice(deviceFromDb);
                Person personFromDb = await _context.Persons.Where(
                    p => p.Id == devicePerson.Person.Id && p.UserId == user.Id
                ).FirstOrDefaultAsync();
                deviceFromDb.PersonDevice = new PersonDevice
                {
                    Person = personFromDb,
                    Device = deviceFromDb
                };
                personFromDb.PersonDevices.Add(deviceFromDb.PersonDevice);
            }
            _context.Devices.Update(deviceFromDb);
            await _context.SaveChangesAsync();
            //continue here
            return null;
        }

        private async void __RemovePersonDevice(Device deviceFromDb)
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
