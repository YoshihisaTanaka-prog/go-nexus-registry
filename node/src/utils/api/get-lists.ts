import { type Ref, ref } from 'vue'
import { axios, convertError, alertError } from './_base'

export type Library = {
  id: string;
  name: string;
  version: string;
  status: 'uploading' | 'uploaded' | 'failed';
  isPublished: boolean;
}

type Params = {
  readonly cursor?: Cursor
  readonly limit: number;
}

type Cursor = {
  readonly name: string;
  readonly v1: number;
  readonly v2: number;
  readonly v3: number;
}

type PaginatedLibrariesResponse = {
  data: Library[];
  cursor?: Cursor;
  limit: number;
}

async function getLibrariesUnit(params: Params): Promise<PaginatedLibrariesResponse> {
  try {
    let url = 'get-libraries?limit=' + params.limit
    if (params.cursor) {
      const { name, v1, v2, v3 } = params.cursor;
      url += `&name=${encodeURI(name)}&v1=${v1}&v2=${v2}&v3=${v3}`
    }
    const result = await axios.get(url);
    return result.data;
  } catch (error) {
    const { status, data } = convertError(error);
    alertError(status, data);
    return { data: [], limit: 0 }
  }
}

export async function getLibraries(currentLibraryList: Ref<Library[]>) {
  const newLibraryList = currentLibraryList.value.length === 0 ? currentLibraryList : ref<Library[]>([]);
  let cursor: Cursor | undefined = undefined;
  let limit: number = 10;
  
  do {
    const { data: libraries, ...cursorData } = await getLibrariesUnit({cursor, limit});
    const currentLibs = currentLibraryList.value;
    currentLibraryList.value = [...currentLibs, ...libraries ]
    cursor = cursorData.cursor;
    limit = cursorData.limit;
  } while (cursor !== undefined);

  currentLibraryList.value = [...newLibraryList.value];
}
