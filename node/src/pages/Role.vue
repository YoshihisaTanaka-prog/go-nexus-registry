<script setup lang="ts">
  import Base from '@/components/base/ContentsBase.vue';
  import { getRoles } from '@/utils/api/role';

  window.addEventListener('load', async () => {
    console.log(await getRoles());
  });
</script>

<template>
  <Base :path="'role'" />
</template>

<style scoped></style>
