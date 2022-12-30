using KeyBook.Constants;
using KeyBook.Database;
using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace KeyBook.Services
{
    public class UserService
    {
        private readonly UserManager<User> __userManager;
        private readonly IHttpContextAccessor __httpContextAccessor;

        public UserService(UserManager<User> userManager, IHttpContextAccessor httpContextAccessor)
        {
            __userManager = userManager;
            __httpContextAccessor = httpContextAccessor;
        }

        public async Task<List<User>?> GetUsersExceptCurrentUser()
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? currentUser = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (currentUser == null) return null;
            if (!(await __userManager.IsInRoleAsync(currentUser, Roles.SuperAdmin.ToString()))) return null;

            List<User> allUsersExceptCurrentUser = await __userManager.Users.Where(a => a.Id != currentUser.Id).ToListAsync();
            return allUsersExceptCurrentUser;
        }
    }
}
