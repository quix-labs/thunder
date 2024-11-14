<template>
  <div tabindex="0" class="text-sm leading-tight"
       @click.prevent.stop="setActivePath('')"
       @keydown.enter.space.prevent.stop="setActivePath('')"
       :class="{'font-bold':getPath('') === (state.activeMappingPath||'')}"
  >
    <div
        class="inline-flex items-baseline leading-normal relative w-full dark:hover:bg-white/5 hover:bg-primary/5"
        :class="{'underline':getPath('') === (state.activeMappingPath||'')}"
    >
      <!-- Collapse -->
      <UIcon class="carret hover:text-sky-600" name="i-heroicons-chevron-down" @click.prevent.stop="collapsed=true"
             v-if="path && !collapsed "/>
      <UIcon class="carret hover:text-sky-600" name="i-heroicons-chevron-right" @click.prevent.stop="collapsed=false"
             v-else-if="path"/>
      <!-- Name-->
      <span v-if="name && name!==as" class="italic opacity-60">{{ name }}:&nbsp;</span>
      <span v-if="as">{{ as }}&nbsp;</span>

      <!-- brackets -->
      <p @click.prevent.stop="collapsed=!collapsed" class="hover:text-sky-600">
        {{ collapsed ? '{...}' : '{' }}
        <span class="text-gray-400 ml-2" v-if="collapsed">
          {{ mapping.length >= 2 ? `// ${mapping.length} items` : `// ${mapping.length} item` }}
        </span>
      </p>

    </div>

    <div
        v-if="!collapsed"
        v-for="field in mapping"
        class="pl-4 cursor-pointer dark:hover:bg-white/5 hover:bg-primary/5"
        :class="{
          'border-l border-dashed':true,
          'border-sky-600 dark:border-green-400/60':isPathActive(''),
          'border-[#bfcbd9]':!isPathActive(''),
        }">


      <p v-if="field._type==='simple'"
         :class="{'font-bold underline':getPath(field.column) === (state.activeMappingPath||'')}"
         @keydown.enter.space.prevent.stop="setActivePath(field.column)"
         @click.prevent.stop="setActivePath(field.column)" tabindex="0">
        <span v-if="field.name && field.name!==field.column" class="italic text-sm">{{ field.name }}:&nbsp;</span>
        <span v-if="field.column">{{ field.column }}</span>
      </p>

      <MappingThree
          :mapping="field.mapping"
          :name="field.table"
          :as="field.name"
          :path="getPath(field.name)"
          v-else-if="field._type==='relation'"
      />

    </div>
    <span class="inline-flex items-baseline relative" v-if="!collapsed">}</span>
  </div>
</template>

<script setup lang="ts">

interface Props {
  mapping: any
  name?: string
  as?: string
  path?: string
}

const {mapping = [], path, name, as} = defineProps<Props>()
const collapsed = ref(false)

const state = useProcessFormState()
const setActivePath = async (path: string) => {
  state.value.activeTabs = "mapping"
  state.value.activeMappingPath = getPath(path)
}

const getPath = (fieldname: string): string => {
  return [path, fieldname].filter(i => i).join('.')
}

const isPathActive = (path: string) => {
  return state.value.activeMappingPath?.startsWith(getPath(path) + '.')
  || state.value.activeMappingPath === getPath(path)
  || getPath(path)===''
}
</script>

<style scoped lang="postcss">
.carret {
  @apply absolute left-0;
  transform: translateX(calc(-100% - 0.13em)) translateY(25%);
}
</style>