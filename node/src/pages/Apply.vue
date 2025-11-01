<script setup lang="ts">
  import { ref, useTemplateRef } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import Field from './Apply/Field.vue';
  import SelectLibKind from '@/components/SelectLibKind.vue'
  import * as api from '@/utils/api';
  import { formatNpm } from '@/utils/format';

  const libKindModel = ref<string>();

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

  function createSpecificLibData(name: string, version: string): LibData {
    const library = createLibData();
    library.name = name;
    const [v1 = 0, v2 = 0, v3 = 0] = version.split('.').map(v => Number(v));
    library.v1 = v1;
    library.v2 = v2;
    library.v3 = v3;
    return library;
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
    api.apply(libKindModel.value, ...libDataArray.value.map(libData => {return {name: libData.name, v1: libData.v1, v2: libData.v2, v3: libData.v3}}));
    if (libDataArray.value.length === 0) {
      setTimeout(async () => {
        libDataArray.value = [createLibData()];
      }, 100);
    }
  }

  const fileInputRef = useTemplateRef('file-input-ref');

  const onChangeFileInputRef = async () => {
    const files =fileInputRef.value?.files;
    if (files == null) {
      return
    }

    const newLibraryDataArray = [...libDataArray.value];
    switch (libKindModel.value) {
      case 'npm':
        const libraries = await formatNpm(files)
        newLibraryDataArray.push(...libraries.map(library => createSpecificLibData(library.name, library.version)))
        break;
    
      default:
        break;
    }
    const filteredLibDataArray: LibData[] = [];
    for(const libData of newLibraryDataArray) {
      const foundLibData = filteredLibDataArray.find((l) => {
        return (l.name === libData.name) &&
          (l.v1 === libData.v1) &&
          (l.v2 === libData.v2) &&
          (l.v3 === libData.v3)
      });
      if (foundLibData == undefined) {
        if (libData.name != '') {
          filteredLibDataArray.push(libData)
        }
      }
    }
    libDataArray.value = filteredLibDataArray;
  }
</script>

<template>
  <Base :path="'apply'">
  <form>
    <h2>ライブラリの申請</h2>
    <SelectLibKind v-model="libKindModel" />
    <div v-if="libKindModel !== undefined">
      <p style="text-align: center;">
        <input ref="file-input-ref" type="file" @change="onChangeFileInputRef" />
      </p>
      <table>
        <tbody>
          <Field
            v-for="libData in libDataArray"
            :key="libData.id"
            :length="libDataArray.length"
            :library="libData"
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
    </div>
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
    position: relative;
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

