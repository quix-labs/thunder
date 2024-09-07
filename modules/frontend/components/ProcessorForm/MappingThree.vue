<template>
  <div tabindex="0" class="text-md leading-tight">
    <span @click.prevent="path?setActivePath(''):void(0)" class="inline-flex items-center">
      {{ prefix ? `${prefix} {` : '{' }}

      <UKbd class="cursor-pointer mx-1 inline-flex items-center" :class="{'opacity-60':collapsed}" @click.prevent="collapsed=!collapsed">
        <span>{{ mapping.length>=2?`${mapping.length} fields`:`${mapping.length} field` }}</span>
        <UIcon name="i-heroicons-chevron-down" class="ml-1" v-if="!collapsed"/>
        <UIcon name="i-heroicons-chevron-right" class="ml-1" v-else/>
      </UKbd>

    </span>

    <div
        v-for="field in mapping"
        class="pl-8 cursor-pointer hover:underline"
        :class="{
          'border-l ':path,
          'border-sky-600 dark:border-sky-300/60':path && state.activeMappingPath?.includes(path),
          'border-gray-200 dark:border-gray-700':path && !state.activeMappingPath?.includes(path),
        }"
        v-if="!collapsed">

      <template v-if="field._type==='simple'">
        <p
            :class="{'font-bold':getPath(field.column) === (state.activeMappingPath||'')}"
            @click.prevent="setActivePath(field.column)"
            tabindex="0"
        >
          {{ field.column }}
        </p>
      </template>

      <template v-else-if="field._type==='relation'">
        <MappingThree
            :class="{'font-bold':getPath(field.name) === (state.activeMappingPath||'')}"
            :mapping="field.mapping"
            :prefix="field.name"
            :path="getPath(field.name)"
        />
      </template>
    </div>
    <span>}</span>
  </div>
</template>

<script setup lang="ts">

interface Props {
  mapping: any
  prefix?: string
  path?: string
}

const {mapping = [], path} = defineProps<Props>()
const collapsed = ref(false)

const state = useProcessFormState()
const setActivePath = async (path: string) => {
  state.value.activeMappingPath = getPath(path)
}

const getPath = (fieldname: string): string => {
  return [path, fieldname].filter(i => i).join('.')
}
</script>