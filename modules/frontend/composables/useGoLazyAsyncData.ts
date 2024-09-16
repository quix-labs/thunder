import {joinURL} from "ufo";
import {type FetchOptions} from "ofetch";


const useGoLazyAsyncData = <T>(uniqueKey: string, request: string, opts?: FetchOptions) => {
    const fullPath = joinURL('/go-api', request)

    return useAsyncData<T>(uniqueKey, () => $fetch(fullPath, opts), {
        server: false,
        lazy: true,
        getCachedData: key => useNuxtApp().payload.data[key],
    })
}
export default useGoLazyAsyncData
