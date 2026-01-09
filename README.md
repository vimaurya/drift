## Drift

Drift is a lightweight, robust database migration CLI tool written in Go. It is designed to ensure "Strong Consistency" between your local migration files and your remote database schema by using SHA256 integrity checks to prevent Migration Drift.

## Features
**Drift Protection:** Detects if a migration file was edited after being applied using SHA256 checksum validation.

Multi-Driver Support: Native support for PostgreSQL and MySQL.

**LIFO Rollbacks:** Intelligent down command that reverts migrations in the correct reverse-chronological order.

**Atomic Migrations:** Supports full transactional migrations for PostgreSQL.

**Configuration Persistence:** Saves project settings in a local configuration file for a seamless developer experience.

## Connection String Patterns
When running drift init, you must provide a database URL. Use the following patterns based on your database type:

1. PostgreSQL
Postgres uses a standard URI format.

```
postgres://<user>:<password>@<host>:<port>/<database_name>?sslmode=disable
```
Example: `postgres://admin:secret@localhost:5432/my_db?sslmode=disable`

2. MySQL
The MySQL driver requires a specific Data Source Name (DSN) format. If your migrations contain multiple SQL statements, you must append multiStatements=true.

```
mysql://<user>:<password>@tcp(<host>:<port>)/<database_name>?multiStatements=true
```
Example: `mysql://root:password@tcp(127.0.0.1:3306)/my_db?multiStatements=true`

## Installation
Download the binary for your operating system from the latest releases.

Move the binary to your project folder (or add it to your PATH).

Rename it to drift (or drift.exe on Windows).

## Usage
1. Initialize
Sets up the local configuration and creates the schema_migrations tracking table in your database.

```Bash
drift init -url="your_connection_string"
```
2. Create Migration
Generates a timestamped pair of .up.sql and .down.sql files.

```Bash
drift create -name=add_users_table
```
3. Run Migrations
Applies all pending migrations.

```Bash
drift up
```
4. Rollback
Reverts the single most recent migration applied to the database.

```Bash
drift down
```
## Project Structure
cmd/: Entry point for the CLI application.

internal/driver/: Database factory and driver implementations (Postgres/MySQL).

internal/core/: The migration engine logic and checksum validation.

internal/config/: Configuration file management.
