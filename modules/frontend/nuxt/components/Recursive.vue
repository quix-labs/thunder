<template>
  <section>
    <h1>TEST</h1>
    <UButton @click.prevent="addRelation">
      Add relation
    </UButton>

    <div class="p-2" v-for="(relation,index) in relations">
      <h1>{{ relation.name }} + {{id}}</h1>
      <LazyRecursive/>
    </div>
  </section>

</template>
<script setup lang="ts">
const id = useId()
const relations = ref<Array<{ name: string }>>([])
const addRelation = () => {
  relations.value.push({name: `Relation ${relations.value?.length + 1}`})
}
</script>