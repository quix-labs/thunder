<template>
  <UFormGroup
      v-for="field in fields"
      :required="field.required"
      :label="field.label"
      :name="`${rootName}${field.name}`"
      :help="field.help">

    <UInput
        autocomplete="off"
        type="number"
        v-if="['number'].includes(field.type)"
        :min="field.min"
        :max="field.max"
        :value="model[field.name]"
        @input="(e:any) => (model[field.name]=e.target.value===''?undefined:parseFloat(e.target.value))"
    />
    <UInput
        autocomplete="off"
        :type="field.type"
        v-else-if="['url','text','password'].includes(field.type)"
        v-model="model[field.name]"/>
    <span v-else>{{ field.type }} not supported!</span>
  </UFormGroup>
</template>
<script setup lang="ts">
import {z, ZodObject, type ZodRawShape} from "zod";

interface Props {
  fields?: any[]
  rootName?: string,
}

const {fields = [], rootName = "config."} = defineProps<Props>()

const model = defineModel<{ [key: string]: any }>({default: {}})
const schema = defineModel<ZodObject<ZodRawShape> | undefined>('schema', {required: false, default: {}})

function refreshSchema() {
  if (fields.length === 0) {
    schema.value = undefined;
    return;
  }

  const driverFields = fields || [];
  const tmpSchema: ZodRawShape = {};

  for (const field of driverFields) {
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

  schema.value = z.object(tmpSchema);
}


function setDefaults() {
  model.value ||= {}
  fields.forEach(field => {
    if (field.default && !model.value?.[field.name]) {
      model.value[field.name] = field.type !== "number" ? field.default : parseFloat(field.default)
    }
  })
}

onBeforeMount(() => {
  refreshSchema()
  setDefaults()
})
</script>