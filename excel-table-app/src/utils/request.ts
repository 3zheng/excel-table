import router from '@/router'
import { ElMessage } from 'element-plus'
import i18n from '@/i18n'

const t = i18n.global.t
const TIMEOUT_MS = 10000 // 10秒超时

export async function request(url: string, options: RequestInit = {}): Promise<Response> {
  const token = localStorage.getItem('token')

  //超时计时器
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), TIMEOUT_MS)

  let res: Response
  try {
    res = await fetch(url, {
      ...options,
      signal: controller.signal,  //超时信号
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    })
  } catch (err: any) {
    if (err.name === 'AbortError') {
      ElMessage.error('Request timeout, please try again')
    } else {
      ElMessage.error(t('errNetwork'))
    }
    throw err
  } finally {
    clearTimeout(timeout)
  }

  // 401：token 失效，直接跳登录页，无需弹窗
  if (res.status === 401) {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
    throw new Error('401')
  }

  if (res.status === 403) {
    ElMessage.error(t('errForbidden'))
    throw new Error('403')
  }

  if (res.status === 500) {
    ElMessage.error(t('errServer'))
    throw new Error('500')
  }

  if (!res.ok) {
    try {
      const data = await res.json()
      ElMessage.error(data.error || `${t('errUnknown')} (${res.status})`)
    } catch {
      ElMessage.error(`${t('errUnknown')} (${res.status})`)
    }
    throw new Error(`${res.status}`)
  }

  return res
}