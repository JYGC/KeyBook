using ExcelDataReader;
using KeyBook.Models;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage;
using System.Data;

namespace KeyBook.Controllers
{
    [Authorize]
    public class DataImportController : Controller
    {
        private const string __UPLOAD_FOLDER = "Uploads";

        private readonly KeyBookDbContext _context;
        private IHostEnvironment __environment;

        public DataImportController(IHostEnvironment environment, KeyBookDbContext context)
        {
            __environment = environment;
            _context = context;
        }

        public IActionResult Index()
        {
            return View();
        }

        [HttpPost]
        //[ValidateAntiForgeryToken] - Add XSRF protection later 
        public IActionResult Excel(IFormFile postedFile)
        {
            User? user = _context.UserTable.FirstOrDefault(u => u.Name == "Administrator"); //replace this - Authentication
            if (postedFile == null || (!postedFile.FileName.EndsWith(".xls") && !postedFile.FileName.EndsWith(".xlsx"))) return NotFound();
            using IDbContextTransaction transaction = _context.Database.BeginTransaction();
            try
            {
                DataTableCollection dataTableCollection = __ConvertUploadToDataTable(postedFile);
                Dictionary<string, Device> inboundDevicesWithIdent = new Dictionary<string, Device>();
                List<DeviceHistory> deviceHistories = new List<DeviceHistory>();
                List<PersonHistory> personHistories = new List<PersonHistory>();
                Dictionary<string, Dictionary<string, PersonDeviceHistory>> personDevicesHistories = new Dictionary<string, Dictionary<string, PersonDeviceHistory>>();
                // Get dates
                Dictionary<int, DateTime> datesWithColNum = new Dictionary<int, DateTime>();
                DateTime[] dateTimeHeading = dataTableCollection[0].Rows[0].ItemArray.Select(d =>
                {
                    DateTime.TryParse(d.ToString(), out DateTime dateTime);
                    return dateTime;
                }).ToArray();
                // Get users
                Dictionary<string, Person> personsWithName = new Dictionary<string, Person>();
                for (int r1 = 1; r1 < dataTableCollection[1].Rows.Count; r1++)
                {
                    string personName = dataTableCollection[1].Rows[r1][0].ToString();
                    personsWithName[personName] = new Person
                    {
                        Name = personName,
                        ImportIdentifier = dataTableCollection[1].Rows[r1][1].ToString(),
                        User = user,
                    };
                }
                Person[] existingPersons = _context.Persons.Where(p => p.UserId == user.Id && personsWithName.Keys.Contains(p.ImportIdentifier)).ToArray();
                foreach (Person existingPerson in existingPersons)
                {
                    personsWithName[existingPerson.Name] = existingPerson;
                }
                // Get device and persondevices
                for (int r0 = 1; r0 < dataTableCollection[0].Rows.Count; r0++)
                {
                    string deviceIdentifier = dataTableCollection[0].Rows[r0][1].ToString();
                    string deviceName = dataTableCollection[0].Rows[r0][0].ToString();
                    if (string.IsNullOrWhiteSpace(deviceIdentifier)) throw new Exception(string.Format("Device with name, {0}, must have identifier", deviceName));
                    if (string.IsNullOrWhiteSpace(deviceName)) throw new Exception(string.Format("Device with identifier, {0}, must have name", deviceIdentifier));
                    inboundDevicesWithIdent[deviceIdentifier] = new Device
                    {
                        Name = deviceName,
                        Identifier = deviceIdentifier,
                        Type = (Device.DeviceType)Enum.Parse(typeof(Device.DeviceType), dataTableCollection[0].Rows[r0][2].ToString()),
                        User = user,
                    };
                    personDevicesHistories[deviceIdentifier] = new Dictionary<string, PersonDeviceHistory>();
                    for (int c = 3; c < dataTableCollection[0].Rows[r0].ItemArray.Length; c++)
                    {
                        string personName = dataTableCollection[0].Rows[r0][c].ToString();
                        if (string.IsNullOrWhiteSpace(personName)) continue;
                        personDevicesHistories[deviceIdentifier][personName] = new PersonDeviceHistory
                        {
                            DeviceId = inboundDevicesWithIdent[deviceIdentifier].Id,
                            PersonId = personsWithName[personName].Id,
                            IsNoLongerHas = false,
                            IsDeleted = false,
                            RecordDateTime = datesWithColNum[c]
                        };
                    }
                }
                Device[] existingDevices = _context.Devices.Include(d => d.PersonDevice).Where(d => d.UserId == user.Id && inboundDevicesWithIdent.Keys.Contains(d.Identifier)).ToArray();
                foreach (Device existingDevice in existingDevices)
                {
                    existingDevice.Name = inboundDevicesWithIdent[existingDevice.Identifier].Name;
                    inboundDevicesWithIdent.Remove(existingDevice.Identifier);
                    deviceHistories.Add(new DeviceHistory
                    {
                        Name = existingDevice.Name,
                        Identifier = existingDevice.Identifier,
                        Status = existingDevice.Status,
                        Type = existingDevice.Type,
                        IsDeleted = existingDevice.IsDeleted,
                        Description = "sync with excel",
                        Device = existingDevice
                    });
                }
                foreach (Device newDevice in inboundDevicesWithIdent.Values)
                {
                    deviceHistories.Add(new DeviceHistory
                    {
                        Name = newDevice.Name,
                        Identifier = newDevice.Identifier,
                        Status = newDevice.Status,
                        Type = newDevice.Type,
                        IsDeleted = newDevice.IsDeleted,
                        Description = "create new device from excel",
                        Device = newDevice
                    });
                }
                _context.Devices.AddRange(inboundDevicesWithIdent.Values);
                _context.DeviceHistories.AddRange(deviceHistories);
                _context.SaveChanges();
                transaction.Commit();
                return RedirectToAction("Index", "Device");
            }
            catch (Exception ex) // Error handling
            {
                transaction.Rollback();
                return NotFound(ex);
            }
        }

        private DataTableCollection __ConvertUploadToDataTable(IFormFile postedFile)
        {
            System.Text.Encoding.RegisterProvider(System.Text.CodePagesEncodingProvider.Instance);
            using MemoryStream stream = new MemoryStream();
            postedFile.CopyTo(stream);
            using IExcelDataReader reader = ExcelReaderFactory.CreateReader(stream);
            return reader.AsDataSet().Tables;
        }
    }
}
