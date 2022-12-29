using KeyBook.Constants;
using Microsoft.AspNetCore.Mvc;

namespace KeyBook.Controllers
{
    public class HomeController : Controller
    {
        [Route("Old")]
        [Route("OldHome")]
        public IActionResult Index()
        {
            if (User.IsInRole(Roles.SuperAdmin.ToString())) return RedirectToAction("Index", "User");
            if (User.IsInRole(Roles.Owner.ToString())) return RedirectToAction("Index", "Device");
            return Redirect("/Identity/Account/Login");
        }
    }
}
