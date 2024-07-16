using Microsoft.AspNetCore.Identity;
using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models;

// Add profile data for application users by adding properties to the AuthUser class
public class User : IdentityUser
{
    [Required]
    public string Name { get; set; }
    [Required]
    public Guid OrganizationId { get; set; }
    public virtual Organization Organization { get; set; }
    public virtual ICollection<UserHistory> UserHistories { get; set; } = new List<UserHistory>();
}

