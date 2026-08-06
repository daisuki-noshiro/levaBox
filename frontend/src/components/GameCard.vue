<script setup lang="ts">
import type { Game } from '../types/game'

withDefaults(defineProps<{
  game: Game
  selected?: boolean
  variant?: 'rail' | 'grid'
  reordering?: boolean
  moving?: boolean
}>(), {
  selected: false,
  variant: 'grid',
  reordering: false,
  moving: false,
})

defineEmits<{
  select: [game: Game]
  pointerDown: [game: Game, event: PointerEvent]
}>()
</script>

<template>
  <button
    class="game-card"
    :class="[`game-card--${variant}`, { 'game-card--selected': selected, 'game-card--reordering': reordering, 'game-card--moving': moving }]"
    :aria-pressed="selected"
    :aria-label="`${game.title}，${game.status}`"
    :data-game-id="game.id"
    @click="$emit('select', game)"
    @pointerdown="$emit('pointerDown', game, $event)"
  >
    <span class="game-card__art" :style="{ background: game.coverGradient }">
      <img class="game-card__cover" :src="game.cover" :alt="`${game.title}封面`" :style="{ objectPosition: game.coverPosition ?? 'center' }" />
      <span class="game-card__shine"></span>
      <span v-if="game.favorite" class="game-card__favorite" aria-label="已收藏">♥</span>
    </span>
    <span class="game-card__body">
      <strong>{{ game.title }}</strong>
      <span>{{ game.developer }} · {{ game.year }}</span>
    </span>
  </button>
</template>

<style scoped>
.game-card {
  --card-accent: v-bind('game.accent');
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 16px;
  padding: 0;
  overflow: hidden;
  color: #fff;
  background: rgba(12, 18, 31, 0.78);
  text-align: left;
  cursor: pointer;
  flex: 0 0 auto;
  transition: transform 220ms ease, border-color 180ms ease, box-shadow 180ms ease, opacity 180ms ease, filter 180ms ease;
}

.game-card--grid:hover,
.game-card--grid:focus-visible {
  transform: translateY(-5px);
  border-color: var(--card-accent);
  outline: none;
}

.game-card--rail {
  width: var(--rail-card-width, 250px);
  overflow: visible;
  border-color: rgba(255, 255, 255, .12);
  border-radius: 8px;
  background: transparent;
  opacity: .76;
  filter: brightness(.88);
  transform: translateY(0) scale(1);
  transform-origin: center center;
  transition:
    transform var(--rail-motion-duration, 360ms) var(--rail-motion-easing, cubic-bezier(.22, .8, .25, 1)),
    border-color var(--rail-motion-duration, 360ms) var(--rail-motion-easing, cubic-bezier(.22, .8, .25, 1)),
    box-shadow var(--rail-motion-duration, 360ms) var(--rail-motion-easing, cubic-bezier(.22, .8, .25, 1)),
    opacity var(--rail-motion-duration, 360ms) var(--rail-motion-easing, cubic-bezier(.22, .8, .25, 1)),
    filter var(--rail-motion-duration, 360ms) var(--rail-motion-easing, cubic-bezier(.22, .8, .25, 1));
}

.game-card--grid { width: 100%; overflow: hidden; }

.game-card--grid.game-card--selected {
  transform: translateY(-8px) scale(1.035);
  border-color: var(--card-accent);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--card-accent) 70%, transparent), 0 16px 38px rgba(0, 0, 0, 0.4);
}

.game-card__art {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 16 / 10;
  overflow: hidden;
}

.game-card--rail .game-card__art {
  aspect-ratio: 5 / 7;
  border-radius: 7px;
  box-shadow: 0 10px 24px rgba(0, 0, 0, .28);
}

.game-card--rail:hover,
.game-card--rail:focus-visible {
  outline: none;
  opacity: .92;
  filter: brightness(.98);
  border-color: rgba(255, 255, 255, .4);
}

.game-card--rail.game-card--selected {
  z-index: 2;
  opacity: 1;
  filter: none;
  border-color: var(--card-accent);
  transform: translateY(calc(var(--rail-selected-lift, 22px) * -1)) scale(var(--rail-selected-scale, 1.08));
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--card-accent) 64%, transparent), 0 16px 34px rgba(0, 0, 0, .42);
}

.game-card--rail.game-card--moving {
  z-index: 4;
  opacity: 1;
  filter: none;
  transform: translateY(calc(var(--rail-reorder-lift, 35px) * -1)) scale(var(--rail-selected-scale, 1.08));
  cursor: grabbing;
}

.game-card--rail.game-card--moving {
  border-color: #fff;
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--card-accent) 72%, white), 0 20px 42px rgba(0, 0, 0, .5);
}

.game-card__cover {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  user-select: none;
  -webkit-user-drag: none;
}

.game-card__shine {
  position: absolute;
  inset: 0;
  background: linear-gradient(145deg, rgba(255, 255, 255, .08), transparent 45%, rgba(0, 0, 0, .1));
}

.game-card__favorite {
  position: absolute;
  top: 10px;
  right: 12px;
  color: #ffb8cf;
  font-size: 1.1rem;
}

.game-card--rail .game-card__favorite,
.game-card--rail .game-card__body { display: none; }

.game-card__body { display: grid; gap: 4px; padding: 12px 14px 14px; }
.game-card__body strong,
.game-card__body span { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.game-card__body strong { font-size: 1rem; }
.game-card__body span { color: #9caac0; font-size: 0.75rem; }

@media (max-height: 760px) {
  .game-card__body { padding: 9px 12px 11px; }
}

@media (prefers-reduced-motion: reduce) {
  .game-card { transition-duration: 1ms; }
}
</style>
