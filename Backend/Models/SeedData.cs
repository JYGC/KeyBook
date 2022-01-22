using Microsoft.EntityFrameworkCore;

namespace Backend.Models
{
    public static class SeedData
    {
        public static void Initialize(IServiceProvider serviceProvider)
        {
            using KeyBookDbContext context = new(serviceProvider.GetRequiredService<DbContextOptions<KeyBookDbContext>>());
            if (context.Devices.Any())
            {
                return;
            }

            context.Devices.AddRange(
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
            );
            context.SaveChanges();
        }
    }
}
