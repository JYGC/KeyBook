using Microsoft.AspNetCore.Authorization;
using System.Security.Claims;

namespace KeyBook.Permission
{
    internal class PermissionAuthorizationHandler : AuthorizationHandler<PermissionRequirement>
    {
        public PermissionAuthorizationHandler() { }

        protected override async Task HandleRequirementAsync(AuthorizationHandlerContext context, PermissionRequirement requirement)
        {
            if (context.User == null)
            {
                return;
            }
            IEnumerable<Claim> permissionss = context.User.Claims.Where(
                x => x.Type == "Permission" &&
                     x.Value == requirement.Permission &&
                     x.Issuer == "LOCAL AUTHORITY");
            if (permissionss.Any())
            {
                context.Succeed(requirement);
                return;
            }
        }
    }
}
