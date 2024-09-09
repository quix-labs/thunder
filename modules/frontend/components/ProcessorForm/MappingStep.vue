<template>
  <div :class="{'short-mode':state.shortMode}">
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

  nextTick(() => {
    document.getElementById('mapping-' + current.activeMappingPath)?.scrollIntoView({
      block: 'start',
      inline: 'start',
      behavior: 'instant',
    })
  })
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