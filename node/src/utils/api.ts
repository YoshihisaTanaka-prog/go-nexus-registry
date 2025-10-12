import Axios from "axios";

const axios = Axios.create({
  baseURL: '/api/v1'
})

export const signUp = async (userId: string, password: string) => {
  const result = await axios.post('sign-up', {userId, password});
  if (result.status < 400) {
    location.href = "/";
  } else {
    alert(`Response Status ${result.status}\n\n${result.data}`)
  }
}

export const signIn = async (userId: string, password: string) => {
  const result = await axios.post('sign-in', {userId, password});
  if (result.status < 400) {
    location.href = "/";
  } else {
    alert(`Response Status ${result.status}\n\n${result.data}`)
  }
}
