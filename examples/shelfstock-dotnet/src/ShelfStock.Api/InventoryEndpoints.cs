using ShelfStock.Core;

namespace ShelfStock.Api;

public sealed record AddRequest(string? Name, int Quantity);

public sealed record RemoveRequest(int Quantity);

/// <summary>HTTP transport for the inventory. Domain exceptions are mapped to status codes ONLY here.</summary>
public static class InventoryEndpoints
{
    public static IEndpointRouteBuilder MapInventory(this IEndpointRouteBuilder routes)
    {
        routes.MapGet("/healthz", () => Results.Ok(new { status = "ok" }));

        routes.MapGet("/items", (InventoryStore store) => Results.Ok(store.List()));

        routes.MapGet(
            "/items/{sku}",
            (string sku, InventoryStore store) => Handle(() => Results.Ok(store.Get(sku))));

        routes.MapPost(
            "/items/{sku}/add",
            (string sku, AddRequest body, InventoryStore store) =>
                Handle(() => Results.Ok(store.Add(sku, body.Name, body.Quantity))));

        routes.MapPost(
            "/items/{sku}/remove",
            (string sku, RemoveRequest body, InventoryStore store) =>
                Handle(() => Results.Ok(store.Remove(sku, body.Quantity))));

        return routes;
    }

    private static IResult Handle(Func<IResult> action)
    {
        try
        {
            return action();
        }
        catch (InventoryException ex)
        {
            return Results.Json(new { error = ex.Message }, statusCode: StatusFor(ex));
        }
    }

    private static int StatusFor(InventoryException ex) => ex switch
    {
        ItemNotFoundException => StatusCodes.Status404NotFound,
        InsufficientStockException => StatusCodes.Status409Conflict,
        InvalidQuantityException or InvalidSkuException => StatusCodes.Status400BadRequest,
        _ => StatusCodes.Status500InternalServerError,
    };
}
