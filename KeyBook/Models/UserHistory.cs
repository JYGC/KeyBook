using System.ComponentModel.DataAnnotations;

namespace KeyBook.Models
{
    public class UserHistory
    {
        [Key]
        public Guid Id { get; set; } = Guid.NewGuid();
        [Required]
        public string? Name { get; set; }
        [Required]
        public string? Email { get; set; }
        [Required]
        public bool IsAdmin { get; set; } = false;
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public bool IsBlocked { get; set; } = false;
        [Required]
        public string? Description { get; set; }
        [Required]
        public DateTime DateTime { get; set; } = DateTime.UtcNow;
        [Required]
        public string UserId { get; set; }
        public virtual ApplicationUser? User { get; set; }
    }
}
