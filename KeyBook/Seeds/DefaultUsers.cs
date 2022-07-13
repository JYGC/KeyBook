using KeyBook.Constants;
using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using System.Security.Claims;

namespace KeyBook.Seeds
{
    public static class DefaultUsers
    {
        public static async Task SeedBasicUserAsync(UserManager<User> userManager, RoleManager<IdentityRole> roleManager)
        {
            User defaultUser = new User
            {
                UserName = "seeduser@app.server",
                Name = "Sean Ulysses",
                Email = "seeduser@app.server",
                EmailConfirmed = true,
                Organization = new Organization
                {
                    Name = "Seed Organization"
                }
            };
            if (userManager.Users.All(u => u.Id != defaultUser.Id))
            {
                User user = await userManager.FindByEmailAsync(defaultUser.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(defaultUser, "&bC12357");
                    await userManager.AddToRoleAsync(defaultUser, Roles.Owner.ToString());
                }
            }
        }

        public static async Task SeedSuperAdminAsync(UserManager<User> userManager, RoleManager<IdentityRole> roleManager)
        {
            User superAdmin = new User
            {
                UserName = "superadmin@app.server",
                Name = "SuperAdmin",
                Email = "superadmin@app.server",
                EmailConfirmed = true,
                Organization = new Organization
                {
                    Name = "SuperAdmin"
                }
            };
            if (userManager.Users.All(u => u.Id != superAdmin.Id))
            {
                User? user = await userManager.FindByEmailAsync(superAdmin.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(superAdmin, "$uperUs3r");
                    await userManager.AddToRoleAsync(superAdmin, Roles.Manager.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.Owner.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.Admin.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.SuperAdmin.ToString());
                }
                await roleManager.SeedClaimsForSuperAdmin();
            }
        }

        private async static Task SeedClaimsForSuperAdmin(this RoleManager<IdentityRole> roleManager)
        {
            IdentityRole? adminRole = await roleManager.FindByNameAsync("SuperAdmin");
            await roleManager.AddPermissionClaim(adminRole, "Products");
        }

        public static async Task AddPermissionClaim(this RoleManager<IdentityRole> roleManager, IdentityRole role, string module)
        {
            IList<Claim>? allClaims = await roleManager.GetClaimsAsync(role);
            List<string>? allPermissions = Permissions.GeneratePermissionsForModule(module);
            foreach (string permission in allPermissions)
            {
                if (!allClaims.Any(a => a.Type == "Permission" && a.Value == permission))
                {
                    await roleManager.AddClaimAsync(role, new Claim("Permission", permission));
                }
            }
        }
    }
}
