<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { applyDisplayPreset } from '../services/windowService'
import { settingsState } from '../state/settingsStore'
import {
  DEFAULT_DISPLAY_PRESET_ID,
  displayPresets,
  findDisplayPreset,
  getPresetAspectRatio,
  getPresetResolution,
  type DisplayPreset,
} from '../types/display'

const pendingPresetId = ref(settingsState.displayPresetId)
const viewportWidth = ref(window.innerWidth)
const viewportHeight = ref(window.innerHeight)
const devicePixelRatio = ref(formatDevicePixelRatio())
const applying = ref(false)
const resultMessage = ref('选择模式后点击应用，窗口才会发生变化。')
const resultKind = ref<'idle' | 'success' | 'error'>('idle')
let viewportUpdateTimer: ReturnType<typeof setTimeout> | undefined

const fullscreenSelected = computed({
  get: () => findDisplayPreset(pendingPresetId.value).fullscreen === true,
  set: (enabled: boolean) => {
    pendingPresetId.value = enabled ? 'fullscreen' : DEFAULT_DISPLAY_PRESET_ID
  },
})

function selectPreset(preset: DisplayPreset): void {
  pendingPresetId.value = preset.id
}

function updateViewportInfo(): void {
  viewportWidth.value = window.innerWidth
  viewportHeight.value = window.innerHeight
  devicePixelRatio.value = formatDevicePixelRatio()
}

function formatDevicePixelRatio(): number {
  return Number(window.devicePixelRatio.toFixed(2))
}

async function applySelectedPreset(): Promise<void> {
  const preset = findDisplayPreset(pendingPresetId.value)
  applying.value = true
  resultKind.value = 'idle'
  const result = await applyDisplayPreset(preset)
  applying.value = false
  resultMessage.value = result.message
  resultKind.value = result.ok ? 'success' : 'error'

  if (result.ok) {
    settingsState.displayPresetId = preset.id
    settingsState.fullscreen = result.fullscreen === true
    if (viewportUpdateTimer) clearTimeout(viewportUpdateTimer)
    viewportUpdateTimer = setTimeout(updateViewportInfo, 180)
  }
}

onMounted(() => {
  updateViewportInfo()
  window.addEventListener('resize', updateViewportInfo)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateViewportInfo)
  if (viewportUpdateTimer) clearTimeout(viewportUpdateTimer)
})
</script>

<template>
  <main class="settings page-shell">
    <header class="page-header">
      <div>
        <p class="page-kicker">PREFERENCES</p>
        <h1>设置</h1>
        <p>调整 levaBox 的显示、声音与控制偏好，设置会保存在当前设备。</p>
      </div>
    </header>

    <div class="settings__layout">
      <section class="settings__group settings__group--display">
        <div class="settings__group-title">
          <span>▣</span>
          <div><h2>显示模式</h2><p>真实调整 Wails 窗口尺寸，用于测试常见掌机屏幕。</p></div>
        </div>

        <div class="display-presets" role="radiogroup" aria-label="显示模式">
          <button
            v-for="preset in displayPresets"
            :key="preset.id"
            type="button"
            class="display-preset"
            :class="{ 'display-preset--selected': pendingPresetId === preset.id }"
            role="radio"
            :aria-checked="pendingPresetId === preset.id"
            @click="selectPreset(preset)"
          >
            <span class="display-preset__indicator"></span>
            <span class="display-preset__copy">
              <strong>{{ preset.name }}</strong>
              <small>{{ getPresetResolution(preset) }}</small>
            </span>
            <em>{{ getPresetAspectRatio(preset) }}</em>
          </button>
        </div>

        <div class="display-actions">
          <div class="viewport-readout" aria-live="polite">
            <span>当前内容区域：<b>{{ viewportWidth }} × {{ viewportHeight }}</b></span>
            <span>设备像素比例：<b>{{ devicePixelRatio }}</b></span>
          </div>
          <button class="apply-display-button" type="button" :disabled="applying" @click="applySelectedPreset">
            {{ applying ? '正在应用…' : '应用显示设置' }}
          </button>
        </div>
        <p class="display-result" :class="`display-result--${resultKind}`" role="status">{{ resultMessage }}</p>
      </section>

      <section class="settings__group">
        <div class="settings__group-title">
          <span>♫</span>
          <div><h2>声音</h2><p>声音播放将在后续阶段接入，偏好状态现已保存。</p></div>
        </div>
        <div class="setting-row">
          <div><strong>背景音乐</strong><span>浏览大厅时播放你选择的音乐。</span></div>
          <label class="toggle"><input v-model="settingsState.backgroundMusic" type="checkbox" aria-label="背景音乐" /><span></span></label>
        </div>
        <div class="setting-row">
          <div><strong>按键音效</strong><span>菜单操作和焦点切换时播放反馈音。</span></div>
          <label class="toggle"><input v-model="settingsState.soundEffects" type="checkbox" aria-label="按键音效" /><span></span></label>
        </div>
        <label class="volume-row">
          <div><strong>音乐音量</strong><span>{{ settingsState.musicVolume }}%</span></div>
          <input v-model="settingsState.musicVolume" type="range" min="0" max="100" aria-label="音乐音量" />
        </label>
      </section>

      <section class="settings__group">
        <div class="settings__group-title">
          <span>◇</span>
          <div><h2>显示与控制</h2><p>适合触摸、键盘与后续手柄焦点操作。</p></div>
        </div>
        <div class="setting-row">
          <div><strong>全屏显示</strong><span>选中后点击“应用显示设置”进入全屏。</span></div>
          <label class="toggle"><input v-model="fullscreenSelected" type="checkbox" aria-label="全屏显示" /><span></span></label>
        </div>
        <div class="setting-row">
          <div><strong>启用手柄</strong><span>允许使用控制器浏览界面。</span></div>
          <label class="toggle"><input v-model="settingsState.gamepadEnabled" type="checkbox" aria-label="启用手柄" /><span></span></label>
        </div>
        <div class="setting-row setting-row--disabled">
          <div><strong>按键映射</strong><span>自定义手柄按键布局。</span></div><em>开发中</em>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.settings { overflow-y: auto; }
.settings__layout { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 20px; margin-top: 26px; padding-bottom: 40px; }
.settings__group { align-self: start; overflow: hidden; border: 1px solid var(--line); border-radius: 22px; background: rgba(14, 22, 37, .8); box-shadow: var(--shadow); }
.settings__group--display { grid-column: 1 / -1; }
.settings__group-title { display: flex; gap: 14px; padding: 19px 23px; border-bottom: 1px solid var(--line); background: linear-gradient(120deg, rgba(103, 205, 255, .09), rgba(151, 117, 255, .06)); }
.settings__group-title > span { display: grid; place-items: center; flex: 0 0 auto; width: 40px; height: 40px; border-radius: 12px; color: #9ceaff; background: rgba(100, 207, 255, .13); }
.settings__group-title h2 { margin: 0 0 4px; font-size: 1rem; }
.settings__group-title p { margin: 0; color: var(--muted); font-size: .74rem; }
.display-presets { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 11px; padding: 18px 22px; }
.display-preset { display: flex; align-items: center; gap: 12px; min-height: 72px; padding: 12px 14px; border: 1px solid rgba(255, 255, 255, .09); border-radius: 13px; color: #fff; background: rgba(255, 255, 255, .035); text-align: left; cursor: pointer; transition: transform 150ms ease, border-color 150ms ease, background 150ms ease; }
.display-preset:hover { transform: translateY(-2px); border-color: rgba(129, 224, 255, .35); background: rgba(99, 190, 224, .08); }
.display-preset:focus-visible, .apply-display-button:focus-visible { outline: 2px solid #8ae6ff; outline-offset: 3px; }
.display-preset--selected { border-color: rgba(119, 223, 255, .58); background: linear-gradient(135deg, rgba(74, 174, 216, .17), rgba(115, 99, 194, .13)); box-shadow: inset 0 0 0 1px rgba(143, 230, 255, .08); }
.display-preset__indicator { flex: 0 0 auto; width: 13px; height: 13px; border: 2px solid rgba(255, 255, 255, .4); border-radius: 50%; }
.display-preset--selected .display-preset__indicator { border: 4px solid #92e8ff; background: #0c1a28; }
.display-preset__copy { display: grid; min-width: 0; gap: 4px; }
.display-preset__copy strong { overflow: hidden; font-size: .82rem; white-space: nowrap; text-overflow: ellipsis; }
.display-preset__copy small { color: var(--muted); font-size: .7rem; font-variant-numeric: tabular-nums; }
.display-preset em { margin-left: auto; padding: 4px 7px; border-radius: 6px; color: #9ee8fb; background: rgba(105, 211, 241, .08); font-size: .62rem; font-style: normal; }
.display-actions { display: flex; align-items: center; justify-content: space-between; gap: 22px; padding: 1px 22px 14px; }
.viewport-readout { display: flex; flex-wrap: wrap; gap: 8px 22px; color: var(--muted); font-size: .72rem; }
.viewport-readout b { color: #eaf8ff; font-weight: 700; font-variant-numeric: tabular-nums; }
.apply-display-button { min-width: 190px; min-height: 50px; padding: 0 21px; border: 1px solid rgba(255, 255, 255, .24); border-radius: 12px; color: #07131d; background: linear-gradient(135deg, #e7faff, #83def7); font-weight: 800; cursor: pointer; }
.apply-display-button:disabled { cursor: wait; opacity: .6; }
.display-result { min-height: 18px; margin: 0; padding: 0 22px 18px; color: var(--muted); font-size: .7rem; }
.display-result--success { color: #9cebc8; }
.display-result--error { color: #ffb4b4; }
.setting-row { display: flex; align-items: center; justify-content: space-between; min-height: 74px; padding: 0 23px; border-bottom: 1px solid rgba(255, 255, 255, .06); }
.setting-row > div { display: grid; gap: 5px; }
.setting-row strong, .volume-row strong { font-size: .88rem; }
.setting-row > div span { color: var(--muted); font-size: .72rem; }
.setting-row em { padding: 5px 9px; border-radius: 7px; color: #e5c684; background: rgba(227, 171, 74, .1); font-size: .65rem; font-style: normal; }
.setting-row--disabled { opacity: .72; }
.toggle { position: relative; flex: 0 0 auto; width: 50px; height: 28px; }
.toggle input { position: absolute; opacity: 0; }
.toggle span { position: absolute; inset: 0; border-radius: 99px; background: #2d394d; cursor: pointer; transition: background 160ms ease; }
.toggle span::after { content: ''; position: absolute; top: 4px; left: 4px; width: 20px; height: 20px; border-radius: 50%; background: #aab5c5; transition: transform 160ms ease, background 160ms ease; }
.toggle input:checked + span { background: linear-gradient(90deg, #348fbd, #776cc6); }
.toggle input:checked + span::after { transform: translateX(22px); background: #fff; }
.toggle input:focus-visible + span { outline: 2px solid #8ae6ff; outline-offset: 3px; }
.volume-row { display: grid; gap: 14px; padding: 19px 23px 23px; }
.volume-row div { display: flex; justify-content: space-between; }
.volume-row div span { color: #9de9ff; font-size: .78rem; }
.volume-row input { width: 100%; min-height: 28px; accent-color: #72dfff; }
@media (max-width: 1100px) { .display-presets { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 900px) {
  .settings__layout { grid-template-columns: 1fr; }
  .display-presets { grid-template-columns: 1fr; }
  .display-actions { align-items: stretch; flex-direction: column; }
  .apply-display-button { width: 100%; }
}
@media (max-height: 760px) {
  .settings__layout { margin-top: 18px; }
  .settings__group-title { padding-top: 15px; padding-bottom: 15px; }
  .display-preset { min-height: 62px; }
}
</style>
