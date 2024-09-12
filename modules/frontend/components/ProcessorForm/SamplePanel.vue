<template>
  <div class="flex flex-col flex-1 gap-y-4 max-h-screen overflow-auto p-4">
    <UCard class="flex flex-col flex-1 overflow-x-hidden overflow-y-auto" :ui="{
      divide:'',
      body: { base: 'flex-1' },
      header:{base:'sticky z-20 top-0 shadow-md bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800'},
      footer:{base:'sticky z-20 bottom-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800'}
    }">
      <template #header>
        <div class="flex items-center justify-between">
          <p class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            Document Sample
          </p>
          <UButton color="sky" :loading="status==='pending'" @click.prevent="refresh">
            Fetch Sample
          </UButton>
        </div>
      </template>

      <div v-if="status==='idle'">
        Please click the "Fetch Sample" button to generate the document.
      </div>
      <div v-else-if="status==='pending'" class="flex items-center gap-x-2 text-lg">
        Generating document
        <UIcon name="i-heroicons-arrow-path" class="animate-spin"/>
      </div>

      <div v-else-if="status==='success'">
        <p class="mx-4 mb-4 ">
          <span class="font-semibold">Document primary keys:</span>
          [{{ data?.primary_keys?.join(', ') }}]
        </p>
        <VueJsonPretty
            show-icon
            show-length
            show-line-number
            :data="data.document" :theme="useColorMode().value"/>
      </div>

      <div v-else-if="status==='error'" class=" text-red-600 flex flex-col gap-y-2 items-center ">
        <p class="font-semibold  w-full">{{ error?.data?.error }}</p>
        <p class="w-full">There was an error retrieving the document:</p>
        <p>If the problem persists, please verify that the data mapping might be incorrect or that the server is
          responding as expected.</p>
        <p>You may need to check the response format or any related configuration.</p>
        <p class="font-semibold text-lg w-full">{{ error?.message }}</p>
        <UButton color="red" @click="refresh">Retry</UButton>
      </div>

      <template #footer>
        <div class="flex items-center justify-between">
          <label class="flex items-center gap-x-4 cursor-pointer">
            <UToggle
                on-icon="i-heroicons-check-20-solid"
                off-icon="i-heroicons-x-mark-20-solid"
                v-model="state.liveReload"
            />
            Live reload
          </label>
          <UButton color="gray" :disabled="!data" @click.prevent.stop="download">
            Download
          </UButton>
        </div>
      </template>
    </UCard>
  </div>
</template>
<script setup lang="ts">
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';

const state = useProcessFormState()
const form = defineModel<any>('form', {required: true})
const {status, error, data, refresh} = await useGoFetch<any>("/processors/test", {
  method: 'post',
  body: form,
  watch: false,
  responseType: 'json',
  immediate: false,
})
const throttledRefresh = throttle(async () => {
  await refresh();
}, 4_000);

watch(form.value, async (current) => {
  if (!state.value.liveReload) return
  await throttledRefresh()
})

const download = () => {
  if (!data.value) return
  downloadJSON(data.value, form.value.index + '_sample.json')
}
</script>
