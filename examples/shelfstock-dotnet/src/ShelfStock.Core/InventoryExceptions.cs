namespace ShelfStock.Core;

/// <summary>Base type for all expected domain failures. The API maps these to HTTP status codes.</summary>
public abstract class InventoryException : Exception
{
    protected InventoryException(string message)
        : base(message)
    {
    }
}

public sealed class InvalidSkuException : InventoryException
{
    public InvalidSkuException()
        : base("sku must not be empty")
    {
    }
}

public sealed class InvalidQuantityException : InventoryException
{
    public InvalidQuantityException()
        : base("quantity must be positive")
    {
    }
}

public sealed class ItemNotFoundException : InventoryException
{
    public ItemNotFoundException(string sku)
        : base($"item not found: {sku}")
    {
    }
}

public sealed class InsufficientStockException : InventoryException
{
    public InsufficientStockException(int have, int want)
        : base($"insufficient stock: have {have}, want {want}")
    {
    }
}
