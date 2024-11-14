<template>
  <div class="flex flex-col flex-1 gap-y-4 p-4">
    <!-- SOURCE -->
    <UCard
        tabindex="0"
        class="cursor-pointer"
        :class="state.activeTabs==='source' ? 'ring-1 ring-[var(--ui-primary)]' : undefined"
        @click.prevent="selectTab('source')"
        @keydown.enter.space="selectTab('source')"
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
      <p>Conditions: {{ form.conditions?.length || '---' }}</p>

    </UCard>

    <!-- MAPPING -->
    <UCard
        tabindex="0"
        class="flex-1 cursor-pointer"
        :class="state.activeTabs==='mapping' ? 'ring-1 ring-[var(--ui-primary)]' : undefined"
        @click.prevent="selectTab('mapping')"
        @keydown.enter.space="selectTab('mapping')"
    >
      <template #header>
        <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Mapping
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Describe how to transform your relational data into unstructured documents
        </p>
      </template>
      <ProcessorFormMappingThree :mapping="form.mapping" :name="form.table" :as="form.index"/>
    </UCard>

    <!-- OUTPUT -->
    <UCard
        tabindex="0"
        class="cursor-pointer"
        :class="state.activeTabs==='output' ? 'ring-1 ring-[var(--ui-primary)]' : undefined"
        @click.prevent="selectTab('output')"
        @keydown.enter.space="selectTab('output')"
    >
      <template #header>
        <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Output
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Describe how your data will be synchronized
        </p>
      </template>
      <p class="font-semibold">
        Targets: <span class="font-normal">{{form.targets?.map(i=>`Target n°${i}`)?.join(', ')}}</span>
      </p>
      <p class="font-semibold">
        Expected index: <span class="font-normal">{{ form.index }}</span>
      </p>

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
