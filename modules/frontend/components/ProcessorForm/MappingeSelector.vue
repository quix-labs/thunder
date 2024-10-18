<template>
  <div class="grid space-y-4">

    <div class="targettable grid space-y-4" :id="!path?'mapping-':undefined"
         @click.stop="!path?setActivePath(''):void(0)"
         @keydown.enter.space.stop="!path?setActivePath(''):void(0)"
         :class="{'targetted':!path?getPath('')===state.activeMappingPath:false}">

      <template v-for="(mapped,index) in mapping">
        <UCard
            :ui="{body:{base:'flex justify-between items-center gap-x-2'}}" v-if="mapped._type==='simple'"
            class="targettable"
            :class="{'targetted':getPath(mapped.column)===state.activeMappingPath}"
            @click.stop="setActivePath(mapped.column)"
            @keydown.enter.space.stop="setActivePath(mapped.column)"
            :id="`mapping-${getPath(mapped.column)}`">
          <p>Simple: {{ mapped.column }}</p>
          <UButton color="red" size="xs" icon="i-heroicons-trash" @click.prevent="removeMapped(index)"/>
        </UCard>

        <UCard
            v-if="mapped._type==='relation'"
            class="targettable"
            :class="{'targetted':getPath(mapped.name)===state.activeMappingPath}"
            @click.stop="setActivePath(mapped.name)"
            @keydown.enter.space.stop="setActivePath(mapped.name)"
            :id="`mapping-${getPath(mapped.name)}`">
          <template #header>
            <div class="flex justify-between items-center">
              <p>Relation: {{ mapped.name || '---' }}</p>
              <UButton color="red" size="xs" icon="i-heroicons-trash" @click.prevent="removeMapped(index)"/>
            </div>

          </template>

          <div class="space-y-2">
            <UFormGroup label="Name" required>
              <UInput v-model="mapped.name"/>
            </UFormGroup>
            <UFormGroup label="Type" required>
              <USelect v-model="mapped.type" :options="['one-to-one','has-many']" required/>
            </UFormGroup>

            <UDivider label="Configuration"/>
            <UFormGroup label="Foreign Table" required>
              <USelectMenu v-model="mapped.table" :options="Object.keys(stats||{})" searchable/>
            </UFormGroup>


            <UFormGroup label="Local key" required>
              <div class="flex">
                <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                  {{ table }}.
                </UKbd>
                <USelectMenu
                    v-model="mapped.local_key"
                    class="w-full"
                    :ui="{rounded:'rounded-md rounded-l-none'}"
                    :options="stats?.[table]?.columns || []"
                    searchable
                />
              </div>
            </UFormGroup>

            <template v-if="mapped.type==='has-many'">
              <UDivider label="Pivot Configuration"/>
              <UCheckbox label="Use pivot table" v-model="mapped.use_pivot_table"/>
              <template v-if="mapped.use_pivot_table">
                <UFormGroup label="Pivot Table" required>
                  <USelectMenu v-model="mapped.pivot_table" :options="Object.keys(stats||{})" searchable/>
                </UFormGroup>

                <UFormGroup label="Local Pivot key" required>
                  <div class="flex">
                    <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                      {{ mapped.pivot_table }}.
                    </UKbd>
                    <USelectMenu
                        v-model="mapped.local_pivot_key"
                        class="w-full"
                        :ui="{rounded:'rounded-md rounded-l-none'}"
                        :options="stats?.[mapped.pivot_table]?.columns || []"
                        searchable
                    />
                  </div>
                </UFormGroup>

                <UFormGroup label="Foreign Pivot key" required>
                  <div class="flex">
                    <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                      {{ mapped.pivot_table }}.
                    </UKbd>
                    <USelectMenu
                        v-model="mapped.foreign_pivot_key"
                        class="w-full"
                        :ui="{rounded:'rounded-md rounded-l-none'}"
                        :options="stats?.[mapped.pivot_table]?.columns || []"
                        searchable
                    />
                  </div>
                </UFormGroup>
              </template>
              <UDivider label="Pivot Configuration End"/>
            </template>


            <UFormGroup label="Foreign key" required>
              <div class="flex">
                <UKbd class="h-auto p-2 pointer-events-none select-none rounded-r-none border-l-0">
                  {{ mapped.table }}.
                </UKbd>
                <USelectMenu
                    v-model="mapped.foreign_key"
                    class="w-full"
                    :ui="{rounded:'rounded-md rounded-l-none'}"
                    :options="stats?.[mapped.table]?.columns || []"
                    searchable
                />
              </div>
            </UFormGroup>

            <UFormGroup label="Primary keys" required name="primaryKeys" v-if="mapped.table">
              <USelectMenu
                  searchable multiple
                  v-model="mapped.primary_keys"
                  :options="[...new Set([...(stats?.[mapped.table]?.columns || []),...(mapped.primary_keys||[])])]"
              >
                <template #label>
                      <span v-if="mapped.primary_keys?.length>0" class="truncate">
                        {{ mapped.primary_keys.join(', ') }}
                      </span>
                  <span v-else>Select primary keys</span>
                </template>
                <template #empty>No columns available for the selected table.</template>
              </USelectMenu>
            </UFormGroup>

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


    <div class="flex items-start justify-start mt-4 gap-8">
      <div class="flex gap-2">
        <USelectMenu
            creatable
            searchable
            multiple
            clear-search-on-close
            :options="availableSimpleFields"
            v-model="tempSimpleField"
            placeholder="Select a field to add"
        />
        <UButton
            trailing-icon="i-heroicons-plus"
            @click.prevent.stop="addSimpleFields"
            :disabled="!tempSimpleField"
        >
          Add Simple Fields
        </UButton>
      </div>


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

import useSourceStats from "~/composables/useSourceStats";

const mapping = defineModel<any>({required: false, default: []})
const props = defineProps<{
  source: number,
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
  @apply outline outline-2 outline-amber-400;
}
</style>