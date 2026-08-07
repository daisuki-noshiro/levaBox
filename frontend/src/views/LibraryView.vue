<script setup lang="ts">
import { computed, ref } from 'vue'
import GameCard from '../components/GameCard.vue'
import type { Game } from '../types/game'

const props = defineProps<{ games: Game[] }>()
defineEmits<{ openDetail: [game: Game] }>()

type Filter = 'all' | 'favorite' | 'recent'
const query = ref('')
const filter = ref<Filter>('all')
const sortBy = ref<'title' | 'year' | 'playtime'>('title')

const visibleGames = computed(() => {
  const keyword = query.value.trim().toLowerCase()
  return props.games
    .filter((game) => filter.value === 'all' || (filter.value === 'favorite' ? game.favorite : game.recentlyPlayed))
    .filter((game) => !keyword || `${game.title} ${game.subtitle} ${game.developer}`.toLowerCase().includes(keyword))
    .toSorted((a, b) => {
      if (sortBy.value === 'year') return b.year - a.year
      if (sortBy.value === 'playtime') return b.playtime - a.playtime
      return a.title.localeCompare(b.title, 'zh-CN')
    })
})
</script>

<template>
  <main class="library page-shell">
    <header class="page-header">
      <div><p class="page-kicker">LIBRARY</p><h1>游戏库</h1><p>整理、筛选并重新发现你的故事。</p></div>
      <div class="library__count"><strong>{{ visibleGames.length }}</strong><span>当前结果</span></div>
    </header>
    <section class="library__toolbar" aria-label="游戏筛选">
      <label class="library__search"><span>⌕</span><input v-model="query" type="search" placeholder="搜索游戏、开发商……" /></label>
      <div class="library__filters">
        <button v-for="item in ([['all', '全部游戏'], ['favorite', '我的收藏'], ['recent', '最近游玩']] as const)" :key="item[0]" :class="{ active: filter === item[0] }" @click="filter = item[0]">{{ item[1] }}</button>
      </div>
      <label class="library__sort"><span>排序</span><select v-model="sortBy"><option value="title">按名称</option><option value="year">按年份</option><option value="playtime">按游玩时间</option></select></label>
    </section>
    <section v-if="visibleGames.length" class="library__grid">
      <GameCard v-for="game in visibleGames" :key="game.id" :game="game" variant="grid" @select="$emit('openDetail', $event)" />
    </section>
    <section v-else class="library__empty"><span>⌕</span><h2>没有找到游戏</h2><p>试试其他关键词或筛选条件。</p></section>
  </main>
</template>

<style scoped>
.library { overflow-y: auto; }
.library__count { display: grid; min-width: 115px; padding: 15px 20px; border: 1px solid var(--line); border-radius: 18px; background: var(--panel); text-align: right; }
.library__count strong { font-size: 1.55rem; }
.library__count span { color: var(--muted); font-size: .72rem; }
.library__toolbar { display: flex; align-items: center; gap: 12px; margin: 29px 0 26px; padding: 13px; border: 1px solid var(--line); border-radius: 20px; background: rgba(14, 22, 37, .72); }
.library__search { display: flex; align-items: center; gap: 10px; width: min(330px, 30cqw); padding: 0 14px; border: 1px solid var(--line); border-radius: 13px; background: rgba(255, 255, 255, .05); }
.library__search span { color: #91a0b6; font-size: 1.5rem; }
.library__search input { width: 100%; height: 43px; border: 0; outline: 0; color: #fff; background: transparent; font: inherit; }
.library__filters { display: flex; gap: 7px; }
.library__filters button { min-height: 43px; padding: 0 15px; border: 1px solid transparent; border-radius: 12px; color: #9eabc0; background: transparent; cursor: pointer; }
.library__filters button:hover, .library__filters button.active { border-color: rgba(107, 214, 255, .25); color: #fff; background: rgba(104, 190, 255, .13); }
.library__sort { display: flex; align-items: center; gap: 9px; margin-left: auto; color: var(--muted); font-size: .8rem; }
.library__sort select { height: 43px; padding: 0 34px 0 13px; border: 1px solid var(--line); border-radius: 12px; color: #fff; background: #172136; }
.library__grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(clamp(180px, 12cqw, 240px), 1fr)); gap: clamp(16px, 1.5cqw, 30px); padding-bottom: clamp(36px, 5cqh, 70px); }
.library__empty { display: grid; place-items: center; padding: 90px 20px; color: var(--muted); text-align: center; }
.library__empty span { font-size: 3rem; }
.library__empty h2 { margin: 12px 0 4px; color: #fff; }
.library__empty p { margin: 0; }
@container levabox-app (max-width: 1050px) { .library__toolbar { flex-wrap: wrap; } .library__search { width: 100%; } .library__sort { margin-left: 0; } }
@container levabox-app (max-height: 760px) { .library__toolbar { margin: 20px 0 20px; } }
</style>
