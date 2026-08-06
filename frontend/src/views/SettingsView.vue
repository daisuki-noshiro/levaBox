<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getFullscreenState, setFullscreen } from '../services/windowService'
import { settingsState } from '../state/settingsStore'

const applyingFullscreen = ref(false)
const fullscreenMessage = ref('')
const fullscreenMessageKind = ref<'idle' | 'success' | 'error'>('idle')

async function handleFullscreenChange(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const requested = input.checked
  applyingFullscreen.value = true
  fullscreenMessageKind.value = 'idle'
  const result = await setFullscreen(requested)
  applyingFullscreen.value = false
  fullscreenMessage.value = result.message
  fullscreenMessageKind.value = result.ok ? 'success' : 'error'
  settingsState.fullscreen = result.ok ? requested : !requested
  if (!result.ok) input.checked = settingsState.fullscreen
}

onMounted(async () => {
  const result = await getFullscreenState()
  if (result.ok) settingsState.fullscreen = result.fullscreen === true
})
</script>

<template>
  <main class="settings page-shell">
    <header class="page-header">
      <div>
        <p class="page-kicker">PREFERENCES</p>
        <h1>设置</h1>
        <p>调整 levaBox 的声音、显示与控制偏好，设置会保存在当前设备。</p>
      </div>
    </header>

    <div class="settings__layout">
      <section class="settings__group">
        <div class="settings__group-title">
          <span aria-hidden="true">♪</span>
          <div><h2>声音</h2><p>声音资源将在后续阶段接入，当前偏好状态会正常保存。</p></div>
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
          <span aria-hidden="true">◇</span>
          <div><h2>显示与控制</h2><p>适合触摸、键盘与后续手柄焦点操作。</p></div>
        </div>
        <div class="setting-row">
          <div><strong>全屏显示</strong><span>在 Wails 应用中切换全屏或窗口模式。</span></div>
          <label class="toggle" :class="{ 'toggle--busy': applyingFullscreen }">
            <input
              :checked="settingsState.fullscreen"
              :disabled="applyingFullscreen"
              type="checkbox"
              aria-label="全屏显示"
              @change="handleFullscreenChange"
            />
            <span></span>
          </label>
        </div>
        <p v-if="fullscreenMessage" class="setting-message" :class="`setting-message--${fullscreenMessageKind}`" role="status">{{ fullscreenMessage }}</p>
        <div class="setting-row">
          <div><strong>启用手柄</strong><span>允许使用控制器浏览界面。</span></div>
          <label class="toggle"><input v-model="settingsState.gamepadEnabled" type="checkbox" aria-label="启用手柄" /><span></span></label>
        </div>
        <div class="setting-row setting-row--disabled">
          <div><strong>界面缩放</strong><span>调整正式界面的文字与交互区域大小。</span></div><em>开发中</em>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.settings { overflow-y: auto; }
.settings__layout { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 20px; margin-top: 26px; padding-bottom: 40px; }
.settings__group { align-self: start; overflow: hidden; border: 1px solid var(--line); border-radius: 22px; background: rgba(14, 22, 37, .8); box-shadow: var(--shadow); }
.settings__group-title { display: flex; gap: 14px; padding: 19px 23px; border-bottom: 1px solid var(--line); background: linear-gradient(120deg, rgba(103, 205, 255, .09), rgba(151, 117, 255, .06)); }
.settings__group-title > span { display: grid; place-items: center; flex: 0 0 auto; width: 40px; height: 40px; border-radius: 12px; color: #9ceaff; background: rgba(100, 207, 255, .13); }
.settings__group-title h2 { margin: 0 0 4px; font-size: 1rem; }
.settings__group-title p { margin: 0; color: var(--muted); font-size: .74rem; }
.setting-row { display: flex; align-items: center; justify-content: space-between; min-height: 82px; padding: 0 23px; border-bottom: 1px solid rgba(255, 255, 255, .06); }
.setting-row > div { display: grid; gap: 5px; }
.setting-row strong, .volume-row strong { font-size: .88rem; }
.setting-row > div span { color: var(--muted); font-size: .72rem; }
.setting-row em { padding: 5px 9px; border-radius: 7px; color: #e5c684; background: rgba(227, 171, 74, .1); font-size: .65rem; font-style: normal; }
.setting-row--disabled { opacity: .72; }
.setting-message { margin: -1px 23px 0; padding: 10px 0; color: var(--muted); font-size: .7rem; }
.setting-message--success { color: #9cebc8; }
.setting-message--error { color: #ffb4b4; }
.toggle { position: relative; flex: 0 0 auto; width: 50px; height: 28px; }
.toggle input { position: absolute; opacity: 0; }
.toggle span { position: absolute; inset: 0; border-radius: 99px; background: #2d394d; cursor: pointer; transition: background 160ms ease; }
.toggle span::after { content: ''; position: absolute; top: 4px; left: 4px; width: 20px; height: 20px; border-radius: 50%; background: #aab5c5; transition: transform 160ms ease, background 160ms ease; }
.toggle input:checked + span { background: linear-gradient(90deg, #348fbd, #776cc6); }
.toggle input:checked + span::after { transform: translateX(22px); background: #fff; }
.toggle input:focus-visible + span { outline: 2px solid #8ae6ff; outline-offset: 3px; }
.toggle--busy { opacity: .55; }
.volume-row { display: grid; gap: 14px; padding: 19px 23px 23px; }
.volume-row div { display: flex; justify-content: space-between; }
.volume-row div span { color: #9de9ff; font-size: .78rem; }
.volume-row input { width: 100%; min-height: 28px; accent-color: #72dfff; }
@container levabox-app (max-width: 900px) {
  .settings__layout { grid-template-columns: 1fr; }
}
@container levabox-app (max-height: 760px) {
  .settings__layout { margin-top: 18px; }
  .settings__group-title { padding-top: 15px; padding-bottom: 15px; }
  .setting-row { min-height: 68px; }
}
</style>
