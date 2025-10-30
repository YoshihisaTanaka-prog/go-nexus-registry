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
    manage: {
      className: "menu-button",
      displayText: "ライブラリ設定",
    },
  };
</script>

<template>
  <BaseBase :path="path" :selections="selections">
    <slot />
  </BaseBase>
  <SSE />
</template>
