namespace ShelfStock.Core;

/// <summary>A stocked product identified by SKU.</summary>
public sealed record Item(string Sku, string Name, int Quantity);
