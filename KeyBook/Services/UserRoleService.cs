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
        private readonly IHttpContextAccessor __httpContextAccessor;

        public UserRoleService(SignInManager<User> signInManager, UserManager<User> userManager, RoleManager<IdentityRole> roleManager, IHttpContextAccessor httpContextAccessor)
        {
            __signInManager = signInManager;
            __userManager = userManager;
            __roleManager = roleManager;
            __httpContextAccessor = httpContextAccessor;
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

        public async Task<(bool, string?)> ManageUserRolesForUser(string id, ManageUserRolesViewModel model)
        {
            try
            {
                User user = await __userManager.FindByIdAsync(id);
                IList<string> roles = await __userManager.GetRolesAsync(user);
                IdentityResult result = await __userManager.RemoveFromRolesAsync(user, roles);
                if (__httpContextAccessor.HttpContext == null) throw new Exception("Cannot find HTTP context");
                User currentUser = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
                await __signInManager.RefreshSignInAsync(currentUser);
                await Seeds.DefaultUsers.SeedSuperAdminAsync(__userManager, __roleManager); // fallback code in case one admin tries to change the roles of another - Find better implementation?
                return (true, null);
            }
            catch (Exception ex)
            {
                return (false, ex.Message);
            }
        }

        public async Task<(bool, string?)> UpdateUserRolesForUser(string id, ManageUserRolesViewModel model)
        {
            try
            {
                User user = await __userManager.FindByIdAsync(id);
                foreach (UserRolesViewModel userRolesViewModel in model.UserRoles)
                {
                    if (userRolesViewModel.Selected)
                    {
                        await __userManager.AddToRoleAsync(user, userRolesViewModel.RoleName);
                    }
                    else
                    {
                        await __userManager.RemoveFromRoleAsync(user, userRolesViewModel.RoleName);
                    }
                }
                return (true, null);
            }
            catch (Exception ex)
            {
                return (false, ex.Message);
            }
        }
    }
}
