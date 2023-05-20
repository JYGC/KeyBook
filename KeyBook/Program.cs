using KeyBook.Services;
using KeyBook.Models;
using KeyBook.Permission;
using KeyBook.Security;
using KeyBook.Seeds;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Migrations;
using System.Text.Json.Serialization;
using KeyBook.Database;

WebApplicationBuilder? builder = WebApplication.CreateBuilder(args);

builder.Services.AddRazorPages();
builder.Services.AddServerSideBlazor();

// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
//builder.Services.AddSingleton<IAuthorizationPolicyProvider, PermissionPolicyProvider>(); // Possible future permissions implementation?
builder.Services.AddScoped<IAuthorizationHandler, PermissionAuthorizationHandler>();

// Database
builder.Services.AddDbContext<KeyBookDbContext>(options => options.UseNpgsql(
    builder.Configuration.GetConnectionString("DefaultConnection"),
    x => x.MigrationsHistoryTable("__efmigrationshistory", "public")
).ReplaceService<IHistoryRepository, LoweredCaseMigrationHistoryRepository>());

// Authentication and roles
builder.Services.AddIdentity<User, IdentityRole>()
    .AddEntityFrameworkStores<KeyBookDbContext>()
    .AddDefaultUI()
    .AddDefaultTokenProviders()
    .AddRoles<IdentityRole>();

// Data services
builder.Services.AddScoped<DeviceService>();
builder.Services.AddScoped<PersonService>();
builder.Services.AddScoped<UserService>();
builder.Services.AddScoped<RoleService>();
builder.Services.AddScoped<PermissionService>();
builder.Services.AddScoped<UserRoleService>();

// Injection for HttpContext
builder.Services.AddHttpContextAccessor();

// Add AntiforgeryToken Provider
builder.Services.AddScoped<TokenProvider>();


// Initialize App
WebApplication? app = builder.Build();

// Seed default data went database is empty
using (IServiceScope scope = app.Services.CreateScope())
{
    IServiceProvider services = scope.ServiceProvider;
    UserManager<User> userManager = services.GetRequiredService<UserManager<User>>();
    RoleManager<IdentityRole> roleManager = services.GetRequiredService<RoleManager<IdentityRole>>();
    await DefaultRoles.SeedAsync(userManager, roleManager);
    await DefaultUsers.SeedBasicUserAsync(userManager, roleManager);
    await DefaultUsers.SeedSuperAdminAsync(userManager, roleManager);
    await DefaultData.SeedAsync(services);
}

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Home/Error");
    // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
    app.UseHsts();
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseStaticFiles();

app.UseRouting();
app.UseAuthentication();

app.UseAuthorization();

app.MapBlazorHub();
app.MapRazorPages();
app.MapFallbackToPage("/_Host");

app.Run();
