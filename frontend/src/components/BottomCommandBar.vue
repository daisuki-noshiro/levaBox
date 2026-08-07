<script setup lang="ts">
defineProps<{ menuOpen: boolean }>()
defineEmits<{ toggleMenu: [] }>()
</script>

<template>
  <footer class="bottom-command-bar" aria-label="主机操作栏">
    <button
      id="global-menu-trigger"
      class="bottom-command-bar__menu"
      type="button"
      aria-controls="global-menu"
      :aria-expanded="menuOpen"
      @click="$emit('toggleMenu')"
    >
      <span class="bottom-command-bar__brand">LEVABOX</span>
      <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 7h16M4 12h16M4 17h16" /></svg>
      <span>菜单</span>
    </button>

    <div class="bottom-command-bar__hints" aria-hidden="true">
      <span><i>A</i>选择</span>
      <span><i>B</i>返回</span>
    </div>
  </footer>
</template>

<style scoped>
.bottom-command-bar {
  position: absolute;
  z-index: 40;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--command-bar-height);
  padding: 0 clamp(20px, 3.2cqw, 72px);
  border-top: 1px solid rgba(255, 255, 255, .09);
  color: rgba(238, 244, 250, .8);
  background: rgba(5, 8, 13, .92);
  box-shadow: 0 -10px 30px rgba(0, 0, 0, .18);
}

.bottom-command-bar__menu {
  display: flex;
  align-items: center;
  gap: clamp(8px, .65cqw, 13px);
  height: calc(100% - 8px);
  padding: 0 clamp(8px, .7cqw, 14px);
  border: 1px solid transparent;
  border-radius: 7px;
  color: inherit;
  background: transparent;
  font-size: clamp(.76rem, 1.3cqh, .96rem);
  font-weight: 700;
  cursor: pointer;
}
.bottom-command-bar__menu svg { width: clamp(17px, 1.25cqw, 23px); height: clamp(17px, 1.25cqw, 23px); fill: none; stroke: currentColor; stroke-width: 1.8; stroke-linecap: round; }
.bottom-command-bar__menu:hover,
.bottom-command-bar__menu:focus-visible { border-color: rgba(255, 255, 255, .2); color: #fff; outline: none; background: rgba(255, 255, 255, .07); }
.bottom-command-bar__brand { padding-right: clamp(9px, .8cqw, 15px); border-right: 1px solid rgba(255, 255, 255, .18); color: #fff; font-size: clamp(.66rem, 1.05cqh, .82rem); font-weight: 900; letter-spacing: .13em; }
.bottom-command-bar__hints { display: flex; align-items: center; gap: clamp(14px, 1.6cqw, 30px); font-size: clamp(.72rem, 1.2cqh, .9rem); font-weight: 700; }
.bottom-command-bar__hints span { display: flex; align-items: center; gap: 7px; }
.bottom-command-bar__hints i { display: grid; place-items: center; width: clamp(20px, 1.5cqw, 26px); aspect-ratio: 1; border: 1px solid rgba(255, 255, 255, .38); border-radius: 50%; color: #fff; font-size: .72em; font-style: normal; }

@media (prefers-reduced-motion: reduce) {
  .bottom-command-bar__menu { transition: none; }
}
</style>
