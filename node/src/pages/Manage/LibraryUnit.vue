<script setup lang="ts">
  import { defineEmits, defineProps, ref } from 'vue';
  import LibrarySwitch from './LibrarySwitch.vue';
  import { type Library } from '@/utils/api';

  const { library } = defineProps<{ library: Library }>();
  const emits = defineEmits<{
    updateIsPublished: [id: string, newIsPublished: boolean, onDone: ()=>void];
  }>();

  const isChanging = ref(false);

  async function updateIsPublished(newValue: boolean) {
    isChanging.value = true
    emits('updateIsPublished', library.id, newValue, () => {
      isChanging.value = false;
    });
  }
</script>

<template>
  <tr>
    <td style="width: 40vw;">
      <div>
        <b>{{ library.fullName }}</b>
        <br />
        Version：{{ library.version }}
      </div>
    </td>
    <td>
      <LibrarySwitch
        :display-texts="['非公開', '公開']"
        :is-changing="isChanging"
        :is-selected="library.isPublished"
        @change="updateIsPublished"
      />
    </td>
  </tr>
</template>