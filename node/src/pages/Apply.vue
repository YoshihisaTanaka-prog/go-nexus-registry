<script setup lang="ts">
  import { ref } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import Field from './Apply/Field.vue';
  import * as api from '@/utils/api';

  const libTypeOptions = ref(['npm']);
  const libTypeModel = ref<string>();

  type LibData = {
    id: string;
    name: string;
    v1: number|undefined;
    v2: number|undefined;
    v3: number|undefined;
  }

  function createLibData(): LibData {
    return {
      id: `${Math.random()}`.slice(2),
      name: '',
      v1: undefined,
      v2: undefined,
      v3: undefined,
    }
  }

  const libDataArray = ref<LibData[]>([createLibData()])

  const addLibData = () => {
    libDataArray.value.push(createLibData())
  }

  const deleteLibData = (id: string) => {
    if (libDataArray.value.length > 1) {
      libDataArray.value = libDataArray.value.filter(libData => libData.id != id)
    }
  }

  const getLibData = (id: string) => {
    return libDataArray.value.find(libData => libData.id === id);
  }

  const setLibName = (id: string, newName: string) => {
    getLibData(id)!.name = newName;
  }

  const setLibV1 = (id: string, newV1: number|undefined) => {
    getLibData(id)!.v1 = newV1;
  }

  const setLibV2 = (id: string, newV2: number|undefined) => {
    getLibData(id)!.v2 = newV2;
  }

  const setLibV3 = (id: string, newV3: number|undefined) => {
    getLibData(id)!.v3 = newV3;
  }

  const apply = () => {
    libDataArray.value = libDataArray.value.filter(libData => libData.name != '');
    api.apply(libTypeModel.value, ...libDataArray.value.map(libData => {return {name: libData.name, v1: libData.v1, v2: libData.v2, v3: libData.v3}}));
    if (libDataArray.value.length === 0) {
      setTimeout(async () => {
        libDataArray.value = [createLibData()];
      }, 100);
    }
  }
</script>

<template>
  <Base :path="'apply'">
  <form>
    <h2>ライブラリの申請</h2>
    <p style="text-align: center;">
      ライブラリの種類：
      <select v-model="libTypeModel">
        <option v-if="libTypeModel === undefined" :value="undefined">選択してください。</option>
        <option
          v-for="(libTypeOption, index) in libTypeOptions"
          :key="index"
        >
          {{ libTypeOption }}
        </option>
      </select>
    </p>
    <table>
      <tbody>
        <Field
          v-for="libData in libDataArray"
          :key="libData.id"
          :length="libDataArray.length"
          :lib-id="libData.id"
          @delete="deleteLibData"
          @set-lib-name="setLibName"
          @set-lib-v1="setLibV1"
          @set-lib-v2="setLibV2"
          @set-lib-v3="setLibV3"
        />
      </tbody>
    </table>
    <p class="button-p">
      <button type="button" @click="addLibData">
        +
      </button>
    </p>
    <p class="button-p">
      <button type="button" class="main-button" @click="apply">
        ライブラリを申請
      </button>
    </p>
  </form>
  </Base>
</template>

<style scoped>
  h2 {
    text-align: center;
  }
  form {
    width: fit-content;
    margin-inline: auto;
    margin-block: auto;
    font-size: 1rem;
    border: 1px solid #b6bfd2;
    background-color: #fff;
    border-radius: 0.5rem;
    padding-block: 1rem;
    padding-inline: 3rem;
    box-sizing: border-box;
  }
  select {
    border: 1px solid #b6bfd2;
    border-radius: 0.2rem;
  }

  .button-p {
    text-align: center;
  }

  .main-button {
    background-color: #158654 !important;
    color: #fff;
    font-weight: bold;
    padding: 0.5em;
    border: none;
    border-radius: 0.3em;
  }

  .main-button:hover {
    background-color: #158654d0 !important;
    color: #fff;
    font-weight: bold;
  }
</style>

