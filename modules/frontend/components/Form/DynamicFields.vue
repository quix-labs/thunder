<template>
  <UForm :state :schema ref="formEl" class="space-y-4">
    <UFormField
        v-for="field in fields"
        :required="field.required"
        :label="field.label"
        :name="`${field.name}`"
        :help="field.help">

      <UInput
          autocomplete="off"
          type="number"
          v-if="['number'].includes(field.type)"
          :min="field.min"
          :max="field.max"
          :value="state[field.name]"
          @input="(e:any) => (state[field.name]=e.target.value===''?undefined:parseFloat(e.target.value))"
      />
      <UInput
          autocomplete="off"
          :type="field.type"
          v-else-if="['url','text','password'].includes(field.type)"
          v-model="state[field.name]"/>
      <span v-else>{{ field.type }} not supported!</span>
    </UFormField>
  </UForm>
</template>
<script setup lang="ts">
import {z, type ZodRawShape} from "zod";

const formEl = useTemplateRef("formEl")

interface Props {
  fields?: any[]
}

const {fields = []} = defineProps<Props>()

const state = defineModel<{ [key: string]: any }>('state', {required: true})

const schema = computed(() => {
  if (fields.length === 0) {
    return undefined;
  }
  const tmpSchema: ZodRawShape = {};
  for (const field of fields) {
    let zodType = {
      email: z.string({required_error: `${field.label} is required`}).email("Invalid email"),
      number: z.number({required_error: `${field.label} is required`}),
      url: z.string({required_error: `${field.label} is required`}).url({message: `${field.label} must be a valid URL`}),
    }[field.type as string] || z.string({required_error: `${field.label} is required`});

    if (field.min) {
      zodType = zodType.min(parseFloat(field.min), {message: `${field.label} must be at least ${field.min}`})
    }
    if (field.max) {
      zodType = zodType.max(parseFloat(field.max), {message: `${field.label} must not exceed ${field.max}`})
    }
    tmpSchema[field.name] = field.required ? zodType.min(1, {message: `${field.label} is required`}) : z.optional(zodType).nullable();
  }
  return z.object(tmpSchema)
})

function setDefaults() {
  state.value ||= {}
  fields.forEach(field => {
    if (field.default && !state.value?.[field.name]) {
      state.value[field.name] = field.type !== "number" ? field.default : parseFloat(field.default)
    }
  })
}

onBeforeMount(() => {
  setDefaults()
})
</script>