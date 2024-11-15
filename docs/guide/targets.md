# Targets

In this application, a **target** represents the endpoint of the data pipeline. It serves as the final destination where data is sent after being processed.

## Key Concepts

### 1. Driver

Each target is tightly coupled with a **driver** that determines the mechanism for data transmission.

A driver is the core component responsible for connecting to and sending data to the target.

By default, the application includes the `thunder.elastic` driver.

However, developers have the ability to create custom drivers tailored to specific needs.

For details on how to develop custom drivers, refer to the [Developer Documentation(TODO)](#).

### 2. Configuration

Each driver requires a **custom configuration** to define how it interacts with the target.

This includes details such as connection parameters, authentication credentials, and any additional settings required by the driver.

## Managing Targets via the Web UI

The configuration of targets can be managed through the **web interface** available at  
`http://localhost:3000/targets`.

From this interface, you can:

- **Add** a new target
- **Edit** an existing target
- **Delete** a target
- **Test** the configuration of a target to ensure its validity


## Managing Targets via the API

Refer to the [API documentation for targets](./modules/api/targets.md) for detailed information.

## Next Steps

You can now initiate data transfer from [sources](./sources) to [targets](./targets) by creating a new [processor](./processors).
