# syntax=docker/dockerfile:1
# Multi-stage build. Targets:
#   test     -> format check + tests (used by `task docker:check`; fails the build on any failure)
#   runtime  -> small ASP.NET runtime image, non-root (default final stage)

ARG DOTNET_VERSION=10.0

FROM mcr.microsoft.com/dotnet/sdk:${DOTNET_VERSION} AS base
WORKDIR /src
# Copy only project metadata first so the restore layer is cached until dependencies change.
COPY global.json Directory.Build.props .editorconfig ShelfStock.slnx ./
COPY src/ShelfStock.Core/ShelfStock.Core.csproj src/ShelfStock.Core/
COPY src/ShelfStock.Api/ShelfStock.Api.csproj src/ShelfStock.Api/
COPY tests/ShelfStock.Tests/ShelfStock.Tests.csproj tests/ShelfStock.Tests/
RUN dotnet restore
COPY . .

FROM base AS test
RUN dotnet format --verify-no-changes --no-restore && dotnet test --no-restore --nologo

FROM base AS publish
RUN dotnet publish src/ShelfStock.Api/ShelfStock.Api.csproj -c Release --no-restore -o /out /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:${DOTNET_VERSION} AS runtime
WORKDIR /app
COPY --from=publish /out ./
# APP_UID is predefined in Microsoft's .NET 8+ images (non-root user).
USER $APP_UID
ENV ASPNETCORE_HTTP_PORTS=8080
EXPOSE 8080
ENTRYPOINT ["dotnet", "ShelfStock.Api.dll"]
