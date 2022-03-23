using Microsoft.EntityFrameworkCore;

namespace KeyBook.Models
{
    public static class SeedData
    {
        private const string __seedAdminName = "Administrator";
        private const string __seedAdminEmail = "admin@company.com";

        public static void Initialize(IServiceProvider serviceProvider)
        {
            using KeyBookDbContext context = new(serviceProvider.GetRequiredService<DbContextOptions<KeyBookDbContext>>());
            // seed users
            User seedAdminUser;
            User? adminInDb = context.Users.FirstOrDefault(u => u.Name == __seedAdminName && u.IsAdmin == true);
            if (adminInDb == null)
            {
                seedAdminUser = new User
                {
                    Name = __seedAdminName,
                    Email = __seedAdminEmail,
                    IsAdmin = true,
                };
                seedAdminUser.UserHistories.Add(new UserHistory
                {
                    Name = seedAdminUser.Name,
                    Email = seedAdminUser.Email,
                    IsAdmin = seedAdminUser.IsAdmin,
                    IsDeleted = seedAdminUser.IsDeleted,
                    IsBlocked = seedAdminUser.IsBlocked,
                    Description = "seeding admin user",
                    User = seedAdminUser
                });
                context.Users.Add(seedAdminUser);
            }
            else
            {
                seedAdminUser = adminInDb;
            }
            // seed devices
            List<Device> seededDevices;
            if (!context.Devices.Any())
            {
                seededDevices = new List<Device>
                {
                    new Device
                    {
                        Name = "Remote 1",
                        Identifier = "`'\"; OR 1=\'1",
                        Type = Device.DeviceType.Remote,
                        User = seedAdminUser
                    },
                    new Device
                    {
                        Name = "Key 2",
                        Identifier = "<script>alert(1);</script>",
                        Type = Device.DeviceType.Key,
                        User = seedAdminUser
                    },
                    new Device
                    {
                        Name = "Fob 2",
                        Identifier = ";D:D:D:D",
                        Type = Device.DeviceType.Fob,
                        User = seedAdminUser
                    },
                    new Device
                    {
                        Name = "Remote 2",
                        Identifier = "8D",
                        Type = Device.DeviceType.Remote,
                        User = seedAdminUser
                    },
                    new Device
                    {
                        Name = "Mail Kay 1",
                        Identifier = "<h1>Hamburger</h1>",
                        Type = Device.DeviceType.MailKey,
                        User = seedAdminUser
                    }
                };
                foreach (Device device in seededDevices)
                {
                    device.DeviceHistories.Add(new DeviceHistory
                    {
                        Name = device.Name,
                        Identifier = device.Identifier,
                        Status = device.Status,
                        Type = device.Type,
                        IsDeleted = device.IsDeleted,
                        Description = "seeding device table",
                        Device = device
                    });
                }
                context.Devices.AddRange(seededDevices);
            }
            else
            {
                seededDevices = context.Devices.Select(device => device).ToList();
            }
            // seed persons
            List<Person> seededPersons;
            if (!context.Persons.Any())
            {
                seededPersons = new List<Person>
                {
                    new Person
                    {
                        Name = "Administrator",
                        Type = PersonType.Owner,
                        User = seedAdminUser
                    },
                    new Person
                    {
                        Name = "CheeJay von Kelhiem",
                        Type = PersonType.Manager,
                        User = seedAdminUser
                    },
                    new Person
                    {
                        Name = "Rachel DeSantis",
                        Type = PersonType.Tenant,
                        User = seedAdminUser
                    },
                    new Person
                    {
                        Name = "William Coldfuchs",
                        Type = PersonType.Tenant,
                        User = seedAdminUser
                    }
                };
                foreach (Person person in seededPersons)
                {
                    person.PersonHistories.Add(new PersonHistory
                    {
                        Name = person.Name,
                        Description = "seeding person table",
                        Person = person
                    });
                }
                context.Persons.AddRange(seededPersons);
            }
            else
            {
                seededPersons = context.Persons.Select(person => person).ToList();
            }
            // seed persondevices
            if (!context.PersonDevices.Any())
            {
                List<PersonDevice> seededPersonDevices = new List<PersonDevice>
                {
                    new PersonDevice
                    {
                        Person = seededPersons[0],
                        Device = seededDevices[0],
                    },
                    new PersonDevice
                    {
                        Person = seededPersons[2],
                        Device = seededDevices[1],
                    },
                    new PersonDevice
                    {
                        Person = seededPersons[2],
                        Device = seededDevices[2],
                    },
                    new PersonDevice
                    {
                        Person = seededPersons[2],
                        Device = seededDevices[3],
                    },
                    new PersonDevice
                    {
                        Person = seededPersons[3],
                        Device = seededDevices[4],
                    }
                };
                List<PersonDeviceHistory> seededPersonDeviceHistories = new List<PersonDeviceHistory>();
                foreach (PersonDevice personDevice in seededPersonDevices)
                {
                    seededPersonDeviceHistories.Add(new PersonDeviceHistory
                    {
                        PersonDeviceId = personDevice.Id,
                        PersonId = personDevice.Person.Id,
                        DeviceId = personDevice.Device.Id,
                        Description = "seeding person device table"
                    });
                }
                context.PersonDeviceHistories.AddRange(seededPersonDeviceHistories);
                context.PersonDevices.AddRange(seededPersonDevices);
            }
            context.SaveChanges();
        }
    }
}
