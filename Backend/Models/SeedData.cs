using Microsoft.EntityFrameworkCore;

namespace Backend.Models
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
                seedAdminUser.UserHistories = new List<UserHistory>
                {
                    new UserHistory
                    {
                        Name = seedAdminUser.Name,
                        Email = seedAdminUser.Email,
                        IsAdmin = seedAdminUser.IsAdmin,
                        IsDeleted = seedAdminUser.IsDeleted,
                        IsBlocked = seedAdminUser.IsBlocked,
                        Description = "seeding admin user",
                        User = seedAdminUser
                    }
                };
                context.Users.Add(seedAdminUser);
            }
            else
            {
                seedAdminUser = adminInDb;
            }
            // seed devices
            if (!context.Devices.Any())
            {
                List<Device> seededDevices = new List<Device>
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
                    }
                };
                foreach (Device device in seededDevices)
                {
                    device.DeviceHistories = new List<DeviceHistory>
                    {
                        new DeviceHistory
                        {
                            Name = device.Name,
                            Identifier = device.Identifier,
                            Status = device.Status,
                            Type = device.Type,
                            IsDeleted = device.IsDeleted,
                            Description = "seeding device table",
                            Device = device
                        }
                    };
                }
                context.Devices.AddRange(seededDevices);
            }
            // seed persons
            if (!context.Persons.Any())
            {
                List<Person> seededPersons = new List<Person>
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
                    person.PersonHistories = new List<PersonHistory>
                    {
                        new PersonHistory
                        {
                            Name = person.Name,
                            Description = "seed person table",
                            Person = person
                        }
                    };
                }
                context.Persons.AddRange(seededPersons);
            }
            context.SaveChanges();
        }
    }
}
