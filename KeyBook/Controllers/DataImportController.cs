using ExcelDataReader;
using Microsoft.AspNetCore.Mvc;
using System.Data;

namespace KeyBook.Controllers
{
    public class DataImportController : Controller
    {
        private const string __UPLOAD_FOLDER = "Uploads";

        private IHostEnvironment __environment;

        public DataImportController(IHostEnvironment environment)
        {
            __environment = environment;
        }

        public IActionResult Index()
        {
            return View();
        }

        [HttpPost]
        //[ValidateAntiForgeryToken] - Add XSRF protection later 
        public IActionResult Excel(IFormFile postedFile)
        {
            if (postedFile == null || (!postedFile.FileName.EndsWith(".xls") && !postedFile.FileName.EndsWith(".xlsx"))) return NotFound();
            System.Text.Encoding.RegisterProvider(System.Text.CodePagesEncodingProvider.Instance);
            using MemoryStream stream = new MemoryStream();
            postedFile.CopyTo(stream);
            using IExcelDataReader reader = ExcelReaderFactory.CreateReader(stream);
            DataTable dt = new DataTable();
            DataRow row;
            DataTable dt_ = new DataTable();
            try
            {
                dt_ = reader.AsDataSet().Tables[0];
                for (int i = 0; i < dt_.Columns.Count; i++) dt.Columns.Add(dt_.Rows[0][i].ToString());
                int rowcounter = 0;
                for (int row_ = 1; row_ < dt.Rows.Count; row_++)
                {
                    row = dt.NewRow();
                    for (int col = 0; col < dt_.Columns.Count; col++)
                    {
                        row[col] = dt_.Rows[row_][col].ToString();
                        rowcounter++;
                    }
                    dt.Rows.Add(row);
                }
            }
            catch (Exception ex) // Error handling
            {
                ModelState.AddModelError("File", "Unable to Upload file!");
                return NotFound();
            }
            DataSet result = new DataSet();
            result.Tables.Add(dt);
            reader.Close();
            DataTable tmp = result.Tables[0];
            //HttpContext.Session. ["tmpdata"] = tmp;  //store datatable into session
            return RedirectToAction("Index", "DataImport");
        }
    }
}
