<template>
  <section class="grid gap-y-4">
    <UFormGroup label="Targets" required name="targets">
      <div class="grid grid-cols-2 gap-4">
        <div v-for="([key,target]) in Object.entries(targets||{})" :key="key">
          <input
              type="checkbox"
              name="targets[]"
              v-model.number="form.targets"
              :value="parseFloat(key)"
              class="sr-only peer"
              :id="`target-${key}`"
              tabindex="-1"
          />
          <label
              tabindex="0"
              @keydown.enter.space.prevent="form.driver=parseFloat(key)"
              :for="`target-${key}`"
              class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
          >
            <p class="font-semibold">Target n°{{ key }}</p>
            <span class="italic text-gray-400" v-if="target.excerpt">{{target.excerpt}}</span>
          </label>
        </div>

        <div
            class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
            @click.prevent="createFormOpened=true"
        >
          + CREATE TARGET
        </div>

      </div>
    </UFormGroup>

    <UFormGroup label="Index" required name="index">
      <UInput v-model="form.index" placeholder="Enter index"/>
    </UFormGroup>

    <TargetForm v-model:opened="createFormOpened" @created="refresh" mode="create"/>
  </section>
</template>
<script setup lang="ts">
const form = defineModel<any>('form', {required: true})
const {data, refresh} = useTargets()
const createFormOpened = ref(false)

const targets = computed(() => {
  const targets = data.value || {};
  if (form.value.targets?.length > 0) {
    form.value.targets.filter((id: number) => targets?.length < id + 1).forEach((id: number) => {
      targets[id] = {} // TODO FAKE INFOS
    })
  }
  return targets;
})


onBeforeMount(() => {
  form.value.index ||= form.value?.table
  form.value.targets ||= []
})
</script>