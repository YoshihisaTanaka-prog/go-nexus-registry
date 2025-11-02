<script setup lang="ts">
  import { ref } from 'vue';
  import Base from '@/components/base/ContentsBase.vue';
  import RoleUnit from './Role/RoleUnit.vue';
  import CreateRole from './Role/CreateRole.vue';
  import { getRoles, createRole, updateRole, type Role } from '@/utils/api/role';

  const roles = ref<Role[]>([]);

  function _updateRole(id: string, name: string, callback:() => void) {
    void updateRole(roles, id, name, callback);
  }

  function _createRole(name: string, callback:() => void) {
    void createRole(roles, name, () => {isCreating.value = false;}, callback);
  }

  window.addEventListener('load', async () => {
    roles.value = await getRoles();
  });

  const isCreating = ref(false);
</script>

<template>
  <Base :path="'role'" style="position: relative;">
    <div class="main-div">
      <div style="position: relative;">
        <p style="position: absolute; right: 1em;">
          <button @click="() => {isCreating = true}">+</button>
        </p>
      </div>
      <h3 style="text-align: center;">ロール一覧</h3>
      <RoleUnit v-for="role in roles" :key="role.id" :role="role" @save="_updateRole" />
    </div>
    <CreateRole v-if="isCreating" @close="() => {isCreating = false}" @create="_createRole" />
  </Base>
</template>

<style scoped>
  .main-div {
    width: 50vw;
    margin: 0 auto;
    padding-inline: 1em;
    border: 1px solid #b6bfd2;
    border-radius: 0.5rem;
    background-color: #fff;
  }
</style>
