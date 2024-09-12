<template>
  <section>
    <UFormGroup label="Source" required name="source">
      <div class="grid grid-cols-2 gap-4">
        <div v-for="(source, key) in sources" :key="key">
          <input
              type="radio"
              name="source"
              v-model="form.source"
              :value="key"
              class="sr-only peer"
              :id="`source-${key}`"
              tabindex="-1"
          />
          <label
              tabindex="0"
              @keydown.enter.space.prevent="form.driver=key"
              :for="`source-${key}`"
              class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
          >
            Source n°{{ key }} (TODO user defined name)
          </label>
        </div>

        <div
            class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
            @click.prevent="createSourceOpened=true"
        >
          + CREATE SOURCE
        </div>

      </div>
    </UFormGroup>
    <UFormGroup label="Table" required name="table" v-if="availableTables">
      <USelectMenu
          creatable
          searchable
          v-model="form.table"
          :options="Object.keys(availableTables)"
          placeholder="Select a table"
      >
        <template #option-create="{ option }">
          Force manual: {{ option }}
        </template>
        <template #empty>No tables available for the selected source.</template>
      </USelectMenu>
    </UFormGroup>

    <UFormGroup label="Primary keys" required name="primaryKeys" v-if="availableTables">
      <USelectMenu
          searchable
          multiple
          v-model="form.primary_keys"
          :options="[...new Set([...availableTables?.[form.table]?.columns || [],...form.primary_keys||[]])]"
      >
        <template #label>
          <span v-if="form.primary_keys?.length>0" class="truncate">{{ form.primary_keys.join(', ') }}</span>
          <span v-else>Select primary keys</span>
        </template>
        <template #empty>No columns available for the selected table.</template>
      </USelectMenu>
    </UFormGroup>
    <CreateSourceDriverForm v-model:opened="createSourceOpened" @created="refresh"/>
  </section>


</template>
<script setup lang="ts">

const form = defineModel<any>('form', {required: true})
const {data: sources, refresh} = useSources()

type Stats = { [key: string]: { columns: string[], primary_keys: string[] } };
const availableTables = ref<Stats>({});


const computeAvailableTables = async () => {
  if (form.value.source === null) return;

  const {data, status, error} = await useGoFetch<Stats>(`/sources/${form.value.source}/stats`);

  availableTables.value = data.value || {};

  if (status.value === 'error') {
    const errorMessage = error.value?.statusCode === 422
        ? error.value.data?.error
        : error.value?.message;

    useToast().add({title: 'Error fetching stats', description: errorMessage, color: 'red'});
    return;
  }

  if (!data.value) {
    useToast().add({title: 'No Data', description: 'Empty stats received', color: 'orange'});
    return;
  }
}
const createSourceOpened = ref(false);

watch(() => form.value.source, computeAvailableTables);
onMounted(computeAvailableTables)
</script>