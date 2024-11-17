<template>
  <UCard :as="UForm" :state="form" @submit.prevent="submit" class="overflow-y-auto grid grid-rows-[auto_1fr_auto]" :ui="
  state.showSamplePanel
    ? { body:'p-0 sm:p-0 grid grid-cols-[1fr_3fr_1.5fr] overflow-y-auto',  footer:'p-0 sm:p-0'}
    : { body:'p-0 sm:p-0 grid grid-cols-[1fr_4.5fr] overflow-y-auto', footer:'p-0 sm:p-0' }
">
    <template #header>

      <ProcessorFormToolbar title="Create processor"/>
    </template>

    <template #default>
      <ProcessorFormLeftPanel v-model:form="form" class="overflow-y-auto"/>
      <ProcessorFormRightPanel v-model:form="form" class="overflow-y-auto"/>
      <LazyProcessorFormSamplePanel v-model:form="form" class="overflow-y-auto" v-if="state.showSamplePanel"/>
    </template>

    <template #footer>
      <UButton class="rounded-none" size="xl" block type="submit">
        Create
      </UButton>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import {UForm} from "#components";

definePageMeta({
  container: false,
  fullHeight: true,
})

const state = useProcessFormState()

onBeforeMount(async () => {
  const from = useRoute().query?.from
  if (!from) return

  const {data: processors, status, error} = await useProcessors({lazy: false})
  if (status.value !== 'success') {
    useToast().add({color: 'error', title: 'Unable to fetch processors', description: error.value?.message})
    return
  }

  const processor = processors.value?.find(i => i.id === from)
  if (!processor) {
    useToast().add({color: 'error', title: `Cannot clone processor ${from}`, description: "Processor not found"})
    return
  }

  // Fill form with default
  for (const key in processor) {
    if (key === 'id') continue;
    if (Object.prototype.hasOwnProperty.call(form, key)) {
      (form as any)[key] = processor[key];
    }
  }
})
const form = reactive({
  source: null as number | null,
  table: null as string | null,
  primary_keys: [] as string[] | null,
  conditions: [] as any[] | null,

  targets: [] as string[],
  index: null as string | null,
  mapping: [] as string[]
});


const submit = async () => {
  const {status, error, data} = await useGoFetch("/processors", {method: 'post', body: form, watch: false})
  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'success', description: serverData.message})

    // Refresh processors list
    await useProcessors().refresh()

    return navigateTo('/processors');
  }

  if (status.value === 'error') {
    useToast().add({title: 'Error', description: error.value?.message, color: 'error'})
  }
}
</script>

