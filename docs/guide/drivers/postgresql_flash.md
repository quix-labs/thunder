# PostgreSQL Flash Driver

:::info
Real-time capabilities are currently under development.  
You can still use the driver for bulk indexing or to generate structured documents via the [Exporters](../exporters).
:::


## Real-time Compatibility Table (Work in Progress, Currently Unavailable)

|      Target      | Insert | Insert rel_col | Update | Update rel_col | Delete | Delete rel col | Truncate |
|:----------------:|:------:|:--------------:|:------:|:--------------:|:------:|:--------------:|:--------:|
|    Base table    |  ❌ ◀   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
|    one_to_one    |  N/A   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
|     has-many     |  ❌ ⏳   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
| has-many (pivot) |  ❌ ◀   |      ❌ ◀       |  ❌ ◀   |      ❌ ◀       |   ❌    |       ❌        |    ❌     |

✅: `Supported`
❌: `Not implemented (but possible)`
◀: `Need to send select request (impact DB)`