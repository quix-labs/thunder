<template>
  <UCard as="section" id="targets"
         :ui="{header:'flex gap-x-2 justify-between items-center',footer:'text-sm leading-5 text-center'}">
    <template #header>
      <h1>Targets</h1>
      <div class="flex space-x-4">
        <KbdButton :kbds="['r']" label="Refresh" color="info" variant="soft" icon="i-heroicons-arrow-path-20-solid"
                   @click="()=>refresh()" :loading="status==='pending' || status==='idle'"/>
        <KbdButton :kbds="['c']" label="Create" variant="soft" leading-icon="i-heroicons-plus" :disabled="slideoverOpen"
                   @click="()=>openCreateForm()"/>
      </div>
    </template>


    <template #footer>
      <span class="font-medium">Total:&nbsp;</span>
      <span>{{ rows?.length || 0 }}&nbsp;targets</span>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import {FormTarget} from "#components";
import KbdButton from "~/components/KbdButton.vue";

const {status, data: targets, refresh} = useTargets()
const slideoverOpen = useSlideover()?.isOpen

const columns = [
  {key: 'id', label: '#', sortable: true, rowClass: 'w-[1px] whitespace-nowrap'},
  {key: 'excerpt', label: 'Excerpt', sortable: true},
  {key: 'driver', label: 'Driver', sortable: true},
  {key: 'actions', sortable: false, rowClass: 'w-[1px] whitespace-nowrap'}
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
    target: {...toRaw(target)},
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