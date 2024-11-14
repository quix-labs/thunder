<script setup lang="ts">
import {inject} from 'vue'
import {drawerInjectionKey, useDrawer} from "~/composables/useDrawer";

const drawerState = inject(drawerInjectionKey)

const {isOpen} = useDrawer()
</script>

<template>
  <component :is="drawerState.component" v-bind="drawerState.props" v-model:open="isOpen" v-if="drawerState"/>
</template>
