<template>
  <section>
    <UCard>
      <template #header>
        <div class="flex gap-x-2 items-center justify-between">
          <h1>Data Sources</h1>
          <UButton @click.prevent="createSlideoverOpen=!createSlideoverOpen" variant="soft">+ Add Source</UButton>
        </div>
      </template>
      <template #default>
        <div class="grid grid-cols-3">
          <div v-for="(source,index) in sources||[]">

          <span v-if="sourceDrivers && sourceDrivers[source.driver]?.config?.image"
                v-html="sourceDrivers[source.driver]?.config?.image" class="nested-svg"/>
            Source n°{{ index }} (TODO user defined name)

            <UButton color="red" icon="i-heroicons-trash" @click.prevent="deleteSource(index)">
              Delete
            </UButton>

          </div>
        </div>

      </template>
    </UCard>
    <CreateSourceDriverForm v-model:opened="createSlideoverOpen" @created="refresh"/>
  </section>


</template>
<script setup lang="ts">

const createSlideoverOpen = ref(false)
const {data: sources, refresh} = useSources()
const {data: sourceDrivers} = useSourceDrivers()

const deleteSource = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/sources/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "red", title: "Unable to delete source", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "green", title: "Successfully deleted source", description: serverData.message})
  }
  await refresh()
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