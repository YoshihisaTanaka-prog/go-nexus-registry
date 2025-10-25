<script setup lang="ts">
  import { ref, watch } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import SelectLibKind from '@/components/SelectLibKind.vue'
  import LibrarySearch from './Menu/LibrarySearch.vue';
  import { type Library, getLibraries, updateIsPublishing } from '@/utils/api';

  const libraries = ref<Library[]>([]);

  const libKindModel = ref<string>();

  watch(libKindModel, () => {
    if(libKindModel.value !== undefined) {
      getLibraries(libKindModel.value, libraries);
    }
  });

  function _updateIsPublishing(id: string, newIsPublishing: boolean, onDone: ()=>void) {
    void updateIsPublishing(id, newIsPublishing, libraries, onDone);
  }
</script>

<template>
  <Base :path="'manage'">
    <SelectLibKind v-model="libKindModel" />
    <LibrarySearch :libraries="libraries" @update-is-publishing="_updateIsPublishing" />
  </Base>
</template>
