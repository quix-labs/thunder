<template>
  <NuxtRouteAnnouncer/>
  <UApp>
    <Navbar/>
    <UContainer class="mt-4" v-if="container">
      <NuxtPage/>
    </UContainer>
    <NuxtPage v-else/>
    <DrawerProvider/>
  </UApp>
</template>

<script setup lang="ts">
const container = computed(() => {
  if (useRoute().meta?.container !== undefined) {
    return useRoute().meta?.container
  }
  return true
})
const fullHeight = computed(() => {
  if (useRoute().meta?.fullHeight !== undefined) {
    return useRoute().meta?.fullHeight
  }
  return false
})

useHead({
  htmlAttrs: {class: 'full-height'},
})
injectWebhookListener()
</script>

<style>
@import "tailwindcss";
@import "@nuxt/ui";

@theme {
  --container-8xl: 90rem;
}

:root {
  --ui-container: var(--container-8xl);
}

@layer base {
  h1 {
    @apply font-semibold text-xl text-gray-900 dark:text-white leading-tight;
  }
}
</style>

<style>
.full-height {
  @apply h-full max-h-full;

  & body, & #__nuxt {
    @apply h-full;
  }

  & #__nuxt {
    display: grid;
    grid-template-rows: auto 1fr;
  }
}
</style>