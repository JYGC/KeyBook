using KeyBook.Models;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace KeyBook.Controllers
{
    [Authorize]
    public class UserController : Controller
    {
        private readonly UserManager<User> __userManager;
        public UserController(UserManager<User> userManager)
        {
            __userManager = userManager;
        }

        public async Task<IActionResult> Index()
        {
            User currentUser = await __userManager.GetUserAsync(HttpContext.User);
            List<User> allUsersExceptCurrentUser = await __userManager.Users.Where(a => a.Id != currentUser.Id).ToListAsync();
            return View(allUsersExceptCurrentUser);
        }
    }
}
