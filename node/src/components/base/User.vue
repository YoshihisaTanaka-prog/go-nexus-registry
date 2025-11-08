<script setup lang="ts">
  import { defineProps, ref } from 'vue';

  const { email } = defineProps<{email?: string}>();

  const remPx = parseFloat(getComputedStyle(document.documentElement).fontSize);

  const iconSize = Math.round(remPx * 2.5);

  const viewBox = `0 0 ${iconSize} ${iconSize}`

  const faceCircleData = { cx: iconSize / 2, cy: iconSize * 0.4, r: iconSize / 3 };
  const bodyCircleData = { cx: iconSize / 2, cy: iconSize * 641 / 480, r: iconSize * 289 / 480 }

  const strokeWidth = Math.round(remPx / 16 * 10);

  const isOpeningModal = ref(false);
</script>

<template>
  <template v-if="email" >
    <div
      v-if="isOpeningModal"
      style="position: fixed; top: 0; left: 0; bottom: 0; right: 0; background-color: #b0c8e040; z-index: 99;"
      @click="() =>{ isOpeningModal = false; }"
    >
      <ul style=""  @click="(e: Event) => {e.stopPropagation()}">
        <li>1</li>
        <li>2</li>
        <li>3</li>
        <li>4</li>
      </ul>
    </div>
    <div style="position: absolute; display: flex; right: 1em;">
      {{ email }}
      <div style="display: inline-block; margin-left: 1em; z-index: 100;" @click="() => { isOpeningModal = !isOpeningModal; }">
        <svg :viewBox="viewBox" :width="`${iconSize}px`" :height="`${iconSize}px`">
          <circle :cx="faceCircleData.cx" :cy="faceCircleData.cy" :r="faceCircleData.r" fill="none" stroke="#158654" :strokeWidth="strokeWidth" />
          <circle :cx="bodyCircleData.cx" :cy="bodyCircleData.cy" :r="bodyCircleData.r" fill="none" stroke="#158654" :strokeWidth="strokeWidth" />
        </svg>
      </div>
    </div>
  </template>
</template>

<style scoped>
  ul {
    width: fit-content;
    margin: 4rem 1.5rem auto auto;
    padding: 1em;
    list-style: none;
    border-radius: 0.5rem;
    border: 1px solid #b6bfd2;
    background-color: #fff;
  }
  svg {
    border-radius: 0.5rem;
    background-color: #f4f5f9;
  }
  svg:hover {
    background-color: #e3eaf3;
  }
</style>