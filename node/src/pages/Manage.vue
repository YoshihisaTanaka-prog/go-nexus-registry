<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import { type Library, getLibraries } from '@/utils/api';

  const libraries = ref<Library[]>([]);

  onMounted(() => {
    getLibraries(libraries);
  });
</script>

<template>
  <Base :path="'manage'">
    <div v-for="library in libraries" :key="library.id">{{ library.name }}</div>
  </Base>
</template>
