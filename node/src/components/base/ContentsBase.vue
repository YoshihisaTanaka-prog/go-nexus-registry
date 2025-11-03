<script setup lang="ts">
  import { defineProps, ref } from 'vue';
  import BaseBase from './BaseBase.vue';
  import SSE from './SSE.vue';

  import type { SelectionObj } from './types';
  const { path } = defineProps<{path?: string, style?: unknown}>();

  const baseSelections = [
    {path: "apply", text: "申請"},
    {path: "manage", text: "ライブラリ設定"},
    {path: "role", text: "ロール設定"},
  ]

  const selections = ref<SelectionObj>({});

  window.addEventListener("load", async () => {
    await Promise.all(baseSelections.map(({path, text}) => {
      return new Promise<void>(async (resolve) => {
        try {
          const response = await fetch(`/${path}`, {redirect: 'error'});
          if (response.ok && response.status != 302) {
            selections.value[path] = {
              className: "menu-button",
              displayText: text,
            }
          }
        } finally {
          resolve();
        }
      });
    }));
    const keys = Object.keys(selections.value);
    if (keys.length == 1 && path != keys[0]) {
      const key = keys[0];
      if (path != key) {
        location.href = `/${key}`
      }
    }
  });
</script>

<template>
  <BaseBase :path="path" :selections="selections">
    <slot />
  </BaseBase>
  <SSE />
</template>
