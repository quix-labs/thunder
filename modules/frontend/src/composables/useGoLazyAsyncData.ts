import {type FetchOptions} from "ofetch";
import type {AsyncDataOptions} from "#app";

const useGoLazyAsyncData = <T>(uniqueKey: string, request: string | (() => string), opts?: FetchOptions, asyncDataOpts?: AsyncDataOptions<T>) => {

    return useAsyncData<T>(uniqueKey, () => {
            const resolvedRequest = typeof request === 'function' ? request() : request;
            return $fetch(resolvedRequest, {...opts, baseURL: '/go-api'})
        },
        {
            server: false,
            lazy: true,
            getCachedData: key => useNuxtApp().payload.data[key],
            ...asyncDataOpts
        }
    );
}
export default useGoLazyAsyncData
