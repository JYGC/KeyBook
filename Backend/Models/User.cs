﻿using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class User
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public string Name { get; set; }
        [Required]
        public string Email { get; set; }
        [Required]
        public bool IsAdmin { get; set; } = false;
        [Required]
        public bool IsDeleted { get; set; } = false;
        [Required]
        public bool IsBlocked { get; set; } = false;
        public virtual ICollection<Person> Persons { get; set; }
        public virtual ICollection<Device> Devices { get; set; }
        public virtual ICollection<UserHistory> UserHistories { get; set; }
    }
}
