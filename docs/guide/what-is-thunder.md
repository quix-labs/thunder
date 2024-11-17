# What is Thunder?

**Thunder** is a lightweight tool that synchronizes structured SQL databases (like Postgres, MySQL, etc.) with search engine indexes (ElasticSearch, OpenSearch, etc.).

Built for real-time updates, Thunder ensures your data is always indexed and searchable, without impacting your database's performance.

<div class="tip custom-block" style="padding-top: 8px">

Want to try it out? Jump straight to the [Quickstart](./installation).

</div>

## Use Cases

Thunder is ideal for scenarios where you need to sync structured data from your database to a search engine:

- **E-commerce**: Make product data instantly searchable (e.g., name, price, stock status).
- **CMS**: Keep articles, blogs, or content dynamically indexed.
- **Analytics & BI**: Sync large data sets for real-time reporting and analysis.
- **CRM**: Index customer data for improved search performance.


## Key Features

- **Real-Time Sync**

  Changes to your database (inserts, updates, deletes) trigger immediate updates to your search engine index, ensuring it stays up-to-date with minimal delay—no need for batch processing or manual reindexing.

- **Non-Intrusive Synchronization**

  Thunder minimizes performance impact by using **Write-Ahead Logs (WAL)** and **replication**.

  It listens for changes recorded in the WAL, avoiding additional queries to the database.

  By streaming changes in near real-time, Thunder ensures minimal resource usage and zero impact on your production database.

- **Cross-Platform Support**

  Thunder supports multiple installation options, including Docker, building from source, or using prebuilt assets, making it flexible for different environments.

- **Data Export**

  Thunder also allows you to export your synchronized data in bulk to formats such as CSV, JSON, and others.
  
  This can be done with a single command, making it easy to back up or transfer large datasets from your database to different systems.

- **Web Interface for Configuration**
 
  Thunder provides a user-friendly web interface to help you configure the synchronization process.
  
  The interface guides you through the setup, automatically suggesting the right columns and fields to sync based on your database schema, making the process smooth and error-free.

## Supported Databases and Search Engines

:::warning Important Notes

Real-time listening is currently only planned for the [PostgreSQL Driver](./drivers/postgresql_flash).  
For details on why it is not implemented for MySQL, refer to the [MySQL Driver](./drivers/mysql) page.

:::


Thunder currently supports:

- **Databases**:
  - PostgreSQL
  - MySQL/MariaDB (Partial, no real-time capabilities)

- **Search Engines**:
  - ElasticSearch
  - ~~OpenSearch~~ (Coming soon)


## Interested?

To get started with Thunder, follow the [setup guide](./installation).
