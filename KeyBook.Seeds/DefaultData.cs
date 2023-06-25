using KeyBook.Database;
using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

namespace KeyBook.Seeds
{
    public static class DefaultData
    {
        public static async Task SeedAsync(IServiceProvider serviceProvider)
        {
            using KeyBookDbContext context = new(serviceProvider.GetRequiredService<DbContextOptions<KeyBookDbContext>>());
            User seedUser = await serviceProvider.GetRequiredService<UserManager<User>>().Users.Where(u => u.UserName == DefaultUsers.SeedUserEmail).FirstAsync();
            seedUser.Organization = await context.Organizations.Where(o => o.Users.Where(u => u.Email == seedUser.Email).Any()).FirstAsync();
            // seed devices
            List <Device> seededDevices;
            if (!context.Devices.Any())
            {
                seededDevices = new List<Device>
                {
                    new Device
                    {
                        Name = "Remote 1",
                        Identifier = "`'\"; OR 1=\'1",
                        Type = Device.DeviceType.Remote,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Device
                    {
                        Name = "Key 2",
                        Identifier = "<script>alert(1);</script>",
                        Type = Device.DeviceType.Key,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Device
                    {
                        Name = "Fob 2",
                        Identifier = ";D:D:D:D",
                        Type = Device.DeviceType.Fob,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Device
                    {
                        Name = "Remote 2",
                        Identifier = "8D",
                        Type = Device.DeviceType.Remote,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Device
                    {
                        Name = "Mail Kay 1",
                        Identifier = "<h1>Hamburger</h1>",
                        Type = Device.DeviceType.MailboxKey,
                        OrganizationId = seedUser.OrganizationId
                    }
                };
                DateTime[] recordDateTime = new DateTime[]
                {
                    new DateTime(2021, 6, 25, 11, 0, 0, DateTimeKind.Utc),
                    new DateTime(2021, 8, 10, 11, 0, 0, DateTimeKind.Utc),
                    new DateTime(2016, 1, 1, 11, 0, 0, DateTimeKind.Utc),
                    new DateTime(2018, 10, 22, 11, 0, 0, DateTimeKind.Utc),
                    new DateTime(2014, 4, 9, 11, 0, 0, DateTimeKind.Utc),
                };
                for (int i = 0; i < seededDevices.Count; i++)
                {
                    seededDevices[i].DeviceHistories.Add(new DeviceHistory
                    {
                        Name = seededDevices[i].Name,
                        Identifier = seededDevices[i].Identifier,
                        DefunctReason = seededDevices[i].DefunctReason,
                        Type = seededDevices[i].Type,
                        IsDeleted = seededDevices[i].IsDeleted,
                        Description = "seeding device table",
                        Device = seededDevices[i],
                        RecordDateTime = recordDateTime[i],
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
                        Name = seedUser.Name,
                        Type = Person.PersonType.Owner,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Person
                    {
                        Name = "CheeJay von Kelhiem",
                        Type = Person.PersonType.Manager,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Person
                    {
                        Name = "Rachel DeSantis",
                        Type = Person.PersonType.Tenant,
                        OrganizationId = seedUser.OrganizationId
                    },
                    new Person
                    {
                        Name = "William Coldfuchs",
                        Type = Person.PersonType.Tenant,
                        OrganizationId = seedUser.OrganizationId
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
                int[] givenDaysAfterCreation = new int[] { 9, 32, 64, 128, 256 };
                for (int i = 0; i < seededPersonDevices.Count; i++)
                {
                    seededPersonDeviceHistories.Add(new PersonDeviceHistory
                    {
                        PersonDeviceId = seededPersonDevices[i].Id,
                        PersonId = seededPersonDevices[i].Person.Id,
                        DeviceId = seededPersonDevices[i].Device.Id,
                        Description = "seeding person device table",
                        RecordDateTime = seededPersonDevices[i].Device.DeviceHistories.ToList()[0].RecordDateTime.AddDays(givenDaysAfterCreation[i]),
                    });
                }
                context.PersonDeviceHistories.AddRange(seededPersonDeviceHistories);
                context.PersonDevices.AddRange(seededPersonDevices);
            }
            context.SaveChanges();
        }
    }
}
