<template>
  <UChip v-if="kbds && kbds?.length>0">
    <!-- Display keyboard shortcuts -->
    <template #content>
      <UButtonGroup>
        <UKbd v-for="(kbd, index) in kbds" :key="index" :value="kbd"/>
      </UButtonGroup>
    </template>

    <!-- Main button with inherited attributes -->
    <UButton v-bind="$attrs" @click="handleClick">
      <slot/>
    </UButton>
  </UChip>
  <UButton v-bind="$attrs" @click="handleClick" v-else>
    <slot/>
  </UButton>
</template>

<script setup lang="ts">
defineOptions({
  inheritAttrs: false,
})

interface Props {
  kbds?: string[],
  onClick?: () => void
}

const props = defineProps<Props>()

const handleClick = () => {
  if (props.onClick) {
    props.onClick()
  }
}

const shortcuts = computed(() => {
  if (!props.kbds?.length || !props.onClick) {
    return {}
  }
  return {
    [props.kbds.join('_')]: () => props.onClick?.()
  }
})

defineShortcuts(shortcuts)
</script>
