<script setup lang="ts">
  import { defineEmits, defineProps } from 'vue';
  import User from './User.vue';
  import IconImage from "/icon.png";

  const { email } = defineProps<{email?: string}>()
  const emits = defineEmits<{
    openUserModal: []
  }>();
</script>

<template>
  <div id="header">
    <div class="header-main-div">
      <img :src="IconImage" />
      <div class="header-sub-div">
        <div style="font-weight: bolder;">
          Golang+Vue <span style="color: #158654;">NePlus Manager</span>
        </div>
        <div style="font-size: 0.5em;">
          <span style="font-weight: bolder; color: #f00">Unofficial</span> <span style="color: #aaa;">Nexus Manager</span>
        </div>
      </div>
      <User :email="email" @open-user-modal="() => { emits('openUserModal') }" />
    </div>
  </div>
</template>

<style scoped>
  .header-main-div {
    display: flex;
    position: relative;
  }
  .header-sub-div {
    display: inline-block;
    height: 100%;
  }
   img {
    width: 2.5rem;
    margin-left: 1rem;
    margin-right: 0.5rem;
  }
</style>