<script setup lang="ts">
import type { MainPage } from '../types/navigation'

defineProps<{ open: boolean; currentPage: MainPage }>()
defineEmits<{ close: []; navigate: [page: MainPage] }>()

const menuItems: Array<{ page: MainPage; label: string; icon: string; hint: string }> = [
  { page: 'home', label: '大厅', icon: '⌂', hint: '沉浸式浏览' },
  { page: 'library', label: '游戏库', icon: '▦', hint: '查看全部游戏' },
  { page: 'import', label: '导入游戏', icon: '＋', hint: '添加本地内容' },
  { page: 'settings', label: '设置', icon: '⚙', hint: '偏好与控制' },
]
</script>

<template>
  <Transition name="fade">
    <button v-if="open" class="sidebar-overlay" aria-label="关闭菜单" @click="$emit('close')"></button>
  </Transition>
  <aside class="sidebar" :class="{ 'sidebar--open': open }" :aria-hidden="!open">
    <header class="sidebar__header">
      <div class="sidebar__brand">
        <span class="sidebar__brand-mark">L</span>
        <span><strong>levaBox</strong><small>GALGAME LIBRARY</small></span>
      </div>
      <button class="sidebar__close" aria-label="关闭菜单" @click="$emit('close')">×</button>
    </header>
    <nav class="sidebar__nav" aria-label="主导航">
      <button
        v-for="item in menuItems"
        :key="item.page"
        class="sidebar__item"
        :class="{ 'sidebar__item--active': currentPage === item.page }"
        @click="$emit('navigate', item.page)"
      >
        <span class="sidebar__icon">{{ item.icon }}</span>
        <span class="sidebar__copy"><strong>{{ item.label }}</strong><small>{{ item.hint }}</small></span>
      </button>
    </nav>
    <footer class="sidebar__footer"><span class="sidebar__status-dot"></span>原型数据 · 本地模式</footer>
  </aside>
</template>

<style scoped>
.sidebar-overlay { position: fixed; inset: 0; z-index: 80; border: 0; background: rgba(2, 6, 15, .66); backdrop-filter: blur(5px); }
.sidebar { position: fixed; inset: 0 auto 0 0; z-index: 90; width: min(340px, 86vw); display: flex; flex-direction: column; padding: 28px 22px 22px; background: linear-gradient(165deg, rgba(18, 27, 45, .98), rgba(7, 12, 24, .99)); border-right: 1px solid rgba(255, 255, 255, .1); box-shadow: 30px 0 80px rgba(0, 0, 0, .45); transform: translateX(-105%); transition: transform 220ms ease; }
.sidebar--open { transform: translateX(0); }
.sidebar__header, .sidebar__brand, .sidebar__item, .sidebar__footer { display: flex; align-items: center; }
.sidebar__header { justify-content: space-between; margin-bottom: 44px; }
.sidebar__brand { gap: 12px; }
.sidebar__brand-mark { display: grid; place-items: center; width: 42px; height: 42px; border-radius: 13px; color: #07101e; background: linear-gradient(135deg, #82e7ff, #9b8cff); font-weight: 900; font-size: 1.35rem; }
.sidebar__brand span:last-child, .sidebar__copy { display: grid; }
.sidebar__brand small, .sidebar__copy small { color: #8795aa; font-size: .68rem; letter-spacing: .08em; }
.sidebar__close { width: 40px; height: 40px; border: 0; border-radius: 12px; color: #b8c2d2; background: rgba(255, 255, 255, .07); font-size: 1.6rem; cursor: pointer; }
.sidebar__nav { display: grid; gap: 9px; }
.sidebar__item { gap: 16px; min-height: 68px; padding: 10px 14px; border: 1px solid transparent; border-radius: 16px; color: #b7c1d1; background: transparent; text-align: left; cursor: pointer; }
.sidebar__item:hover, .sidebar__item--active { color: #fff; border-color: rgba(132, 220, 255, .3); background: linear-gradient(90deg, rgba(93, 189, 255, .18), rgba(154, 123, 255, .08)); }
.sidebar__item--active::after { content: ''; width: 4px; height: 30px; margin-left: auto; border-radius: 99px; background: #79ddff; box-shadow: 0 0 14px #79ddff; }
.sidebar__icon { display: grid; place-items: center; width: 40px; height: 40px; font-size: 1.4rem; }
.sidebar__copy { gap: 4px; }
.sidebar__footer { gap: 9px; margin-top: auto; padding: 14px; color: #7f8da3; font-size: .78rem; }
.sidebar__status-dot { width: 8px; height: 8px; border-radius: 50%; background: #63e6a8; box-shadow: 0 0 10px #63e6a8; }
.fade-enter-active, .fade-leave-active { transition: opacity 180ms ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
