<script setup lang="ts">
  import { defineProps } from 'vue';
  import BaseBase from './BaseBase.vue';
  import SSE from './SSE.vue';

  import type { SelectionObj } from './types';
  const { path } = defineProps<{path?: string}>();

  const selections: SelectionObj = {
    apply: {
      className: "menu-button",
      displayText: "申請",
    },
    allow: {
      className: "menu-button",
      displayText: "認可（要管理者権限）",
    },
  };
</script>

<template>
  <BaseBase :path="path" :selections="selections">
    <slot />
  </BaseBase>
  <SSE />
</template>
