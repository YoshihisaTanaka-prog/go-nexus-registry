import { axios, convertError, alertError } from './_base'

const urlParams = new URLSearchParams(window.location.search);
const redirectTo = urlParams.get('redirect');

export const signUp = async (email: string, password: string) => {
  const onFailedSignUp = (statusCode: number, data: unknown) => {
    if ( statusCode == 409 ) {
      alert('そのユーザーは既に登録されています。');
      location.href = redirectTo ? `/sign-in?redirect=${encodeURIComponent(redirectTo)}` : '/sign-in';
    } else {
      alertError(statusCode, data);
    }
  }
  try {
    const result = await axios.post('sign-up', {email, password});
    if (result.status < 400) {
      location.href = redirectTo || '/';
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
      alert('パスワードが正しくありません。');
    } else {
      alertError(statusCode, data);
    }
  }
  try {
    const result = await axios.post('sign-in', {email, password});
    if (result.status < 400) {
      location.href =  redirectTo || '/';
    } else {
      onFailedSignIn(result.status, result.data)
    }
  } catch (error) {
    const { status, data } = convertError(error);
    onFailedSignIn(status, data)
  }
}
