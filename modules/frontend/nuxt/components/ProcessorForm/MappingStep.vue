<template>
  <div :class="{'short-mode':shortMode}">
    <label class="flex items-center justify-end mb-2 gap-x-4 cursor-pointer">
      Short Mode (beta):
      <UToggle
          on-icon="i-heroicons-check-20-solid"
          off-icon="i-heroicons-x-mark-20-solid"
          v-model="shortMode"
      />
    </label>

    <ProcessorFormMappingeSelector
        :source="form.source"
        :table="form.table"
        v-model="form.mapping"/>
  </div>

</template>

<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
const shortMode = ref(false)

const state = useProcessFormState();
let actualMappingPath: string | undefined = ""

watch(state.value, (current, old) => {
  if (!import.meta.client) return;
  if (actualMappingPath === current.activeMappingPath) return;
  actualMappingPath = current.activeMappingPath

  if (state.value.preventScroll) return;

  console.log('here')
  setTimeout(() => {
    document.getElementById('mapping-' + current.activeMappingPath)?.scrollIntoView({
      block: 'start',
      inline: 'start',
      behavior: 'smooth'
    })
  }, 100)
})
</script>

<style scoped lang="postcss">
:global(.short-mode  .targettable:not(.targetted):not(:has(.targetted))) {
  @apply hidden;
}

:global(.short-mode .targettable.targetted .targettable) {
  @apply block;
}
</style>