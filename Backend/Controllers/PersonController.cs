#nullable disable
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Backend.Models;

namespace Backend.Controllers
{
    [Route("[controller]")]
    [ApiController]
    public class PersonController : ControllerBase
    {
        private readonly KeyBookDbContext _context;

        public PersonController(KeyBookDbContext context)
        {
            _context = context;
        }

        // GET: Person
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Person>>> GetPersons()
        {
            return await _context.Persons.ToListAsync();
        }

        // GET: Person/allforuser
        [HttpGet("allforuser")]
        public async Task<ActionResult<IEnumerable<Person>>> GetAllForUser()
        {
            User user = _context.Users.FirstOrDefault(u => u.Name == "Administrator");
            return await _context.Persons.Where(
                person => person.UserId == user.Id && !person.IsDeleted
            ).ToListAsync(); // Add filter user when adding auth
        }

        // GET: Person/view/id/149BE541-9271-4E3B-8766-08D9D36C9255
        [HttpGet("view/id/{id}")]
        public async Task<ActionResult<Person>> GetPerson(Guid id)
        {
            Person person = await _context.Persons.FindAsync(id);

            if (person == null)
            {
                return NotFound();
            }
            User user = _context.Users.FirstOrDefault(u => u.Name == "Administrator");
            person.PersonDevices = (from device in _context.Devices
                                    join personDevice in _context.PersonDevices on device.Id equals personDevice.DeviceId
                                    where device.UserId == user.Id && personDevice.PersonId == person.Id
                                    select new PersonDevice
                                    {
                                        Id = personDevice.Id,
                                        PersonId = personDevice.PersonId,
                                        DeviceId = personDevice.DeviceId,
                                        Device = device,
                                        IsDeleted = personDevice.IsDeleted,
                                    }).ToList();

            return person;
        }

        // POST: Person/add
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost("add")]
        public async Task<ActionResult<Person>> PostPerson(Person person)
        {
            User user = await _context.Users.FirstOrDefaultAsync(u => u.Name == "Administrator"); //replace this - Authentication
            person.User = user;
            person.PersonHistories.Add(new PersonHistory
            {
                Name = person.Name,
                IsGone = person.IsGone,
                Type = person.Type,
                IsDeleted = person.IsDeleted,
                Description = "create new person",
                Person = person
            });
            _context.Persons.Add(person);
            await _context.SaveChangesAsync();

            return CreatedAtAction("GetPerson", new { id = person.Id }, person);
        }

        // PUT: api/Person/5
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPut("{id}")]
        public async Task<IActionResult> PutPerson(Guid id, Person person)
        {
            if (id != person.Id)
            {
                return BadRequest();
            }

            _context.Entry(person).State = EntityState.Modified;

            try
            {
                await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException)
            {
                if (!PersonExists(id))
                {
                    return NotFound();
                }
                else
                {
                    throw;
                }
            }

            return NoContent();
        }

        // DELETE: api/Person/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeletePerson(Guid id)
        {
            var person = await _context.Persons.FindAsync(id);
            if (person == null)
            {
                return NotFound();
            }

            _context.Persons.Remove(person);
            await _context.SaveChangesAsync();

            return NoContent();
        }

        private bool PersonExists(Guid id)
        {
            return _context.Persons.Any(e => e.Id == id);
        }
    }
}
