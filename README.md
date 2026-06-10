## How to Run

Requires Docker and Docker Compose installed.

```bash
# 1. Copy example environment file
cp .env.example .env

# 2. Configure variables in .env

# 3. Start containers
docker compose up -d --build
```

### Local PostgreSQL Setup

If you prefer to use a local PostgreSQL instance instead of Docker, you can set up the database using the provided `init.sql` script:

1. Connect to your PostgreSQL instance (e.g., using `psql`):
   ```bash
   psql -h localhost -U <your_username> -d <your_database_name>
   ```
2. Run the initialization script:
   ```sql
   \i init.sql
   ```

## Environment Variables
...

Configure these in your `.env` file:

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT signing

## API Usage

- **Documentation:** `http://localhost:8080/swagger/index.html`
- **Auth:**
  - `POST /auth/register`: Create user
  - `POST /auth/login`: Authenticate and get JWT
- **Tasks (requires Authorization header: `Bearer <token>`):**
  - `POST /tasks`: Create task
  - `GET /tasks`: List user tasks
  - `GET /tasks/{id}`: Get task
  - `PUT /tasks/{id}`: Update task
  - `DELETE /tasks/{id}`: Delete task
