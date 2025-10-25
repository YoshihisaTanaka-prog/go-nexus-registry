<script setup lang="ts">
  import { defineEmits, defineProps } from 'vue';

  const { displayTexts, isSelected, isChanging } = defineProps<{
    displayTexts: [string, string];
    isSelected: boolean;
    isChanging: boolean;
  }>();

  const emits = defineEmits<{
    change: [newValue: boolean]
  }>();

  function change(newValue: boolean) {
    if (isSelected != newValue) {
      emits('change', newValue);
    } 
  }
</script>

<template>
  <div v-if="isChanging">
    変更中...
  </div>
  <div v-else class="main-div">
    <div :class="isSelected ? 'item' : 'item selected-item'" @click="() => {change(false)}">
      {{ displayTexts[0] }}
    </div>
    <div :class="isSelected ? 'item selected-item' : 'item'" @click="() => {change(true)}">
      {{ displayTexts[1] }}
    </div>
  </div>
</template>

<style scoped>
  .main-div {
    background-color: #fff;
    border: 1px solid #b6bfd2;
    border-radius: 0.3em;
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }
  .item {
    width: 5em;
  }
  .selected-item {
    border: 1px solid #b6bfd2;
    border-radius: 0.3em;
    padding: 0.5em;
    box-sizing: border-box;
    background-color: #f4f5f9;
  }
</style>