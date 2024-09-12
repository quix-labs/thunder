<template>
  <USlideover v-model="opened" class="overflow-y-auto" @keydown.ctrl.enter.prevent="submit">
    <UForm class="flex flex-col flex-1" :state="form" @submit.prevent="submit" ref="formEl">
      <UCard class="flex flex-col flex-1" :ui="{body: { base: 'flex-1' },rounded:''}">
        <template #header>
          <h2 class="text-center text-2xl font-semibold">
            {{ {create: `New target`, edit: `Edit target: TODO id`, read: 'Target: TODO id'}[mode] }}
          </h2>
        </template>
        <fieldset :disabled="mode==='read'" class="grid gap-y-4">
          <UFormGroup label="Driver" required name="driver">
            <div class="grid grid-cols-2 gap-4">
              <div v-for="([key,driver]) in Object.entries(targetDrivers || {})">
                <input type="radio" name="driver" v-model="form.driver" :value="key" class="sr-only peer"
                       :id="`driver-${key}`"
                       @change="formEl?.clear()"
                       tabindex="-1"
                >
                <label :for="`driver-${key}`"
                       tabindex="0"
                       @keydown.enter.space.prevent="form.driver=key"
                       class="cursor-pointer peer-disabled:cursor-not-allowed flex flex-col gap-y-1 text-center items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500">
                  <span v-html="driver.config.image" v-if="driver.config.image" class="nested-svg"/>
                  <span v-text="driver.config.name"/>
                </label>
              </div>
            </div>
          </UFormGroup>

          <template v-if="targetDrivers && form.driver">
            <!-- Driver notes-->
            <div v-if="targetDrivers[form.driver]?.config?.notes?.length>0">
              <UDivider label="Driver Notes"/>
              <ul class="list-disc list-outside ml-4">
                <li v-for="note in targetDrivers[form.driver]?.config?.notes" v-text="note"/>
              </ul>
            </div>
            <UDivider label="Driver Configuration"/>
            <DynamicFieldsForm v-model="form.config" :fields="targetDrivers[form.driver].fields"
                               v-if="targetDrivers && form.driver"/>
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
interface Props {
  mode: "create" | "edit" | "read"
  target?: any
  targetId?: number
}

const {mode = "create", target, targetId} = defineProps<Props>()

const opened = defineModel('opened')
const $emit = defineEmits(['created', 'updated'])
const formEl = ref<HTMLFormElement>()

const {data: targetDrivers} = useTargetDrivers()

const defaultForm = reactive({
  driver: null as string | null,
  config: {} as { [key: string]: any } | null,
})

const form = computed(() => target || defaultForm)

const submit = async () => {
  if (mode === "create") {
    const {status, error, data} = await useGoFetch("/targets", {method: 'post', body: form, watch: false})
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
    //TODO
    // const {status, error, data} = await useGoFetch(`/processors/${processorId}`, {
    //   method: 'put',
    //   body: form,
    //   watch: false
    // })
    // if (status.value === 'success') {
    //   const serverData = data.value as { message?: string }
    //   useToast().add({title: 'Success', color: 'green', description: serverData.message})
    //   $emit('updated')
    //   opened.value = false
    //   return
    // }
    //
    // if (status.value === 'error') {
    //   useToast().add({title: 'Error', description: error.value?.message, color: 'red'})
    // }
  }

}

const sendTestRequest = async () => {
  if (await formEl.value?.validate(undefined, {silent: true}) === false) {
    return
  }

  const {status, error, data} = await useGoFetch("/target-drivers/test", {method: 'post', body: form, watch: false})

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