<script setup lang="ts">
defineEmits<{ notify: [message: string] }>()

const actions = [
  { icon: '▣', title: '扫描文件夹', description: '选择一个目录，自动寻找可能的游戏程序。', action: '扫描文件夹' },
  { icon: 'EXE', title: '手动选择 EXE', description: '直接指定游戏的主启动程序。', action: '选择 EXE' },
  { icon: 'ZIP', title: '导入压缩包', description: '从本地压缩文件创建待安装项目。', action: '导入压缩包', developing: true },
]
</script>

<template>
  <main class="import-view page-shell">
    <header class="page-header">
      <div><p class="page-kicker">ADD GAMES</p><h1>导入游戏</h1><p>把本地收藏加入 levaBox。当前阶段仅展示操作流程。</p></div>
      <span class="import-view__badge">原型模式</span>
    </header>
    <section class="import-view__actions">
      <article v-for="item in actions" :key="item.title" class="import-card">
        <span class="import-card__icon">{{ item.icon }}</span>
        <div><span v-if="item.developing" class="import-card__tag">开发中</span><h2>{{ item.title }}</h2><p>{{ item.description }}</p></div>
        <button @click="$emit('notify', `${item.action}将在后续阶段接入 Windows 功能`)" >{{ item.action }} <span>→</span></button>
      </article>
    </section>
    <section class="import-view__result">
      <div class="import-view__result-title"><div><h2>已识别游戏</h2><p>扫描结果会显示在这里，确认后再加入游戏库。</p></div><span>0 个项目</span></div>
      <div class="import-view__empty"><span class="import-view__empty-icon">＋</span><h3>还没有待导入的游戏</h3><p>从上方选择一种方式开始。此页面不会访问你的文件。</p></div>
    </section>
  </main>
</template>

<style scoped>
.import-view { overflow-y: auto; }
.import-view__badge { padding: 9px 14px; border: 1px solid rgba(115, 222, 255, .25); border-radius: 999px; color: #aeeeff; background: rgba(80, 185, 230, .1); font-size: .75rem; }
.import-view__actions { display: grid; grid-template-columns: repeat(3, 1fr); gap: clamp(16px, 1.3cqw, 28px); margin: clamp(24px, 3cqh, 44px) 0 clamp(22px, 2.8cqh, 38px); }
.import-card { display: grid; grid-template-rows: auto 1fr auto; min-height: clamp(235px, 25cqh, 340px); padding: clamp(20px, 1.7cqw, 32px); border: 1px solid var(--line); border-radius: 22px; background: linear-gradient(145deg, rgba(24, 36, 58, .92), rgba(13, 20, 34, .92)); box-shadow: var(--shadow); }
.import-card__icon { display: grid; place-items: center; width: 54px; height: 54px; border-radius: 16px; color: #9de9ff; background: linear-gradient(135deg, rgba(86, 204, 255, .18), rgba(153, 120, 255, .18)); font-size: .85rem; font-weight: 900; }
.import-card h2 { margin: 18px 0 7px; font-size: 1.15rem; }
.import-card p { margin: 0; color: var(--muted); font-size: .85rem; line-height: 1.65; }
.import-card__tag { float: right; margin-top: -42px; padding: 4px 8px; border-radius: 7px; color: #f6d58d; background: rgba(239, 179, 75, .12); font-size: .65rem; }
.import-card button { display: flex; align-items: center; justify-content: space-between; min-height: 46px; margin-top: 20px; padding: 0 14px; border: 1px solid rgba(255, 255, 255, .1); border-radius: 12px; color: #fff; background: rgba(255, 255, 255, .06); cursor: pointer; }
.import-card button:hover { border-color: rgba(117, 222, 255, .4); background: rgba(107, 205, 255, .13); }
.import-view__result { border: 1px solid var(--line); border-radius: 22px; background: rgba(12, 19, 32, .72); }
.import-view__result-title { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px; border-bottom: 1px solid var(--line); }
.import-view__result-title h2 { margin: 0 0 4px; font-size: 1rem; }
.import-view__result-title p, .import-view__result-title span { margin: 0; color: var(--muted); font-size: .75rem; }
.import-view__empty { display: grid; place-items: center; padding: 34px 20px 42px; text-align: center; }
.import-view__empty-icon { display: grid; place-items: center; width: 58px; height: 58px; border: 1px dashed #52627a; border-radius: 18px; color: #8fa0b8; font-size: 1.8rem; }
.import-view__empty h3 { margin: 15px 0 5px; }
.import-view__empty p { margin: 0; color: var(--muted); font-size: .8rem; }
@container levabox-app (max-width: 900px) { .import-view__actions { grid-template-columns: 1fr; } }
@container levabox-app (max-height: 760px) { .import-view__actions { margin-top: 20px; } .import-card { min-height: 215px; } }
</style>
