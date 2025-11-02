<script setup lang="ts">
  import { defineProps } from 'vue';
  import Toast from '@/components/Toast.vue';
  import IconImage from "/icon.png";
  import type { SelectionObj } from './types';
  const { path, selections, urlSuffix = '' } = defineProps<{path?: string, selections: SelectionObj, urlSuffix?: string}>();

  if ( path ) {
    if (Object.keys(selections).includes(path) && selections[path]) {
      selections[path].className = "menu-button selected-menu-button"
    }
  }

  const onClickItem = (path: string) => {
    location.href = (import.meta.env.MODE === 'development' ? `/htmls/${path}.html` :  `/${path}`) + urlSuffix;
  }
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
    </div>
  </div>
  <div id="main">
    <div class="menu">
      <ul class="menu-ul">
        <li
          v-for="[key, value] in Object.entries(selections)"
          :key="key"
          :class="value.className"
          @click="()=>{if (path !== key) {onClickItem(key)}}"
        >
          {{ value.displayText }}
        </li>
      </ul>
    </div>
    <div class="contents">
      <slot />
      <Toast />
    </div>
  </div>
  <div id="footer">
    <div>
      NePlusはSonatype社とは一切関係のない非公式ツールです。SonatypeおよびNexusはSonatype, Inc.の商標です。
    </div>
    <div style="margin-right: 1em;">
      &copy; Yoshihisa Tanaka 2025
    </div>
  </div>
</template>

<style scoped>
  .header-main-div {
    display: flex;
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

  .menu {
    padding: 0.9rem;
    box-sizing: border-box;
    border-right: 1px solid #b6bfd2;
  }

  .contents {
    background-color: #f4f5f9;
    flex: 1;
    padding: 1.5rem;
    position: relative;
    overflow-y: scroll;
  }

  .selected-menu-button {
    background-color: #158654 !important;
    color: #fff;
    font-weight: bold;
  }

  .selected-menu-button:hover {
    background-color: #158654d0 !important;
    color: #fff;
    font-weight: bold;
  }

  .menu-button {
    background-color: #00000000;
    width: 14rem;
    text-align: left;
    font-size: 0.9rem;
    padding: 1em;
    padding-bottom: 1.1em;
    border-radius: 0.3em;
    border: none;
  }

  .menu-button:hover {
    background-color: #dadfe7;
  }

  .menu-ul {
    list-style: none;
    margin-block-start: 0;
    padding-inline-start: 0;
  }

  .menu-ul li {
    margin-block-end: 0.9rem;
  }

</style>
