import type { Ref } from 'vue'
import { axios, convertError, alertError } from '@/utils/api/_base'
import type { Library } from './types'

export async function updateIsPublishing(id: string, isPublished: boolean, currentLibraryList: Ref<Library[]>, onDone: ()=>void) {
  try {
    const result = await axios.post('update-is-publishing', {id, isPublished});
    const newLibraryData: Library = result.data;
    const newLibraries = [...currentLibraryList.value.map(lib => lib.id === id ? newLibraryData: lib)];
    currentLibraryList.value = newLibraries;
  } catch (error) {
    const { status, data } = convertError(error);
    alertError(status, data);
  } finally {
    onDone();
  }
}
