<template>
  <div class="flex flex-col flex-1 gap-y-4 p-4">
    <!-- SOURCE -->
    <ProcessorFormSourceStep v-model:form="form" v-if="state.activeTabs==='source'"/>

    <!-- MAPPING -->
    <ProcessorFormMappingStep v-model:form="form" v-if="state.activeTabs==='mapping'"/>

    <!-- OUTPUT -->
    <ProcessorFormOutputStep v-model:form="form" v-if="state.activeTabs==='output'"/>

  </div>
</template>
<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
const state = useProcessFormState()

</script>
