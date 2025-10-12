<script setup lang="ts">
  import { ref } from 'vue';
  import Base from '@/components/base/SignBase.vue';
  import Form from '@/components/sign/Form.vue';
  import { signUp } from '@/utils/api';

  const confirmPasswordModel = ref('')

  const onClickButton = (userId: string, password: string) => {
    if ( password === confirmPasswordModel.value ) {
      signUp(userId, password)
    } else {
      alert('パスワードと確認用のパスワードが異なる値です。')
    }
  }
</script>

<template>
  <Base :path="'sign-up'">
    <Form
      :page-title="'サインアップ'"
      :password-autocomplete="'new-password'"
      @on-clicked-button="onClickButton"
    >
      <input type="password" name="password-confirm" autocomplete="new-password" v-model="confirmPasswordModel" />
    </Form>
  </Base>
</template>

<style scoped></style>
