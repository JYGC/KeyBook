﻿using KeyBook.Models;
using Microsoft.AspNetCore.Identity;

namespace KeyBook.Seeds
{
    public class DefaultRoles
    {
        public static async Task SeedAsync(UserManager<ApplicationUser> userManager, RoleManager<IdentityRole> roleManager)
        {
            await roleManager.CreateAsync(new IdentityRole(UserRoles.SuperAdmin.ToString()));
            await roleManager.CreateAsync(new IdentityRole(UserRoles.Admin.ToString()));
            await roleManager.CreateAsync(new IdentityRole(UserRoles.Owner.ToString()));
            await roleManager.CreateAsync(new IdentityRole(UserRoles.Manager.ToString()));
        }
    }
}
