<template>
  <UCard class="overflow-y-auto grid grid-rows-[auto_1fr_auto] " :ui="
    { body:'p-0 sm:p-0 grid grid-cols-[1fr_3fr] overflow-y-auto',  header:'flex gap-4 justify-between',footer:'p-0 sm:p-0'}"
         v-if="status==='success' && processor">
    <template #header>
      <h1 class="text-2xl font-semibold">
        Processor n°{{ id }}
      </h1>
      <div class="flex gap-4">
        <UButton color="primary" :to="`/processors/${id}/edit`">Edit</UButton>
        <UButton color="error" loading-auto @click="()=>remove()">Delete</UButton>
      </div>
    </template>

    <template #default>
      <ProcessorFormLeftPanel v-model:form="processor" class="overflow-y-auto"/>
      <LazyProcessorFormSamplePanel v-model:form="processor" class="overflow-y-auto"/>
    </template>
    <template #footer>
      <UButton class="rounded-none" size="xl" block to="/processors">
        Return back
      </UButton>
    </template>
  </UCard>
</template>

<script setup lang="ts">

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

const {id} = useRoute().params
const {data: processors, status, refresh} = useProcessors()
const processor = computed(() => processors.value?.find(i => i.id === id))

const remove = async () => {
  const {data, error, status} = await useGoFetch(`/processors/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "error", title: "Unable to delete processor", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "success", title: "Successfully deleted processor", description: serverData.message})

    // Refresh processors list
    await refresh();
    await navigateTo('/processors')
  }
}
</script>