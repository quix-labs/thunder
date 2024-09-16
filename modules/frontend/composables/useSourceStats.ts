export type Stats = { [key: string]: { columns: string[], primary_keys: string[] } };

export default (sourceId: number) => useGoLazyAsyncData<Stats>(`source-stats-${sourceId}`, `/sources/${sourceId}/stats`)
