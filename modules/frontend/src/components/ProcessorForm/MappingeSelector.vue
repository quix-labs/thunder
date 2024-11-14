<template>
  <div class="grid space-y-4">

    <div class="targettable grid space-y-4" :id="!path?'mapping-':undefined"
         @click.stop="!path?setActivePath(''):void(0)"
         @keydown.enter.space.stop="!path?setActivePath(''):void(0)"
         :class="{'targetted':!path?getPath('')===state.activeMappingPath:false}">

      <template v-for="(mapped,index) in mapping">
        <UCard
            v-if="mapped._type==='simple'"
            :ui="{body:'flex justify-between items-center gap-x-2 sm:p-4'}"
            class="targettable"
            :class="{'targetted':getPath(mapped.column)===state.activeMappingPath}"
            @click.stop="setActivePath(mapped.column)"
            @keydown.enter.space.stop="setActivePath(mapped.column)"
            :id="`mapping-${getPath(mapped.column)}`">
          <p>Simple: {{ mapped.column }}</p>
          <UButton color="error" variant="ghost" size="md" icon="i-heroicons-trash"
                   @click.prevent="removeMapped(index)"/>
        </UCard>

        <UCard
            v-if="mapped._type==='relation'"
            :ui="{body:'sm:p-4',header:'flex justify-between items-center gap-x-2 sm:p-4'}"
            class="targettable"
            :class="{'targetted':getPath(mapped.name)===state.activeMappingPath}"
            @click.stop="setActivePath(mapped.name)"
            @keydown.enter.space.stop="setActivePath(mapped.name)"
            :id="`mapping-${getPath(mapped.name)}`">
          <template #header>
            <p>Relation: {{ mapped.name || '---' }}</p>
            <UButton color="error" variant="ghost" size="md" icon="i-heroicons-trash"
                     @click.prevent="removeMapped(index)"/>
          </template>

          <div class="space-y-4">
            <UFormField name="name" label="Name" required>
              <UInput v-model="mapped.name"  class="w-full" />
            </UFormField>
            <UFormField name="type" label="Type" required>
              <USelect v-model="mapped.type" :items="['one-to-one','has-many']" class="w-full" required/>
            </UFormField>

            <USeparator label="Configuration"/>
            <UFormField name="table" label="Foreign Table" required>
              <USelectMenu v-model="mapped.table" :items="Object.keys(stats||{})"  class="w-full" />
            </UFormField>


            <UFormField name="local_key" label="Local key" required>
              <div class="flex">
                <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                  {{ table }}.
                </UKbd>
                <USelectMenu
                    v-model="mapped.local_key"
                    class="w-full rounded-md rounded-l-none"
                    :items="stats?.[table]?.columns || []"
                />
              </div>
            </UFormField>

            <template v-if="mapped.type==='has-many'">
              <USeparator label="Pivot Configuration"/>
              <UCheckbox label="Use pivot table" v-model="mapped.use_pivot_table"/>
              <template v-if="mapped.use_pivot_table">
                <UFormField name="pivot_table" label="Pivot Table" required>
                  <USelectMenu v-model="mapped.pivot_table" :items="Object.keys(stats||{})" class="w-full" />
                </UFormField>

                <UFormField name="local_pivot_key" label="Local Pivot key" required>
                  <div class="flex">
                    <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                      {{ mapped.pivot_table }}.
                    </UKbd>
                    <USelectMenu
                        v-model="mapped.local_pivot_key"
                        class="w-full rounded-md rounded-l-none"
                        :items="stats?.[mapped.pivot_table]?.columns || []"
                    />
                  </div>
                </UFormField>

                <UFormField name="foreign_pivot_key" label="Foreign Pivot key" required>
                  <div class="flex">
                    <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                      {{ mapped.pivot_table }}.
                    </UKbd>
                    <USelectMenu
                        v-model="mapped.foreign_pivot_key"
                        class="w-full rounded-md rounded-l-none"
                        :items="stats?.[mapped.pivot_table]?.columns || []"
                    />
                  </div>
                </UFormField>
              </template>
              <USeparator label="Pivot Configuration End"/>
            </template>


            <UFormField name="foreign_key" label="Foreign key" required>
              <div class="flex">
                <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                  {{ mapped.table }}.
                </UKbd>
                <USelectMenu
                    v-model="mapped.foreign_key"
                    class="w-full rounded-md rounded-l-none"
                    :items="stats?.[mapped.table]?.columns || []"
                />
              </div>
            </UFormField>

            <UFormField label="Primary keys" required name="primary_keys" v-if="mapped.table">
              <UInputMenu
                  class="w-full"
                  :ui="{tagsInput:'flex-1'}"
                  multiple
                  v-model="mapped.primary_keys"
                  :items="[...new Set([...(stats?.[mapped.table]?.columns || []),...(mapped.primary_keys||[])])]"
              >
                <template #empty>No columns available for the selected table.</template>
              </UInputMenu>
            </UFormField>

          </div>

          <template #footer v-if="mapped.table">
            <MappingeSelector
                :source :table="mapped.table" v-model="mapped.mapping"
                :path="getPath(mapped.name)"
            />
          </template>
        </UCard>
      </template>
    </div>


    <div class="flex items-start justify-start mt-4 gap-2 max-w-full overflow-hidden">
        <UInputMenu
            multiple
            :ui="{tagsInput:'flex-1'}"
            class="flex-1 overflow-hidden"
            :items="availableSimpleFields"
            v-model="tempSimpleField"
            placeholder="Select a field to add"
        />
        <UButton
            trailing-icon="i-heroicons-plus"
            class="mr-4"
            @click.prevent.stop="addSimpleFields"
            :disabled="!tempSimpleField"
        >Add</UButton>


      <UButton trailing-icon="i-heroicons-plus" @click.prevent.stop="addAllFields">
        Add all fields
      </UButton>


      <UButton trailing-icon="i-heroicons-plus" @click.prevent.stop="addRelation">
        Add Relation
      </UButton>

    </div>
  </div>
</template>

<script setup lang="ts">

const mapping = defineModel<any>({required: false, default: []})
const props = defineProps<{
  source: string,
  table: string,
  path?: string,
}>()

const {data: stats, status, error} = await useSourceStats(props.source);
//TODO ERROR HANDLING

const manualSimpleFields = ref<string[]>([])
const tempSimpleSelected = ref<string[]>([])

const tempSimpleField = computed<string[]>({
  get: () => tempSimpleSelected.value,
  set: async (labels) => {
    for (const field of labels) {
      if ([...stats.value?.[props.table]?.columns || [], ...manualSimpleFields.value || []]?.includes(field)) {
        continue
      }
      manualSimpleFields.value.push(field?.label || field)
    }

    tempSimpleSelected.value = labels.map(i => i?.label || i).filter(i => !mapping.value?.includes(i))
  }
});

const availableSimpleFields = computed(() => [
  ...stats.value?.[props.table]?.columns || [],
  ...manualSimpleFields.value || []
].filter(col => !mapping.value.find(mapped => mapped?._type === 'simple' && mapped?.column == col)))

const addSimpleFields = () => {
  if (tempSimpleField.value?.length == 0) return
  mapping.value = [...mapping.value || [], ...tempSimpleField.value?.map(i => ({_type: 'simple', column: i}))]
  tempSimpleField.value = []
};

const addAllFields = () => {
  mapping.value = [...mapping.value || [], ...availableSimpleFields.value?.map(i => ({_type: 'simple', column: i}))]

}
const addRelation = () => {
  mapping.value.push({
    _type: 'relation',
    mapping: [],
    type: null,
    name: null,
    table: null,

    foreign_key: null,
    local_key: null,

  })
};

const removeMapped = (index: number) => {
  mapping.value = mapping.value.filter((_, key) => key !== index)
};


const state = useProcessFormState();


const getPath = (fieldname: string): string => {
  return [props.path, fieldname].filter(i => i).join('.')
}

const setActivePath = async (path: string) => {
  state.value.preventScroll = true;
  await nextTick(() => {
    state.value.activeTabs = "mapping"
    state.value.activeMappingPath = getPath(path)
  });
  await nextTick(() => state.value.preventScroll = false);
}

</script>

<style scoped lang="postcss">
.targettable {
  scroll-margin-top: 1em;
}

.targettable.targetted {
  outline: 1px solid var(--ui-color-primary-500);
}
</style>