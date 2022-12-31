using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;

namespace KeyBook.Models
{
    public class KeyBookDbContext : IdentityDbContext<User>
    {
        public KeyBookDbContext(DbContextOptions<KeyBookDbContext> options) : base(options)
        {
        }

        public DbSet<User> UserTable { get; set; }
        public DbSet<Organization> Organizations { get; set; }
        public DbSet<UserHistory> UserHistory { get; set; }
        public DbSet<Device> Devices { get; set; }
        public DbSet<DeviceActivityHistory> DeviceActivityHistory { get; set; }
        public DbSet<DeviceHistory> DeviceHistories { get; set; }
        public DbSet<Person> Persons { get; set; }
        public DbSet<PersonHistory> PersonHistories { get; set; }
        public DbSet<PersonDevice> PersonDevices { get; set; }
        public DbSet<PersonDeviceHistory> PersonDeviceHistories { get; set; }

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        {
            optionsBuilder.UseNpgsql(DAL.ConfigSettings.DefaultConnection);
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);
            modelBuilder.HasDefaultSchema("KeyBook");
            modelBuilder.Entity<Organization>().ToTable("Organization");
            modelBuilder.Entity<User>().HasOne(u => u.Organization).WithMany(o => o.Users).HasForeignKey(u => u.OrganizationId).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<UserHistory>().ToTable("UserHistory");
            modelBuilder.Entity<Device>().HasOne(d => d.Organization).WithMany(o => o.Devices).HasForeignKey(d => d.OrganizationId).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Device>().HasOne(d => d.PersonDevice).WithOne(pd => pd.Device).IsRequired(false).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Device>().ToTable("Device");
            modelBuilder.Entity<DeviceActivityHistory>().HasNoKey().ToView("DeviceActionHistory"); // Will not be a table
            modelBuilder.Entity<DeviceHistory>().ToTable("DeviceHistory");
            modelBuilder.Entity<Person>().HasOne(p => p.Organization).WithMany(o => o.Persons).HasForeignKey(p => p.OrganizationId).OnDelete(DeleteBehavior.Restrict);
            modelBuilder.Entity<Person>().ToTable("Person");
            modelBuilder.Entity<PersonHistory>().ToTable("PersonHistory");
            modelBuilder.Entity<PersonDevice>().ToTable("PersonDevice");
            modelBuilder.Entity<PersonDeviceHistory>().ToTable("PersonDeviceHistory");
        }
    }
}
