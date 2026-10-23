using System.Net;
using System.Net.Http.Json;
using System.Text;
using Microsoft.AspNetCore.Mvc.Testing;
using ShelfStock.Core;

namespace ShelfStock.Tests;

public sealed class ApiTests
{
    // Each test builds its own in-process server so the singleton InventoryStore is never shared.
    private static async Task<HttpResponseMessage> PostJsonAsync(HttpClient client, string path, string json)
    {
        using var content = new StringContent(json, Encoding.UTF8, "application/json");
        return await client.PostAsync(path, content);
    }

    [Fact]
    public async Task Add_then_get_returns_the_item()
    {
        await using var factory = new WebApplicationFactory<Program>();
        using var client = factory.CreateClient();

        using var add = await PostJsonAsync(client, "/items/A/add", """{"name":"Alpha","quantity":3}""");
        var item = await client.GetFromJsonAsync<Item>("/items/A");

        Assert.Equal(HttpStatusCode.OK, add.StatusCode);
        Assert.NotNull(item);
        Assert.Equal(3, item.Quantity);
        Assert.Equal("Alpha", item.Name);
    }

    [Fact]
    public async Task Healthz_returns_ok()
    {
        await using var factory = new WebApplicationFactory<Program>();
        using var client = factory.CreateClient();

        using var response = await client.GetAsync("/healthz");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
    }

    [Theory]
    [InlineData("GET", "/items/nope", "", 404)]
    [InlineData("POST", "/items/A/add", "{", 400)]
    [InlineData("POST", "/items/A/add", """{"quantity":0}""", 400)]
    [InlineData("POST", "/items/A/add", """{"quantity":1,"unknown":true}""", 400)]
    [InlineData("POST", "/items/A/remove", """{"quantity":9}""", 409)]
    public async Task Returns_expected_status_code(string method, string path, string body, int expected)
    {
        await using var factory = new WebApplicationFactory<Program>();
        using var client = factory.CreateClient();
        using var seed = await PostJsonAsync(client, "/items/A/add", """{"quantity":1}""");

        using var request = new HttpRequestMessage(new HttpMethod(method), path);
        if (method == "POST")
        {
            request.Content = new StringContent(body, Encoding.UTF8, "application/json");
        }

        using var response = await client.SendAsync(request);

        Assert.Equal(expected, (int)response.StatusCode);
    }
}
