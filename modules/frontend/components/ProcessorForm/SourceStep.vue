<template>
  <section class="grid gap-y-4">
    <UFormGroup label="Source" required name="source">
      <div class="grid grid-cols-2 gap-4">
        <div v-for="source in sources||[]" :key="`source-${source.id}`">
          <input
              type="radio"
              name="source"
              v-model="form.source"
              :value="source.id"
              class="sr-only peer"
              :id="`source-${source.id}`"
              tabindex="-1"
          />
          <label
              tabindex="0"
              @keydown.enter.space.prevent="form.driver=source.id"
              :for="`source-${source.id}`"
              class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
          >
            <p class="font-semibold">Source n°{{ source.id }}</p>
            <span class="italic text-gray-400" v-if="source.excerpt">{{ source.excerpt }}</span>
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

    <UDivider label="Conditions"/>

    <UCard v-for="(condition,index) in form.conditions"
           :ui="{body:{base:'grid grid-cols-3 gap-x-4'},header:{base:'flex justify-between items-center'}}">
      <template #header>
        Condition {{ index }}
        <UButton color="red" size="xs" icon="i-heroicons-trash" @click.prevent="removeCondition(index)"/>
      </template>
      <UFormGroup label="Column" required :name="`conditions.${index}.column`" v-if="availableTables">
        <USelectMenu searchable creatable v-model="form.conditions[index].column"
                     :options="availableTables?.[form.table]?.columns || []" placeholder="Select a column"/>
      </UFormGroup>

      <UFormGroup label="Operator" required :name="`conditions.${index}.operator`"
                  v-if="form.conditions[index].column||null">
        <USelectMenu
            v-model="form.conditions[index].operator"
            :options="[
                {label:'IS',value:'='},
                {label:'IS NULL',value:'is null'},
                {label:'IS NOT NULL',value:'is not null'},
                {label:'IS TRUE',value:'is true'},
                {label:'IS FALSE',value:'is false'},
            ]"
            @change="form.conditions[index].value=undefined"
            value-attribute="value"
            option-attribute="label"
            placeholder="Select an operator"/>
      </UFormGroup>

      <UFormGroup label="Value" required :name="`conditions.${index}.value`"
                  v-if="form.conditions[index].operator==='='">
        <UInput type="text" v-model="form.conditions[index].value" placeholder="Select an operator"/>
      </UFormGroup>
    </UCard>

    <UButton class="text-center" block @click.prevent="form.conditions.push({})">+ Add condition</UButton>


    <SourceForm v-model:opened="createSourceOpened" @created="refresh" mode="create"/>
  </section>
</template>
<script setup lang="ts">
import {type Stats} from "~/composables/useSourceStats";

const form = defineModel<any>('form', {required: true})
const {data: sources, refresh} = useSources()

const availableTables = ref<Stats>({});


const computeAvailableTables = async () => {
  if (form.value.source === null) return;

  const {data, status, error} = await useSourceStats(form.value.source);

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
onBeforeMount(() => {
  form.value.conditions ||= []
})

function removeCondition(index: number) {
  form.value.conditions = form.value.conditions?.filter((_, key) => key !== index)

}
</script>