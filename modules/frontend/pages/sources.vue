<template>
  <UCard as="section" id="sources"
         :ui="{header:'flex gap-x-2 justify-between items-center',footer:'text-sm leading-5 text-center'}">
    <template #header>
      <h1>Sources</h1>
      <div class="flex space-x-4">
        <KbdButton :kbds="['r']" label="Refresh" color="info" variant="soft" icon="i-heroicons-arrow-path-20-solid"
                   @click="()=>refresh()" :loading="status==='pending' || status==='idle'"/>
        <KbdButton :kbds="['c']" label="Create" variant="soft" leading-icon="i-heroicons-plus" :disabled="slideoverOpen"
                   @click="()=>openCreateForm()"/>
      </div>
    </template>

    <template #footer>
      <span class="font-medium">Total:&nbsp;</span>
      <span>{{ rows?.length || 0 }}&nbsp;sources</span>
    </template>
    <UTable v-model:sorting="sorting" class="flex-1"
            :columns="columns" :data="rows" :loading="['idle','pending'].includes(status)"/>
  </UCard>

</template>


<script setup lang="ts">
import {FormSource, UButton} from "#components";
import KbdButton from "~/components/KbdButton.vue";
import type { TableColumn } from '@nuxt/ui'

const {status, data: sources, refresh} = useSources()
const slideoverOpen = useSlideover()?.isOpen

const columns: TableColumn<any>[] = [
  {
    accessorKey: 'id',
    header: ({column}) => {
      const isSorted = column.getIsSorted()

      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: '#',
        icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up-20-solid' : 'i-heroicons-bars-arrow-down-20-solid' : 'i-heroicons-arrows-up-down-20-solid',
        class: '-mx-2.5',
        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc')
      })
    }
  },
  {accessorKey: 'excerpt', header: 'Excerpt'},
  {accessorKey: 'driver', header: 'Driver'},
  {accessorKey: 'actions', header: ''}
]
const sorting = ref([{
  id: 'id',
  desc: true
}])

const rows = computed(() => sources.value?.map(source => ({
  id: source.id,
  excerpt: source.excerpt,
  driver: source.driver,
})) || [])

// Form

const openCreateForm = () => {
  useSlideover().open(FormSource, {mode: "create", onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openShowForm = (id: number) => {
  const source = sources.value?.find(s => s.id === id)
  useSlideover().open(FormSource, {mode: "read", source, onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openEditForm = (id: number) => {
  const source = sources.value?.find(s => s.id === id)
  useSlideover().open(FormSource, {mode: "edit", source, onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openCloneForm = (id: number) => {
  const source = sources.value?.find(s => s.id === id)
  useSlideover().open(FormSource, {
    mode: "create",
    source: {...toRaw(source)},
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}

const remove = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/sources/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "error", title: "Unable to delete source", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "success", title: "Successfully deleted source", description: serverData.message})
  }
  await refresh()
}
</script>