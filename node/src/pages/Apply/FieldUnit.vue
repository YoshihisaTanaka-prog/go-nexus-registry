<script setup lang="ts">
  import { defineProps, defineEmits, ref } from 'vue';

  const { inputType = 'number', defaultValue } = defineProps<
    { inputType: 'text'|'number', defaultValue?: string } |
    { inputType?: 'number', defaultValue?: number }
  >()

  const model = inputType ==='number' ? ref<number|undefined>(defaultValue as number | undefined): ref(defaultValue as string);
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
