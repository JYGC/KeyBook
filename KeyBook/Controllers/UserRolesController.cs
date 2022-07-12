using KeyBook.Constants;
using KeyBook.Models;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;

namespace KeyBook.Controllers
{
    [Authorize(Roles = "SuperAdmin")]
    public class UserRolesController : Controller
    {
        private readonly SignInManager<ApplicationUser> __signInManager;
        private readonly UserManager<ApplicationUser> __userManager;
        private readonly RoleManager<IdentityRole> __roleManager;

        public UserRolesController(SignInManager<ApplicationUser> signInManager, UserManager<ApplicationUser> userManager, RoleManager<IdentityRole> roleManager)
        {
            __signInManager = signInManager;
            __userManager = userManager;
            __roleManager = roleManager;
        }

        public async Task<IActionResult> Index(string userId)
        {
            List<UserRolesViewModel> userRolesViewModels = new List<UserRolesViewModel>();
            ApplicationUser user = await __userManager.FindByIdAsync(userId);
            foreach (IdentityRole role in __roleManager.Roles.Where(r => r.Name != Roles.SuperAdmin.ToString()).ToList())
            {
                UserRolesViewModel userRolesViewModel = new UserRolesViewModel
                {
                    RoleName = role.Name,
                };
                userRolesViewModel.Selected = await __userManager.IsInRoleAsync(user, role.Name);
                userRolesViewModels.Add(userRolesViewModel);
            }
            ManageUserRolesViewModel model = new ManageUserRolesViewModel()
            {
                UserId = userId,
                UserRoles = userRolesViewModels
            };
            return View(model);
        }

        public async Task<IActionResult> Manage(string id, ManageUserRolesViewModel model)
        {
            ApplicationUser user = await __userManager.FindByIdAsync(id);
            IList<string> roles = await __userManager.GetRolesAsync(user);
            IdentityResult result = await __userManager.RemoveFromRolesAsync(user, roles);
            ApplicationUser currentUser = await __userManager.GetUserAsync(User);
            await __signInManager.RefreshSignInAsync(currentUser);
            await Seeds.DefaultUsers.SeedSuperAdminAsync(__userManager, __roleManager); // fallback code in case one admin tries to change the roles of another - Find better implementation?
            return RedirectToAction("Index", new { userId = id });
        }
    }
}
