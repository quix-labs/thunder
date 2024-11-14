<template>
  <section class="grid gap-y-4">
    <UFormField label="Source" required name="source">
      <FormCardsInput :options="sources?.map(source=>({value:source.id,item:source}))||[]" v-model="form.source">
        <template #default="{item}">
          <div class="text-center">
            <p class="font-semibold">Source n°{{ item.id }}</p>
            <span class="italic text-gray-400" v-if="item.excerpt">{{ item.excerpt }}</span>
          </div>
        </template>

        <template #create>
          <div
              class="col-span-full flex justify-center items-center gap-x-1 cursor-pointer  rounded-lg p-4 text-lg text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800"
              @click="()=>openCreateForm()">
            <UIcon name="heroicons-plus"/>
            <span>Create</span>
          </div>

        </template>
      </FormCardsInput>
    </UFormField>

    <UFormField label="Table" required name="table">
      <UInputMenu
          v-model="form.table"
          :items="Object.keys(availableTables||{})"
          placeholder="Select a table"
          :loading="statsStatus==='idle'||statsStatus==='pending'"
      >
        <template #empty>No tables available for the selected source.</template>
      </UInputMenu>
    </UFormField>

    <UFormField label="Primary keys" required name="primaryKeys" v-if="availableTables">
      <UInputMenu
          multiple
          v-model="form.primary_keys"
          :loading="statsStatus==='idle'||statsStatus==='pending'"
          :items="[...new Set([...availableTables?.[form.table]?.columns || [],...form.primary_keys||[]])]"
      >
        <template #empty>No columns available for the selected table.</template>
      </UInputMenu>
    </UFormField>

    <USeparator label="Conditions"/>

    <UCard v-for="(condition,index) in form.conditions"
           :key="'conditions:'+index"
           :ui="{body:'grid grid-cols-3 gap-x-4',header:'flex justify-between items-center p-2 sm:p-2 font-bold'}">
      <template #header>
        Condition {{ index }}
        <UButton color="error" class="cursor-pointer" variant="ghost" size="sm" icon="i-heroicons-trash" @click.prevent="removeCondition(index)"/>
      </template>

      <UFormField label="Column" required :name="`conditions.${index}.column`" v-if="availableTables">
        <UInputMenu
            :loading="statsStatus==='idle'||statsStatus==='pending'"
            class="w-full"
            searchable
            v-model="form.conditions[index].column"
            :items="availableTables?.[form.table]?.columns || []" placeholder="Select a column"/>
      </UFormField>

      <UFormField label="Operator" required :name="`conditions.${index}.operator`"
                  v-if="form.conditions[index].column||null">
        <USelectMenu
            v-model="form.conditions[index].operator"
            class="w-full"
            :items="[
                    {label:'IS',value:'='},
                    {label:'IS NULL',value:'is null'},
                    {label:'IS NOT NULL',value:'is not null'},
                    {label:'IS TRUE',value:'is true'},
                    {label:'IS FALSE',value:'is false'},
                ]"
            value-key="value"
            placeholder="Select an operator"/>
      </UFormField>
      <UFormField label="Value" required :name="`conditions.${index}.value`"
                  v-if="form.conditions[index].operator==='='">
        <UInput type="text" v-model="form.conditions[index].value" class="w-full" placeholder="Select an operator"/>
      </UFormField>
    </UCard>

    <UButton class="text-center" block @click.prevent="form.conditions.push({})">+ Add condition</UButton>

    <FormSource v-model:open="formOpened" mode="create" @created="()=>refresh()"/>
  </section>
</template>
<script setup lang="ts">
import {FormSource} from "#components";

const form = defineModel<any>('form', {required: true})
const {data: sources, refresh} = useSources()

const {data: availableTables, status: statsStatus} = useSourceStats(toRef(form.value, 'source'))

const formOpened = ref(false)
const openCreateForm = () => {
  formOpened.value = true
}
onBeforeMount(() => {
  form.value.conditions ||= []
})

function removeCondition(index: number) {
  form.value.conditions = form.value.conditions?.filter((_, key) => key !== index)
}
</script>