using System.Text.Json.Serialization;
using ShelfStock.Api;
using ShelfStock.Core;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton<InventoryStore>();

// Reject unknown JSON fields (mirrors the Go service) and always answer malformed bodies with 400.
builder.Services.ConfigureHttpJsonOptions(options =>
{
    options.SerializerOptions.UnmappedMemberHandling = JsonUnmappedMemberHandling.Disallow;
});
builder.Services.Configure<RouteHandlerOptions>(options =>
{
    options.ThrowOnBadRequest = false;
});

// Cap request bodies at 1 MiB.
builder.WebHost.ConfigureKestrel(options =>
{
    options.Limits.MaxRequestBodySize = 1024 * 1024;
});

var app = builder.Build();

app.MapInventory();

app.Run();

// Makes the entry point visible to WebApplicationFactory<Program> in the test project.
public partial class Program
{
}
