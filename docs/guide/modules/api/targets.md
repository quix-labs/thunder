# Managing Targets via the API

The following API routes are available to manage sources programmatically:

- **List all sources**: `GET /go-api/sources`  
  Retrieves a list of all configured sources.

- **Create a new source**: `POST /go-api/sources`  
  Adds a new source to the system.

  The body of the request should include the source configuration.

- **Update an existing source**: `PUT /go-api/sources/{id}`  
  Updates the configuration of a specific source identified by its `{id}`.

- **Delete a source**: `DELETE /go-api/sources/{id}`  
  Removes a source from the system based on its `{id}`.

- **Get source statistics**: `GET /go-api/sources/{id}/stats`  
  Retrieves metrics or usage statistics for a specific source.

- **Get all sources drivers**: `GET /go-api/sources-drivers`  
  Retrieves a list of all installed sources drivers.

- **Test driver configuration**: `POST /go-api/sources-drivers/test`  
  Tests a specific driver configuration without creating it.
