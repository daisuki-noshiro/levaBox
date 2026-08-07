<script setup lang="ts">
import type { Game } from '../types/game'

defineProps<{ game: Game }>()
defineEmits<{ back: []; notify: [message: string] }>()
</script>

<template>
  <main class="detail" :style="{ background: game.backgroundGradient }">
    <div class="detail__veil"></div>
    <button class="detail__back" @click="$emit('back')">← 返回</button>
    <section class="detail__content">
      <div class="detail__cover" :style="{ background: game.coverGradient }"><span>{{ game.title.slice(0, 1) }}</span><small>{{ game.subtitle }}</small></div>
      <article class="detail__info">
        <p class="page-kicker">GAME DETAIL</p>
        <h1>{{ game.title }}</h1>
        <p class="detail__subtitle">{{ game.subtitle }}</p>
        <div class="detail__tags"><span v-for="genre in game.genres" :key="genre">{{ genre }}</span></div>
        <dl class="detail__facts">
          <div><dt>开发商</dt><dd>{{ game.developer }}</dd></div>
          <div><dt>状态</dt><dd>{{ game.status }}</dd></div>
          <div><dt>游玩时间</dt><dd>{{ game.playtime ? `${game.playtime} 小时` : '尚未开始' }}</dd></div>
          <div><dt>发行年份</dt><dd>{{ game.year }}</dd></div>
        </dl>
        <div class="detail__summary"><h2>游戏简介</h2><p>{{ game.description }}</p></div>
        <div class="detail__actions">
          <button class="button button--primary" @click="$emit('notify', `“${game.title}”的真实启动功能将在后续阶段接入`)" >▶ 启动游戏</button>
          <button class="button button--glass" @click="$emit('notify', '编辑资料功能将在后续阶段接入')">编辑资料</button>
          <button class="button button--ghost" @click="$emit('back')">返回</button>
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.detail { position: relative; height: 100%; min-height: 0; overflow-y: auto; background-size: cover !important; }
.detail__veil { position: fixed; inset: 0; background: linear-gradient(90deg, rgba(4, 8, 17, .94), rgba(5, 9, 18, .72) 55%, rgba(5, 9, 18, .35)), linear-gradient(0deg, rgba(4, 8, 17, .9), transparent); }
.detail__back { position: relative; z-index: 2; margin: 28px 0 0 clamp(40px, 4cqw, 96px); padding: 10px 14px; border: 1px solid rgba(255, 255, 255, .12); border-radius: 12px; color: #d3dbe7; background: rgba(4, 8, 17, .32); cursor: pointer; }
.detail__content { position: relative; z-index: 2; display: grid; grid-template-columns: minmax(220px, 28cqw) minmax(0, 760px); align-items: center; gap: clamp(34px, 4.5cqw, 94px); width: min(1600px, calc(100% - clamp(80px, 8cqw, 200px))); min-height: calc(100% - 84px); margin: 0 auto; padding: 30px 0 clamp(54px, 6cqh, 88px); }
.detail__cover { position: relative; display: flex; flex-direction: column; align-items: center; justify-content: center; aspect-ratio: 3 / 4; overflow: hidden; border: 1px solid rgba(255, 255, 255, .22); border-radius: 26px; box-shadow: 0 30px 70px rgba(0, 0, 0, .45); }
.detail__cover::after { content: ''; position: absolute; inset: 0; background: radial-gradient(circle at 72% 18%, rgba(255, 255, 255, .45), transparent 25%), linear-gradient(150deg, transparent 30%, rgba(255, 255, 255, .14)); }
.detail__cover span { position: relative; z-index: 1; font-size: clamp(5rem, 11cqw, 10rem); font-weight: 900; text-shadow: 0 15px 35px rgba(0, 0, 0, .3); }
.detail__cover small { position: absolute; z-index: 1; right: 20px; bottom: 22px; left: 20px; padding-top: 13px; border-top: 1px solid rgba(255, 255, 255, .38); font-size: .7rem; letter-spacing: .12em; text-align: center; }
.detail__info h1 { margin: 3px 0 0; font-size: clamp(2.8rem, 5cqw, 5.6rem); line-height: 1.05; letter-spacing: -.055em; }
.detail__subtitle { margin: 8px 0 18px; color: rgba(255, 255, 255, .54); letter-spacing: .1em; }
.detail__tags { display: flex; gap: 8px; }
.detail__tags span { padding: 6px 10px; border: 1px solid rgba(255, 255, 255, .13); border-radius: 999px; color: #d6deea; background: rgba(0, 0, 0, .16); font-size: .7rem; }
.detail__facts { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1px; margin: 22px 0; overflow: hidden; border: 1px solid rgba(255, 255, 255, .1); border-radius: 15px; background: rgba(255, 255, 255, .1); }
.detail__facts div { padding: 13px 15px; background: rgba(8, 14, 25, .7); }
.detail__facts dt { color: #8997aa; font-size: .66rem; }
.detail__facts dd { margin: 5px 0 0; font-size: .8rem; font-weight: 700; }
.detail__summary h2 { margin: 0 0 8px; font-size: .92rem; }
.detail__summary p { margin: 0; color: #c3cedd; font-size: .88rem; line-height: 1.75; }
.detail__actions { display: flex; gap: 10px; margin-top: 24px; }
@container levabox-app (max-width: 850px) { .detail__content { grid-template-columns: 190px 1fr; width: calc(100% - 60px); gap: 28px; } .detail__facts { grid-template-columns: repeat(2, 1fr); } }
@container levabox-app (max-height: 760px) { .detail__content { align-items: start; padding-top: 18px; } .detail__cover { max-height: 490px; justify-self: center; } .detail__info h1 { font-size: clamp(2.5rem, 4.5cqw, 4.4rem); } .detail__summary p { line-height: 1.55; } }
@container levabox-app (max-width: 1100px) and (max-height: 760px) { .detail__actions { flex-wrap: wrap; } }
</style>
