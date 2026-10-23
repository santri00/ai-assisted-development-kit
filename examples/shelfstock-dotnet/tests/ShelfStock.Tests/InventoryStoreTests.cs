using ShelfStock.Core;

namespace ShelfStock.Tests;

public sealed class InventoryStoreTests
{
    [Fact]
    public void Add_accumulates_quantity_and_keeps_existing_name()
    {
        var store = new InventoryStore();

        store.Add("SKU-1", "Widget", 5);
        var item = store.Add("SKU-1", null, 3);

        Assert.Equal(8, item.Quantity);
        Assert.Equal("Widget", item.Name);
    }

    [Fact]
    public void Add_rejects_empty_sku()
    {
        var store = new InventoryStore();

        Assert.Throws<InvalidSkuException>(() => store.Add(string.Empty, "n", 1));
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-2)]
    public void Add_rejects_non_positive_quantity(int quantity)
    {
        var store = new InventoryStore();

        Assert.Throws<InvalidQuantityException>(() => store.Add("A", "n", quantity));
    }

    [Fact]
    public void Remove_rejects_overselling_and_unknown_sku()
    {
        var store = new InventoryStore();
        store.Add("A", "Alpha", 4);

        Assert.Throws<InsufficientStockException>(() => store.Remove("A", 5));
        Assert.Throws<ItemNotFoundException>(() => store.Remove("missing", 1));
    }

    [Fact]
    public void Remove_can_take_stock_to_exactly_zero()
    {
        var store = new InventoryStore();
        store.Add("A", "Alpha", 4);

        var item = store.Remove("A", 4);

        Assert.Equal(0, item.Quantity);
    }

    [Fact]
    public void List_is_sorted_by_sku()
    {
        var store = new InventoryStore();
        store.Add("B", "b", 1);
        store.Add("A", "a", 1);

        var items = store.List();

        Assert.Equal(new[] { "A", "B" }, items.Select(i => i.Sku).ToArray());
    }
}
