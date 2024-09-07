<template>
  <div class="flex flex-col flex-1 gap-y-4 max-h-screen overflow-auto p-4">

    <!-- SOURCE -->
    <UCard
        tabindex="0"
        class="cursor-pointer"
        @click.prevent="selectTab('source')"
        @keydown.enter.space="selectTab('source')"
        :ui="state.activeTabs==='source' ? {ring:'ring-1 ring-primary dark:ring-primary'} : undefined"
    >
      <template #header>
        <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Source
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Describe which data will be synchronized
        </p>
      </template>


      <p>Source: Source n°{{ form.source }}</p>
      <p>Table: {{ form.table || '---' }}</p>

    </UCard>

    <!-- MAPPING -->
    <UCard
        tabindex="0"
        class="flex-1 cursor-pointer"
        @click.prevent="selectTab('mapping')"
        @keydown.enter.space="selectTab('mapping')"
        :ui="state.activeTabs==='mapping' ? {ring:'ring-1 ring-primary dark:ring-primary'} : undefined"
    >
      <template #header>
        <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Mapping
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Describe how to transform your relational data into unstructured documents
        </p>
      </template>
      <ProcessorFormMappingThree :mapping="form.mapping" :prefix="form.table"/>
    </UCard>

    <!-- OUTPUT -->
    <UCard
        tabindex="0"
        class="cursor-pointer"
        @click.prevent="selectTab('output')"
        @keydown.enter.space="selectTab('output')"
        :ui="state.activeTabs==='output' ? {ring:'ring-1 ring-primary dark:ring-primary'} : undefined"
    >
      <template #header>
        <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Output
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Describe how your data will be synchronized
        </p>
      </template>
      <p>Targets: "TODO Implement target and target drivers"[]</p>
      <p>Index: {{ form.index}}</p>
    </UCard>


  </div>
</template>
<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
const state = useProcessFormState()

function selectTab(tab: typeof state.value.activeTabs) {
  state.value.activeTabs = tab
  // TODO EMIT EVENT FOR RIGHT PANEL
}
</script>
