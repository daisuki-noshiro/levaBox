import {
  WindowFullscreen,
  WindowIsFullscreen,
  WindowUnfullscreen,
} from '../../wailsjs/runtime/runtime'

export interface WindowOperationResult {
  ok: boolean
  message: string
  fullscreen?: boolean
}

type WailsRuntimeWindow = Window & {
  runtime?: {
    WindowFullscreen?: () => void
    WindowUnfullscreen?: () => void
    WindowIsFullscreen?: () => Promise<boolean>
  }
}

export function isWailsWindow(): boolean {
  if (typeof window === 'undefined') return false
  const runtime = (window as WailsRuntimeWindow).runtime
  return typeof runtime?.WindowFullscreen === 'function'
    && typeof runtime.WindowUnfullscreen === 'function'
    && typeof runtime.WindowIsFullscreen === 'function'
}

export async function setFullscreen(enabled: boolean): Promise<WindowOperationResult> {
  if (!isWailsWindow()) {
    return { ok: false, message: '当前不是 Wails 窗口，无法切换全屏模式。' }
  }

  try {
    const currentlyFullscreen = await WindowIsFullscreen()
    if (enabled && !currentlyFullscreen) WindowFullscreen()
    if (!enabled && currentlyFullscreen) WindowUnfullscreen()
    return {
      ok: true,
      fullscreen: enabled,
      message: enabled ? '已切换为全屏模式。' : '已切换为窗口模式。',
    }
  } catch (error) {
    const reason = error instanceof Error ? error.message : String(error)
    return { ok: false, message: `切换全屏模式失败：${reason}` }
  }
}

export async function getFullscreenState(): Promise<WindowOperationResult> {
  if (!isWailsWindow()) {
    return { ok: false, message: '当前不是 Wails 窗口，无法读取全屏状态。' }
  }
  try {
    const fullscreen = await WindowIsFullscreen()
    return { ok: true, fullscreen, message: fullscreen ? '当前为全屏模式。' : '当前为窗口模式。' }
  } catch (error) {
    const reason = error instanceof Error ? error.message : String(error)
    return { ok: false, message: `读取全屏状态失败：${reason}` }
  }
}
