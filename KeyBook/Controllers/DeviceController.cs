using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Mvc;

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
            return View(new DeviceListViewModel
            {
                Devices = _context.Devices.Where(
                    device => device.Status == Device.DeviceStatus.NotUsed || device.Status == Device.DeviceStatus.WithManager || device.Status == Device.DeviceStatus.Used
                ).ToList()
            });
        }

        public IActionResult Welcome(string name, int Id = 1)
        {
            ViewData["name"] = "Hello" + name;
            ViewData["id"] = Id;
            return View();
        }

        public IActionResult Edit(Guid id)
        {
            Device device = _context.Devices.Find(id);
            device.PersonDevice = _context.PersonDevices.FirstOrDefault(pd => pd.DeviceId == device.Id);

            if (device == null)
            {
                return NotFound();
            }

            return View(device);
        }
    }
}
