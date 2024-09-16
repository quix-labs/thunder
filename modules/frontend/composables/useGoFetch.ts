import {joinURL} from "ufo";

const useGoFetch: typeof useFetch = (path, options = {}) => {
    const fullPath = joinURL('/go-api', path.toString())
    return useFetch(fullPath, {server:false,...options,lazy:true})
}
export default useGoFetch