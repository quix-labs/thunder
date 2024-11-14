<template>
  <fieldset class="grid grid-cols-2 gap-4" :id="id" :disabled="disabled">
    <div v-for="option in options" :key="option.value">
      <input
          :type="multiple?'checkbox':'radio'"
          :name="multiple?`${name}[]`:name"
          v-model="model"
          @change="onUpdate"
          :value="asNumber?parseFloat(option.value):option.value"
          class="sr-only peer"
          :id="`${name}:${option.value}`"
          tabindex="-1"
          :disabled="disabled"
      />
      <label
          tabindex="0"
          @keydown.enter.space.prevent="model=asNumber?parseFloat(option.value):option.value"
          :for="`${name}:${option.value}`"
          class="flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800  [.peer:checked+&]:ring-2 [.peer:checked+&]:ring-[var(--ui-primary)]"
          :class="{
            'cursor-pointer':!disabled,
            '[.peer:disabled+&]:cursor-not-allowed':!disabled,
            'cursor-not-allowed':disabled,
          }"
      >
        <slot name="default" :item="option.item"/>
      </label>
    </div>
    <slot name="create"/>
  </fieldset>
</template>
<script setup lang="ts" generic="T">
interface Props<T> {
  options: Array<{ value: any, item: T }>
  name?: string,
  required?: boolean,
  multiple?: boolean,
  asNumber?: boolean,
}

const props = withDefaults(defineProps<Props<T>>(), {
  required: false,
  multiple: false,
  asNumber: false
})
const model = defineModel()
const slots = defineSlots<{ default: { item: T }, create: {} }>()
const emits = defineEmits<{ change: [payload: Event] }>()

const {emitFormChange, emitFormInput, color, name, id: _id, disabled} = useFormField<Props<T>>(props, {bind: false})

const id = _id.value ?? useId()

function onUpdate(value: any) {
  // @ts-expect-error - 'target' does not exist in type 'EventInit'
  const event = new Event('change', {target: {value}})
  emits('change', event)
  emitFormChange()
  emitFormInput()
}

</script>