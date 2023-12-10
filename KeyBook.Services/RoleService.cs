using KeyBook.Constants;
using KeyBook.Models;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace KeyBook.Services
{
    public class RoleService
    {
        private readonly RoleManager<IdentityRole> __roleManager;
        private readonly UserManager<User> __userManager;
        private readonly IHttpContextAccessor __httpContextAccessor;
        public RoleService(RoleManager<IdentityRole> roleManager, UserManager<User> userManager, IHttpContextAccessor httpContextAccessor)
        {
            __roleManager = roleManager;
            __userManager = userManager;
            __httpContextAccessor = httpContextAccessor;
        }

        public async Task<List<IdentityRole>?> GetAllRoles()
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? currentUser = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (currentUser == null) return null;
            if (!(await __userManager.IsInRoleAsync(currentUser, Roles.SuperAdmin.ToString()))) return null;

            List<IdentityRole> roles = await __roleManager.Roles.ToListAsync();
            return roles;
        }

        public async Task<IdentityResult?> AddRole(string roleName)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? currentUser = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (currentUser == null) return null;
            if (!(await __userManager.IsInRoleAsync(currentUser, Roles.SuperAdmin.ToString()))) return null;

            IdentityResult? result = null;
            if (roleName != null) result = await __roleManager.CreateAsync(new IdentityRole(roleName.Trim()));
            return result;
        }
    }
}
