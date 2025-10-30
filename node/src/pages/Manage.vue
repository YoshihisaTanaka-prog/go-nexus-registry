<script setup lang="ts">
  import { ref, watch } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import SelectLibKind from '@/components/SelectLibKind.vue'
  import LibrarySearch from './Manage/LibrarySearch.vue';
  import { type Library, getLibraries, updateIsPublished } from '@/utils/api';

  const libraries = ref<Library[]>([]);

  const libKindModel = ref<string>();

  watch(libKindModel, () => {
    if(libKindModel.value !== undefined) {
      getLibraries(libKindModel.value, libraries);
    }
  });

  function _updateIsPublished(id: string, newIsPublished: boolean, onDone: ()=>void) {
    void updateIsPublished(id, newIsPublished, libraries, onDone);
  }
</script>

<template>
  <Base :path="'manage'">
    <SelectLibKind v-model="libKindModel" />
    <LibrarySearch :libraries="libraries" @update-is-published="_updateIsPublished" />
  </Base>
</template>
