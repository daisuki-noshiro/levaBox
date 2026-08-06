export interface DisplayPreset {
  id: string
  name: string
  width?: number
  height?: number
  fullscreen?: boolean
}

export const DEFAULT_DISPLAY_PRESET_ID = 'handheld-1280x800'

export const displayPresets: readonly DisplayPreset[] = [
  { id: DEFAULT_DISPLAY_PRESET_ID, name: '16:10 标准掌机', width: 1280, height: 800 },
  { id: 'performance-1280x720', name: '16:9 性能模式', width: 1280, height: 720 },
  { id: 'hd-1920x1080', name: '16:9 高清', width: 1920, height: 1080 },
  { id: 'hd-1920x1200', name: '16:10 高清', width: 1920, height: 1200 },
  { id: 'ultra-2560x1600', name: '16:10 超高清', width: 2560, height: 1600 },
  { id: 'fullscreen', name: '全屏模式', fullscreen: true },
]

export function findDisplayPreset(id: string): DisplayPreset {
  return displayPresets.find((preset) => preset.id === id) ?? displayPresets[0]!
}

export function getPresetResolution(preset: DisplayPreset): string {
  return preset.fullscreen ? '使用当前显示器' : `${preset.width} × ${preset.height}`
}

export function getPresetAspectRatio(preset: DisplayPreset): string {
  if (preset.fullscreen || !preset.width || !preset.height) return '自适应'
  if (preset.width * 9 === preset.height * 16) return '16:9'
  if (preset.width * 10 === preset.height * 16) return '16:10'
  return `${preset.width}:${preset.height}`
}
