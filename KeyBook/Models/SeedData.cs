using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;

namespace KeyBook.Models
{
    public static class SeedData
    {
        private const string __seedUserName = "Seed";
        private const string __seedUserEmail = "seed@app.server";
        private const string __seedUserPassword = "$Eed3D";

        public static void Initialize(IServiceProvider serviceProvider)
        {
            using KeyBookDbContext context = new(serviceProvider.GetRequiredService<DbContextOptions<KeyBookDbContext>>());
            // seed roles
            //foreach (UserRoles userRole in Enum.GetValues(typeof(UserRoles)))
            //{
            //    bool roleExist = await roleManager.RoleExistsAsync(userRole.ToString());
            //    if (!roleExist) await roleManager.CreateAsync(new IdentityRole(userRole.ToString()));
            //}
            // seed users
            //ApplicationUser? superuserInDb = context.ApplicationUsers.FirstOrDefault(u => u.UserName == __seedUserName && u.Email == __seedUserEmail);
            //if (superuserInDb)
            //{
            //AuthUser superuser = new()
            //{
            //    Id = Guid.NewGuid().ToString(),
            //    UserName = __seedAdminName,
            //    Email = __seedAdminEmail,
            //    NormalizedUserName = __seedAdminEmail,
            //    NormalizedEmail = __seedAdminEmail,
            //    LockoutEnabled = false,
            //    EmailConfirmed = true,
            //    PhoneNumberConfirmed = true,
            //    SecurityStamp = Guid.NewGuid().ToString("D")
            //};
            ////context.AuthUsers.Add(superuser);
            //PasswordHasher<AuthUser> passwordHasher = new PasswordHasher<AuthUser>();
            //superuser.SecurityStamp = Guid.NewGuid().ToString();
            //superuser.PasswordHash = passwordHasher.HashPassword(superuser, __seedSuperuserPassword);
            //UserStore<ApplicationUser> userStore = new UserStore<ApplicationUser>(context);
            //UserManager<ApplicationUser> userManager = serviceProvider.GetRequiredService<UserManager<ApplicationUser>>();
            //if (!userManager.CreateAsync(new ApplicationUser { UserName = __seedUserName, Email = __seedUserEmail }, __seedUserPassword).Result.Succeeded) throw new Exception("Superuser could not be created");
            //}
            User seedUser;
            User? seedUserInDb = context.UserTable.FirstOrDefault(u => u.Name == __seedUserName);
            if (seedUserInDb == null)
            {
                seedUser = new User
                {
                    Name = __seedUserName,
                    Email = __seedUserEmail,
                };
                seedUser.UserHistories.Add(new UserHistory
                {
                    Name = seedUser.Name,
                    Email = seedUser.Email,
                    IsAdmin = seedUser.IsAdmin,
                    IsDeleted = seedUser.IsDeleted,
                    IsBlocked = seedUser.IsBlocked,
                    Description = "seeding admin user",
                    User = seedUser
                });
                context.UserTable.Add(seedUser);
            }
            else
            {
                seedUser = seedUserInDb;
            }
            // seed devices
            List<Device> seededDevices;
            if (!context.Devices.Any())
            {
                var MakeDEviceHistoryWithRecordDate = (Device device, DateTime recordDateTime) =>
                {
                    return new DeviceHistory
                    {
                        Name = device.Name,
                        Identifier = device.Identifier,
                        Status = device.Status,
                        Type = device.Type,
                        IsDeleted = device.IsDeleted,
                        Description = "seeding device table",
                        RecordDateTime = recordDateTime,
                        Device = device,
                    };
                };
                seededDevices = new List<Device>
                {
                    new Device
                    {
                        Name = "Remote 1",
                        Identifier = "`'\"; OR 1=\'1",
                        Type = Device.DeviceType.Remote,
                        User = seedUser
                    },
                    new Device
                    {
                        Name = "Key 2",
                        Identifier = "<script>alert(1);</script>",
                        Type = Device.DeviceType.Key,
                        User = seedUser
                    },
                    new Device
                    {
                        Name = "Fob 2",
                        Identifier = ";D:D:D:D",
                        Type = Device.DeviceType.Fob,
                        User = seedUser
                    },
                    new Device
                    {
                        Name = "Remote 2",
                        Identifier = "8D",
                        Type = Device.DeviceType.Remote,
                        User = seedUser
                    },
                    new Device
                    {
                        Name = "Mail Kay 1",
                        Identifier = "<h1>Hamburger</h1>",
                        Type = Device.DeviceType.MailboxKey,
                        User = seedUser
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
                        Status = seededDevices[i].Status,
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
                        Name = __seedUserName,
                        Type = Person.PersonType.Owner,
                        User = seedUser
                    },
                    new Person
                    {
                        Name = "CheeJay von Kelhiem",
                        Type = Person.PersonType.Manager,
                        User = seedUser
                    },
                    new Person
                    {
                        Name = "Rachel DeSantis",
                        Type = Person.PersonType.Tenant,
                        User = seedUser
                    },
                    new Person
                    {
                        Name = "William Coldfuchs",
                        Type = Person.PersonType.Tenant,
                        User = seedUser
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
