using Microsoft.EntityFrameworkCore;

namespace Backend.Models
{
    public static class SeedData
    {
        public static void Initialize(IServiceProvider serviceProvider)
        {
            using KeyBookDbContext context = new(serviceProvider.GetRequiredService<DbContextOptions<KeyBookDbContext>>());
            // seed devices
            if (!context.Devices.Any())
            {
                List<Device> seededDevices = new List<Device>
                {
                    new Device
                    {
                        Name = "Remote 1",
                        Identifier = "3285 09 2017",
                        Type = Device.DeviceType.Remote
                    },
                    new Device
                    {
                        Name = "Key 2",
                        Identifier = "104 3055 02 2017",
                        Type = Device.DeviceType.Key
                    },
                    new Device
                    {
                        Name = "Fob 2",
                        Identifier = "01354 SR",
                        Type = Device.DeviceType.Fob
                    }
                };

                foreach (Device device in seededDevices)
                {
                    device.DeviceHistory = new List<DeviceHistory>
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
                        Type = PersonType.Owner
                    },
                    new Person
                    {
                        Name = "CheeJay von Kelhiem",
                        Type = PersonType.Manager
                    },
                    new Person
                    {
                        Name = "Rachel DeSantis",
                        Type = PersonType.Tenant
                    },
                    new Person
                    {
                        Name = "William Coldfuchs",
                        Type = PersonType.Tenant
                    }
                };
                foreach (Person person in seededPersons)
                {
                    person.PersonHistory = new List<PersonHistory>
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
