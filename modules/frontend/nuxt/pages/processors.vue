<template>
  <section>
    <UCard>
      <template #header>
        <div class="flex gap-x-2 items-center justify-between">
          <h1>Processors</h1>
          <UButton @click.prevent="createProcessor" variant="soft">+ Add Processor</UButton>
        </div>
      </template>
      <template #default>


        <div class="grid grid-cols-3">
          <div v-for="(processor,index) in processors||[]" v-if="status==='success'" class="grid grid-cols-3">
            <div class="col-span-full">
              Processor {{processor.index||index}} (TODO user defined name)
            </div>

            <UButton color="sky" icon="i-heroicons-eye" @click.prevent="showProcessor(index)">
              Show
            </UButton>
            <UButton color="sky" icon="i-heroicons-pencil" @click.prevent="editProcessor(index)">
              Edit
            </UButton>
            <UButton color="orange" icon="i-heroicons-trash" @click.prevent="cloneProcessor(index)">
              Duplicate
            </UButton>
            <UButton color="red" icon="i-heroicons-trash" @click.prevent="deleteProcessor(index)">
              Delete
            </UButton>
          </div>
          <div class="flex gap-2 items-center" v-else>
            <UIcon size="large" name="i-heroicons-arrow-path" class="animate-spin"/>
            <span>Loading</span>
          </div>
        </div>
      </template>
    </UCard>

    <ProcessorFormComplete
        @updated="refresh"
        @created="refresh"
        :mode="formMode"
        v-model:opened="formOpened"
        :processor="formProcessor"
    />
  </section>

</template>

<script setup lang="ts">
const {status,data: processors, refresh} = useProcessors()

const formOpened = ref(false);
const formMode = ref<"create" | "edit" | "read">("create");
const formProcessor = ref<any>();

const createProcessor = async (id: number) => {
  formMode.value="create"
  formProcessor.value = undefined
  formOpened.value = true
}
const showProcessor = async (id: number) => {
  formMode.value="read"
  formProcessor.value = processors.value?.at(id)
  formOpened.value = true
}
const editProcessor = async (id: number) => {
  formMode.value="edit"
  formProcessor.value = processors.value?.at(id)
  formOpened.value = true
}
const cloneProcessor = async (id: number) => {
  formMode.value="create"
  formProcessor.value = processors.value?.at(id)
  formOpened.value = true
}

const deleteProcessor = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/processors/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "red", title: "Unable to delete processor", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "green", title: "Successfully deleted processor", description: serverData.message})
  }
  await refresh()
}
</script>