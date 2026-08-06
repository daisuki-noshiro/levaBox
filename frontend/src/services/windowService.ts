import {
  WindowCenter,
  WindowFullscreen,
  WindowIsFullscreen,
  WindowSetSize,
  WindowUnfullscreen,
} from '../../wailsjs/runtime/runtime'
import type { DisplayPreset } from '../types/display'

export interface WindowOperationResult {
  ok: boolean
  message: string
  fullscreen?: boolean
}

type WailsRuntimeWindow = Window & {
  runtime?: {
    WindowSetSize?: (width: number, height: number) => void
    WindowIsFullscreen?: () => Promise<boolean>
  }
}

export function isWailsWindow(): boolean {
  if (typeof window === 'undefined') return false
  const runtimeWindow = window as WailsRuntimeWindow
  return typeof runtimeWindow.runtime?.WindowSetSize === 'function'
    && typeof runtimeWindow.runtime?.WindowIsFullscreen === 'function'
}

export async function applyDisplayPreset(preset: DisplayPreset): Promise<WindowOperationResult> {
  if (!isWailsWindow()) {
    return { ok: false, message: '当前不是 Wails 窗口，无法调整应用尺寸。' }
  }

  try {
    const currentlyFullscreen = await WindowIsFullscreen()
    if (preset.fullscreen) {
      if (!currentlyFullscreen) WindowFullscreen()
      return { ok: true, fullscreen: true, message: '已切换为全屏模式。' }
    }

    if (!preset.width || !preset.height) {
      return { ok: false, message: '显示预设缺少有效的窗口尺寸。' }
    }

    if (currentlyFullscreen) WindowUnfullscreen()
    WindowSetSize(preset.width, preset.height)
    WindowCenter()
    return {
      ok: true,
      fullscreen: false,
      message: `已应用 ${preset.width} × ${preset.height}，窗口已居中。`,
    }
  } catch (error) {
    const reason = error instanceof Error ? error.message : String(error)
    return { ok: false, message: `调整窗口失败：${reason}` }
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
