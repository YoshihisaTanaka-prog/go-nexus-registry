<script setup lang="ts">
  import Base from '@/components/base/SignBase.vue';
  import Form from '@/components/sign/Form.vue';
  import { signIn } from '@/utils/api';

  const onClickButton = (userId: string, password: string) => {
    signIn(userId, password)
  }
</script>

<template>
  <Base :path="'sign-in'">
    <Form
      :page-title="'サインイン'"
      :password-autocomplete="'current-password'"
      @on-clicked-button="onClickButton"
    />
  </Base>
</template>

<style scoped></style>
