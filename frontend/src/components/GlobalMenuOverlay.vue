<script setup lang="ts">
import { nextTick, ref } from 'vue'
import type { MainPage } from '../types/navigation'

const props = defineProps<{ currentPage: MainPage }>()
const emit = defineEmits<{ close: []; navigate: [page: MainPage] }>()

const menuButtons = ref<HTMLButtonElement[]>([])
const focusedIndex = ref(0)

const items: Array<{ page: MainPage; label: string; icon: 'home' | 'library' | 'import' | 'settings' }> = [
  { page: 'home', label: '首页', icon: 'home' },
  { page: 'library', label: '游戏库', icon: 'library' },
  { page: 'import', label: '导入游戏', icon: 'import' },
  { page: 'settings', label: '设置', icon: 'settings' },
]

function setMenuButton(element: unknown, index: number): void {
  if (element instanceof HTMLButtonElement) menuButtons.value[index] = element
}

function focusIndex(index: number): void {
  focusedIndex.value = (index + items.length) % items.length
  nextTick(() => menuButtons.value[focusedIndex.value]?.focus({ preventScroll: true }))
}

function focusCurrent(): void {
  menuButtons.value = []
  focusedIndex.value = Math.max(0, items.findIndex((item) => item.page === props.currentPage))
  nextTick(() => menuButtons.value[focusedIndex.value]?.focus({ preventScroll: true }))
}

function moveFocus(direction: -1 | 1): void {
  focusIndex(focusedIndex.value + direction)
}

function confirmFocused(): void {
  const item = items[focusedIndex.value]
  if (item) emit('navigate', item.page)
}

function handleKeydown(event: KeyboardEvent): void {
  if (event.key === 'ArrowUp' || event.key === 'ArrowDown') {
    event.preventDefault()
    moveFocus(event.key === 'ArrowDown' ? 1 : -1)
    return
  }
  if (event.key === 'Enter') {
    event.preventDefault()
    confirmFocused()
    return
  }
  if (event.key === 'Tab') {
    event.preventDefault()
    moveFocus(event.shiftKey ? -1 : 1)
  }
}

defineExpose({ focusCurrent, moveFocus, confirmFocused })
</script>

<template>
  <div class="global-menu-layer" @click.self="$emit('close')">
    <aside
      id="global-menu"
      class="global-menu"
      role="dialog"
      aria-modal="true"
      aria-label="全局菜单"
      @keydown="handleKeydown"
    >
      <header><span>LEVABOX</span><strong>主菜单</strong></header>
      <nav aria-label="页面导航">
        <button
          v-for="(item, index) in items"
          :key="item.page"
          :ref="(element) => setMenuButton(element, index)"
          type="button"
          :class="{ 'global-menu__item--active': currentPage === item.page }"
          :aria-current="currentPage === item.page ? 'page' : undefined"
          @focus="focusedIndex = index"
          @click="$emit('navigate', item.page)"
        >
          <svg v-if="item.icon === 'home'" viewBox="0 0 24 24" aria-hidden="true"><path d="M3.5 10.5 12 3.8l8.5 6.7v9a1 1 0 0 1-1 1h-5.2v-6h-4.6v6H4.5a1 1 0 0 1-1-1z" /></svg>
          <svg v-else-if="item.icon === 'library'" viewBox="0 0 24 24" aria-hidden="true"><rect x="3.5" y="4" width="7" height="7" rx="1.2" /><rect x="13.5" y="4" width="7" height="7" rx="1.2" /><rect x="3.5" y="14" width="7" height="7" rx="1.2" /><rect x="13.5" y="14" width="7" height="7" rx="1.2" /></svg>
          <svg v-else-if="item.icon === 'import'" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 4v16M4 12h16" /></svg>
          <svg v-else viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="3.2" /><path d="M12 2.8v2.1M12 19.1v2.1M21.2 12h-2.1M4.9 12H2.8M18.5 5.5 17 7M7 17l-1.5 1.5M18.5 18.5 17 17M7 7 5.5 5.5" /></svg>
          <span>{{ item.label }}</span>
          <i v-if="currentPage === item.page">当前</i>
        </button>
      </nav>
      <footer>Esc / B 返回</footer>
    </aside>
  </div>
</template>

<style scoped>
.global-menu-layer { position: absolute; z-index: 100; inset: 0; background: rgba(2, 5, 10, .48); backdrop-filter: blur(8px) brightness(.68); }
.global-menu { display: flex; flex-direction: column; width: clamp(280px, 24cqw, 400px); height: 100%; padding: clamp(34px, 5cqh, 76px) clamp(20px, 2.2cqw, 38px) clamp(22px, 3cqh, 42px); border-right: 1px solid rgba(255, 255, 255, .11); color: #f5f8fb; background: linear-gradient(100deg, rgba(10, 13, 18, .99), rgba(13, 17, 23, .97)); box-shadow: 18px 0 44px rgba(0, 0, 0, .3); }
.global-menu header { display: grid; gap: 5px; margin: 0 12px clamp(28px, 4cqh, 54px); }
.global-menu header span { color: rgba(255, 255, 255, .55); font-size: clamp(.66rem, 1.05cqh, .82rem); font-weight: 900; letter-spacing: .18em; }
.global-menu header strong { font-size: clamp(1.65rem, 3.2cqh, 2.65rem); letter-spacing: -.04em; }
.global-menu nav { display: grid; gap: clamp(6px, .8cqh, 11px); }
.global-menu button { position: relative; display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: clamp(12px, 1.2cqw, 19px); width: 100%; min-height: clamp(54px, 7cqh, 76px); padding: 0 clamp(14px, 1.35cqw, 22px); border: 1px solid transparent; border-radius: 7px; color: rgba(232, 238, 245, .64); background: transparent; font-size: clamp(.92rem, 1.75cqh, 1.18rem); font-weight: 700; text-align: left; cursor: pointer; transition: color 190ms ease, background 190ms ease, border-color 190ms ease, transform 190ms ease; }
.global-menu button::before { content: ''; position: absolute; top: 50%; left: -1px; width: 3px; height: 0; border-radius: 3px; background: #d9f3fb; transform: translateY(-50%); transition: height 190ms ease; }
.global-menu button:hover,
.global-menu button:focus-visible { border-color: rgba(255, 255, 255, .12); color: #fff; outline: none; background: rgba(255, 255, 255, .075); transform: translateX(4px); }
.global-menu button:hover::before,
.global-menu button:focus-visible::before,
.global-menu__item--active::before { height: 46%; }
.global-menu__item--active { color: rgba(255, 255, 255, .92) !important; background: rgba(255, 255, 255, .045) !important; }
.global-menu svg { width: clamp(22px, 1.7cqw, 30px); height: clamp(22px, 1.7cqw, 30px); fill: none; stroke: currentColor; stroke-width: 1.7; stroke-linecap: round; stroke-linejoin: round; }
.global-menu button i { color: rgba(255, 255, 255, .47); font-size: .62em; font-style: normal; font-weight: 600; }
.global-menu footer { margin-top: auto; padding: 0 12px; color: rgba(255, 255, 255, .45); font-size: clamp(.7rem, 1.15cqh, .86rem); font-weight: 700; }

@media (prefers-reduced-motion: reduce) {
  .global-menu button { transition-duration: 1ms; }
}
</style>
