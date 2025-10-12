<script setup lang="ts">
  import { defineProps } from 'vue';
  import BaseBase from './BaseBase.vue';
  import type { SelectionObj } from './types';
  const { path } = defineProps<{path?: string}>();

  const selections: SelectionObj = {
    "sign-in": {
      className: "menu-button",
      displayText: "サインイン",
    },
    "sign-up": {
      className: "menu-button",
      displayText: "サインアップ",
    },
  };
</script>

<template>
  <BaseBase :path="path" :selections="selections">
    <slot />
  </BaseBase>
</template>
