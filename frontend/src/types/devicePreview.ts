export interface PreviewDevice {
  id: string
  name: string
  width: number
  height: number
  screenInches?: number
}

export const DEFAULT_PREVIEW_DEVICE_ID = 'handheld-1280x800'

export const previewDevices: readonly PreviewDevice[] = [
  { id: DEFAULT_PREVIEW_DEVICE_ID, name: '16:10 标准掌机', width: 1280, height: 800 },
  { id: 'handheld-1280x720', name: '16:9 性能掌机', width: 1280, height: 720 },
  { id: 'handheld-1920x1080', name: '16:9 高清掌机', width: 1920, height: 1080 },
  { id: 'handheld-1920x1200', name: '16:10 高清掌机', width: 1920, height: 1200 },
]

export function findPreviewDevice(id: string): PreviewDevice {
  return previewDevices.find((device) => device.id === id) ?? previewDevices[0]!
}

export function getPreviewAspectRatio(device: PreviewDevice): string {
  return device.width * 9 === device.height * 16 ? '16:9' : '16:10'
}
