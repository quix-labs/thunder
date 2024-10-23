<template>
  <section class="grid gap-y-4">
    <UFormField label="Targets" required name="targets">
      <FormCardsInput multiple :options="targets?.map(target=>({value:target.id,item:target}))||[]" v-model="form.targets">
        <template #default="{item}">
          <div class="text-center">
            <p class="font-semibold">Target n°{{ item.id }}</p>
            <span class="italic text-gray-400" v-if="item.excerpt">{{ item.excerpt }}</span>
          </div>
          <!--TODO CREATE TARGET-->
        </template>
      </FormCardsInput>
    </UFormField>

    <UFormField label="Index" required name="index">
      <UInput v-model="form.index" placeholder="Enter index"/>
    </UFormField>
  </section>
</template>
<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
const {data: targets, refresh} = useTargets()

onBeforeMount(() => {
  form.value.index ||= form.value?.table
  form.value.targets ||= []
})
</script>