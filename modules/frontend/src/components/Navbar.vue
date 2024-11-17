<template>
  <UNavigationMenu color="primary" variant="pill" :items="items">
    <template #item-trailing="{item}">
      <UButtonGroup>
        <UKbd v-for="key in item.shortcuts||[]" :value="key"/>
      </UButtonGroup>
    </template>
  </UNavigationMenu>

</template>
<script setup lang="ts">
import type {ShortcutsConfig} from "#ui/composables/defineShortcuts";

const items = computed(() => [
  [
    {
      label: 'Dashboard',
      icon: 'i-heroicons-home',
      to: '/',
      shortcuts: ['alt', 'd']
    }
  ],
  [
    {
      label: 'Sources',
      icon: 'i-heroicons-circle-stack',
      to: '/sources',
      shortcuts: ['alt', 's']
    },
    {
      label: 'Processors',
      icon: 'i-heroicons-briefcase',
      to: '/processors',
      shortcuts: ['alt', 'p'],
      active: useRoute().path.startsWith('/processors'),
    },
    {
      label: 'Targets',
      icon: 'i-heroicons-arrow-up-tray',
      to: '/targets',
      shortcuts: ['alt', 't']
    }
  ],
  [
    {
      slot: 'colorMode'
    }
  ],
])
const colorModeItems = [
  {
    label: 'System',
    value: 'system',
    icon: 'i-heroicons-computer-desktop'
  },
  {
    label: 'Light',
    value: 'light',
    icon: 'i-heroicons-sun',
  },
  {
    label: 'Dark',
    value: 'dark',
    icon: 'i-heroicons-moon',

  }
]
const shortcuts = computed<ShortcutsConfig>(() => Object.fromEntries(
    items.value.flat(1)
        .filter(item => 'shortcuts' in item)
        .map(item => {
          const key = item.shortcuts.join('_');
          console.log(key)
          return [key, () => navigateTo(item.to)];
        })
));
defineShortcuts(shortcuts)
</script>