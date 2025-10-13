import { ref } from "vue";

type ToastType = "success"|"failed"|"error";

type ToastData = {
  type: ToastType;
  message: string;
  id: string;
}

export const toasts = ref<ToastData[]>([]);

type ToastProps = [
  type: ToastType,
  message: string,
  seconds?: number,
]

const toast = (...[type, message, seconds = 5]: ToastProps) => {
  const id = `${Date.now()}-${`${Math.random()}`.slice(2)}`;
  toasts.value = [...toasts.value, { type, message, id}];
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id != id);
  }, seconds * 1000);
}

export default toast;
