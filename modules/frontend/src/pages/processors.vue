<template>
  <UCard as="section" id="processors"
         :ui="{body:'p-0 sm:p-0', header:'flex gap-x-2 justify-between items-center',footer:'text-sm leading-5 text-center'}">
    <template #header>
      <h1>Processors</h1>
      <div class="flex space-x-4">
        <KbdButton :kbds="['r']" label="Refresh" color="info" variant="soft" icon="i-heroicons-arrow-path-20-solid"
                   @click="()=>refresh()" :loading="status==='pending' || status==='idle'"/>
        <KbdButton :kbds="['c']" label="Create" variant="soft" leading-icon="i-heroicons-plus" :disabled="slideoverOpen"
                   @click="()=>openCreateForm()"/>
      </div>
    </template>

    <CustomTable :columns :rows :loading="['idle','pending'].includes(status)" :sorting="[{desc:true,id:'id'}]">
      <template #cell-source="{ row }">
        <div class="flex gap-1">
          <UBadge size="sm" :label="`Source n°${row.source}`" color="secondary" variant="subtle"/>
        </div>
      </template>
      <template #cell-targets="{row}">
        <div class="flex gap-1">
          <UBadge size="sm" :label="`Target n°${target}`" color="secondary" variant="subtle"
                  v-for="target in row.targets"/>
        </div>
      </template>


      <template #cell-stats="{ row }">
        <div class="flex gap-1">
          <UBadge size="sm" :label="`${row.stats.total} fields`" color="neutral" variant="subtle"/>
          <UBadge size="sm" :label="`${row.stats.relations} relations`" color="neutral" variant="subtle"/>
          <UBadge size="sm" :label="`${row.stats.conditions} conditions`" color="neutral" variant="subtle"/>
        </div>
      </template>
      <template #cell-indexing="{ row }">
        <UBadge :label="row.indexing?'Indexing':'Not indexing'"
                :color="row.indexing?'success':'error'" variant="subtle" size="sm"
        />
      </template>
      <template #cell-listening="{ row }">
        <UBadge :label="row.listening?'Listening':'Not listening'"
                :color="row.listening?'success':'error'" variant="subtle" size="sm"/>
      </template>

      <template #cell-actions="{ row }">
        <div class="flex gap-2 justify-end">
          <UDropdownMenu :items="[[
                  {label:'Replicate',onSelect:()=>openCloneForm(row.id)},
                  {label:'Claim indexing',onSelect:()=>claimIndex(row.id),disabled:row.indexing},
                  {label:'Start listening',onSelect:()=>claimStart(row.id),disabled:row.listening},
                  {label:'Stop listening',onSelect:()=>claimStop(row.id),disabled:!row.listening},
              ]]" @click.stop>
            <UButton icon="i-heroicons-ellipsis-horizontal" variant="link" color="neutral" class="p-0" size="xl"/>
          </UDropdownMenu>

          <UButton icon="i-heroicons-eye" variant="link" color="neutral" class="p-0" size="xl"
                   @click.stop.prevent="openShowForm(row.id)"/>
          <UButton icon="i-heroicons-pencil-square" variant="link" color="neutral" class="p-0" size="xl"
                   @click.stop.prevent="openEditForm(row.id)"/>
          <UButton icon="i-heroicons-document-arrow-down" variant="link" color="primary" class="p-0" size="xl"
                   @click.stop.prevent="openDownloadForm(row.id)"/>
          <UButton icon="i-heroicons-trash" variant="link" color="error" class="p-0" size="xl"
                   @click.stop.prevent="remove(row.id)"/>
        </div>
      </template>
    </CustomTable>
    <template #footer>
      <span class="font-medium">Total:&nbsp;</span>
      <span>{{ rows?.length || 0 }}&nbsp;processors</span>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import KbdButton from "../components/KbdButton.vue";
import {LazyDownloadProcessorForm, ProcessorFormComplete, UButton, UDropdownMenu} from "#components";

const {status, data: processors, refresh} = useProcessors()
const slideoverOpen = useSlideover()?.isOpen

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
    processor: reactive({...JSON.parse(JSON.stringify(processor))}),
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}

const openDownloadForm = (id: number) => {
  useModal().open(LazyDownloadProcessorForm, {processorId: id})
}

const remove = async (id: number) => {
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
  useToast().add({color: "error", title: "Not implemented yet!", description: "Currently working on!"})
  return

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
  useToast().add({color: "error", title: "Not implemented yet!", description: "Currently working on!"})
  return

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