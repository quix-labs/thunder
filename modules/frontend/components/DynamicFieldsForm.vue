<template>
  <UFormGroup
      v-for="field in fields"
      :required="field.required"
      :label="field.label"
      :name="`config.${field.name}`"
      :help="field.help">

    <UInput
        autocomplete="off"
        type="number"
        :placeholder="field.default"
        v-if="['number'].includes(field.type)"
        :min="field.min"
        :max="field.max"
        :value="model[field.name]"
        @input="(e:any) => (model[field.name]=e.target.value===''?undefined:parseFloat(e.target.value))"
    />
    <UInput
        autocomplete="off"
        :type="field.type"
        :placeholder="field.default"
        v-else-if="['text','password'].includes(field.type)"
        v-model="model[field.name]"/>
    <span v-else>{{ field.type }} not supported!</span>
  </UFormGroup>
</template>
<script setup lang="ts">
interface Props {
  fields?: any[]
}

const {fields = []} = defineProps<Props>()
const model = defineModel<{ [key: string]: any }>({default: {}})
</script>