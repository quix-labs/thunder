<template>
  <div>
    <NuxtRouteAnnouncer/>
    <Navbar/>
    <UContainer class="mt-4">
      <NuxtPage/>
    </UContainer>
    <UNotifications :ui="{  position: 'start-1/2 -translate-x-1/2'}"/>
  </div>
</template>