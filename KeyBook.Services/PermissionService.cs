using KeyBook.Constants;
using KeyBook.Models;
using KeyBook.Permission;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Identity;
using System.Security.Claims;

namespace KeyBook.Services
{
    public class PermissionService
    {
        private readonly RoleManager<IdentityRole> __roleManager;
        private readonly UserManager<User> __userManager;
        private readonly IHttpContextAccessor __httpContextAccessor;

        public PermissionService(RoleManager<IdentityRole> roleManager, UserManager<User> userManager, IHttpContextAccessor httpContextAccessor)
        {
            __roleManager = roleManager;
            __userManager = userManager;
            __httpContextAccessor = httpContextAccessor;
        }

        public async Task<List<RoleClaimsViewModel>?> GetAllPermissionsByRoleId(string roleId)
        {
            if (__httpContextAccessor.HttpContext == null) return null;
            User? currentUser = await __userManager.GetUserAsync(__httpContextAccessor.HttpContext.User);
            if (currentUser == null) return null;
            if (!(await __userManager.IsInRoleAsync(currentUser, Roles.SuperAdmin.ToString()))) return null;

            PermissionViewModel model = new PermissionViewModel();
            List<RoleClaimsViewModel> allPermissions = new List<RoleClaimsViewModel>();
            allPermissions.GetPermissions(typeof(Permissions.Products), roleId);
            IdentityRole role = await __roleManager.FindByIdAsync(roleId);
            model.RoleId = roleId;
            IList<Claim> claims = await __roleManager.GetClaimsAsync(role);
            List<string> allClaimValues = allPermissions.Select(a => a.Value).ToList();
            List<string> roleClaimValues = claims.Select(a => a.Value).ToList();
            List<string> authorizedClaims = allClaimValues.Intersect(roleClaimValues).ToList();
            foreach (var permission in allPermissions)
            {
                if (authorizedClaims.Any(a => a == permission.Value))
                {
                    permission.Selected = true;
                }
            }

            return allPermissions;
        }

        public async Task<string> UpdatePermissionForRoleId(PermissionViewModel model)
        {
            var role = await __roleManager.FindByIdAsync(model.RoleId);
            var claims = await __roleManager.GetClaimsAsync(role);
            foreach (var claim in claims)
            {
                await __roleManager.RemoveClaimAsync(role, claim);
            }
            var selectedClaims = model.RoleClaims.Where(a => a.Selected).ToList();
            foreach (var claim in selectedClaims)
            {
                await __roleManager.AddPermissionClaim(role, claim.Value);
            }
            return model.RoleId;
        }
    }
}
