import type {AsyncDataOptions} from "#app";

export type Exporters = {
    [key: string]: string
}

export default (opts?: AsyncDataOptions<Exporters>) => useGoLazyAsyncData<Exporters>('exporters', '/exporters', undefined, opts)
