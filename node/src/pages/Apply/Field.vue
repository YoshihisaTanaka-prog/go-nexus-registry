<script setup lang="ts">
  import FieldUnit from './FieldUnit.vue';
  import { defineEmits, defineModel, defineProps } from 'vue';

  const { vmId, length } = defineProps<{vmId: string, length: number}>()

  const emits = defineEmits<{'delete': [id: string]}>()

  const nameModel = defineModel<string>('name', {default: ''});
  const v1Model = defineModel<number>('v1');
  const v2Model = defineModel<number>('v2');
  const v3Model = defineModel<number>('v3');
</script>

<template>
  <tr>
    <FieldUnit
      :input-name="'name'"
      :input-type="'text'"
      v-model="nameModel"
    />
    <td>@</td>
    <FieldUnit
      :input-name="'v1'"
      v-model="v1Model"
    />
    <td>.</td>
    <FieldUnit
      :input-name="'v2'"
      v-model="v2Model"
    />
    <td>.</td>
    <FieldUnit
      :input-name="'v3'"
      v-model="v3Model"
    />
    <td>
      <button v-if="length > 1" type="button" @click="() => {emits('delete', vmId)}">X</button>
    </td>
  </tr>
</template>
