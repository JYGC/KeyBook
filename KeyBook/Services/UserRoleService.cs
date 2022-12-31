using KeyBook.Constants;
using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Identity;

namespace KeyBook.Services
{
    public class UserRoleService
    {
        private readonly SignInManager<User> __signInManager;
        private readonly UserManager<User> __userManager;
        private readonly RoleManager<IdentityRole> __roleManager;

        public UserRoleService(SignInManager<User> signInManager, UserManager<User> userManager, RoleManager<IdentityRole> roleManager)
        {
            __signInManager = signInManager;
            __userManager = userManager;
            __roleManager = roleManager;
        }

        public async Task<List<UserRolesViewModel>?> GetUserRolesByUserId(string userId)
        {
            List<UserRolesViewModel> userRolesViewModels = new List<UserRolesViewModel>();
            User user = await __userManager.FindByIdAsync(userId);
            foreach (IdentityRole role in __roleManager.Roles.Where(r => r.Name != Roles.SuperAdmin.ToString()).ToList())
            {
                UserRolesViewModel userRolesViewModel = new UserRolesViewModel
                {
                    RoleName = role.Name,
                };
                userRolesViewModel.Selected = await __userManager.IsInRoleAsync(user, role.Name);
                userRolesViewModels.Add(userRolesViewModel);
            }
            return userRolesViewModels;
        }
    }
}
