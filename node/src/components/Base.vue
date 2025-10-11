<script setup lang="ts">
  import { defineProps } from 'vue';
  const { path } = defineProps({
    path: {
      type: String
    }
  });

  type SelectionObj = {
    className: string;
    displayText: string
  }

  const selections: {[key: string]: SelectionObj} = {
    apply: {
      className: "menu-button",
      displayText: "申請",
    },
    allow: {
      className: "menu-button",
      displayText: "認可（要管理者権限）",
    },
  }

  if (Object.keys(selections).includes(path)) {
    selections[path].className = "menu-button selected-menu-button"
  }
</script>

<template>
  <div id="header">
    Header
  </div>
  <div id="main">
    <div class="menu">
      <ul class="menu-ul">
        <li
          v-for="[key, value] in Object.entries(selections)"
          key="key"
          :class="value.className"
        >
          {{ value.displayText }}
        </li>
      </ul>
    </div>
    <div class="contents">
      <slot></slot>
    </div>
  </div>
  <div id="footer">
    &copy; Yoshihisa Tanaka 2025
  </div>
</template>

<style scoped>
  #header {
    border-bottom: 1px solid #b6bfd2;
  }

  #footer {
    border-top: 1px solid #b6bfd2;
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
  }

  .selected-menu-button {
    background-color: #158654 !important;
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
