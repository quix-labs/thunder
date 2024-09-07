import mitt from 'mitt'

export type ApplicationEvents = {};

const emitter = mitt<ApplicationEvents>()

export default function useEventBus() {
    return emitter
}
