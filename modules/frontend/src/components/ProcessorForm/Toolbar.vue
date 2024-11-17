<template>
  <div class="flex justify-between">
    <h2 class="col-span-full  text-center text-2xl font-semibold">
      {{ title }}
    </h2>

    <div class="flex gap-x-4">
      <label class="flex items-center gap-x-4 cursor-pointer">
        Short Mode (beta):
        <USwitch
            checked-icon="i-heroicons-check-20-solid"
            unchecked-icon="i-heroicons-x-mark-20-solid"
            v-model="state.shortMode"
        />
      </label>
      <label class="flex items-center gap-x-4 cursor-pointer ">
        Sample Tab:
        <USwitch
            checked-icon="i-heroicons-check-20-solid"
            unchecked-icon="i-heroicons-x-mark-20-solid"
            v-model="state.showSamplePanel"
        />
      </label>

      <div v-if="$slots.default" class="flex gap-4">
        <slot/>
      </div>

    </div>
  </div>
</template>
<script setup lang="ts">
interface Props {
  title: string
}
defineProps<Props>()
const state = useProcessFormState()

</script>