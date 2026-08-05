<script setup lang="ts">
import type { MainPage } from '../types/navigation'

defineProps<{ currentPage: MainPage }>()
defineEmits<{ navigate: [page: MainPage] }>()

const navigationItems: Array<{ page: MainPage; label: string }> = [
  { page: 'home', label: '游戏大厅' },
  { page: 'library', label: '游戏库' },
  { page: 'import', label: '添加游戏' },
  { page: 'settings', label: '设置' },
]
</script>

<template>
  <nav class="global-navigation" aria-label="主导航">
    <button
      v-for="item in navigationItems"
      :key="item.page"
      class="global-navigation__item"
      :class="{ 'global-navigation__item--active': currentPage === item.page }"
      :aria-label="item.label"
      :aria-current="currentPage === item.page ? 'page' : undefined"
      @click="$emit('navigate', item.page)"
    >
      <svg v-if="item.page === 'home'" viewBox="0 0 24 24" aria-hidden="true">
        <path d="M3.5 10.5 12 3.8l8.5 6.7v9a1 1 0 0 1-1 1h-5.2v-6h-4.6v6H4.5a1 1 0 0 1-1-1z" />
      </svg>
      <svg v-else-if="item.page === 'library'" viewBox="0 0 24 24" aria-hidden="true">
        <rect x="3.5" y="4" width="7" height="7" rx="1.2" />
        <rect x="13.5" y="4" width="7" height="7" rx="1.2" />
        <rect x="3.5" y="14" width="7" height="7" rx="1.2" />
        <rect x="13.5" y="14" width="7" height="7" rx="1.2" />
      </svg>
      <svg v-else-if="item.page === 'import'" viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 4v16M4 12h16" />
      </svg>
      <svg v-else viewBox="0 0 24 24" aria-hidden="true">
        <circle cx="12" cy="12" r="3.2" />
        <path d="M12 2.8v2.1M12 19.1v2.1M21.2 12h-2.1M4.9 12H2.8M18.5 5.5 17 7M7 17l-1.5 1.5M18.5 18.5 17 17M7 7 5.5 5.5" />
      </svg>
    </button>
  </nav>
</template>

<style scoped>
.global-navigation {
  position: fixed;
  z-index: 24;
  top: 50%;
  left: clamp(16px, 1.7vw, 34px);
  display: grid;
  gap: clamp(9px, .8vw, 14px);
  transform: translateY(-50%);
}

.global-navigation__item {
  position: relative;
  display: grid;
  place-items: center;
  width: clamp(46px, 3vw, 58px);
  height: clamp(46px, 3vw, 58px);
  padding: 0;
  border: 1px solid transparent;
  border-radius: 9px;
  color: rgba(235, 243, 251, .48);
  background: transparent;
  cursor: pointer;
  transition: color 150ms ease, background 150ms ease, border-color 150ms ease, transform 150ms ease;
}

.global-navigation__item svg {
  width: clamp(21px, 1.45vw, 28px);
  height: clamp(21px, 1.45vw, 28px);
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.global-navigation__item:hover,
.global-navigation__item:focus-visible {
  color: rgba(255, 255, 255, .94);
  border-color: rgba(255, 255, 255, .2);
  outline: none;
  background: rgba(8, 14, 23, .34);
  transform: scale(1.04);
}

.global-navigation__item--active {
  color: #fff;
  border-color: rgba(158, 222, 245, .32);
  background: rgba(8, 14, 23, .46);
}

.global-navigation__item--active::before {
  content: '';
  position: absolute;
  left: -5px;
  top: 50%;
  width: 2px;
  height: 20px;
  border-radius: 2px;
  background: #b8e9f8;
  transform: translateY(-50%);
}

@media (max-width: 1100px), (max-height: 760px) {
  .global-navigation { left: 12px; gap: 7px; }
  .global-navigation__item { width: 43px; height: 43px; }
}

@media (prefers-reduced-motion: reduce) {
  .global-navigation__item { transition: none; }
}
</style>
