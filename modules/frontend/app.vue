<template>
  <NuxtRouteAnnouncer/>

  <UApp >
    <Navbar/>
    <UContainer class="mt-4">
      <NuxtPage/>
    </UContainer>
    <DrawerProvider/>
  </UApp>
</template>

<script setup lang="ts">

injectWebhookListener()
</script>

<style>
@import "tailwindcss";
@import "@nuxt/ui";

@layer base {
  h1 {
    @apply font-semibold text-xl text-gray-900 dark:text-white leading-tight;
  }
}
</style>