<script setup lang="ts">
  import { defineProps, ref } from 'vue';
  import BaseBase from './BaseBase.vue';
  import SSE from './SSE.vue';

  import { getMyProfile } from '@/utils/api';
  import type { SelectionObj } from './types';
  const { path } = defineProps<{path?: string, style?: unknown}>();

  const selectionConfig = {
    allUsers: [{targetPath: 'apply', text: "申請"}],
    forEditor: [{targetPath: 'manage', text: "ライブラリ設定"}],
    forAdmin: [{targetPath: 'role', text: "ロール設定"}],
  }

  const selections = ref<SelectionObj>({});

  const email = ref<string>()
  window.addEventListener('load', async () => {
    for (const {targetPath, text} of selectionConfig.allUsers){
      selections.value[targetPath] = {
        className: 'menu-button',
        displayText: text
      }
    }
    const profile = await getMyProfile();
    if (profile) {
      try {
        email.value = profile.email;
        const { isAdmin, isEditor } = profile;
        if ( isAdmin || isEditor ) {
          for (const {targetPath, text} of selectionConfig.forEditor){
            selections.value[targetPath] = {
              className: 'menu-button',
              displayText: text
            }
          }
        }
        if (isAdmin) {
          for (const {targetPath, text} of selectionConfig.forAdmin){
            selections.value[targetPath] = {
              className: 'menu-button',
              displayText: text
            }
          }
        }
      } finally {
        const keys = Object.keys(selections.value);
        if (keys.length == 1 && path !== undefined) {
          const key = keys[0];
          if (path != key) {
            location.href = import.meta.env.MODE === 'development' ? `/htmls/${key}.html` :  `/${key}`;
          }
        }
      }
    }
  });
</script>

<template>
  <BaseBase :path="path" :selections="selections" :email="email">
    <slot />
  </BaseBase>
  <SSE />
</template>
