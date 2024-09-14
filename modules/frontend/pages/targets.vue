<template>
  <section class="flex flex-col flex-1">
    <UCard class="w-full" :ui="{
      divide: 'divide-y divide-gray-200 dark:divide-gray-700',
      body: { padding: '' },
      header:{base:'flex gap-x-2 justify-between items-center'},
      footer:{base:'text-sm leading-5 text-center'}
    }">
      <template #header>
        <h1 class="font-semibold text-xl text-gray-900 dark:text-white leading-tight">Targets</h1>
        <UButton @click.prevent="create" variant="soft">+ Add Target</UButton>
      </template>
      <template #default>
        <UTable @select="(row:any)=>show(row.id)" :columns="columns" :rows="rows"
                :sort="{column:'id',direction:'desc'}"
                :loading="status==='pending' || status==='idle'">
          <template #driver-data="{ row }">
            <UBadge size="xs" :label="row.driver" color="sky" variant="subtle"/>
          </template>
          <template #actions-data="{ row }">
            <div class="flex gap-1">
              <UDropdown :items="[[{label:'Replicate',click:()=>clone(row.id)}]]" @click.stop>
                <UButton icon="i-heroicons-ellipsis-horizontal" variant="link" color="gray" size="xl" :padded="false"/>
              </UDropdown>
              <UButton icon="i-heroicons-eye" variant="link" color="gray" size="xl" :padded="false"
                       @click.stop.prevent="show(row.id)"/>
              <UButton icon="i-heroicons-pencil-square" variant="link" color="gray" size="xl" :padded="false"
                       @click.stop.prevent="edit(row.id)"/>
              <UButton icon="i-heroicons-trash" variant="link" color="red" size="xl" :padded="false"
                       @click.stop.prevent="remove(row.id)"/>
            </div>
          </template>
        </UTable>
      </template>
      <template #footer>
        <span class="font-medium">Total:&nbsp;</span>
        <span>{{ rows?.length || 0 }}&nbsp;targets</span>
      </template>
    </UCard>
    <TargetForm
        @updated="refresh"
        @created="refresh"
        :mode="formMode"
        v-model:opened="formOpened"
        :target="formTarget"
    />
  </section>
</template>

<script setup lang="ts">
const {status, data: targets, refresh} = useTargets()

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

const formOpened = ref(false);
const formMode = ref<"create" | "edit" | "read">("create");
const formTarget = ref<any>();

const create = () => {
  formMode.value = "create"
  formTarget.value = undefined
  formOpened.value = true
}
const show = (id: number) => {
  formMode.value = "read"
  formTarget.value = targets.value?.filter(i => i.id === id)?.at(0)
  formOpened.value = true
}
const edit = (id: number) => {
  formMode.value = "edit"
  formTarget.value = targets.value?.filter(i => i.id === id)?.at(0)
  formOpened.value = true
}
const clone = (id: number) => {
  formMode.value = "create"
  formTarget.value = {...targets.value?.filter(i => i.id === id)?.at(0)}
  formOpened.value = true
}
const remove = async (id: number) => {
  const {data, error, status} = await useGoFetch(`/targets/${id}`, {method: "DELETE"})
  if (status.value === "error") {
    useToast().add({color: "red", title: "Unable to delete target", description: error.value?.message})
  } else if (status.value === "success") {
    const serverData = data.value as { message?: string }
    useToast().add({color: "green", title: "Successfully deleted target", description: serverData.message})
  }
  await refresh()
}
</script>