<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import {
  DEFAULT_PREVIEW_DEVICE_ID,
  findPreviewDevice,
  getPreviewAspectRatio,
  previewDevices,
} from '../types/devicePreview'

type ScaleMode = 'fit' | 'physical' | '50' | '75' | '100'

interface PreviewPreferences {
  deviceId: string
  scaleMode: ScaleMode
  showFrame: boolean
  showInfo: boolean
}

const STORAGE_KEY = 'levabox.device-preview.v1'
const defaultPreferences: PreviewPreferences = {
  deviceId: DEFAULT_PREVIEW_DEVICE_ID,
  scaleMode: 'fit',
  showFrame: true,
  showInfo: true,
}

const stored = readPreferences()
const deviceId = ref(stored.deviceId)
const scaleMode = ref<ScaleMode>(stored.scaleMode)
const showFrame = ref(stored.showFrame)
const showInfo = ref(stored.showInfo)
const fitScale = ref(1)
const devicePixelRatio = ref(readDevicePixelRatio())
const workspace = ref<HTMLElement | null>(null)
let resizeObserver: ResizeObserver | undefined

const device = computed(() => findPreviewDevice(deviceId.value))
const previewScale = computed(() => {
  if (scaleMode.value === 'fit') return fitScale.value
  if (scaleMode.value === 'physical') return 1 / devicePixelRatio.value
  return Number(scaleMode.value) / 100
})
const framePadding = computed(() => showFrame.value ? 18 : 0)
const visualWidth = computed(() => device.value.width * previewScale.value + framePadding.value * 2)
const visualHeight = computed(() => device.value.height * previewScale.value + framePadding.value * 2)
const scaleLabel = computed(() => `${(previewScale.value * 100).toFixed(2)}%`)

function setScaleMode(mode: ScaleMode): void {
  scaleMode.value = mode
  if (mode === 'fit') updateFitScale()
}

function updateFitScale(): void {
  const container = workspace.value
  if (!container) return
  devicePixelRatio.value = readDevicePixelRatio()
  const safeMargin = 64
  const infoHeight = showInfo.value ? 34 : 0
  const frameSpace = framePadding.value * 2
  const availableWidth = Math.max(1, container.clientWidth - safeMargin - frameSpace)
  const availableHeight = Math.max(1, container.clientHeight - safeMargin - frameSpace - infoHeight)
  fitScale.value = Math.max(.1, Math.min(
    availableWidth / device.value.width,
    availableHeight / device.value.height,
  ))
}

function readDevicePixelRatio(): number {
  return Math.max(.1, window.devicePixelRatio || 1)
}

function readPreferences(): PreviewPreferences {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return { ...defaultPreferences }
    const value = JSON.parse(raw) as Partial<PreviewPreferences>
    const modes: ScaleMode[] = ['fit', 'physical', '50', '75', '100']
    return {
      deviceId: previewDevices.some((item) => item.id === value.deviceId)
        ? value.deviceId!
        : defaultPreferences.deviceId,
      scaleMode: modes.includes(value.scaleMode as ScaleMode)
        ? value.scaleMode as ScaleMode
        : defaultPreferences.scaleMode,
      showFrame: typeof value.showFrame === 'boolean' ? value.showFrame : defaultPreferences.showFrame,
      showInfo: typeof value.showInfo === 'boolean' ? value.showInfo : defaultPreferences.showInfo,
    }
  } catch {
    return { ...defaultPreferences }
  }
}

function persistPreferences(): void {
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify({
      deviceId: deviceId.value,
      scaleMode: scaleMode.value,
      showFrame: showFrame.value,
      showInfo: showInfo.value,
    } satisfies PreviewPreferences))
  } catch {
    // The preview remains usable when localStorage is unavailable.
  }
}

function handleWindowResize(): void {
  devicePixelRatio.value = readDevicePixelRatio()
  if (scaleMode.value === 'fit') updateFitScale()
}

watch([deviceId, scaleMode, showFrame, showInfo], () => {
  persistPreferences()
  nextTick(updateFitScale)
})

onMounted(() => {
  resizeObserver = new ResizeObserver(updateFitScale)
  if (workspace.value) resizeObserver.observe(workspace.value)
  window.addEventListener('resize', handleWindowResize)
  updateFitScale()
})

onUnmounted(() => {
  resizeObserver?.disconnect()
  window.removeEventListener('resize', handleWindowResize)
})
</script>

<template>
  <div class="device-preview-shell">
    <header class="device-preview-toolbar" @keydown.stop>
      <div class="device-preview-toolbar__brand">
        <strong>DEVICE PREVIEW</strong>
        <span>开发测试工具</span>
      </div>

      <label class="device-picker">
        <span>设备</span>
        <select v-model="deviceId" aria-label="设备预设">
          <option v-for="item in previewDevices" :key="item.id" :value="item.id">
            {{ item.name }} · {{ item.width }} × {{ item.height }}
          </option>
        </select>
      </label>

      <div class="scale-modes" role="group" aria-label="预览缩放">
        <button type="button" :class="{ active: scaleMode === 'fit' }" @click="setScaleMode('fit')">适应窗口</button>
        <button type="button" :class="{ active: scaleMode === 'physical' }" @click="setScaleMode('physical')">物理像素 1:1</button>
        <button v-for="value in (['50', '75', '100'] as const)" :key="value" type="button" :class="{ active: scaleMode === value }" @click="setScaleMode(value)">{{ value }}%</button>
      </div>

      <label class="preview-check"><input v-model="showFrame" type="checkbox" />掌机外框</label>
      <label class="preview-check"><input v-model="showInfo" type="checkbox" />尺寸信息</label>
      <div class="device-preview-toolbar__metrics">
        <span>DPR {{ devicePixelRatio.toFixed(2) }}</span>
        <b>{{ scaleLabel }}</b>
      </div>
    </header>

    <main ref="workspace" class="device-preview-workspace">
      <div class="device-preview-canvas" :style="{ width: `${visualWidth}px` }">
        <div
          class="device-frame"
          :class="{ 'device-frame--hidden': !showFrame }"
          :style="{ width: `${visualWidth}px`, height: `${visualHeight}px`, padding: `${framePadding}px` }"
        >
          <div
            class="device-stage"
            :style="{
              width: `${device.width}px`,
              height: `${device.height}px`,
              transform: `scale(${previewScale})`,
            }"
          >
            <slot />
          </div>
        </div>

        <p v-if="showInfo" class="device-preview-info">
          <strong>{{ device.width }} × {{ device.height }}</strong>
          <span>{{ getPreviewAspectRatio(device) }}</span>
          <span>缩放 {{ scaleLabel }}</span>
          <span>DPR {{ devicePixelRatio.toFixed(2) }}</span>
        </p>
      </div>
    </main>
  </div>
</template>

<style scoped>
.device-preview-shell { display: grid; grid-template-rows: auto minmax(0, 1fr); width: 100%; height: 100%; overflow: hidden; color: #e8edf4; background: #24272d; }
.device-preview-toolbar { position: relative; z-index: 200; display: flex; align-items: center; gap: 12px; min-height: 68px; padding: 10px 16px; border-bottom: 1px solid #3d4149; background: #191b20; box-shadow: 0 4px 18px rgba(0, 0, 0, .22); font-family: 'Segoe UI', sans-serif; }
.device-preview-toolbar__brand { display: grid; flex: 0 0 auto; gap: 2px; padding-right: 12px; border-right: 1px solid #3a3d44; }
.device-preview-toolbar__brand strong { color: #9fe5f8; font-size: .7rem; letter-spacing: .12em; }
.device-preview-toolbar__brand span, .device-picker > span { color: #8f969f; font-size: .65rem; }
.device-picker { display: grid; flex: 0 1 270px; gap: 3px; }
.device-picker select { width: 100%; min-height: 36px; padding: 0 30px 0 10px; border: 1px solid #454a54; border-radius: 7px; color: #eef2f7; background: #292c33; }
.scale-modes { display: flex; gap: 5px; }
.scale-modes button { min-height: 36px; padding: 0 10px; border: 1px solid #414650; border-radius: 7px; color: #b9c0ca; background: #272a31; cursor: pointer; white-space: nowrap; }
.scale-modes button.active { border-color: #75cce5; color: #eefbff; background: #214553; }
.preview-check { display: flex; align-items: center; gap: 6px; color: #b9c0ca; font-size: .72rem; white-space: nowrap; }
.preview-check input { width: 16px; height: 16px; accent-color: #69c9e5; }
.device-preview-toolbar button:focus-visible, .device-preview-toolbar select:focus-visible, .device-preview-toolbar input:focus-visible { outline: 2px solid #8be2fa; outline-offset: 2px; }
.device-preview-toolbar__metrics { display: grid; flex: 0 0 auto; gap: 2px; margin-left: auto; text-align: right; font-size: .68rem; font-variant-numeric: tabular-nums; }
.device-preview-toolbar__metrics span { color: #8f969f; }
.device-preview-toolbar__metrics b { color: #dff9ff; }
.device-preview-workspace { min-width: 0; min-height: 0; overflow: auto; padding: 32px; background: radial-gradient(circle at 50% 35%, #343943, #22252a 70%); }
.device-preview-canvas { margin: 0 auto; }
.device-frame { border-radius: 24px; background: #0d0f12; box-shadow: 0 22px 60px rgba(0, 0, 0, .42); }
.device-frame--hidden { border-radius: 0; background: transparent; box-shadow: none; }
.device-stage { position: relative; overflow: hidden; border-radius: 10px; transform-origin: top left; background: #080d17; }
.device-frame--hidden .device-stage { border-radius: 0; }
.device-preview-info { display: flex; justify-content: center; gap: 14px; margin: 10px 0 0; color: #9da4ae; font-family: 'Segoe UI', sans-serif; font-size: .7rem; font-variant-numeric: tabular-nums; }
.device-preview-info strong { color: #e9eef5; }
@media (max-width: 1180px) {
  .device-preview-toolbar { align-items: stretch; flex-wrap: wrap; }
  .device-preview-toolbar__metrics { margin-left: 0; }
}
@media (prefers-reduced-motion: reduce) {
  .scale-modes button { transition: none; }
}
</style>
