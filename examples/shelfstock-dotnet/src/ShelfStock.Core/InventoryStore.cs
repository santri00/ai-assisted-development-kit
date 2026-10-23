namespace ShelfStock.Core;

/// <summary>In-memory, thread-safe inventory. All access to the dictionary happens under <c>_gate</c>.</summary>
public sealed class InventoryStore
{
    private readonly object _gate = new();
    private readonly Dictionary<string, Item> _items = new(StringComparer.Ordinal);

    /// <summary>Increases stock for <paramref name="sku"/>, creating the item if needed.</summary>
    public Item Add(string sku, string? name, int quantity)
    {
        if (string.IsNullOrEmpty(sku))
        {
            throw new InvalidSkuException();
        }

        if (quantity <= 0)
        {
            throw new InvalidQuantityException();
        }

        lock (_gate)
        {
            _items.TryGetValue(sku, out var existing);
            var resolvedName = string.IsNullOrEmpty(name) ? existing?.Name ?? string.Empty : name;
            var updated = new Item(sku, resolvedName, (existing?.Quantity ?? 0) + quantity);
            _items[sku] = updated;
            return updated;
        }
    }

    /// <summary>Decreases stock. Never lets stock go negative.</summary>
    public Item Remove(string sku, int quantity)
    {
        if (quantity <= 0)
        {
            throw new InvalidQuantityException();
        }

        lock (_gate)
        {
            if (!_items.TryGetValue(sku, out var existing))
            {
                throw new ItemNotFoundException(sku);
            }

            if (existing.Quantity < quantity)
            {
                throw new InsufficientStockException(existing.Quantity, quantity);
            }

            var updated = existing with { Quantity = existing.Quantity - quantity };
            _items[sku] = updated;
            return updated;
        }
    }

    public Item Get(string sku)
    {
        lock (_gate)
        {
            return _items.TryGetValue(sku, out var item) ? item : throw new ItemNotFoundException(sku);
        }
    }

    /// <summary>All items sorted by SKU (ordinal).</summary>
    public IReadOnlyList<Item> List()
    {
        lock (_gate)
        {
            return _items.Values.OrderBy(i => i.Sku, StringComparer.Ordinal).ToList();
        }
    }
}
