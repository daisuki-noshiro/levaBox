<script setup lang="ts">
import GameCard from '../components/GameCard.vue'
import type { Game } from '../types/game'
import type { MainPage } from '../types/navigation'

defineProps<{ games: Game[]; selectedGame: Game }>()
defineEmits<{
  selectGame: [game: Game]
  openDetail: [game: Game]
  notify: [message: string]
  openMenu: []
  navigate: [page: MainPage]
}>()
</script>

<template>
  <section class="home">
    <div class="home__media" aria-hidden="true">
      <Transition name="media-fade" mode="out-in">
        <video
          v-if="selectedGame.backgroundType === 'video' && selectedGame.backgroundVideo"
          :key="`${selectedGame.id}-video`"
          class="home__backdrop"
          :src="selectedGame.backgroundVideo"
          :poster="selectedGame.backgroundImage"
          :style="{ objectPosition: selectedGame.backgroundPosition ?? 'center' }"
          autoplay
          muted
          loop
          playsinline
          preload="metadata"
        ></video>
        <img
          v-else
          :key="`${selectedGame.id}-image`"
          class="home__backdrop"
          :src="selectedGame.backgroundImage"
          alt=""
          :style="{ objectPosition: selectedGame.backgroundPosition ?? 'center' }"
        />
      </Transition>
    </div>
    <div class="home__shade"></div>

    <nav class="home-rail" aria-label="大厅快捷导航">
      <button class="home-rail__menu" aria-label="展开主菜单" @click="$emit('openMenu')">
        <span class="home-rail__menu-lines"><i></i><i></i><i></i></span><strong>levaBox</strong>
      </button>
      <div class="home-rail__divider"></div>
      <button class="home-rail__item home-rail__item--active" @click="$emit('navigate', 'home')"><span>⌂</span><strong>大厅</strong></button>
      <button class="home-rail__item" @click="$emit('navigate', 'library')"><span>▦</span><strong>游戏库</strong></button>
      <button class="home-rail__item" @click="$emit('notify', '收藏筛选可在游戏库中查看')"><span>♥</span><strong>收藏</strong></button>
      <button class="home-rail__item" @click="$emit('notify', '最近游玩记录将在后续阶段接入')"><span>◷</span><strong>最近游玩</strong></button>
      <button class="home-rail__item" @click="$emit('navigate', 'import')"><span>＋</span><strong>导入游戏</strong></button>
      <button class="home-rail__item home-rail__item--bottom" @click="$emit('navigate', 'settings')"><span>⚙</span><strong>设置</strong></button>
    </nav>

    <Transition name="hero" mode="out-in">
      <article :key="selectedGame.id" class="home__hero">
        <p class="home__eyebrow">
          <span class="home__live-dot"></span>
          {{ selectedGame.backgroundType === 'video' ? '动态背景' : '精选游戏' }} · {{ selectedGame.genres.join(' / ') }}
        </p>
        <img v-if="selectedGame.logo" class="home__logo" :src="selectedGame.logo" :alt="selectedGame.title" />
        <button v-else class="home__title" @click="$emit('openDetail', selectedGame)">{{ selectedGame.title }}</button>
        <p class="home__subtitle">{{ selectedGame.subtitle }}</p>
        <p class="home__description">{{ selectedGame.shortDescription }}</p>
        <div class="home__meta">
          <span>{{ selectedGame.developer }}</span><i></i><span>{{ selectedGame.year }}</span><i></i><span>{{ selectedGame.status }}</span>
        </div>
        <div class="home__actions">
          <button class="button button--primary" @click="$emit('notify', `“${selectedGame.title}”的启动功能将在后续阶段接入`)" >▶ 开始游戏</button>
          <button class="button button--glass" @click="$emit('openDetail', selectedGame)">查看详情</button>
        </div>
      </article>
    </Transition>

    <section class="home__carousel" aria-label="游戏封面轮播">
      <header class="home__carousel-header">
        <div><span></span><strong>游戏大厅</strong><small>选择游戏以切换背景</small></div>
        <p><b>{{ String(games.findIndex((game) => game.id === selectedGame.id) + 1).padStart(2, '0') }}</b> / {{ String(games.length).padStart(2, '0') }}</p>
      </header>
      <div class="home__cards">
        <GameCard
          v-for="game in games"
          :key="game.id"
          :game="game"
          variant="rail"
          :selected="game.id === selectedGame.id"
          @select="$emit('selectGame', $event)"
        />
      </div>
    </section>
  </section>
</template>

<style scoped>
.home { position: relative; height: 100%; min-height: 0; overflow: hidden; background: #070b12; }
.home__media, .home__backdrop, .home__shade { position: absolute; inset: 0; }
.home__backdrop { width: 100%; height: 100%; object-fit: cover; }
.home__shade { background: linear-gradient(90deg, rgba(3, 7, 13, .93) 0%, rgba(3, 7, 13, .7) 27%, rgba(3, 7, 13, .12) 63%, rgba(3, 7, 13, .08) 100%), linear-gradient(0deg, rgba(3, 7, 13, .96) 0%, rgba(3, 7, 13, .58) 26%, transparent 57%), linear-gradient(180deg, rgba(3, 7, 13, .25), transparent 22%); }

.home-rail { position: absolute; z-index: 10; top: 18px; bottom: 18px; left: 18px; display: flex; flex-direction: column; width: 58px; padding: 7px; overflow: hidden; border: 1px solid rgba(255, 255, 255, .11); border-radius: 18px; background: rgba(5, 10, 18, .7); box-shadow: 0 18px 50px rgba(0, 0, 0, .28); backdrop-filter: blur(16px); transition: width 180ms ease, background 180ms ease; }
.home-rail:hover, .home-rail:focus-within { width: 178px; background: rgba(5, 10, 18, .9); }
.home-rail button { display: flex; align-items: center; flex: 0 0 auto; width: 162px; height: 44px; padding: 0; border: 0; border-radius: 12px; color: #9da9b9; background: transparent; cursor: pointer; }
.home-rail button:hover, .home-rail button:focus-visible { color: #fff; outline: none; background: rgba(255, 255, 255, .09); }
.home-rail__menu { color: #fff !important; }
.home-rail__menu-lines, .home-rail__item > span { display: grid; place-items: center; flex: 0 0 42px; width: 42px; }
.home-rail__menu-lines { gap: 4px; }
.home-rail__menu-lines i { display: block; width: 18px; height: 2px; border-radius: 9px; background: #fff; }
.home-rail strong { opacity: 0; white-space: nowrap; transition: opacity 100ms ease 30ms; }
.home-rail:hover strong, .home-rail:focus-within strong { opacity: 1; }
.home-rail__menu strong { font-size: .85rem; letter-spacing: .04em; }
.home-rail__divider { height: 1px; margin: 7px 5px 9px; background: rgba(255, 255, 255, .09); }
.home-rail__item { margin-bottom: 4px; font-size: .78rem; text-align: left; }
.home-rail__item > span { font-size: 1.2rem; }
.home-rail__item--active { color: #fff !important; background: linear-gradient(90deg, rgba(91, 206, 255, .22), rgba(91, 206, 255, .05)) !important; }
.home-rail__item--active::before { content: ''; position: absolute; left: 0; width: 3px; height: 24px; border-radius: 99px; background: #72dfff; box-shadow: 0 0 12px #72dfff; }
.home-rail__item--bottom { margin-top: auto; }

.home__hero { position: relative; z-index: 3; width: min(520px, 45vw); padding: clamp(80px, 11vh, 118px) 0 0 clamp(100px, 8.2vw, 156px); }
.home__eyebrow { display: flex; align-items: center; gap: 8px; margin: 0 0 12px; color: #c3cfda; font-size: .68rem; font-weight: 700; letter-spacing: .12em; }
.home__live-dot { width: 7px; height: 7px; border-radius: 50%; background: #67e2ff; box-shadow: 0 0 12px #67e2ff; }
.home__logo { display: block; max-width: min(430px, 38vw); max-height: 125px; object-fit: contain; object-position: left center; }
.home__title { display: block; max-width: 100%; padding: 0; border: 0; color: #fff; background: transparent; font: inherit; font-size: clamp(2.7rem, 4.6vw, 5.3rem); font-weight: 800; line-height: 1.05; letter-spacing: -.06em; text-align: left; text-shadow: 0 10px 35px rgba(0, 0, 0, .48); cursor: pointer; }
.home__title:hover { color: #dff8ff; }
.home__subtitle { margin: 6px 0 16px; color: rgba(255, 255, 255, .58); font-size: clamp(.72rem, 1vw, .92rem); letter-spacing: .11em; }
.home__description { max-width: 440px; margin: 0; color: #d8e0e9; font-size: clamp(.82rem, 1.04vw, .95rem); line-height: 1.6; text-shadow: 0 2px 16px rgba(0, 0, 0, .65); }
.home__meta { display: flex; align-items: center; gap: 9px; margin-top: 13px; color: #b6c0ce; font-size: .68rem; }
.home__meta i { width: 3px; height: 3px; border-radius: 50%; background: #6e7a89; }
.home__actions { display: flex; gap: 10px; margin-top: 18px; }
.home__actions .button { min-height: 43px; padding: 0 17px; font-size: .8rem; }

.home__carousel { position: absolute; z-index: 4; right: 0; bottom: 0; left: 0; padding: 0 clamp(25px, 3.4vw, 58px) clamp(15px, 2vh, 24px) clamp(96px, 7.7vw, 148px); }
.home__carousel-header { display: flex; align-items: end; justify-content: space-between; margin-bottom: 10px; }
.home__carousel-header div { display: flex; align-items: center; gap: 9px; }
.home__carousel-header div > span { width: 3px; height: 20px; border-radius: 9px; background: #70ddff; box-shadow: 0 0 10px #70ddff; }
.home__carousel-header strong { font-size: .88rem; }
.home__carousel-header small { color: #8491a1; font-size: .65rem; }
.home__carousel-header p { margin: 0; color: #7f8b9b; font-size: .7rem; }
.home__carousel-header b { color: #fff; font-size: .95rem; }
.home__cards { display: flex; gap: clamp(12px, 1.25vw, 20px); padding: 8px 4px 5px; overflow-x: auto; overflow-y: hidden; scrollbar-width: none; }
.home__cards::-webkit-scrollbar { display: none; }

.media-fade-enter-active, .media-fade-leave-active { transition: opacity 300ms ease; }
.media-fade-enter-from, .media-fade-leave-to { opacity: 0; }
.hero-enter-active, .hero-leave-active { transition: opacity 170ms ease, transform 170ms ease; }
.hero-enter-from { opacity: 0; transform: translateY(10px); }
.hero-leave-to { opacity: 0; transform: translateY(-7px); }

@media (max-height: 760px) {
  .home__hero { padding-top: 72px; }
  .home__title { font-size: clamp(2.5rem, 4.2vw, 4.4rem); }
  .home__subtitle { margin-bottom: 10px; }
  .home__actions { margin-top: 13px; }
  .home__carousel { padding-bottom: 13px; }
}
</style>
