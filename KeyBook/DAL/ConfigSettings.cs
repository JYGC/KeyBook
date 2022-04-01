namespace KeyBook.DAL
{
    public class ConfigSettings
    {
        public static string DefaultConnection { get; }
        static ConfigSettings()
        {
            string appsettings;
#if DEBUG
            appsettings = "appsettings.Development.json";
#else
            appsettings = "appsettings.json";
#endif
            ConfigurationBuilder configurationBuilder = new ConfigurationBuilder();
            string path = Path.Combine(Directory.GetCurrentDirectory(), appsettings);
            configurationBuilder.AddJsonFile(path, false);
            DefaultConnection = configurationBuilder.Build().GetSection("ConnectionStrings:DefaultConnection").Value;
        }
    }
}
