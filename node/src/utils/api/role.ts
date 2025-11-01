import { axios, convertError, alertError } from './_base'

type Role = {
  id: string;
  name: string;
  privileges: string[]
}

export async function getRoles():Promise<Role[]> {
  const result = await axios.get('get-roles');
  return result.data;
}