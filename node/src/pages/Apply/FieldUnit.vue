<script setup lang="ts">
  import { defineProps, defineModel } from 'vue';

  const { inputName, inputType = 'number' } = defineProps<{ inputName: string, inputType?: 'text'|'number' }>()

  const model = defineModel<string>({default: ''});
</script>

<template>
  <td>
    <input :type="inputType" :name="inputName" v-model="model" :class="inputType === 'number' ? 'num-field' : ''" />
  </td>
</template>

<style scoped>
  .num-field {
    width: 3em;
  }
</style>
