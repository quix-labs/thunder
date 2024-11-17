import type {AsyncDataOptions} from "#app";

export type ProcessorsResponse = Array<any>;

export default (asyncDataOpts?: AsyncDataOptions<ProcessorsResponse>) => useGoLazyAsyncData<ProcessorsResponse>('processors', '/processors', undefined, asyncDataOpts)
