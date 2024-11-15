<template>
  <UCard as="section" id="targets"
         :ui="{body:'p-0 sm:p-0', header:'flex gap-x-2 justify-between items-center',footer:'text-sm leading-5 text-center'}">
    <template #header>
      <h1>Targets</h1>
      <div class="flex space-x-4">
        <KbdButton :kbds="['r']" label="Refresh" color="info" variant="soft" icon="i-heroicons-arrow-path-20-solid"
                   @click="()=>refresh()" :loading="status==='pending' || status==='idle'"/>
        <KbdButton :kbds="['c']" label="Create" variant="soft" leading-icon="i-heroicons-plus" :disabled="slideoverOpen"
                   @click="()=>openCreateForm()"/>
      </div>
    </template>
    <CustomTable :columns :rows :loading="['idle','pending'].includes(status)"  :sorting="[{desc:true,id:'id'}]">
      <template #cell-driver="{ row }">
        <UBadge size="sm" :label="row.driver" color="secondary" variant="subtle"/>
      </template>
      <template #cell-actions="{ row }">
        <div class="flex gap-2 justify-end">
          <UDropdownMenu :items="[[{label:'Replicate',onSelect:()=>openCloneForm(row.id)}]]" @click.stop>
            <UButton icon="i-heroicons-ellipsis-horizontal" variant="link" color="neutral" class="p-0" size="xl"/>
          </UDropdownMenu>
          <UButton icon="i-heroicons-eye" variant="link" color="neutral" size="xl"
                   @click.stop.prevent="openShowForm(row.id)" class="p-0"/>
          <UButton icon="i-heroicons-pencil-square" variant="link" color="neutral" size="xl"
                   @click.stop.prevent="openEditForm(row.id)" class="p-0"/>
          <UButton icon="i-heroicons-trash" variant="link" color="error" size="xl" @click.stop.prevent="remove(row.id)"
                   class="p-0"/>
        </div>
      </template>
    </CustomTable>
    <template #footer>
      <span class="font-medium">Total:&nbsp;</span>
      <span>{{ rows?.length || 0 }}&nbsp;targets</span>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import {FormTarget} from "#components";

const {status, data: targets, refresh} = useTargets()
const slideoverOpen = useSlideover()?.isOpen

const columns: any[] = [
  {key: 'id', label: '#', sortable: true},
  {key: 'excerpt', label: 'Excerpt', sortable: true},
  {key: 'driver', label: 'Driver', sortable: true},
  {key: 'actions'}
]

const rows = computed(() => targets.value?.map(target => ({
  id: target.id,
  excerpt: target.excerpt,
  driver: target.driver,
})) || [])

// Form

const openCreateForm = () => {
  useSlideover().open(FormTarget, {mode: "create", onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openShowForm = (id: number) => {
  const target = targets.value?.find(s => s.id === id)
  useSlideover().open(FormTarget, {mode: "read", target, onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openEditForm = (id: number) => {
  const target = targets.value?.find(s => s.id === id)
  useSlideover().open(FormTarget, {mode: "edit", target, onCreated: () => refresh(), onUpdated: () => refresh()})
}
const openCloneForm = (id: number) => {
  const target = targets.value?.find(s => s.id === id)
  useSlideover().open(FormTarget, {
    mode: "create",
    target: reactive({...JSON.parse(JSON.stringify(target))}),
    onCreated: () => refresh(),
    onUpdated: () => refresh()
  })
}

const remove = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/targets/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "error", title: "Unable to delete target", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "success", title: "Successfully deleted target", description: serverData.message})
  }
  await refresh()
}
</script>