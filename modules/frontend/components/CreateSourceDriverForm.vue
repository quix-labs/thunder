<template>
  <USlideover v-model="opened" class="overflow-y-auto">
    <UForm :state="form" class="flex flex-col flex-1" :schema="schema" @submit="onSubmit" ref="formEl"
           v-if="status==='success' && inDrivers">
      <UCard class="flex flex-col flex-1" :ui="{ rounded:'', body: { base: 'flex-1 space-y-2' }}" rounded>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">Create Data Source</h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" class="-my-1"
                     @click="opened = false"/>
          </div>
        </template>

        <template #default>
          <!-- Driver selection -->
          <UFormGroup label="Driver" required name="driver">
            <div class="grid grid-cols-2 gap-4">
              <div v-for="([key,driver]) in Object.entries(inDrivers)">
                <input type="radio" name="driver" v-model="form.driver" :value="key" class="sr-only peer"
                       :id="`driver-${key}`"
                       @change="formEl.clear()"
                       tabindex="-1"
                >
                <label :for="`driver-${key}`"
                       tabindex="0"
                       @keydown.enter.space.prevent="form.driver=key"
                       class="cursor-pointer flex flex-col gap-y-1 text-center items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500">
                  <span v-html="driver.config.image" v-if="driver.config.image" class="nested-svg"/>
                  <span v-text="driver.config.name"/>
                </label>
              </div>
            </div>
          </UFormGroup>
          <!-- Driver configuration -->
          <template v-if="form.driver">

            <!-- Driver notes-->
            <div v-if="inDrivers[form.driver]?.config?.notes?.length>0">
              <UDivider label="Driver Notes"/>
              <ul class="list-disc list-outside ml-4">
                <li v-for="note in inDrivers[form.driver]?.config?.notes" v-text="note"/>
              </ul>
            </div>

            <!-- Driver Configuration field-->
            <UDivider label="Driver Configuration"/>
            <UFormGroup
                v-for="field in inDrivers[form.driver].fields"
                :required="field.required"
                :label="field.label"
                :name="`config.${field.name}`"
                :help="field.help">

              <UInput
                  type="number" :placeholder="field.default" v-if="['number'].includes(field.type)"
                  :value="form.config[field.name]"
                  @input="(e:any) => (form.config[field.name]=e.target.value===''?undefined:parseFloat(e.target.value))"
              />
              <UInput :type="field.type" :placeholder="field.default"
                      v-else-if="['text','password'].includes(field.type)" v-model="form.config[field.name]"/>
              <span v-else>{{ field.type }} not supported!</span>
            </UFormGroup>
          </template>
        </template>

        <template #footer>
          <div class="flex justify-between">
            <UButton type="button" @click.prevent="sendTestRequest" variant="soft" color="sky">
              Test configuration
            </UButton>
            <UButton type="submit">Create</UButton>
          </div>
        </template>
      </UCard>
    </UForm>


  </USlideover>
</template>
<script setup lang="ts">
const opened = defineModel('opened')
const $emit = defineEmits(['created'])

import {z, type ZodRawShape} from 'zod'
import type {FormSubmitEvent} from '#ui/types'

const {data: inDrivers, status} = useSourceDrivers()
const formEl = ref()
const form = reactive({
  driver: null as (keyof SourceDrivers) | null,
  config: {} as { [key: string]: any }
})

const schema = computed(() => {
  let driverConfig: any = undefined;

  if (!inDrivers.value) return z.object({
    driver: z.enum([]),
  })
  if (!form.driver) return z.object({
    driver: z.enum(Object.keys(inDrivers.value)),
  })


  const driverFields = inDrivers.value[form.driver].fields || []
  const schema: ZodRawShape = {}
  for (const field of driverFields) {
    let zodType = {
      text: z.string({required_error: `Field ${field.name} is required`}),
      email: z.string({required_error: `Field ${field.name} is required`}).email("Invalid email"),
      password: z.string({required_error: `Field ${field.name} is required`}),
      number: z.number({required_error: `Field ${field.name} is required`}),
    }[field.type as string] || z.string()

    schema[field.name] = !field.required ? z.optional(zodType) : zodType
  }
  driverConfig = z.object(schema)


  return z.object({
    driver: z.enum(Object.keys(inDrivers.value)),
    ...(driverConfig ? {config: driverConfig} : {})
  })
})
type Schema = z.output<typeof schema.value>


async function onSubmit(event: FormSubmitEvent<Schema>) {

  const {status, error, data} = await useGoFetch("/sources", {method: 'post', body: form, watch: false})

  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'green', description: serverData.message})
    $emit('created')
    opened.value = false;
    return
  }

  if (status.value === 'error') {
    useToast().add({title: 'Error', description: error.value?.message, color: 'red'})
  }
}


const sendTestRequest = async () => {
  if (await formEl.value.validate(undefined, {silent: true}) === false) {
    return
  }

  const {status, error, data} = await useGoFetch("/source-drivers/test", {method: 'post', body: form, watch: false})

  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'green', description: serverData.message})
    return
  }

  if (status.value === 'error' && error.value?.statusCode === 400) {
    const serverData = error.value.data as { error?: string }
    useToast().add({title: 'Error', description: serverData.error, color: 'red'})
  }
}
</script>
<style scoped>
.nested-svg {
  display: block;
  height: 3em;
}

.nested-svg :global(svg) {
  width: 100%;
  height: 100%;
}
</style>