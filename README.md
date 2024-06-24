
# Literary Lions Forum

The **Literary Lions Forum** is a web-based platform designed to facilitate lively discussions among book club members. Built with **Golang** and **SQLite**, this forum offers a seamless online experience for book enthusiasts. The project includes a **Dockerfile** for easy containerization, ensuring efficient management and deployment of the application.

## Key Features

- **User Authentication**: Users can register with a unique email address, username, and password. Login sessions are managed using cookies.
- **Post Creation and Commenting**: Registered users can create posts and comments, and associate categories with posts.
- **Like/Dislike Functionality**: Users can like or dislike posts, and like comments. This popular feature enhances user engagement.
- **Post Filtering**: Users can filter posts by category, created posts, and liked posts.
- **Search Feature**: Users can search for posts and users.
- **User Profiles**: Users can edit their profiles, posts, and delete their comments and posts.
- **Admin Management**: An admin user has extended rights to delete posts and comments from other users.

## Technical Details

- **Backend**: Golang
- **Database**: SQLite, using `go-sqlite3` driver
- **Frontend**: HTML, CSS, and minimal JavaScript
- **Password Encryption**: Securely stored passwords using encryption
- **Session Management**: Managed with UUIDs

## Dockerization

The forum application is containerized using Docker, ensuring ease of deployment across different environments. The Dockerfile defines the application's environment and dependencies within a container. Key steps include:

1. **Building the Docker Image**: Encapsulates the application and its required components.
2. **Running the Application**: A container can be spun up from the created image, effectively running the application.
3. **Applying Metadata**: Enhances organization and management by applying metadata to Docker objects such as images and containers.
4. **Optimizing Resource Usage**: Maintains a clean environment by addressing unused objects.

## Requirements

- **Golang** version 1.22
- **go-sqlite3**: Install the required package using:
  ```sh
  go get github.com/mattn/go-sqlite3
  ```

## Usage

To run the forum application, execute the following command:

```sh
go run .
```

The server will start on port `8000`.

## Summary

The **Literary Lions Forum** provides a digital haven for bookworms to engage in discussions, share insights, and preserve valuable literary conversations. With features like user authentication, post creation, commenting, and search functionality, it fosters a thriving online book club environment. The use of Docker ensures efficient deployment and management, making it easy to set up and run the forum on any server.

## Note

Comments can only be liked, not disliked, as this is a more common feature in many forums, promoting positive engagement and feedback.
