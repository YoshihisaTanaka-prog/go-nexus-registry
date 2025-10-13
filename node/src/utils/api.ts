import Axios from "axios";

const axios = Axios.create({
  baseURL: '/api/v1'
})

export const signUp = async (email: string, password: string) => {
  const result = await axios.post('sign-up', {email, password});
  if (result.status < 400) {
    location.href = "/";
  } else {
    alert(`Response Status ${result.status}\n\n${result.data}`)
  }
}

export const signIn = async (email: string, password: string) => {
  const result = await axios.post('sign-in', {email, password});
  if (result.status < 400) {
    location.href = "/";
  } else {
    alert(`Response Status ${result.status}\n\n${result.data}`)
  }
}
