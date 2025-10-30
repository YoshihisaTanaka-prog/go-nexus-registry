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
    margin-top: 1rem;
    width: 15vw;
    background-color: #fff;
    border-radius: 0.2rem;
    border: 1px solid #b6bfd2;
    box-sizing: border-box;
    padding-block: 0.5rem;
  }
</style>
