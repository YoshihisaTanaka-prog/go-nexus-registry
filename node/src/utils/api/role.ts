import { type Ref } from 'vue';
import { axios, convertError, alertError } from './_base'

export type Role = {
  id: string;
  name: string;
  privileges: string[]
}

export async function getRoles():Promise<Role[]> {
  return new Promise(async (resolve, reject) => {
    try {
      const result = await axios.get('get-roles');
      resolve(result.data);
    } catch (error) {
      const { status, data } = convertError(error);
      alertError(status, data)
      reject()
    }
  })
}

export async function createRole(currentRoleRef: Ref<Role[]>, name: string, onSuccess:() => void, onFailed:() => void): Promise<void> {
  try {
    const keptRoles = [...currentRoleRef.value];
    const result = await axios.post('create-role', {name});
    const newRole = result.data as Role;
    keptRoles.push(newRole);
    keptRoles.sort((a, b) => {
      if (a.name > b.name) {
        return 1;
      }
      if (a.name < b.name) {
        return -1;
      }
      return 0;
    });
    currentRoleRef.value = keptRoles;
    onSuccess();
  } catch (error) {
    const { status, data } = convertError(error);
    alertError(status, data);
    onFailed();
  }
}

export async function updateRoleName(currentRoleRef: Ref<Role[]>, id: string, name: string, callback:() => void): Promise<void> {
  const keptRoles = [...currentRoleRef.value];
  const targetRole = keptRoles.find(r => r.id == id);
  try {
    if (targetRole!.name == name) {
      callback();
      return;
    }
    await axios.put('update-role', {id, name});
    targetRole!.name = name;
    keptRoles.sort((a, b) => {
      if (a.name > b.name) {
        return 1;
      }
      if (a.name < b.name) {
        return -1;
      }
      return 0;
    });
  } catch (error) {
    const { status, data } = convertError(error);
    alertError(status, data)
  } finally {
    currentRoleRef.value = [];
    setTimeout(() => {
      currentRoleRef.value = keptRoles;
      callback();
    }, 0);
  }
}