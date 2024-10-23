export type Stats = { [key: string]: { columns: string[], primary_keys: string[] } };

export default (sourceId: MaybeRef<number>) => {
    return useGoLazyAsyncData<Stats>(
        `source-stats-${unref(sourceId)}`,
        () => `/sources/${unref(sourceId)}/stats`,
        undefined,
        isRef(sourceId) ? {watch: [sourceId]} : undefined,
    )
}
