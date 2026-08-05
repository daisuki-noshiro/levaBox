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
  dragStart: [game: Game, event: DragEvent]
  dragOver: [game: Game, event: DragEvent]
  drop: [game: Game, event: DragEvent]
  dragEnd: [event: DragEvent]
}>()
</script>

<template>
  <button
    class="game-card"
    :class="[`game-card--${variant}`, { 'game-card--selected': selected, 'game-card--reordering': reordering, 'game-card--moving': moving }]"
    :aria-pressed="selected"
    :aria-label="`${game.title}，${game.status}`"
    :data-game-id="game.id"
    :draggable="variant === 'rail' && moving"
    @click="$emit('select', game)"
    @dragstart="$emit('dragStart', game, $event)"
    @dragover="$emit('dragOver', game, $event)"
    @drop="$emit('drop', game, $event)"
    @dragend="$emit('dragEnd', $event)"
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
  transition: transform 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.game-card--grid:hover,
.game-card--grid:focus-visible {
  transform: translateY(-5px);
  border-color: var(--card-accent);
  outline: none;
}

.game-card--rail {
  width: clamp(190px, 20vw, 250px);
  overflow: visible;
  border-color: transparent;
  border-radius: 14px;
  background: transparent;
  opacity: .58;
  transform: scale(.92);
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
  aspect-ratio: 16 / 9;
  border: 1px solid rgba(255, 255, 255, .16);
  border-radius: 13px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, .22);
}

.game-card--rail:hover,
.game-card--rail:focus-visible {
  outline: none;
  opacity: .86;
  transform: translateY(-4px) scale(.96);
}

.game-card--rail.game-card--selected {
  z-index: 2;
  opacity: 1;
  transform: translateY(-11px) scale(1.02);
}

.game-card--rail.game-card--selected .game-card__art {
  border-color: var(--card-accent);
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--card-accent) 70%, transparent), 0 12px 30px rgba(0, 0, 0, .38), 0 0 24px color-mix(in srgb, var(--card-accent) 24%, transparent);
}

.game-card--rail.game-card--reordering:not(.game-card--moving) { opacity: .4; }
.game-card--rail.game-card--moving {
  z-index: 4;
  opacity: 1;
  transform: translateY(-20px) scale(1.075);
  cursor: grabbing;
}

.game-card--rail.game-card--moving .game-card__art {
  border-color: #fff;
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--card-accent) 85%, white), 0 18px 42px rgba(0, 0, 0, .5), 0 0 32px color-mix(in srgb, var(--card-accent) 45%, transparent);
}

.game-card__cover {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
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
</style>
