import {joinURL} from "ufo";

export default function useEvents() {
    const eventSource = import.meta.client ? useState('events-source', () => new EventSource('/go-api/events')) : {value: undefined};
    function on(event: string, closure: (data: any) => void) {
        eventSource.value?.addEventListener(event, (e) => {
            const data = JSON.parse(e.data)
            closure(data)
        })
    }

    return {on}
}