<script setup lang="ts">
  import { defineEmits, ref } from 'vue';

  const emits = defineEmits<{
    close: [],
    create: [name: string, callback: () => void]
  }>()

  function close() {
    emits('close');
  }

  function noClose(e: Event) {
    e.stopPropagation();
  }

  const nameModel = ref('');
  const isSubmitting = ref(false);

  function create() {
    if (!isSubmitting.value) {
      isSubmitting.value = true;
      emits('create', nameModel.value, () => {
        isSubmitting.value = false;
      });
    }
  }
</script>

<template>
  <div class="modal-back-div" @click="close">
    <div @click="noClose" class="modal-main-div">
      <div style="position: relative;">
        <button class="close-button" @click="close">&times;</button>
      </div>
      <h3>ロール作成</h3>
      <p>
        <input type="text" v-model="nameModel" />
      </p>
      <p>
        <span v-if="isSubmitting">実行中...</span>
        <button v-else @click="create">作成</button>
      </p>
    </div>
  </div>
</template>

<style scoped>
  .modal-back-div {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 1.5rem;
    background-color: #b0c8e040;
  }
  .modal-main-div {
    width: calc(50vw - 2em);
    margin: 1em auto;
    border: 1px solid #b6bfd2;
    border-radius: 0.5rem;
    background-color: #fff;
    text-align: center;
  }
  .close-button {
    position: absolute;
    top: -1px;
    right: -1px;
    box-sizing: content-box;
  }
</style>
