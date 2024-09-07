<template>
  <USlideover v-model="opened" class="overflow-y-auto"
              :ui="mode==='read'?{width: 'w-screen max-w-[25vw]'}:{width: 'w-screen max-w-[50vw]'}">
    <div class="grid min-h-screen grid-rows-[auto_1fr_auto]" :class="{'grid-cols-2':mode!=='read'}">

      <h2 @click.prevent="submit" class="col-span-full mx-4 mt-4 text-center text-2xl font-semibold">
        {{ {create: `New processor`, edit: `Edit processor: TODO id`, read: 'Processor: TODO id'}[mode] }}
      </h2>

      <ProcessorFormLeftPanel v-model:form="form" :class="{'':mode==='read'}"/>
      <ProcessorFormRightPanel v-model:form="form" v-if="mode!=='read'"/>

      <template v-if="mode!=='read'">
        <!--        <UDivider label="Actions" class="col-span-full "/>-->
        <UButton @click.prevent="submit" size="xl" block class="col-span-full" :ui="{rounded:''}">
          {{ {create: 'Create', edit: 'Save changes'}[mode] }}
        </UButton>
      </template>


    </div>

    <!-- TODO SUBMIT -->
  </USlideover>
</template>

<script setup lang="ts">

const opened = defineModel('opened')
const $emit = defineEmits(['created', 'updated'])

interface Props {
  mode: "create" | "edit" | "read"
  processor?: any
}

const {mode = "create", processor} = defineProps<Props>()

const defaultForm = reactive({
  source: null as number | null,
  table: null as string | null,
  index: null as string | null,
  mapping: [] as string[]
})

const form = computed(() => processor || defaultForm)
const submit = async () => {
  if (mode === "edit") {
    useToast().add({title: 'Error', description: 'Edit endpoint not implemented', color: 'red'})
    return
  }
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
}
</script>