<template>
  <section>
    <UCard class="w-full" :ui="{
      divide: 'divide-y divide-gray-200 dark:divide-gray-700',
      body: { padding: '' },
      header:{base:'flex gap-x-2 justify-between items-center'},
      footer:{base:'text-sm leading-5 text-center'}
    }">
      <template #header>
        <h1 class="font-semibold text-xl text-gray-900 dark:text-white leading-tight">Processors</h1>
        <UButton @click.prevent="createProcessor" variant="soft">+ Add Processor</UButton>
      </template>
      <template #default>
        <UTable @select="(row:any)=>showProcessor(row.id)" :columns="columns" :rows="rows"
                :sort="{column:'id',direction:'desc'}"
                :loading="status==='pending' || status==='idle'">
          <template #actions-data="{ row }">
            <div class="flex gap-1">
              <UDropdown :items="[[
                  {label:'Replicate',click:()=>cloneProcessor(row.id)},
                  {label:'Claim indexing',click:()=>claimIndex(row.id),disabled:row.indexing},
                  {label:'Start listening',click:()=>claimStart(row.id),disabled:row.listening},
                  {label:'Stop listening',click:()=>claimStop(row.id),disabled:!row.listening},
              ]]" @click.stop>
                <UButton icon="i-heroicons-ellipsis-horizontal" variant="link" color="gray" size="xl" :padded="false"/>
              </UDropdown>
              <UButton icon="i-heroicons-eye" variant="link" color="gray" size="xl" :padded="false"
                       @click.stop.prevent="showProcessor(row.id)"/>
              <UButton icon="i-heroicons-pencil-square" variant="link" color="gray" size="xl" :padded="false"
                       @click.stop.prevent="editProcessor(row.id)"/>
              <UButton icon="i-heroicons-trash" variant="link" color="red" size="xl" :padded="false"
                       @click.stop.prevent="deleteProcessor(row.id)"/>
            </div>
          </template>
          <template #targets-data="{ row }">
            <div class="flex gap-1">
              <UBadge size="xs" :label="`Target n°${target}`" color="sky" variant="subtle"
                      v-for="target in row.targets"/>
            </div>
          </template>
          <template #source-data="{ row }">
            <div class="flex gap-1">
              <UBadge size="xs" :label="`Source n°${row.source}`" color="sky" variant="subtle"/>
            </div>
          </template>
          <template #stats-data="{ row }">
            <div class="flex gap-1">
              <UBadge size="xs" :label="`${row.stats.total} fields`" color="gray"/>
              <UBadge size="xs" :label="`${row.stats.relations} relations`" color="gray"/>
              <UBadge size="xs" :label="`${row.stats.conditions} conditions`" color="gray"/>
            </div>
          </template>
          <template #indexing-data="{ row }">
            <UBadge size="xs" :label="row.indexing?'Indexing':'Not indexing'" :color="row.indexing?'green':'red'"/>
          </template>
          <template #listening-data="{ row }">
            <UBadge size="xs" :label="row.listening?'Listening':'Not listening'" :color="row.listening?'green':'red'"/>
          </template>
        </UTable>
      </template>
      <template #footer>
        <span class="font-medium">Total:&nbsp;</span>
        <span>{{ rows?.length || 0 }}&nbsp;processors</span>
      </template>
    </UCard>
    <ProcessorFormComplete :mode="formMode" v-model:opened="formOpened" :processor="formProcessor"/>
  </section>

</template>

<script setup lang="ts">

const {status, data: processors} = useProcessors()

// Table
const columns = [
  {key: 'id', label: '#', sortable: true, rowClass: 'w-[1px] whitespace-nowrap'},
  {key: 'source', label: 'Source', sortable: true},
  {key: 'index', label: 'Index', sortable: true},
  {key: 'targets', label: 'Targets', sortable: false},
  {key: 'stats', label: 'Stats', sortable: false},
  {key: 'listening', label: 'Listening', sortable: false},
  {key: 'indexing', label: 'Indexing', sortable: false},
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

const formOpened = ref(false);
const formMode = ref<"create" | "edit" | "read">("create");
const formProcessor = ref<any>();

const createProcessor = async () => {
  formMode.value = "create"
  formProcessor.value = undefined
  formOpened.value = true
}
const showProcessor = async (id: number) => {
  formMode.value = "read"
  formProcessor.value = processors.value?.filter(p => p.id === id)?.at(0)
  formOpened.value = true
}
const editProcessor = async (id: number) => {
  formMode.value = "edit"
  formProcessor.value = processors.value?.filter(p => p.id === id)?.at(0)
  formOpened.value = true
}
const cloneProcessor = async (id: number) => {
  formMode.value = "create"
  formProcessor.value = {...processors.value?.filter(p => p.id === id)?.at(0)}
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
}

const claimIndex = (id: number) => {
  useGoFetch(`/processors/${id}/index`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "red",
        title: "Unable to indexing",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "green", title: "Successfully indexed", description: serverData.message})
    }
  })
}
const claimStart = (id: number) => {
  useGoFetch(`/processors/${id}/start`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "red",
        title: "Unable to start",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "green", title: "Successfully started", description: serverData.message})
    }
  })
}
const claimStop = (id: number) => {
  useGoFetch(`/processors/${id}/stop`, {method: "POST"}).then(({data, error, status}) => {
    if (status.value === "error") {
      useToast().add({
        color: "red",
        title: "Unable to stop",
        description: error.value?.data?.error || error.value?.message
      })
    } else if (status.value === "success") {
      const serverData = data.value as { message?: string }
      useToast().add({color: "green", title: "Successfully stopped", description: serverData.message})
    }
  })
}
</script>