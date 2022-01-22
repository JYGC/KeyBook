using System.ComponentModel.DataAnnotations;

namespace Backend.Models
{
    public class PersonDevice
    {
        [Key]
        public Guid Id { get; set; }
        [Required]
        public Person Person { get; set; }
        [Required]
        public Device Device { get; set; }
    }
}
