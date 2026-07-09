import { ref, readonly } from 'vue'

const toasts = ref([])
let id = 0

export function useToast() {
  function show(message, type = 'info', duration = 3000) {
    const tid = ++id
    toasts.value.push({ id: tid, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== tid)
    }, duration)
  }

  return { toasts: readonly(toasts), show }
}
