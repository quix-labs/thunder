<template>
  <section class="grid gap-y-4">
    <UFormGroup label="Targets" required name="targets">
      <div class="grid grid-cols-2 gap-4">
        <div v-for="target in targets || []" :key="`target-${target.id}`">
          <input
              type="checkbox"
              name="targets[]"
              v-model.number="form.targets"
              :value="target.id"
              class="sr-only peer"
              :id="`target-${target.id}`"
              tabindex="-1"
          />
          <label
              tabindex="0"
              @keydown.enter.space.prevent="form.driver=target.id"
              :for="`target-${target.id}`"
              class="cursor-pointer flex flex-col gap-y-2 items-center rounded-lg p-4 text-gray-900 dark:text-white bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 peer-checked:ring-2 peer-checked:ring-primary-500"
          >
            <p class="font-semibold">Target n°{{ target.id }}</p>
            <span class="italic text-gray-400" v-if="target.excerpt">{{ target.excerpt }}</span>
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
const {data: targets, refresh} = useTargets()
const createFormOpened = ref(false)


onBeforeMount(() => {
  form.value.index ||= form.value?.table
  form.value.targets ||= []
})
</script>