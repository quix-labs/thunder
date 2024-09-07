<template>
  <section>
    <UFormGroup label="Index" required name="index">
      <UInput v-model="form.index" placeholder="Enter index"/>
    </UFormGroup>
  </section>
</template>
<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
onBeforeMount(()=>{
  form.value.index ||= form.value?.table
})
</script>