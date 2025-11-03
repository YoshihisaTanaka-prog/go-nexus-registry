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

  const urlParams = new URLSearchParams(window.location.search);
  const redirectTo = urlParams.get('redirect');
  const urlSuffix = redirectTo ? `?redirect=${encodeURIComponent(redirectTo)}` : undefined;
</script>

<template>
  <BaseBase :path="path" :selections="selections" :url-suffix="urlSuffix" >
    <slot />
  </BaseBase>
</template>
