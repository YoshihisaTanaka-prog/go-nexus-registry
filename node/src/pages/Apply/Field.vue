<script setup lang="ts">
  import FieldUnit from './FieldUnit.vue';
  import { defineEmits, defineProps } from 'vue';

  type LibData = {
    id: string;
    name: string;
    v1: number|undefined;
    v2: number|undefined;
    v3: number|undefined;
  }

  const { library, length } = defineProps<{
    library: LibData,
    length: number,
  }>();

  const libId = library.id;

  const emits = defineEmits<{
    'delete': [id: string];
    'setLibName': [id: string, name: string];
    'setLibV1': [id: string, v1: number|undefined];
    'setLibV2': [id: string, v2: number|undefined];
    'setLibV3': [id: string, v3: number|undefined];
  }>();
</script>

<template>
  <tr>
    <FieldUnit
      :input-type="'text'"
      :default-value="library.name"
      @set="(newValue: string | number | undefined)=>{emits('setLibName', libId, newValue as string)}"
    />
    <td>@</td>
    <FieldUnit
      :default-value="library.v1"
      @set="(newValue: string | number | undefined)=>{emits('setLibV1', libId, newValue as number | undefined)}"
    />
    <td>.</td>
    <FieldUnit
      :default-value="library.v2"
      @set="(newValue: string | number | undefined)=>{emits('setLibV2', libId, newValue as number | undefined)}"
    />
    <td>.</td>
    <FieldUnit
      :default-value="library.v3"
      @set="(newValue: string | number | undefined)=>{emits('setLibV3', libId, newValue as number | undefined)}"
    />
    <td>
      <button v-if="length > 1" type="button" @click="() => {emits('delete', libId)}">X</button>
    </td>
  </tr>
</template>
