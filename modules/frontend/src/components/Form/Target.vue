<template>
  <USlideover
      :title="{create: `New target`, edit: `Edit target n°${target?.id}`, read: `Target n° ${target?.id}`}[mode]"
      description="Configure data target" v-model:open="open">
    <template #body>
      <UForm class="space-y-4" :state="form" :schema @submit.prevent="submit" ref="formEl" :disabled="mode==='read'"
             id="targetForm">
        <UFormField name="driver" label="Driver" required>
          <FormCardsInput
              :options="Object.entries(targetDrivers || {}).map(([ID,item])=>({value:ID,item}))"
              v-model="form.driver">
            <template #default="{item}">
              <span v-html="item.config.image" v-if="item.config.image" class="nested-svg"/>
              <span v-text="item.config.name"/>
            </template>
          </FormCardsInput>
        </UFormField>
        <template v-if="targetDrivers && form.driver">
          <!-- Driver notes-->
          <UCollapsible v-if="(targetDrivers?.[form.driver]?.config?.notes?.length || 0) >0">
            <USeparator label="Driver Notes"/>
            <template #content>
              <ul class="list-disc list-outside ml-4">
                <li v-for="note in targetDrivers[form.driver].config.notes" v-text="note"/>
              </ul>
            </template>
          </UCollapsible>
          <UCollapsible>
            <USeparator label="Driver Configuration"/>
            <LazyFormDynamicFields :state="form.config" :disabled="mode==='read'" :fields="targetDrivers[form.driver].fields"/>
          </UCollapsible>
        </template>
      </UForm>
    </template>
    <template #footer>
      <UButton type="button" @click.prevent="sendTestRequest" variant="soft" color="secondary">
        Test configuration
      </UButton>
      <UButton type="submit" class="ml-auto" v-if="mode!=='read'" form="targetForm">
        {{ {create: 'Create', edit: 'Save changes'}[mode] }}
      </UButton>
    </template>
  </USlideover>
</template>
<script setup lang="ts">
import {z} from "zod";

interface Props {
  mode: "create" | "edit" | "read"
  target?: any
}

const {mode = "create", target} = defineProps<Props>()
const open = defineModel('open', {type: Boolean})

const $emit = defineEmits(['created', 'updated'])
const formEl = ref<HTMLFormElement>()

const {data: targetDrivers} = useTargetDrivers()

const defaultForm = reactive({
  driver: null as string | null,
  config: {} as { [key: string]: any } | null,
})

const form = computed(() => target || defaultForm)

const schema = computed(() => {
  const driverIDS = Object.keys(targetDrivers.value || {})
  return z.object({
    driver: (driverIDS.length > 0) ? z.enum([driverIDS[0], ...driverIDS.slice(1)], {
      required_error: "Driver is required",
    }) : z.string({required_error: "Driver is required"})
  })
})

const submit = async () => {
  if (mode === "create") {
    const {status, error, data} = await useGoFetch("/targets", {method: 'post', body: form, watch: false})
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'success', description: serverData.message})
      $emit('created')
    } else if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.data?.error || error.value?.message, color: 'red'})
    }
  } else if (mode === "edit") {
    const {status, error, data} = await useGoFetch(`/targets/${target.id}`, {
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

  const {status, error, data} = await useGoFetch("/target-drivers/test", {method: 'post', body: form, watch: false})

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