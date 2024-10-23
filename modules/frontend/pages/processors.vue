<template>
  <UCard as="section" id="processors"
         :ui="{header:'flex gap-x-2 justify-between items-center',footer:'text-sm leading-5 text-center'}">
    <template #header>
      <h1>Processors</h1>
      <div class="flex space-x-4">
        <KbdButton :kbds="['r']" label="Refresh" color="info" variant="soft" icon="i-heroicons-arrow-path-20-solid"
                   @click="()=>refresh()" :loading="status==='pending' || status==='idle'"/>
        <KbdButton :kbds="['c']" label="Create" variant="soft" leading-icon="i-heroicons-plus" :disabled="slideoverOpen"
                   @click="()=>openCreateForm()"/>
      </div>
    </template>
    {{ focusedId }}
    <div v-for="processor in processors" class="flex justify-between w-full" @focusin.passive="focusedId=processor.id"
         tabindex="0">
      {{ processor.id }} - {{ processor.index }}
      <div class="flex gap-x-4">
        <KbdButton :kbds="focusedId===processor.id ? ['s'] : undefined" label="Show" color="neutral" variant="soft"
                   icon="i-heroicons-eye"
                   @click="()=>openShowForm(processor.id)"/>
        <KbdButton :kbds="focusedId===processor.id ? ['e'] : undefined"  label="Edit" color="info" variant="soft" icon="i-heroicons-pencil-square"
                   @click="()=>openEditForm(processor.id)"/>
        <KbdButton :kbds="focusedId===processor.id ? ['d'] : undefined"  label="Download" variant="soft" color="primary"
                   leading-icon="i-heroicons-document-arrow-down"
                   @click="()=>openDownloadForm(processor.id)"/>
        <KbdButton :kbds="focusedId===processor.id ? ['meta','d'] : undefined" label="Delete" variant="soft" color="error" leading-icon="i-heroicons-trash"
                   @click="()=>deleteProcessor(processor.id)"/>
      </div>
    </div>

    <template #footer>
      <span class="font-medium">Total:&nbsp;</span>
      <span>{{ rows?.length || 0 }}&nbsp;processors</span>
    </template>
  </UCard>

</template>

<script setup lang="ts">
import KbdButton from "../components/KbdButton.vue";
import {LazyDownloadProcessorForm, ProcessorFormComplete} from "#components";

const {status, data: processors, refresh} = useProcessors()
const slideoverOpen = useSlideover()?.isOpen

const focusedId = ref<number>()

// Table
const columns = [
  {key: 'id', label: '#', sortable: true, rowClass: 'w-[1px] whitespace-nowrap'},
  {key: 'source', label: 'Source', sortable: true},
  {key: 'index', label: 'Index', sortable: true},
  {key: 'targets', label: 'Targets', sortable: true},
  {key: 'stats', label: 'Stats', sortable: false},
  {key: 'listening', label: 'Listening', sortable: true},
  {key: 'indexing', label: 'Indexing', sortable: true},
  {key: 'actions', sortable: false, rowClass: 'w-[1px] whitespace-nowrap'}
]

const rows = computed(() => processors.value?.map(processor => ({
  id: processor.id,
  targets: processor.targets,
  index: processor.index || '---',
  source: processor.source !== undefined ? processor.source : '---',
  stats: {
    conditions: processor.conditions?.length || 0,
    ...getMappingStats(processor.mapping || [])
  },
  indexing: processor.indexing || false,
  listening: processor.listening || false,
})) || [])

const getMappingStats = (mapping: any[]) => {
  const relations = mapping.filter(i => i._type === 'relation');
  const stats = {
    total: mapping.length,
    relations: relations.length,
  }

  mapping.filter(i => i._type === 'relation').forEach(i => {
    const {total: relTotal, relations: relRelations} = getMappingStats(i.mapping || [])
    stats.total += relTotal
    stats.relations += relRelations
  })
  return stats
}

// Form

const openCreateForm = () => {
  useSlideover().open(ProcessorFormComplete, {mode: "create", onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openShowForm = (id: number) => {
  const processor = processors.value?.find(s => s.id === id)
  useSlideover().open(ProcessorFormComplete, {
    mode: "read",
    processor,
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}
const openEditForm = (id: number) => {
  const processor = processors.value?.find(s => s.id === id)
  useSlideover().open(ProcessorFormComplete, {
    mode: "edit",
    processor,
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}

const openCloneForm = (id: number) => {
  const processor = processors.value?.find(s => s.id === id)
  useSlideover().open(ProcessorFormComplete, {
    mode: "create",
    processor: {...toRaw(processor)},
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}

const openDownloadForm = (id: number) => {
  useModal().open(LazyDownloadProcessorForm,{processorId:id})
}

const deleteProcessor = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/processors/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "error", title: "Unable to delete processor", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "success", title: "Successfully deleted processor", description: serverData.message})
  }
}

const claimIndex = (id: number) => {
  useGoFetch(`/processors/${id}/index`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "error",
        title: "Unable to indexing",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "success", title: "Successfully indexed", description: serverData.message})
    }
  })
}
const claimStart = (id: number) => {
  useGoFetch(`/processors/${id}/start`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "error",
        title: "Unable to start",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "success", title: "Successfully started", description: serverData.message})
    }
  })
}
const claimStop = (id: number) => {
  useGoFetch(`/processors/${id}/stop`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "error",
        title: "Unable to stop",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "success", title: "Successfully stopped", description: serverData.message})
    }
  })
}
</script>