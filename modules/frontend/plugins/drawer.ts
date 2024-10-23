import { shallowRef } from "vue";
import { defineNuxtPlugin } from "#imports";
import {drawerInjectionKey} from "~/composables/useDrawer";

export default defineNuxtPlugin((nuxtApp) => {
    const drawerState = shallowRef({
        component: "div",
        props: {}
    });
    nuxtApp.vueApp.provide(drawerInjectionKey, drawerState);
});
