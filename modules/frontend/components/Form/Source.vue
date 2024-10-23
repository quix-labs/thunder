<template>
  <USlideover
      :title="{create: `New source`, edit: `Edit source n°${source?.id}`, read: `Source n° ${source?.id}`}[mode]"
      description="Configure data source" v-model:open="open">

    <template #body>
      <UForm class="space-y-4" :state="form" :schema @submit.prevent="submit" ref="formEl" :disabled="mode==='read'"
             id="sourceForm">
        <UFormField name="driver" label="Driver" required>
          <FormCardsInput
              :options="Object.entries(sourceDrivers || {}).map(([ID,item])=>({value:ID,item}))"
              v-model="form.driver">
            <template #default="{item}">
              <span v-html="item.config.image" v-if="item.config.image" class="nested-svg"/>
              <span v-text="item.config.name"/>
            </template>
          </FormCardsInput>
        </UFormField>

        <template v-if="sourceDrivers && form.driver">
          <!-- Driver notes-->
          <UCollapsible v-if="(sourceDrivers?.[form.driver]?.config?.notes?.length || 0) >0" default-open>
            <USeparator label="Driver Notes" type="dashed"/>
            <template #content>
              <ul class="list-disc list-outside ml-4">
                <li v-for="note in sourceDrivers[form.driver].config.notes" v-text="note"/>
              </ul>
            </template>
          </UCollapsible>

          <UCollapsible default-open>
            <USeparator label="Driver Configuration"/>
            <template #content>
              <LazyFormDynamicFields :state="form.config" :fields="sourceDrivers[form.driver].fields"/>
            </template>
          </UCollapsible>

        </template>
      </UForm>
    </template>
    <template #footer>
      <UButton type="button" @click.prevent="sendTestRequest" variant="soft" color="secondary">
        Test configuration
      </UButton>
      <UButton type="submit" class="ml-auto" v-if="mode!=='read'" form="sourceForm">
        {{ {create: 'Create', edit: 'Save changes'}[mode] }}
      </UButton>
    </template>
  </USlideover>
</template>
<script setup lang="ts">
import {z} from "zod";
import {UForm} from "#components";

interface Props {
  mode: "create" | "edit" | "read"
  source?: any
}

const open = defineModel('open', {type: Boolean})
const {mode = "create", source} = defineProps<Props>()

const $emit = defineEmits(['created', 'updated'])
const formEl = ref<HTMLFormElement>()

const {data: sourceDrivers} = useSourceDrivers()

const defaultForm = reactive({
  driver: null as string | null,
  config: {} as { [key: string]: any } | null,
})

const form = computed(() => source || defaultForm)

const schema = computed(() => {
  const driverIDS = Object.keys(sourceDrivers.value || {})
  return z.object({
    driver: (driverIDS.length > 0) ? z.enum([driverIDS[0], ...driverIDS.slice(1)], {
      required_error: "Driver is required",
    }) : z.string({required_error: "Driver is required"})
  })
})

const submit = async () => {
  if (mode === "create") {
    const {status, error, data} = await useGoFetch("/sources", {method: 'post', body: form, watch: false})
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'success', description: serverData.message})
      $emit('created')
    } else if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.data?.error || error.value?.message, color: 'error'})
    }
  } else if (mode === "edit") {
    const {status, error, data} = await useGoFetch(`/sources/${source.id}`, {
      method: 'put',
      body: form,
      watch: false
    })
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'success', description: serverData.message})
      $emit('updated')
    } else if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.message, color: 'error'})
    }
  }
  open.value = false
}

const sendTestRequest = async () => {
  if (await formEl.value?.validate(undefined, {silent: true}) === false) {
    return
  }

  const {status, error, data} = await useGoFetch("/source-drivers/test", {method: 'post', body: form, watch: false})

  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'success', description: serverData.message})
    return
  }

  if (status.value === 'error' && [400, 422, 500]?.includes(error.value?.statusCode || 0)) {
    const serverData = error.value?.data as { error?: string } | undefined
    useToast().add({title: 'Error', description: serverData?.error, color: 'error'})
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