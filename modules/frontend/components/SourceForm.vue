<template>
  <USlideover v-model="opened" class="overflow-y-auto" @keydown.ctrl.enter.prevent="submit">
    <UForm class="flex flex-col flex-1" :state="form" :schema @submit.prevent="submit" ref="formEl">
      <UCard class="flex flex-col flex-1" :ui="{body: { base: 'flex-1' },rounded:''}">
        <template #header>
          <h2 class="text-center text-2xl font-semibold">
            {{ {create: `New source`, edit: `Edit source n°${source?.id}`, read: `Source n° ${source?.id}`}[mode] }}
          </h2>
        </template>
        <fieldset :disabled="mode==='read'" class="grid gap-y-4">
          <FormCardsInput
              name="driver" label="Driver" required
              :options="Object.values(sourceDrivers || {}).map((item)=>({value:item.config.ID,item}))"
              v-model="form.driver">
            <template #default="{item}">
              <span v-html="item.config.image" v-if="item.config.image" class="nested-svg"/>
              <span v-text="item.config.name"/>
            </template>
          </FormCardsInput>

          <template v-if="sourceDrivers && form.driver">
            <!-- Driver notes-->
            <div v-if="(sourceDrivers?.[form.driver]?.config?.notes?.length || 0) >0">
              <UDivider label="Driver Notes"/>
              <ul class="list-disc list-outside ml-4">
                <li v-for="note in sourceDrivers[form.driver].config.notes" v-text="note"/>
              </ul>
            </div>
            <UDivider label="Driver Configuration"/>
            <LazyFormDynamicFields
                v-model="form.config"
                v-model:schema="configSchema"
                :fields="sourceDrivers[form.driver].fields"
                v-if="sourceDrivers && form.driver"/>
          </template>
        </fieldset>
        <template #footer>
          <div class="flex">
            <UButton type="button" @click.prevent="sendTestRequest" variant="soft" color="sky">
              Test configuration
            </UButton>
            <UButton type="submit" class="ml-auto" v-if="mode!=='read'">
              {{ {create: 'Create', edit: 'Save changes'}[mode] }}
            </UButton>
          </div>
        </template>
      </UCard>
    </UForm>
  </USlideover>
</template>
<script setup lang="ts">
import {z, type ZodObject, type ZodRawShape} from "zod";

interface Props {
  mode: "create" | "edit" | "read"
  source?: any
}

const {mode = "create", source} = defineProps<Props>()

const opened = defineModel('opened')
const $emit = defineEmits(['created', 'updated'])
const formEl = ref<HTMLFormElement>()

const {data: sourceDrivers} = useSourceDrivers()

const defaultForm = reactive({
  driver: null as string | null,
  config: {} as { [key: string]: any } | null,
})

const form = computed(() => source || defaultForm)

const configSchema = ref<ZodObject<ZodRawShape> | undefined>()
const schema = computed(() => {

  let base = z.object({
    driver: z.enum(Object.keys(sourceDrivers.value || {}), {
      required_error: "Driver is required",
    }),
  })

  if (form.value.driver && configSchema.value) {
    base = base.merge(z.object({
      config: configSchema.value
    }))
  }

  return base;
})
const submit = async () => {
  if (mode === "create") {
    const {status, error, data} = await useGoFetch("/sources", {method: 'post', body: form, watch: false})
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'green', description: serverData.message})
      $emit('created')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.data?.error || error.value?.message, color: 'red'})
    }
  } else if (mode === "edit") {
    const {status, error, data} = await useGoFetch(`/sources/${source.id}`, {
      method: 'put',
      body: form,
      watch: false
    })
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'green', description: serverData.message})
      $emit('updated')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.message, color: 'red'})
    }
  }

}

const sendTestRequest = async () => {
  if (await formEl.value?.validate(undefined, {silent: true}) === false) {
    return
  }

  const {status, error, data} = await useGoFetch("/source-drivers/test", {method: 'post', body: form, watch: false})

  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'green', description: serverData.message})
    return
  }

  if (status.value === 'error' && [400, 422, 500]?.includes(error.value?.statusCode || 0)) {
    const serverData = error.value?.data as { error?: string } | undefined
    useToast().add({title: 'Error', description: serverData?.error, color: 'red'})
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