# Real-Time Forum

## Description

* Create a new forum with registration, login, post creation, commenting, and private messaging features
* Use SQLite for data storage, Golang for backend and WebSockets, JavaScript for frontend events and client WebSockets, HTML for page organization, and CSS for styling
* Implement a single-page application with JavaScript handling page changes
* Ensure real-time functionality for private messages using WebSockets

## Table of Contents (Optional)

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Licence](#license)
- [Acknowledgments](#acknowledgments)

## Installation

## Usage

## Configuration

- Basic commands or features.
- Examples of how to run the application.
- Screenshots or GIFs to illustrate usage (if applicable).

## API Documentation

- Endpoint URLs.
- Request methods (GET, POST, etc.).
- Request and response formats.
- Example requests and responses.

## License

## Acknowledgments

***

*work in progress*

## Techno

- Database Layer: SQLite
- Business Logic Layer: Golang
- ~~API Layer: RESTful API using Golang and Gorilla WebSocket~~
- Frontend Layer: HTML, CSS, and JavaScript
- WebSocket Layer: Gorilla WebSocket and JavaScript WebSocket API
  (optional)
- ~~Caching Layer: Redis or Memcached~~

## Folder Structure

### BACK

```
/real-time-forum
│
├── /cmd
│   └── main.go              # Entry point of the application
│
├── /config                  # Configuration handling
│   └── config.go            # Load and manage configuration
├── /db
│   ├── db.go                # Code for database connection and operations
│   ├── forum.db             # The actual database file (SQLite)
│   └── schema.sql           # Defines the structure of the database (tables, columns, etc.)
│
├── /models                   # Data models (structs)
│   └── user.go              # Example model for a user
│   └── post.go              # Example model for a post
│
├── /repository               # Database access
│   └── user_repository.go    # User data access methods
│   └── post_repository.go    # Post data access methods
│
├── /service                  # Business logic
│   └── user_service.go       # User-related business logic
│   └── post_service.go       # Post-related business logic
│
├── /handler                  # HTTP handlers (controllers)
│   └── user_handler.go       # User-related HTTP handlers
│   └── post_handler.go       # Post-related HTTP handlers
│
├── /middleware               # Middleware functions
│   └── auth_middleware.go    # Example authentication middleware
│
├── /api                      # API route definitions
│   └── routes.go            # Define API routes
│
├── /utils                    # Utility functions
│   └── helpers.go           # Helper functions
│
├── go.mod                    # Go module file
└── go.sum                    # Go module dependencies
```

### Explanation of Each Folder

- **`/cmd`**: Contains the main application entry point. This is where your application starts.
- **`/config`**: A simple directory for managing configuration settings, such as loading environment variables or configuration files.
- **`/db`**: This directory is related to database management.
- **`/models`**: Define your data structures here. For a forum, you might have models for users, posts, comments, etc.
- **`/repository`**: This is where you handle data access. You can create separate files for different entities (e.g., users, posts) to keep things organized.
- **`/service`**: Contains the business logic for your application. This is where you implement the core functionality related to users and posts.
- **`/handler`**: HTTP handlers that respond to incoming requests. Each handler file can correspond to a specific resource (e.g., user-related handlers, post-related handlers).
- **`/middleware`**: Custom middleware functions for handling requests, such as authentication or logging.
- ~~**`/api`**: This directory contains the API route definitions. You can define your routes in a single file for simplicity.~~ ? utility?
- **`/utils`**: Utility functions that can be reused across the application.

CMD

```bash
mkdir -p ./{cmd,config,models,repository,service,handler,middleware,api,utils} && touch ./{cmd/main.go,config/config.go,models/{user.go,post.go},repository/{user_repository.go,post_repository.go},service/{user_service.go,post_service.go},handler/{user_handler.go,post_handler.go},middleware/auth_middleware.go,api/routes.go,utils/helpers.go}
mkdir -p my-forum-api/db/migrations && touch my-forum-api/db/schema.sql
```

### FRONT

todo

## Lexical

- **Cross-Origin Resource Sharing (CORS)**: is a security feature implemented in web browsers that allows or restricts web applications running at one origin (domain) to make requests to resources on a different origin. It uses HTTP headers to inform the browser whether to allow or deny the request based on the origin of the request. CORS is essential for enabling secure interactions between different web applications while preventing unauthorized access to resources.
- **WebSocket**: A protocol for bidirectional, real-time communication between a client and a server over the web.
- **Concurrent connections**: Multiple clients connected to a server simultaneously.
- **Message persistence**: Storing messages in a data store to ensure they are not lost in case of a connection drop or server restart.
- **Concurrency**: To handle multiple clients and messages concurrently, you'll need to use Go's concurrency features, such as goroutines and channels. You can use channels to communicate between goroutines and handle messages.
- **Message broadcasting**: To broadcast messages to all connected clients, you'll need to maintain a list of active connections and send messages to each client. You can use a map to store active connections and iterate over it to send messages.
- **Message queue**: To handle messages efficiently, you can use a message queue to store incoming messages and process them in a separate goroutine.
