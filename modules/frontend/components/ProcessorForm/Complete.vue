<template>
  <section class="flex flex-col flex-1 relative">

    <USlideover v-model="opened" class="overflow-y-auto"
                @keydown.ctrl.enter.prevent="submit"
                :ui="state.showSamplePanel
                ?(mode==='read'?{width: 'w-screen max-w-[40vw]'}:{width: 'w-screen max-w-[70vw]'})
                :(mode==='read'?{width: 'w-screen max-w-[25vw]'}:{width: 'w-screen max-w-[55vw]'})">

      <div class="grid min-h-screen grid-rows-[auto_1fr_auto]"
           :class="{
              'grid-cols-[2fr_3fr]':mode!=='read' && !state.showSamplePanel,
              'grid-cols-[2fr_3fr_15vw]':mode!=='read' && state.showSamplePanel,
              'grid-cols-[1fr_15vw]':state.showSamplePanel && mode=='read'

          }">

        <!--Toolbar-->
        <div class="col-span-full flex justify-between border-b p-4">
          <h2 @click.prevent="submit" class="col-span-full  text-center text-2xl font-semibold">
            {{ {create: `New processor`, edit: `Edit processor: TODO id`, read: 'Processor: TODO id'}[mode] }}
          </h2>
          <div class="flex gap-x-4">
            <label class="flex items-center gap-x-4 cursor-pointer" v-if="mode!=='read'">
              Short Mode (beta):
              <UToggle
                  on-icon="i-heroicons-check-20-solid"
                  off-icon="i-heroicons-x-mark-20-solid"
                  v-model="state.shortMode"
              />
            </label>
            <label class="flex items-center gap-x-4 cursor-pointer ">
              Sample Tab:
              <UToggle
                  on-icon="i-heroicons-check-20-solid"
                  off-icon="i-heroicons-x-mark-20-solid"
                  v-model="state.showSamplePanel"
              />
            </label>
          </div>

        </div>


        <ProcessorFormLeftPanel v-model:form="form"/>
        <LazyProcessorFormRightPanel v-model:form="form" v-if="mode!=='read'"/>
        <LazyProcessorFormSamplePanel v-model:form="form" v-if="state.showSamplePanel"/>

        <!--Full Width save button-->
        <UButton class="col-span-full" @click.prevent="submit" size="xl" block :ui="{rounded:''}" v-if="mode!=='read'">
          {{ {create: 'Create', edit: 'Save changes'}[mode] }}
        </UButton>

      </div>
    </USlideover>
  </section>
</template>

<script setup lang="ts">


const opened = defineModel('opened')
const $emit = defineEmits(['created', 'updated'])
const state = useProcessFormState()

interface Props {
  mode: "create" | "edit" | "read"
  processor?: any
  processorId?: number
}

const {mode = "create", processor, processorId} = defineProps<Props>()

const defaultForm = reactive({
  source: null as number | null,
  table: null as string | null,
  primary_keys: [] as string[] | null,

  targets: [] as string[],
  index: null as string | null,
  mapping: [] as string[]
})

const form = computed(() => processor || defaultForm)

const submit = async () => {
  if (mode === "create") {
    const {status, error, data} = await useGoFetch("/processors", {method: 'post', body: form, watch: false})
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'green', description: serverData.message})
      $emit('created')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.message, color: 'red'})
    }
  } else if (mode === "edit") {
    const {status, error, data} = await useGoFetch(`/processors/${processorId}`, {
      method: 'put',
      body: form,
      watch: false
    })
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'green', description: serverData.message})
      $emit('updated')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error?.value?.data?.error || error.value?.message, color: 'red'})
    }
  }
}
</script>