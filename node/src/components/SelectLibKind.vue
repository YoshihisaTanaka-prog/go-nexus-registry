<script setup lang="ts">
  import { ref, defineModel } from 'vue';

  const libKindOptions = ref(['npm']);
  const libKindModel = defineModel<string>();
</script>

<template>
  <p style="text-align: center;">
    ライブラリの種類：
    <select v-model="libKindModel">
      <option v-if="libKindModel === undefined" :value="undefined">選択してください。</option>
      <option
        v-for="(libKindOption, index) in libKindOptions"
        :key="index"
      >
        {{ libKindOption }}
      </option>
    </select>
  </p>
</template>

<style scoped>
  select {
    border: 1px solid #b6bfd2;
    border-radius: 0.2rem;
  }
</style>

