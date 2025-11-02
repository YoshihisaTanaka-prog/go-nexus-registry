<script setup lang="ts">
  import { defineEmits, defineProps, ref } from 'vue';
  import type { Role } from '@/utils/api/role';
  
  const { role } = defineProps<{role: Role}>();

  const roleName = ref(role.name);
  const isEditing = ref(false);
  const isSubmitting = ref(false);

  const emits = defineEmits<{save: [id: string, name: string, callback:() => void]}>()

  function save() {
    if (!isSubmitting.value) {
      isSubmitting.value = true;
      emits('save', role.id, roleName.value, () => {
        isEditing.value = false;
        isSubmitting.value = false;
      });
    }
  }
</script>

<template>
  <p>
    <span v-if="isEditing">
      <div class="button-div">
        <span v-if="isSubmitting">
          実行中...
        </span>
        <button v-else @click="save">保存</button>
      </div>
      <input v-model="roleName" />
    </span>
    <span v-else>
      <div class="button-div">
        <button @click="()=>{isEditing = true}">編集</button>
      </div>
      {{ roleName }}
    </span>
  </p>
</template>

<style scoped>
  .button-div {
    display: inline-block;
    width: 5em;
    margin-right: 1em;
    text-align: left;
  }
</style>
