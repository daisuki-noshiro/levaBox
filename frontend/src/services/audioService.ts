export type UiSoundType = 'select' | 'open' | 'confirm' | 'cancel' | 'launch'

const volumeState = {
  master: 1,
  ui: 1,
  background: 1,
}

/**
 * UI 声音统一入口。当前没有音频资源，因此保持为空操作。
 * 后续只需在这里接入 Web Audio 或 Wails 提供的本地资源。
 */
export function playUiSound(_type: UiSoundType): void {
  // Intentionally empty during the prototype stage.
}

export function setMasterVolume(value: number): void {
  volumeState.master = clampVolume(value)
}

export function setUiVolume(value: number): void {
  volumeState.ui = clampVolume(value)
}

export function setBackgroundVolume(value: number): void {
  volumeState.background = clampVolume(value)
}

/** 预留首次用户交互后解锁视频原声或背景音乐的位置。 */
export function unlockMediaAudio(): void {
  // Keep background videos muted until real audio settings are implemented.
}

function clampVolume(value: number): number {
  return Math.min(1, Math.max(0, value))
}
