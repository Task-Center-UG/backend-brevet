# Backend Brevet

Backend Brevet is a REST API for managing Brevet training programs at a Tax Center. It covers student registration, email verification, course and batch management, purchases and payment confirmation, meetings, attendance, assignments, quizzes, score tracking, testimonials, blogs, dashboards, and certificate generation.

The application is built with Go, Fiber, GORM, PostgreSQL, Redis, Docker, and Docker Compose.

## Tech Stack

- Go 1.23+
- Fiber v2
- GORM
- PostgreSQL 16
- Redis 7
- JWT authentication
- Docker and Docker Compose
- Excel import/export with `excelize`
- PDF and document generation for certificates and receipts

## Main Features

- Authentication
  - Register, login, logout
  - Email verification code
  - Access token and refresh token flow
  - Redis-based access token blacklist
  - User session tracking

- User and profile management
  - Student, teacher, and admin roles
  - User profile data
  - Group type support: Gunadarma student, non-Gunadarma student, general public

- Course and batch management
  - Course CRUD
  - Batch CRUD
  - Batch schedule, registration period, quota, room, course type, allowed group types, and active days
  - Batch quota checking

- Purchase and payment
  - Student batch purchase
  - Group-based pricing
  - Unique transfer amount generation
  - Payment proof submission
  - Admin payment status confirmation
  - Receipt generation and email attachment
  - Expired payment cleanup scheduler

- Meeting management
  - Batch meetings
  - Assign teachers to meetings
  - Online/offline and exam/basic meeting types
  - Meeting content access control

- Materials
  - Teacher/admin material upload and management
  - Student access after successful purchase

- Assignments and submissions
  - Teacher/admin assignment creation
  - Student submission with notes, essay text, and files
  - Grade and feedback
  - Excel export/import for grading
  - Sequential meeting rules before submitting next content

- Quizzes
  - Quiz metadata creation
  - Excel question import
  - Multiple choice and true/false support
  - Attempt tracking
  - Temporary answer saving
  - Manual and automatic submission
  - Result calculation

- Attendance
  - Bulk attendance update
  - Attendance used for progress and certificate eligibility

- Certificates
  - Certificate generation after course completion
  - PDF certificate generation
  - QR code verification
  - Public certificate verification endpoint

- Dashboard
  - Admin dashboard
  - Teacher dashboard
  - Student dashboard
  - Revenue chart
  - Pending payments
  - Batch progress
  - Teacher workload
  - Certificate statistics
  - Recent activities

- Blog and testimonials
  - Public blog and testimonial reads
  - Admin blog management
  - Student testimonial management

## Project Structure

```text
.
|-- cmd/                 # CLI commands, including database migration
|-- config/              # Database, Redis, and environment configuration
|-- controllers/         # HTTP request handlers
|-- docker/              # Docker init scripts
|-- dto/                 # Request and response structs
|-- helpers/             # Formatting, logging, and document helpers
|-- middlewares/         # Auth, role guard, logging, and request validation
|-- models/              # GORM database models and enum types
|-- repository/          # Database access layer
|-- routes/              # API route registration
|-- scheduler/           # Background jobs
|-- seed/                # Default seed data
|-- services/            # Business logic
|-- tests/               # Unit tests
|-- utils/               # Shared utility functions
|-- validators/          # Custom request validators
|-- DockerFile
|-- docker-compose.yml
|-- go.mod
`-- main.go
```

## Requirements

For Docker-based development:

- Docker
- Docker Compose

For local development without Docker:

- Go 1.23+
- PostgreSQL 16+
- Redis 7+
- LibreOffice, required for receipt document conversion

## Environment Variables

Create a `.env` file in the project root. Docker Compose reads this file automatically.

Example for Docker development:

```env
APP_ENV=development
APP_PORT=8083
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:8083

DB_HOST=db
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=brevetdb
DB_PORT=5432
DB_SSLMODE=disable

REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

ACCESS_TOKEN_SECRET=change-this-access-secret
REFRESH_TOKEN_SECRET=change-this-refresh-secret
VERIFICATION_TOKEN_SECRET=change-this-verification-secret
ACCESS_TOKEN_EXPIRY_HOURS=24
REFRESH_TOKEN_EXPIRY_HOURS=24
VERIFICATION_TOKEN_EXPIRY_MINUTES=15
TOKEN_BLACKLIST_TTL=86400
CLEANUP_INTERVAL_HOURS=1

UPLOAD_DIR=/root/public/uploads

SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=dev@example.test
SMTP_PASS=dev-password
```

Notes:

- When running inside Docker, use `DB_HOST=db` and `REDIS_ADDR=redis:6379`.
- When running directly on your machine, use `DB_HOST=localhost` and `REDIS_ADDR=localhost:6379`.
- Replace all JWT and SMTP values for staging or production.
- `.env` should not be committed if it contains real secrets.

## Run with Docker Compose

Build and start all services:

```bash
docker compose up -d --build
```

This starts:

- PostgreSQL
- Redis
- Migration and seeding job
- Backend API

Check running containers:

```bash
docker compose ps
```

View API logs:

```bash
docker compose logs -f api
```

Stop services:

```bash
docker compose down
```

Stop services and remove volumes:

```bash
docker compose down -v
```

Use `down -v` only when you intentionally want to delete local database and Redis data.

## Health Check

After the API starts, open:

```text
http://localhost:8083/hello
```

Expected response:

```json
{
  "data": null,
  "message": "Backend Brevet API is running",
  "success": true
}
```

## Database Migration and Seed

Migration and seeding are handled by the `migrate` service in Docker Compose:

```bash
docker compose up --build --force-recreate migrate
```

The migration command:

```bash
go run ./cmd/migrate.go
```

The seed command:

```bash
go run ./seed/main.go
```

The seed process creates default users and default prices.

## Default Seeded Users

After running the seeder, these accounts are available:

```text
Admin
Email:    admin@brevet.local
Password: Admin123!

Teacher
Email:    guru@brevet.local
Password: Guru123!

Student
Email:    siswa@brevet.local
Password: Siswa123!
```

## Run Locally Without Docker

Install dependencies:

```bash
go mod download
```

Make sure PostgreSQL and Redis are running, then update `.env`:

```env
DB_HOST=localhost
REDIS_ADDR=localhost:6379
UPLOAD_DIR=./public/uploads
```

Run migration:

```bash
go run ./cmd/migrate.go
```

Run seed:

```bash
go run ./seed/main.go
```

Start the API:

```bash
go run main.go
```

Or use the Makefile:

```bash
make build
make run
```

For hot reload, install `air` first, then run:

```bash
make dev
```

## Tests

Run all tests:

```bash
go test ./...
```

The current test suite includes service and repository unit tests.

## API Overview

Base path:

```text
/api/v1
```

Main route groups:

```text
POST   /auth/register
POST   /auth/login
POST   /auth/verify
POST   /auth/resend-verification
POST   /auth/refresh-token
DELETE /auth/logout

GET    /courses
GET    /courses/:slug
POST   /courses
PUT    /courses/:id
DELETE /courses/:id
GET    /courses/:courseSlug/batches
POST   /courses/:courseId/batches

GET    /batches
GET    /batches/:slug
PUT    /batches/:id
DELETE /batches/:id
GET    /batches/:batchSlug/meetings
POST   /batches/:batchID/meetings
GET    /batches/:batchSlug/students
GET    /batches/:batchSlug/quota

GET    /me
PATCH  /me
GET    /me/purchases
POST   /me/purchases
PATCH  /me/purchases/:id/pay
PATCH  /me/purchases/:id/cancel
GET    /me/batches
GET    /me/batches/:batchID/progress
POST   /me/batches/:batchID/certificate
GET    /me/batches/:batchID/certificate
GET    /me/batches/:batchID/scores
GET    /me/assignments/upcoming
GET    /me/quizzes/upcoming

GET    /purchases
GET    /purchases/:id
PATCH  /purchases/:id/status

GET    /meetings
GET    /meetings/:id
PATCH  /meetings/:id
DELETE /meetings/:id
GET    /meetings/:meetingID/teachers
POST   /meetings/:meetingID/teachers
PUT    /meetings/:meetingID/teachers
DELETE /meetings/:meetingID/teachers/:teacherID
GET    /meetings/:meetingID/assignments
POST   /meetings/:meetingID/assignments
GET    /meetings/:meetingID/materials
POST   /meetings/:meetingID/materials
GET    /meetings/:meetingID/quizzes
POST   /meetings/:meetingID/quizzes

GET    /assignments
GET    /assignments/:assignmentID
PATCH  /assignments/:assignmentID
DELETE /assignments/:assignmentID
GET    /assignments/:assignmentID/submissions
POST   /assignments/:assignmentID/submissions
GET    /assignments/:assignmentID/grades/excel
PUT    /assignments/:assignmentID/grades/import

GET    /submissions/:submissionID
PATCH  /submissions/:submissionID
DELETE /submissions/:submissionID
GET    /submissions/:submissionID/grade
PUT    /submissions/:submissionID/grade

GET    /quizzes/:quizID
GET    /quizzes/:quizID/questions
POST   /quizzes/:quizID/import-questions
POST   /quizzes/:quizID/start
GET    /quizzes/:quizID/attempts
GET    /quizzes/:quizID/attempts/active
POST   /quizzes/attempts/:attemptID/temp-submissions
POST   /quizzes/attempts/:attemptID/submissions
GET    /quizzes/attempts/:attemptID
GET    /quizzes/attempts/:attemptID/result

GET    /certificates/number/:number
GET    /certificates/:certificateID
GET    /certificates/:certificateID/verify

GET    /blogs
GET    /blogs/:slug
POST   /blogs
PUT    /blogs/:id
DELETE /blogs/:id

GET    /testimonials
GET    /testimonials/:testimonialID
PATCH  /testimonials/:testimonialID
DELETE /testimonials/:testimonialID

GET    /dashboard/admin
GET    /dashboard/teacher
GET    /dashboard/student
```

Some routes are public, while most write operations require authentication and specific roles.

## Roles and Access Control

Roles:

- `admin`
- `guru`
- `siswa`

General access rules:

- Admin can manage master data, users, payments, courses, batches, dashboards, and certificates.
- Teacher can access meetings assigned to them, manage materials, assignments, quizzes, and grades.
- Student can purchase batches, access paid batch content, submit assignments, take quizzes, view scores, and request certificates.

## Query Parameters

Most list endpoints support:

```text
q       Search keyword
sort    Sort field
order   asc or desc
select  Comma-separated selected fields
limit   Page size
page    Page number
```

Additional query parameters are treated as filters when supported by the repository.

Example:

```text
GET /api/v1/courses?q=tax&sort=created_at&order=desc&page=1&limit=10
```

## Background Jobs

Schedulers start automatically with the API:

- Expired session cleanup
- Expired purchase marking
- Quiz auto-submit every minute

Cleanup interval is controlled by:

```env
CLEANUP_INTERVAL_HOURS=1
```

## File Uploads and Static Files

Uploaded files are served from:

```text
/uploads
```

The upload directory is configured by:

```env
UPLOAD_DIR=/root/public/uploads
```

In Docker Compose, uploads are stored in the `uploads_data` volume.

## Docker Image and Deployment

The GitHub Actions workflow builds and pushes:

```text
zidanindratama/backend-brevet:latest
```

Deployment workflow:

1. Checkout source.
2. Build and push Docker image.
3. SSH into VPS.
4. Pull latest repository code.
5. Pull latest API image.
6. Start database and Redis.
7. Run migration and seeding.
8. Start the API.
9. Prune unused Docker images.

Expected VPS project path:

```text
~/backend-brevet
```

Make sure the production `.env` file exists on the VPS before deployment.

## Common Commands

```bash
# Start all services
docker compose up -d --build

# Run migration and seed only
docker compose up --build --force-recreate migrate

# Show logs
docker compose logs -f api

# Run tests
go test ./...

# Stop services
docker compose down
```

## Notes for Contributors

- Keep business logic in `services/`.
- Keep database queries in `repository/`.
- Keep request validation in `dto/`, `validators/`, and route middleware.
- Run `gofmt` before committing Go files.
- Run `go test ./...` before opening a pull request.
- Avoid committing real secrets in `.env`.

