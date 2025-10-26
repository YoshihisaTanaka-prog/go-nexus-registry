export type Library = {
  id: string;
  name: string;
  version: string;
  v1: number;
  v2: number;
  v3: number;
  status: 'uploading' | 'uploaded' | 'failed' | 'updating';
  isPublished: boolean;
}

export type PaginationParams = {
  readonly kind: string,
  readonly cursor?: Cursor
  readonly limit: number;
}

export type Cursor = {
  readonly name: string;
  readonly v1: number;
  readonly v2: number;
  readonly v3: number;
}

export type PaginatedLibrariesResponse = {
  data: Library[];
  cursor?: Cursor | null;
  limit: number;
}