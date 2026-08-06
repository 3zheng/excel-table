import { onMounted, onUnmounted } from 'vue'

/**
 * 监听 Ctrl+S 快捷键，触发保存回调
 * 在组件 setup() 中调用，自动管理事件监听的挂载和卸载
 *
 * @param callback 按下 Ctrl+S 时执行的回调函数
 *
 * @example
 * useCtrlS(() => {
 *   if (canEdit.value || canReview.value) {
 *     saveInvoice()
 *   }
 * })
 */
// 定义：接收一个函数作为参数，参数名叫 callback， 类型是() => void函数
export function useCtrlS(callback: () => void) {
  const handleKeydown = (e: KeyboardEvent) => {
    if (e.ctrlKey && e.key === 's') {
      e.preventDefault() //拦截原来的ctrl+s的保存网页行为
      callback()  // 按下 Ctrl+S 时，调用传进来的函数
    }
  }

  onMounted(() => window.addEventListener('keydown', handleKeydown))
  onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
}
