using KeyBook.Constants;
using KeyBook.Models;
using Microsoft.AspNetCore.Identity;
using System.Security.Claims;

namespace KeyBook.Seeds
{
    public static class DefaultUsers
    {
        // Default seed user passwords - CHANGE THESE ON WEB IN PRODUCTION!!!!
        private const string __SeedUserPassword = "&bC12357";
        private const string __DefaultSuperAdminPassword = "$uperUs3r";

        // Seed user details
        public const string SeedUserEmail = "seeduser@app.server";
        public const string SeedUsersName = "Sean Ulysses";
        public const string SeedUsersOrganizationName = "Seed Organization";
        // SuperAdmin details
        public const string SuperAdminEmail = "superadmin@app.server";
        public const string SuperAdminName = "SuperAdmin";

        public static async Task SeedBasicUserAsync(UserManager<User> userManager, RoleManager<IdentityRole> roleManager)
        {
            User seedUser = new User
            {
                UserName = SeedUserEmail,
                Name = SeedUsersName,
                Email = SeedUserEmail,
                EmailConfirmed = true,
                Organization = new Organization
                {
                    Name = SeedUsersOrganizationName // Normal users use same name
                }
            };
            if (userManager.Users.All(u => u.Id != seedUser.Id))
            {
                User user = await userManager.FindByEmailAsync(seedUser.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(seedUser, __SeedUserPassword);
                    await userManager.AddToRoleAsync(seedUser, Roles.Owner.ToString());
                }
            }
        }

        public static async Task SeedSuperAdminAsync(UserManager<User> userManager, RoleManager<IdentityRole> roleManager)
        {
            User superAdmin = new User
            {
                UserName = SuperAdminEmail,
                Name = SuperAdminName,
                Email = SuperAdminEmail,
                EmailConfirmed = true,
                Organization = new Organization
                {
                    Name = SuperAdminName
                }
            };
            if (userManager.Users.All(u => u.Id != superAdmin.Id))
            {
                User? user = await userManager.FindByEmailAsync(superAdmin.Email);
                if (user == null)
                {
                    await userManager.CreateAsync(superAdmin, __DefaultSuperAdminPassword);
                    await userManager.AddToRoleAsync(superAdmin, Roles.Manager.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.Owner.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.Admin.ToString());
                    await userManager.AddToRoleAsync(superAdmin, Roles.SuperAdmin.ToString());
                }
                await roleManager.SeedClaimsForSuperAdmin(); // NOT USED - Saved for possible implementation of permissions
            }
        }

        private async static Task SeedClaimsForSuperAdmin(this RoleManager<IdentityRole> roleManager)
        {
            IdentityRole? adminRole = await roleManager.FindByNameAsync(SuperAdminName);
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
