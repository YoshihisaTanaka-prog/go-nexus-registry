<script setup lang="ts">
  import { ref, watch } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import SelectLibKind from '@/components/SelectLibKind.vue'
  import { type Library, getLibraries } from '@/utils/api';

  const libraries = ref<Library[]>([]);

  const libKindModel = ref<string>();


  watch(libKindModel, () => {
    if(libKindModel.value !== undefined) {
      getLibraries(libKindModel.value, libraries);
    }
  });
</script>

<template>
  <Base :path="'manage'">
    <SelectLibKind v-model="libKindModel" />
    <div v-for="library in libraries" :key="library.id">{{ library.name }} &nbsp; {{ library.version }}</div>
  </Base>
</template>
