# Rema Backend

This project is the backend for a Rema application. It provides a GraphQL API to manage users, tasks, and shop items. The backend is built using Go and Prisma for database management.


## Getting Started

### Prerequisites

- Go 1.21.5 or later
- PostgreSQL or MySQL database

### Setup

1. **Clone the repository:**

    ```sh
    git clone https://github.com/youssefyr/rema-backend.git
    cd rema-backend
    ```

2. **Copy the example environment file and update it with your database credentials:**

    ```sh
    cp .env.example .env
    ```

    Update the `.env` file with your database credentials:

    ```env
    DB_HOST=localhost
    DB_USER=YOUR_DB_USER
    DB_PASSWORD=YOUR_DB_PASSWORD
    DB_NAME=YOUR_DB_NAME
    ```

3. **Install dependencies:**

    ```sh
    go mod tidy
    ```

4. **Generate Prisma client:**

    ```sh
    go run github.com/steebchen/prisma-client-go generate
    ```

### Running the Backend

1. **Start the backend server:**

    ```sh
    go run main.go
    ```

    The server will start at `http://localhost:8686/graphql`.

### GraphQL API

You can access the GraphQL playground at `http://localhost:8686/graphql` to interact with the API.

## License

This project is licensed under the GNU General Public License. See the `LICENSE` file for details.