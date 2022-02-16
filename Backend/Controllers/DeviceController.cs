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

            return CreatedAtAction("GetDevice", new { id = device.Id }, device);
        }

        // POST: Device/save
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost("save")]
        public async Task<ActionResult<Device>> DeviceSave(Device device, Person person)
        {
            // Continue here - getting empty objects
            Console.WriteLine(Newtonsoft.Json.JsonConvert.SerializeObject(device));
            Console.WriteLine(Newtonsoft.Json.JsonConvert.SerializeObject(person));
            return null;
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
