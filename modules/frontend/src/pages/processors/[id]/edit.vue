<template>
  <UCard :as="UForm" :state="form" @submit.prevent="submit" class="overflow-y-auto grid grid-rows-[auto_1fr_auto]" :ui="
  state.showSamplePanel
    ? { body:'p-0 sm:p-0 grid grid-cols-[1fr_3fr_1.5fr] overflow-y-auto',  footer:'p-0 sm:p-0'}
    : { body:'p-0 sm:p-0 grid grid-cols-[1fr_4.5fr] overflow-y-auto', footer:'p-0 sm:p-0' }
">
    <template #header>
      <ProcessorFormToolbar :title="`Edit processor n°${id}`">
        <UButton color="secondary" to="/processors">Go back</UButton>
      </ProcessorFormToolbar>
    </template>

    <template #default>
      <ProcessorFormLeftPanel v-model:form="form" class="overflow-y-auto"/>
      <ProcessorFormRightPanel v-model:form="form" class="overflow-y-auto"/>
      <LazyProcessorFormSamplePanel v-model:form="form" class="overflow-y-auto" v-if="state.showSamplePanel"/>
    </template>

    <template #footer>
      <UButton loading-auto class="rounded-none" size="xl" block type="submit">
        Save changes
      </UButton>
    </template>
  </UCard>
</template>

<script setup lang="ts">

import {UForm} from "#components";

definePageMeta({
  container: false,
  fullHeight: true,

  validate: async (route) => {
    const {id} = route.params;
    const {data: processors, error, status} = await useProcessors({lazy: false, server: true})
    if (status.value !== 'success') {
      console.error(error.value)
      return false
    }
    return !!processors.value?.find(i => i.id === id)
  }
})
const state = useProcessFormState()

const {id} = useRoute().params
const {data: processors, status,refresh} = useProcessors()

const form = computed(() => processors.value?.find(i => i.id === id))

const submit = async () => {
  const {status, error, data} = await useGoFetch(`/processors/${id}`, {
    method: 'put',
    body: form,
    watch: false
  })

  if (status.value === 'success') {
    const serverData = data.value as { message?: string }
    useToast().add({title: 'Success', color: 'success', description: serverData.message})

    // Refresh processors list
    await refresh()
  }

  if (status.value === 'error') {
    useToast().add({title: 'Error', description: error.value?.message, color: 'error'})
  }
}
</script>