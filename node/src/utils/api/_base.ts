import Axios from 'axios';

export const axios = Axios.create({
  baseURL: '/api/v1'
});

axios.interceptors.response.use((response) => {
  if (import.meta.env.DEV) {
    return response;
  }
  if (typeof response.data === 'string') {
    if (response.data.startsWith('<!DOCTYPE html')) {
      const url = new URL(location.href);
      if (!url.pathname.startsWith('/sign-')) {
        location.href = `/sign-in?redirect=${encodeURIComponent(url.pathname + url.search)}`
      }
    }
  }
  return response;
});

export const convertError = (error: unknown) => {
  if (Axios.isAxiosError(error) && error.response) {
    const { response } = error;
    const { status, data } = response;
    return { status, data };
  } else {
    return { status: 500, data: 'Internal Server Error' };
  }
}

export const alertError = (statusCode: number, data: unknown) => {
  alert(`Response Status ${statusCode}\n\n${typeof data === 'string' ? data : JSON.stringify(data)}`);
}
