<template>
  <USlideover v-model="opened" class="overflow-y-auto "
              @keydown.ctrl.enter.prevent="submit"
              :ui="{body:' sm:p-0 p-0',footer:'sm:p-0 p-0',...state.showSamplePanel
                ?(mode==='read'?{content: 'w-screen max-w-[40vw]'}:{content: 'w-screen max-w-[70vw]'})
                :(mode==='read'?{content: 'w-screen max-w-[25vw]'}:{content: 'w-screen max-w-[55vw]'})}">

    <template #header>
      <div class="col-span-full flex justify-between ">
        <h2 @click.prevent="submit" class="col-span-full  text-center text-2xl font-semibold">
          {{
            {
              create: `New processor`,
              edit: `Edit processor: ${processor?.id}`,
              read: `Processor: ${processor?.id}`
            }[mode]
          }}
        </h2>
        <div class="flex gap-x-4">
          <label class="flex items-center gap-x-4 cursor-pointer" v-if="mode!=='read'">
            Short Mode (beta):
            <USwitch
                checked-icon="i-heroicons-check-20-solid"
                unchecked-icon="i-heroicons-x-mark-20-solid"
                v-model="state.shortMode"
            />
          </label>
          <label class="flex items-center gap-x-4 cursor-pointer ">
            Sample Tab:
            <USwitch
                checked-icon="i-heroicons-check-20-solid"
                unchecked-icon="i-heroicons-x-mark-20-solid"
                v-model="state.showSamplePanel"
            />
          </label>
        </div>
      </div>
    </template>
    <template #body>
      <div class="grid h-full"
           :class="{
              'grid-cols-[2fr_3fr]':mode!=='read' && !state.showSamplePanel,
              'grid-cols-[2fr_3fr_15vw]':mode!=='read' && state.showSamplePanel,
              'grid-cols-[1fr_15vw]':state.showSamplePanel && mode=='read'

          }">
        <ProcessorFormLeftPanel v-model:form="form" class="overflow-y-auto"/>
        <ProcessorFormRightPanel v-model:form="form"  class="overflow-y-auto"/>
        <LazyProcessorFormSamplePanel v-model:form="form" v-if="state.showSamplePanel" class="overflow-y-auto"/>


      </div>
    </template>
    <template #footer>
      <!--Full Width save button-->
      <UButton class="rounded-none" @click.prevent="submit" size="xl" block v-if="mode!=='read'">
        {{ {create: 'Create', edit: 'Save changes'}[mode] }}
      </UButton>
    </template>

  </USlideover>
</template>

<script setup lang="ts">

const opened = defineModel('opened')
const $emit = defineEmits(['created', 'updated'])
const state = useProcessFormState()

interface Props {
  mode: "create" | "edit" | "read"
  processor?: any
}

const {mode = "create", processor} = defineProps<Props>()

const defaultForm = reactive({
  source: null as number | null,
  table: null as string | null,
  primary_keys: [] as string[] | null,
  conditions: [] as any[] | null,

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
      useToast().add({title: 'Success', color: 'success', description: serverData.message})
      $emit('created')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error.value?.message, color: 'error'})
    }
  } else if (mode === "edit") {
    const {status, error, data} = await useGoFetch(`/processors/${processor.id}`, {
      method: 'put',
      body: form,
      watch: false
    })
    if (status.value === 'success') {
      const serverData = data.value as { message?: string }
      useToast().add({title: 'Success', color: 'success', description: serverData.message})
      $emit('updated')
      opened.value = false
      return
    }

    if (status.value === 'error') {
      useToast().add({title: 'Error', description: error?.value?.data?.error || error.value?.message, color: 'error'})
    }
  }
}
</script>