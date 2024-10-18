<template>
  <UModal ref="modal">
    <UForm @submit.prevent="submit" :state="form" method="GET" target="_blank">
      <UCard :ui="{body:{base:'grid gap-y-4'},footer:{padding:''}}">
        <template #header>
          <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            Download data for processor n°{{ processorId }}
          </p>
        </template>
        <UFormGroup label="Format" name="format" required>
          <USelect v-model="form.format" :options="['csv', 'json']" required/>
        </UFormGroup>
        <UFormGroup label="Filename" name="filename">
          <UInput :placeholder="defaultFilename" v-model="form.filename"/>
        </UFormGroup>
        <template #footer>
          <UButton type="submit" class="p-3 rounded-t-none" block> Download</UButton>
        </template>
      </UCard>
    </UForm>
  </UModal>
</template>
<script lang="ts" setup>

interface Props {
  processorId: number
}

const props = defineProps<Props>()
const modal = useTemplateRef("modal")
const form = reactive({
  format: 'csv' as "csv" | "json",
  filename: null as string | null
})
const defaultFilename = computed(() => `processor-${props.processorId}.${form.format}`)

const submit = async () => {

  const filename = form.filename || defaultFilename.value;

  await navigateTo(getGoApiUrl(`/processors/${props.processorId}/download?format=${form.format}&filename=${filename}`), {
    external: true,
    open: {target: "_self"}
  })
  modal.value?.close()
}
</script>
