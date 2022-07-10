using KeyBook.Constants;
using KeyBook.Permission;
using KeyBook.ViewModels;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using System.Security.Claims;

namespace KeyBook.Controllers
{
    [Authorize(Roles = "SuperAdmin")]
    public class PermissionController : Controller
    {
        private readonly RoleManager<IdentityRole> __roleManager;

        public PermissionController(RoleManager<IdentityRole> roleManager)
        {
            __roleManager = roleManager;
        }

        public async Task<ActionResult> Index(string roleId)
        {
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
            model.RoleClaims = allPermissions;
            return View(model);
        }

        public async Task<IActionResult> Update(PermissionViewModel model)
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
            return RedirectToAction("Index", new { roleId = model.RoleId });
        }
    }
}
