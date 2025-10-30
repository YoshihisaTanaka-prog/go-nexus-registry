<script setup lang="ts">
  import { defineProps, defineEmits, ref } from 'vue';

  const { inputType = 'number' } = defineProps<{ inputType?: 'text'|'number' }>()

  const model = inputType ==='number' ? ref<number|undefined>(): ref('');
  const emits = defineEmits<{
    'set': [value: string|number|undefined]
  }>();

  const setValue = () => {
    emits('set', model.value);
  }
</script>

<template>
  <td>
    <input
      :class="inputType === 'number' ? 'num-field' : ''"
      :type="inputType"
      v-model="model"
      @input="setValue"
    />
  </td>
</template>

<style scoped>
  input {
    border: 1px solid #b6bfd2;
    border-radius: 0.2rem;
  }
  .num-field {
    width: 3em;
  }
</style>
