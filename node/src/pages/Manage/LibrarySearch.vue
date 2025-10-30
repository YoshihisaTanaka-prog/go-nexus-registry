<script setup lang="ts">
  import { computed, defineEmits, defineProps, ref, watch } from 'vue';
  import LibraryUnit from './LibraryUnit.vue';
  import { type Library } from '@/utils/api';

  const { libraries } = defineProps<{ libraries: Library[] }>();
  const emits = defineEmits<{
    updateIsPublished: [id: string, newIsPublished: boolean, onDone: ()=>void]
  }>();

  const page = ref(0);
  const numOfDisplayedLibraries = ref(10);
  const searchWord = ref<string>("");

  const filteredLibraries = computed(() => searchWord.value === "" ? libraries : libraries.filter((l) => l.name.includes(searchWord.value)));
  watch(filteredLibraries, () => {
    page.value = 0;
  });
  const length = computed(() => Math.ceil(filteredLibraries.value.length / numOfDisplayedLibraries.value));
  const selectedLibraries = computed(() => {
    const sliceStart = page.value * numOfDisplayedLibraries.value;
    const sliceEnd = sliceStart + numOfDisplayedLibraries.value;
    return filteredLibraries.value.slice(sliceStart, sliceEnd);
  });

  function goToPrevPage() {
    page.value = Math.max(0, page.value - 1);
  }
  function goToNextPage() {
    page.value = Math.min(length.value-1, page.value + 1);
  }

  function updateIsPublished(id: string, newIsPublished: boolean, onDone: ()=>void) {
    emits('updateIsPublished', id, newIsPublished, onDone);
  }
</script>

<template>
  <div style="text-align: center;">
    <div v-if="length > 0">
      <label>
        検索：
        <input type="text" v-model="searchWord" />
      </label>
    </div>
    <table style="display: inline-block; text-align: left;">
      <tbody>
        <LibraryUnit v-for="library in selectedLibraries" :key="library.id" :library="library" @update-is-published="updateIsPublished" />
      </tbody>
    </table>
    <div style="display: flex; align-items: center; justify-content: center; vertical-align: middle;">
      <div v-if="length > 1">
        <button @click="goToPrevPage">←</button>
        <button @click="goToNextPage">→</button>
      </div>
      <div v-if="length == 0">
        none / none
      </div>
      <div v-else>
        {{ length === 0 ? 0 : page + 1 }} / {{ length }}
      </div>
    </div>
  </div>
</template>

<style scoped>
  table {
    border-collapse: separate;
    border-spacing: 0 1em;
  }
  button {
    margin-right: 1.5em;
  }
</style>