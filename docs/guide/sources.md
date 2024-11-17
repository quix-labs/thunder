# Sources

In this application, a **source** represents the starting point of the data pipeline. It serves as the entry point where
data is extracted from a specific system or database before being processed further.

## Key Concepts

### 1. Driver

Each source is tightly coupled with a **driver** that determines the mechanism for data extraction.

A driver is the core component responsible for connecting to and extracting data from the source.

By default, the application includes the [PostgreSQL Driver](./drivers/postgresql_flash) and [MySQL/MariaDB](./drivers/mysql) driver.

However, developers have the ability to create custom drivers tailored to specific needs.

For details on how to develop custom drivers, refer to the [Developer Documentation(TODO)](#).

### 2. Configuration

Each driver requires a **custom configuration** to define how it interacts with the source.

This includes details such as connection parameters, authentication credentials, and any additional settings required by
the driver.

## Managing Sources via the Web UI

The configuration of sources can be managed through the **web interface** available at
`http://localhost:3000/sources`.

From this interface, you can:

- **Add** a new source
- **Edit** an existing source
- **Delete** a source
- **Test** the configuration of a source to ensure its validity
## Managing Sources via the API

Refer to the [API documentation for sources](./modules/api/sources.md) for detailed information.

## Next Steps

Once you have configured your source(s), the next step is to set up [your target](./targets)
