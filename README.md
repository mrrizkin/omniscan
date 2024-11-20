# Boot

## Overview

Boot is a sophisticated, full-stack web application framework that combines the power of Go with modern frontend technologies. Built on top of go-fiber and Uber's go-fx dependency injection framework, it follows a Laravel-inspired architecture while maintaining Go idioms. Boot provides a robust foundation for building scalable, maintainable web applications with a focus on developer productivity.

## Key Features

-   **Dependency Injection**: Powered by Uber's go-fx for clean, testable architecture
-   **Laravel-Inspired Structure**: Familiar directory organization for rapid development
-   **Modern Frontend**: Integrated Vite bundler with TailwindCSS and Alpine.js
-   **HTMX Integration**: Seamless server-side interactions
-   **Flexible Database**: Support for MySQL, PostgreSQL, and SQLite via GORM
-   **Advanced Logging**: Structured logging with zerolog
-   **Session Management**: Multiple storage backends (Database, File, Memory)
-   **View Rendering**: Gonja templating engine with Laravel-like blade syntax
-   **Form Validation**: Built-in request validation
-   **Task Scheduling**: Integrated scheduler for background jobs
-   **Asset Management**: Vite integration with hot module replacement

## Project Structure

```
.
├── app/                     # Application core
│   ├── console/             # Console commands and scheduled tasks
│   ├── controllers/         # HTTP request handlers
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Database models and schemas
│   ├── providers/           # Service providers and core components
│   ├── repositories/        # Data access layer
│   └── services/            # Business logic layer
├── bootstrap/               # Application bootstrapping
├── config/                  # Configuration files
├── pkg/                     # Reusable packages
├── public/                  # Static assets
├── resources/               # Frontend resources
│   ├── css/                 # Stylesheets
│   ├── js/                  # JavaScript files
│   └── views/               # Template files
├── routes/                  # Route definitions
├── storage/                 # Application storage
└── tests/                   # Test suites
```

## Getting Started

### Prerequisites

-   Go (version 1.22.2 or later)
-   Node.js
-   PNPM (Package Manager)
-   Air (Go live reload)
-   Swag (Swagger documentation generator)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/mrrizkin/omniscan.git
    cd boot
    ```

2. Run the setup script to install all dependencies:
    ```bash
    pnpm run setup
    ```
    This will:
    - Install Node.js dependencies via PNPM
    - Download Go module dependencies

### Development Scripts

```bash
# Setup project
pnpm run setup                # Install all dependencies

# Development
pnpm run dev                  # Start development server with hot reload
pnpm run dev:assets           # Start Vite dev server for assets
pnpm run dev:app              # Start Go server with Air for hot reload

# Build
pnpm run build               # Build both frontend and backend
pnpm run build:assets        # Build frontend assets with Vite
pnpm run build:app           # Build Go application

# Documentation
pnpm run generate:swagger    # Generate Swagger/OpenAPI documentation
```

## Frontend Tooling

Boot comes with a modern frontend development stack:

-   **Vite**: Fast frontend tooling and HMR
-   **TailwindCSS**: Utility-first CSS framework
-   **Alpine.js**: Lightweight JavaScript framework
-   **Prettier**: Code formatting with Django/Alpine.js support
-   **PostCSS**: CSS processing and autoprefixer

### Frontend Dependencies

```json
{
    "dependencies": {
        "axios": "^1.7.7" // HTTP client
    },
    "devDependencies": {
        "autoprefixer": "^10.4.20", // PostCSS autoprefixer
        "concurrently": "^9.1.0", // Run multiple commands
        "postcss": "^8.4.49", // CSS processing
        "prettier": "^3.3.3", // Code formatting
        "prettier-plugin-django-alpine": "^1.3.0", // Template formatting
        "tailwindcss": "^3.4.15", // Utility CSS framework
        "vite": "^5.4.11", // Frontend tooling
        "vite-plugin-backend": "^1.0.0", // Backend integration
        "vite-plugin-full-reload": "^1.2.0" // Full page reload
    }
}
```

## Development Workflow

1. **Start Development Server**

    ```bash
    pnpm run dev
    ```

    This will:

    - Start Vite dev server for frontend assets
    - Wait 5 seconds for Vite to initialize
    - Start Go server with Air for hot reloading

2. **Working with Assets**

    - Frontend assets are in `resources/`:
        ```
        resources/
        ├── css/
        │   └── index.css      # Main CSS file
        ├── js/
        │   ├── app.js         # Main JavaScript
        │   └── bootstrap.js   # JS initialization
        └── views/             # Template files
        ```
    - Changes to CSS/JS files trigger hot reload
    - Template changes trigger full page reload

3. **API Documentation**

    ```bash
    pnpm run generate:swagger
    ```

    - Generates OpenAPI documentation in `public/docs`
    - Parses internal types and comments
    - Outputs in JSON format

4. **Building for Production**
    ```bash
    pnpm run build
    ```
    This will:
    - Build and optimize frontend assets
    - Compile Go application
    - Output production-ready files

## Configuration Files

-   `vite.config.js`: Vite configuration
-   `tailwind.config.js`: TailwindCSS configuration
-   `postcss.config.js`: PostCSS plugins
-   `.prettierrc`: Prettier formatting rules

## Architecture Overview

Boot follows a clean architecture pattern with clear separation of concerns, leveraging go-fx for dependency injection and automatic component loading.

### Controllers

Controllers in Boot are automatically loaded using go-fx. Here's an example of a typical controller:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mrrizkin/omniscan/app/providers/app"
)

type WelcomeController struct {
    *app.App
}

func (*WelcomeController) Construct() interface{} {
    return func(app *app.App) (*WelcomeController, error) {
        return &WelcomeController{
            App: app,
        }, nil
    }
}

func (c *WelcomeController) Index(ctx *fiber.Ctx) error {
    return c.Render(ctx, "pages/welcome", nil)
}
```

Controllers are automatically registered using the `constructor.Load` helper:

```go
package controllers

import (
    "go.uber.org/fx"
    "github.com/mrrizkin/omniscan/pkg/boot/constructor"
)

func New() fx.Option {
    return constructor.Load(
        &WelcomeController{},
        &UserController{},
    )
}
```

### Application Bootstrap

The application is bootstrapped using go-fx, which handles dependency injection and component lifecycle:

```go
func App() *fx.App {
    return fx.New(
        config.New(),
        controllers.New(),
        middleware.New(),
        models.New(),
        providers.New(),
        repositories.New(),
        services.New(),

        fx.Invoke(
            app.Boot,
            console.Schedule,
            models.AutoMigrate,
            routes.ApiRoutes,
            routes.WebRoutes,
            serveHTTP,
            startScheduler,
        ),

        fx.WithLogger(useLogger),
    )
}
```

### Component Organization

1. **Controllers** (`app/controllers/`)

    - Handle HTTP requests
    - Implement `Construct() interface{}` for dependency injection
    - Automatically loaded via go-fx

2. **Services** (`app/services/`)

    - Business logic implementation
    - Injected into controllers via go-fx

3. **Repositories** (`app/repositories/`)

    - Data access layer
    - Automatically loaded and injected

4. **Middleware** (`app/middleware/`)

    - HTTP middleware components
    - Loaded via go-fx

5. **Providers** (`app/providers/`)
    - Core service providers
    - Application bootstrapping
    - Database connections
    - Logging, caching, etc.

### Creating New Components

1. **Create a New Controller**:

    ```go
    type CustomController struct {
        *app.App
    }

    func (*CustomController) Construct() interface{} {
        return func(app *app.App) (*CustomController, error) {
            return &CustomController{App: app}, nil
        }
    }

    func (c *CustomController) Custom(ctx *fiber.Ctx) error {
       return c.Render(ctx, "pages/custom", nil)
    }
    ```

2. **Register the Controller**:

    ```go
    // app/controllers/controllers.go
    func New() fx.Option {
        return constructor.Load(
            &CustomController{},
            // ... other controllers
        )
    }
    ```

3. **Add Routes**:
    ```go
    // routes/web.go
    func WebRoutes(app *app.App, controller *controllers.CustomController) {
        app.Get("/custom", controller.Custom)
    }
    ```

## Development Workflow

-   `go run main.go`: Start the application
-   `pnpm run dev`: Start development server with Vite for frontend assets
-   `pnpm run build`: Build frontend assets for production

## Best Practices

-   Use the `Construct() interface{}` pattern for dependency injection
-   Keep controllers thin, delegating business logic to services
-   Use the provided `app.App` base for common functionality
-   Leverage go-fx for automatic dependency management
-   Follow the established naming conventions:
    -   Controllers: `*Controller`
    -   Services: `*Service`
    -   Repositories: `*Repository`
-   Use middleware for cross-cutting concerns
-   Implement proper error handling using the provided logger

## Documentation

-   API documentation is available at `/docs` after running `pnpm run generate:swagger`
-   Frontend assets are automatically versioned in production
-   Hot Module Replacement (HMR) is enabled in development
-   Full page reload is triggered for template changes

## Contributing

We welcome contributions! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for details on the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
