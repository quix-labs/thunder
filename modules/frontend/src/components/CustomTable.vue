<template>
  <UTable @update:rowSelection="()=>{console.log('ok')}" :columns="computedColumns" v-model:sorting="sorting"
          :data="rows" :loading="loading"/>
</template>

<script setup lang="ts" generic="T">
import type {TableColumn} from "@nuxt/ui";
import {UButton} from "#components";
import type {SortingState} from "@tanstack/table-core";

interface Props {
  rows: T[];
  columns: Array<{
    key: string;
    label?: string;
    sortable?: boolean;
  }>;
  loading?: boolean;
}

const {rows, columns, loading = false} = defineProps<Props>()
const sorting = defineModel<SortingState>('sorting', {required: false})

const slots = defineSlots()


const computedColumns = computed<TableColumn<T>[]>(() => {
  return columns.map(col => ({
    accessorKey: col.key,
    header: col.label,

    // Sortable
    ...(!!col.sortable ? {
      header: ({column}) => {
        const isSorted = column.getIsSorted()
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: col.label,
          icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up-20-solid' : 'i-heroicons-bars-arrow-down-20-solid' : 'i-heroicons-arrows-up-down-20-solid',
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc')
        })
      }
    } : {}),

    // Slot
    ...(`cell-${col.key}` in slots ? {
      cell: ({row}) => h(slots[`cell-${col.key}`], {row: row.original})
    } : {})
  }))
})

</script>