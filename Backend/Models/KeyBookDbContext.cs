using Microsoft.EntityFrameworkCore;

namespace Backend.Models
{
    public class KeyBookDbContext : DbContext
    {
        public KeyBookDbContext(DbContextOptions<KeyBookDbContext> options) : base(options)
        {
        }
    
        public DbSet<User> Users { get; set; }
        public DbSet<Device> Devices { get; set; }
        public DbSet<DeviceHistory> DeviceHistories { get; set; }
        public DbSet<Person> Persons { get; set; }
        public DbSet<PersonHistory> PersonHistories { get; set; }
        public DbSet<PersonDevice> PersonDevices { get; set; }
        public DbSet<PersonDeviceHistory> PersonDeviceHistories { get; set; }

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        {
            optionsBuilder.UseNpgsql(@"Host=192.168.56.102;Port=5432;Database=KeyBook;Username=KeyBook;Password=keybook");
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<User>().ToTable("User");
            modelBuilder.Entity<Device>().HasOne(u => u.User).WithMany(u => u.Devices).HasForeignKey(d => d.UserId).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Device>().HasOne(d => d.PersonDevice).WithOne(pd => pd.Device).IsRequired(false).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Device>().ToTable("Device");
            modelBuilder.Entity<DeviceHistory>().ToTable("DeviceHistory");
            modelBuilder.Entity<Person>().HasOne(p => p.User).WithMany(u => u.Persons).HasForeignKey(p => p.UserId).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Person>().ToTable("Person");
            modelBuilder.Entity<PersonHistory>().ToTable("PersonHistory");
            modelBuilder.Entity<PersonDevice>().ToTable("PersonDevice");
            modelBuilder.Entity<PersonDeviceHistory>().ToTable("PersonDeviceHistory");
        }
    }
}
