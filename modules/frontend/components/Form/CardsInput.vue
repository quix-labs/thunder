<template>
  <UFormGroup :label :required :name>
    <div class="grid grid-cols-2 gap-4">
      <div v-for="option in options" :key="option.value">
        <input
            :type="multiple?'checkbox':'radio'"
            :name="multiple?`${name}[]`:name"
            v-model="model"
            :value="asNumber?parseFloat(option.value):option.value"
            class="sr-only peer"
            :id="`${name}-${option.value}`"
            tabindex="-1"
        />
        <label
            tabindex="0"
            @keydown.enter.space.prevent="model=asNumber?parseFloat(option.value):option.value"
            :for="`${name}-${option.value}`"
            class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
        >
          <slot name="default" :item="option.item"/>
        </label>
      </div>
    </div>
  </UFormGroup>
</template>
<script setup lang="ts">
interface Props {
  options: Array<{ value: any, item: any }>
  name: string,
  label: string,
  required?: boolean,
  multiple?: boolean,
  asNumber?: boolean,
}

const {options, name, label, required = false, multiple = false, asNumber = false} = defineProps<Props>()
const model = defineModel()
const slots = defineSlots<{
  default(props: { item: any }): any
}>()
</script>