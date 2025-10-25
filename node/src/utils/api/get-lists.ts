import { type Ref, ref } from 'vue'
import { axios, convertError, alertError } from './_base'

export type Library = {
  id: string;
  name: string;
  version: string;
  v1: number;
  v2: number;
  v3: number;
  status: 'uploading' | 'uploaded' | 'failed';
  isPublished: boolean;
}

type Params = {
  readonly kind: string,
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
  cursor?: Cursor | null;
  limit: number;
}

function sortByName(a: Library, b: Library) {
   if (a.name > b.name) {
    return 1;
   } else if (a.name < b.name) {
    return -1;
   } else {
    return 0;
   }
}

function sortByV1(a: Library, b: Library) {
   return a.v1 - b.v1
}

function sortByV2(a: Library, b: Library) {
   return a.v2 - b.v2
}

function sortByV3(a: Library, b: Library) {
   return a.v3 - b.v3
}

async function getLibrariesUnit(params: Params): Promise<PaginatedLibrariesResponse> {
  try {
    const { limit, kind, cursor } = params;
    let url = `get-libraries?limit=${limit}&kind=${kind}`
    if (cursor) {
      const { name, v1, v2, v3 } = cursor;
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

export async function getLibraries(kind: string, currentLibraryList: Ref<Library[]>) {
  currentLibraryList.value = [];
  let cursor: Cursor | null | undefined = undefined;
  let limit: number = 10;
  
  do {
    const { data: libraries, ...cursorData } = await getLibrariesUnit({kind, cursor, limit});
    const currentLibs = currentLibraryList.value;
    currentLibs.push(...libraries);
    currentLibs.sort((a,b) => {
      let res = sortByName(a, b);
      if (res === 0) {
        res = sortByV1(a,b);
        if (res === 0) {
          res = sortByV2(a,b);
          if (res === 0) {
            return sortByV3(a,b);
          }
        }
      }
      return res;
    });
    currentLibraryList.value = [...currentLibs]
    cursor = cursorData.cursor;
    limit = cursorData.limit;
    await new Promise((resolve) => setTimeout(resolve, 500))
  } while ((cursor != null) && limit > 0);
}
