using Microsoft.EntityFrameworkCore;

namespace Backend.Models
{
    public class KeyBookDbContext : DbContext
    {
        public KeyBookDbContext(DbContextOptions<KeyBookDbContext> options) : base(options)
        {
        }
    
        public DbSet<Device> Devices { get; set; }
        public DbSet<DeviceHistory> DeviceHistories { get; set; }
        public DbSet<Person> Persons { get; set; }
        public DbSet<PersonDevice> PersonDevices { get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<Device>().ToTable("Device");
            modelBuilder.Entity<DeviceHistory>().ToTable("DeviceHistory");
            modelBuilder.Entity<Person>().ToTable("Person");
            modelBuilder.Entity<PersonDevice>().ToTable("PersonDevice");
        }
    }
}
