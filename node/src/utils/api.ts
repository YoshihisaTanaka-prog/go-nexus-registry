import Axios from "axios";

const axios = Axios.create({
  baseURL: '/api/v1'
});

const convertError = (error: unknown) => {
  if (Axios.isAxiosError(error) && error.response) {
    const { response } = error;
    const { status, data } = response;
    return { status, data };
  } else {
    return { status: 500, data: "Internal Server Error" };
  }
}

const alertError = (statusCode: number, data: unknown) => {
  alert(`Response Status ${statusCode}\n\n${typeof data === 'string' ? data : JSON.stringify(data)}`);
}

export const signUp = async (email: string, password: string) => {
  const onFailedSignUp = (statusCode: number, data: unknown) => {
    if ( statusCode == 409 ) {
      alert("そのユーザーは既に登録されています。");
      location.href = "/sign-in";
    } else {
      alertError(statusCode, data);
    }
  }
  try {
    const result = await axios.post('sign-up', {email, password});
    if (result.status < 400) {
      location.href = "/";
    } else {
      onFailedSignUp(result.status, result.data)
    }
  } catch (error) {
    const { status, data } = convertError(error);
    onFailedSignUp(status, data)
  }
}

export const signIn = async (email: string, password: string) => {
  const onFailedSignIn = (statusCode: number, data: unknown) => {
    if ( statusCode == 401 ) {
      alert("パスワードが正しくありません。");
    } else {
      alertError(statusCode, data);
    }
  }
  try {
    const result = await axios.post('sign-in', {email, password});
    if (result.status < 400) {
      location.href = "/";
    } else {
      onFailedSignIn(result.status, result.data)
    }
  } catch (error) {
    const { status, data } = convertError(error);
    onFailedSignIn(status, data)
  }
}
