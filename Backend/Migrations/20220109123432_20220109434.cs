using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Migrations
{
    public partial class _20220109434 : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "IsDamaged",
                table: "Device");

            migrationBuilder.DropColumn(
                name: "IsLost",
                table: "Device");

            migrationBuilder.DropColumn(
                name: "IsRetired",
                table: "Device");

            migrationBuilder.AddColumn<int>(
                name: "Status",
                table: "Device",
                type: "int",
                nullable: false,
                defaultValue: 0);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "Status",
                table: "Device");

            migrationBuilder.AddColumn<bool>(
                name: "IsDamaged",
                table: "Device",
                type: "bit",
                nullable: false,
                defaultValue: false);

            migrationBuilder.AddColumn<bool>(
                name: "IsLost",
                table: "Device",
                type: "bit",
                nullable: false,
                defaultValue: false);

            migrationBuilder.AddColumn<bool>(
                name: "IsRetired",
                table: "Device",
                type: "bit",
                nullable: false,
                defaultValue: false);
        }
    }
}
