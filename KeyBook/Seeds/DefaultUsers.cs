using KeyBook.Constants;
using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using System.Security.Claims;

namespace KeyBook.Seeds
{
    public static class DefaultUsers
    {
        public static async Task SeedBasicUserAsync(UserManager<ApplicationUser> userManager, RoleManager<IdentityRole> roleManager)
        {
            ApplicationUser defaultUser = new ApplicationUser
            {
                UserName = "seeduser@app.server",
                Email = "seeduser@app.server",
                EmailConfirmed = true
            };
            if (userManager.Users.All(u => u.Id != defaultUser.Id))
            {
                ApplicationUser user = await userManager.FindByEmailAsync(defaultUser.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(defaultUser, "&bC12357");
                    await userManager.AddToRoleAsync(defaultUser, Roles.Owner.ToString());
                }
            }
        }

        public static async Task SeedSuperAdminAsync(UserManager<ApplicationUser> userManager, RoleManager<IdentityRole> roleManager)
        {
            var defaultUser = new ApplicationUser
            {
                UserName = "superadmin@gmail.com",
                Email = "superadmin@gmail.com",
                EmailConfirmed = true
            };
            if (userManager.Users.All(u => u.Id != defaultUser.Id))
            {
                var user = await userManager.FindByEmailAsync(defaultUser.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(defaultUser, "Test3r$");
                    await userManager.AddToRoleAsync(defaultUser, Roles.Manager.ToString());
                    await userManager.AddToRoleAsync(defaultUser, Roles.Owner.ToString());
                    await userManager.AddToRoleAsync(defaultUser, Roles.Admin.ToString());
                    await userManager.AddToRoleAsync(defaultUser, Roles.SuperAdmin.ToString());
                }
                await roleManager.SeedClaimsForSuperAdmin();
            }
        }

        private async static Task SeedClaimsForSuperAdmin(this RoleManager<IdentityRole> roleManager)
        {
            var adminRole = await roleManager.FindByNameAsync("SuperAdmin");
            await roleManager.AddPermissionClaim(adminRole, "Products");
        }

        public static async Task AddPermissionClaim(this RoleManager<IdentityRole> roleManager, IdentityRole role, string module)
        {
            var allClaims = await roleManager.GetClaimsAsync(role);
            var allPermissions = Permissions.GeneratePermissionsForModule(module);
            foreach (var permission in allPermissions)
            {
                if (!allClaims.Any(a => a.Type == "Permission" && a.Value == permission))
                {
                    await roleManager.AddClaimAsync(role, new Claim("Permission", permission));
                }
            }
        }
    }
}
