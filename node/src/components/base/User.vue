<script setup lang="ts">
  import { defineEmits, defineProps } from 'vue';

  const { email } = defineProps<{email?: string}>();
  const emits = defineEmits<{
    openUserModal: []
  }>();

  const remPx = parseFloat(getComputedStyle(document.documentElement).fontSize);

  const iconSize = Math.round(remPx * 2.5);

  const viewBox = `0 0 ${iconSize} ${iconSize}`

  const faceCircleData = { cx: iconSize / 2, cy: iconSize * 0.4, r: iconSize / 3 };
  const bodyCircleData = { cx: iconSize / 2, cy: iconSize * 641 / 480, r: iconSize * 289 / 480 }

  const strokeWidth = Math.round(remPx / 16 * 10);
</script>

<template>
  <div v-if="email" style="position: absolute; display: flex; right: 1em;">
    {{ email }}
    <div style="display: inline-block; margin-left: 1em;" @click="() => { emits('openUserModal') }">
      <svg :viewBox="viewBox" :width="`${iconSize}px`" :height="`${iconSize}px`">
        <circle :cx="faceCircleData.cx" :cy="faceCircleData.cy" :r="faceCircleData.r" fill="none" stroke="#158654" :strokeWidth="strokeWidth" />
        <circle :cx="bodyCircleData.cx" :cy="bodyCircleData.cy" :r="bodyCircleData.r" fill="none" stroke="#158654" :strokeWidth="strokeWidth" />
      </svg>
    </div>
  </div>
</template>

<style scoped>
  svg {
    border-radius: 0.5rem;
    background-color: #f4f5f9;
  }
  svg:hover {
    background-color: #e3eaf3;
  }
</style>