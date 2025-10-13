<script setup lang="ts">
  import { defineProps, defineEmits, ref } from 'vue';
  const { pageTitle, passwordAutocomplete } = defineProps<{pageTitle: string, passwordAutocomplete: string}>();

  const emits = defineEmits<{onClickedButton: [string, string]}>()

  const userId = ref<string|undefined>('yoshihisa-tanaka@example.com')
  const password = ref<string>('')
</script>

<template>
  <div v-if="userId === undefined"></div>
  <form v-else>
    <h2>{{ pageTitle }}</h2>
    <p>
      <label>
        ユーザーID<br />
        <input type="hidden" name="email" :value="userId!" />
        <input type="text" :value="userId!" autocomplete="email" readonly disabled />
      </label>
    </p>
    <p>
      <label>
        パスワード<br />
        <input type="password" name="password" :autocomplete="passwordAutocomplete" v-model="password" />
      </label>
    </p>
    <slot />
    <p class="button-p">
      <button type="button" @click="() => {emits('onClickedButton', userId!, password)}">
        {{ pageTitle }}
      </button>
    </p>
  </form>
</template>

<style scoped>
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

  .button-p button {
    background-color: #158654 !important;
    color: #fff;
    font-weight: bold;
    padding: 0.5em;
    border: none;
    border-radius: 0.3em;
  }

  .button-p button:hover {
    background-color: #158654d0 !important;
    color: #fff;
    font-weight: bold;
  }
</style>
