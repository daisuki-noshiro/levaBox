<script setup lang="ts">
import { ref } from 'vue'

const backgroundMusic = ref(true)
const volume = ref(65)
const soundEffects = ref(true)
const fullscreen = ref(false)
const gamepad = ref(true)

const settings = [
  { key: 'music', title: '背景音乐', description: '浏览大厅时播放你选择的音乐。', model: backgroundMusic },
  { key: 'sound', title: '按键音效', description: '菜单操作和焦点切换时播放反馈音。', model: soundEffects },
  { key: 'fullscreen', title: '全屏显示', description: '启动后占满掌机屏幕。', model: fullscreen },
  { key: 'gamepad', title: '启用手柄', description: '允许使用控制器浏览界面。', model: gamepad },
]
</script>

<template>
  <main class="settings page-shell">
    <header class="page-header"><div><p class="page-kicker">PREFERENCES</p><h1>设置</h1><p>调整 levaBox 的显示与控制偏好。本阶段刷新后会恢复默认值。</p></div></header>
    <div class="settings__layout">
      <section class="settings__group">
        <div class="settings__group-title"><span>♫</span><div><h2>声音</h2><p>声音功能将在后续阶段真正接入。</p></div></div>
        <div v-for="setting in settings.slice(0, 2)" :key="setting.key" class="setting-row">
          <div><strong>{{ setting.title }}</strong><span>{{ setting.description }}</span></div>
          <label class="toggle"><input v-model="setting.model.value" type="checkbox" /><span></span></label>
        </div>
        <label class="volume-row"><div><strong>音乐音量</strong><span>{{ volume }}%</span></div><input v-model="volume" type="range" min="0" max="100" /></label>
      </section>
      <section class="settings__group">
        <div class="settings__group-title"><span>▣</span><div><h2>显示与控制</h2><p>为横屏掌机预留的大尺寸控制项。</p></div></div>
        <div v-for="setting in settings.slice(2)" :key="setting.key" class="setting-row">
          <div><strong>{{ setting.title }}</strong><span>{{ setting.description }}</span></div>
          <label class="toggle"><input v-model="setting.model.value" type="checkbox" /><span></span></label>
        </div>
        <div class="setting-row setting-row--disabled"><div><strong>按键映射</strong><span>自定义手柄按键布局。</span></div><em>开发中</em></div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.settings { overflow-y: auto; }
.settings__layout { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 20px; margin-top: 30px; padding-bottom: 40px; }
.settings__group { align-self: start; overflow: hidden; border: 1px solid var(--line); border-radius: 22px; background: rgba(14, 22, 37, .8); box-shadow: var(--shadow); }
.settings__group-title { display: flex; gap: 14px; padding: 21px 23px; border-bottom: 1px solid var(--line); background: linear-gradient(120deg, rgba(103, 205, 255, .09), rgba(151, 117, 255, .06)); }
.settings__group-title > span { display: grid; place-items: center; width: 40px; height: 40px; border-radius: 12px; color: #9ceaff; background: rgba(100, 207, 255, .13); }
.settings__group-title h2 { margin: 0 0 4px; font-size: 1rem; }
.settings__group-title p { margin: 0; color: var(--muted); font-size: .74rem; }
.setting-row { display: flex; align-items: center; justify-content: space-between; min-height: 74px; padding: 0 23px; border-bottom: 1px solid rgba(255, 255, 255, .06); }
.setting-row > div { display: grid; gap: 5px; }
.setting-row strong, .volume-row strong { font-size: .88rem; }
.setting-row > div span { color: var(--muted); font-size: .72rem; }
.setting-row em { padding: 5px 9px; border-radius: 7px; color: #e5c684; background: rgba(227, 171, 74, .1); font-size: .65rem; font-style: normal; }
.setting-row--disabled { opacity: .72; }
.toggle { position: relative; flex: 0 0 auto; width: 50px; height: 28px; }
.toggle input { position: absolute; opacity: 0; }
.toggle span { position: absolute; inset: 0; border-radius: 99px; background: #2d394d; cursor: pointer; transition: background 160ms ease; }
.toggle span::after { content: ''; position: absolute; width: 20px; height: 20px; top: 4px; left: 4px; border-radius: 50%; background: #aab5c5; transition: transform 160ms ease, background 160ms ease; }
.toggle input:checked + span { background: linear-gradient(90deg, #348fbd, #776cc6); }
.toggle input:checked + span::after { transform: translateX(22px); background: #fff; }
.toggle input:focus-visible + span { outline: 2px solid #8ae6ff; outline-offset: 3px; }
.volume-row { display: grid; gap: 14px; padding: 19px 23px 23px; }
.volume-row div { display: flex; justify-content: space-between; }
.volume-row div span { color: #9de9ff; font-size: .78rem; }
.volume-row input { width: 100%; accent-color: #72dfff; }
@media (max-width: 900px) { .settings__layout { grid-template-columns: 1fr; } }
</style>
