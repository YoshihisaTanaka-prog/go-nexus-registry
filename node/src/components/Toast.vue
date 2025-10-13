<script setup lang="ts">
  import { toasts } from '@/utils/toast';
</script>

<template>
  <div class="main">
    <div
      class="item"
      v-for="(toast, index) in toasts"
      :key="toast.id"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<style scoped>
  .main {
    position: absolute;
    bottom: 1rem;
    right: 1rem;
  }
  .item {
    margin-bottom: 1rem;
    width: 15vw;
    background-color: #fff;
    border-radius: 0.2rem;
    border: 1px solid #000;
  }
</style>
