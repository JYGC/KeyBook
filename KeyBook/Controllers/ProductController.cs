using Microsoft.AspNetCore.Mvc;

namespace KeyBook.Controllers
{
    public class ProductController : Controller
    {
        public IActionResult Index()
        {
            return View();
        }
    }
}
