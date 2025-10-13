<script setup lang="ts">
  import { ref } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import Field from './Apply/Field.vue';
  import * as api from '@/utils/api';

  function createVModels() {
    return {
      id: `${Math.random()}`.slice(2),
      name: ref(''),
      v1: ref<number>(),
      v2: ref<number>(),
      v3: ref<number>(),
    }
  }

  const vModelsArray = ref([createVModels()])

  const addVModel = () => {
    vModelsArray.value.push(createVModels())
  }

  const deleteVModel = (id: string) => {
    if (vModelsArray.value.length > 1) {
      vModelsArray.value = vModelsArray.value.filter(vm => vm.id != id)
    }
  }

  const apply = () => {
    api.apply(...vModelsArray.value.map(vm => {return {type: 'npm', name: vm.name, v1: vm.v1, v2: vm.v2, v3: vm.v3}}))
  }
</script>

<template>
  <Base :path="'apply'">
  <form>
    <h2>ライブラリの申請</h2>
    <table>
      <tbody>
        <Field
          v-for="vm in vModelsArray"
          :key="vm.id"
          :length="vModelsArray.length"
          :vm-id="vm.id"
          v-model:name="vm.name"
          v-model:v1="vm.v1"
          v-model:v2="vm.v2"
          v-model:v3="vm.v3"
          @delete="deleteVModel"
        />
      </tbody>
    </table>
    <p class="button-p">
      <button type="button" @click="addVModel">
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

