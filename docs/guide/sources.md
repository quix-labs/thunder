# Sources

Sources define your database connections

Go to `/sources` to configure them



## Realtime compatibility Table

|      Target      | Insert | Insert rel_col | Update | Update rel_col | Delete | Delete rel col | Truncate |
|:----------------:|:------:|:--------------:|:------:|:--------------:|:------:|:--------------:|:--------:|
|    Base table    |  ❌ ◀   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
|    one_to_one    |  N/A   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
|     has-many     |  ❌ ⏳   |      ❌ ◀       |   ✅    |      ❌ ◀       |   ✅    |       ❌        |    ✅     |
| has-many (pivot) |  ❌ ◀   |      ❌ ◀       |  ❌ ◀   |      ❌ ◀       |   ❌    |       ❌        |    ❌     |

✅: `Supported`
❌: `Not implemented (but possible)`
◀: `Need to send select request (impact DB)`
