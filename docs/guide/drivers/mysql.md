# MySQL Driver

:::danger Important Notes

This driver is incompatible with real-time synchronization.  
However, you can use the driver for exporting data and batch filling your targeted index.  
Support for the driver is not yet planned. I need to explore alternatives to `pg_notify` or similar processes like `WAL replication`.  
Under the hood, [Flash](https://flash.quix-labs.com) is used to track events from the database in real-time, and any improvements for MySQL must be integrated into Flash.

:::

## Realtime compatibility Table (currently not planned)

|      Target      | Insert | Insert rel_col | Update | Update rel_col | Delete | Delete rel col | Truncate |
|:----------------:|:------:|:--------------:|:------:|:--------------:|:------:|:--------------:|:--------:|
|    Base table    |   ❌    |       ❌        |   ❌    |       ❌        |   ❌    |       ❌        |    ❌     |
|    one_to_one    |   ❌    |       ❌        |   ❌    |       ❌        |   ❌    |       ❌        |    ❌     |
|     has-many     |   ❌    |       ❌        |   ❌    |       ❌        |   ❌    |       ❌        |    ❌     |
| has-many (pivot) |   ❌    |       ❌        |   ❌    |       ❌        |   ❌    |       ❌        |    ❌     |

❌: `Not Planned`
