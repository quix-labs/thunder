import type {Component, InjectionKey, ShallowRef} from "vue";
import {createSharedComposable} from "@vueuse/shared";
import type {ComponentProps} from "#ui/types/component";
import type {DrawerRootProps} from "vaul-vue";

export type DrawerProps = DrawerRootProps

export interface DrawerState {
    component: Component | string;
    props: DrawerProps;
}

export const drawerInjectionKey: InjectionKey<ShallowRef<DrawerState>> = Symbol("custom.drawer");

function _useDrawer() {
    const drawerState = inject(drawerInjectionKey)

    const isOpen = shallowRef(false);

    function open<T extends Component>(component: T, props?: DrawerProps & ComponentProps<T>) {
        if (!drawerState) {
            throw new Error("useDrawer() is called without provider");
        }
        drawerState.value = {
            component,
            props: props ?? {}
        };
        isOpen.value = true;
    }

    async function close() {
        if (!drawerState) return;
        isOpen.value = false;
    }

    function reset() {
        if (!drawerState) return;
        drawerState.value = {
            component: "div",
            props: {}
        };
    }

    function patch<T extends Component = Record<string, never>>(props: Partial<DrawerProps & ComponentProps<T>>) {
        if (!drawerState) return;
        drawerState.value = {
            ...drawerState.value,
            props: {
                ...drawerState.value.props,
                ...props
            }
        };
    }

    return {open, close, reset, patch, isOpen}
}

export const useDrawer = createSharedComposable(_useDrawer)