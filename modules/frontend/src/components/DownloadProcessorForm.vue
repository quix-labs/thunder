<template>
  <UModal :title="`Download data for processor n°${processorId}`">
    <template #body>
      <UForm :state="form" @submit="submit" class="space-y-4">
        <UFormField label="Exporter" name="exporter" required>
          <USelect v-model="form.exporter"
                   :items="Object.entries(exporters||{})?.map(([ID,label])=>({label,value:ID,description:'test'}))"
                   required/>
        </UFormField>
        <UFormField label="Filename" name="filename">
          <UInput :placeholder="defaultFilename" v-model="form.filename"/>
        </UFormField>
        <UButton type="submit" class="p-3 cursor-pointer" block loading-auto>
          Generate file
        </UButton>
      </UForm>
    </template>
  </UModal>
</template>
<script lang="ts" setup>

import {UForm} from "#components";


interface Props {
  processorId: number
}

const {data: exporters, status, error} = await useExporters({lazy: false})

const props = defineProps<Props>()
const form = reactive({
  exporter: (status.value === "success" && exporters.value) ? Object.keys(exporters.value)?.at(0) : null as string | null,
  filename: null as string | null
})
const defaultFilename = computed(() => `processor-${props.processorId}.${form.exporter?.replace('thunder.', '')}`)

const submit = async () => {
  const data = {...toRaw(form)}
  data.filename ||= defaultFilename.value

  const {data: blob, error, status} = await useGoFetch<Blob>(`processors/${props.processorId}/download`, {
    params: data,
    responseType: "blob",
  })

  useModal()?.close()

  if (status.value === "success" && blob.value) {
    const download = () => {
      const url = window.URL.createObjectURL(blob.value)
      const a = document.createElement('a');
      a.href = url;
      a.download = data.filename as string;
      a.click();
      setTimeout(() => {
        window.URL.revokeObjectURL(url);
      }, 0)
    }

    download()

    useToast().add({
      color: 'success',
      title: 'File Ready for Download',
      close: false,
      actions: [{
        icon: 'i-heroicons-arrow-down-tray',
        label: 'Download file',
        color: 'primary',
        variant: 'outline',
        onClick: download,
      }],
    });
  } else if (error.value) {
    if (error.value.data instanceof Blob) {
      const reader = new FileReader();
      reader.onload = (event) => {
        try {
          const json = JSON.parse((event.target?.result || "null") as string);
          useToast().add({
            title: 'Failed to generate file',
            description: json.error || 'An unknown error occurred.',
            color: 'error',
          });
        } catch (parseError) {
          useToast().add({
            title: 'Failed to generate file',
            description: 'Error parsing the error response.',
            color: 'error',
          });
        }
      };

      reader.readAsText(error.value.data);
    } else {
      useToast().add({
        title: 'Failed to generate file',
        description: error.value.message,
        color: 'error',
      });
    }
  }
}
</script>
