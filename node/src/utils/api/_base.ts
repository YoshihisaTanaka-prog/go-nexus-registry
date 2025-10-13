import Axios from 'axios';

export const axios = Axios.create({
  baseURL: '/api/v1'
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
