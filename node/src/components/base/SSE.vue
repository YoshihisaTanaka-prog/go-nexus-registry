<script setup lang="ts">
  import toast from '@/utils/toast';

  window.addEventListener("load", () => {
    const eventSource = new EventSource('/sse');
    eventSource.onopen = () => {
      // 接続の確立
      console.log('connected')
    };
    eventSource.onerror = () => {
      // 接続の切断、など
      console.error('disconnected')
    }
    eventSource.onmessage = ({ data }) => {
      // サーバからのメッセージを受信（データは data にある）
      toast("success", JSON.stringify(data));
    };
  });
</script>

<template>
  <div />
</template>
