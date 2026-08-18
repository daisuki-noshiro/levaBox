<script setup lang="ts">
import { ref } from 'vue'

import { selectExecutable, startImport } from '../services/importService'
import type { service } from '../../wailsjs/go/models'

const emit = defineEmits<{ notify: [message: string] }>()

const draft = ref<service.ImportDraft | null>(null)
const loading = ref(false)
const errorMessage = ref('')
const alreadyExists = ref(false)
const existingGameId = ref('')

const actions = [
  { icon: '▣', title: '扫描文件夹', description: '选择一个目录，自动寻找可能的游戏程序。', action: '扫描文件夹', developing: true },
  { icon: 'EXE', title: '手动选择 EXE', description: '直接指定游戏的主启动程序。', action: '选择 EXE', executable: true },
  { icon: 'ZIP', title: '导入压缩包', description: '从本地压缩文件创建待安装项目。', action: '导入压缩包', developing: true },
]

function errorReason(error: unknown): string {
  return error instanceof Error ? error.message : String(error)
}

function notifyDeveloping(title: string): void {
  emit('notify', `${title}功能正在开发中`)
}

async function handleSelectExecutable(): Promise<void> {
  if (loading.value) return

  loading.value = true
  errorMessage.value = ''

  try {
    let executablePath: string
    try {
      executablePath = await selectExecutable()
    } catch (error) {
      errorMessage.value = `选择游戏失败：${errorReason(error)}`
      return
    }

    if (!executablePath) return

    draft.value = null
    alreadyExists.value = false
    existingGameId.value = ''

    let result: service.StartImportResult
    try {
      result = await startImport(executablePath)
    } catch (error) {
      errorMessage.value = `导入准备失败：${errorReason(error)}`
      return
    }

    if (result.status === 'ready' && result.draft) {
      draft.value = result.draft
      return
    }

    if (result.status === 'already_exists') {
      alreadyExists.value = true
      existingGameId.value = result.existingGameId ?? ''
      return
    }

    errorMessage.value = result.status === 'ready'
      ? '导入准备失败：后端未返回导入草稿。'
      : `导入准备失败：无法识别的状态“${result.status}”。`
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="import-view page-shell">
    <header class="page-header">
      <div><p class="page-kicker">ADD GAMES</p><h1>导入游戏</h1><p>把本地 Galgame 加入游戏库。</p></div>
      <span class="import-view__badge">导入游戏</span>
    </header>
    <section class="import-view__actions">
      <article v-for="item in actions" :key="item.title" class="import-card">
        <span class="import-card__icon">{{ item.icon }}</span>
        <div><span v-if="item.developing" class="import-card__tag">开发中</span><h2>{{ item.title }}</h2><p>{{ item.description }}</p></div>
        <button
          :disabled="item.executable && loading"
          @click="item.executable ? handleSelectExecutable() : notifyDeveloping(item.title)"
        >
          {{ item.executable && loading ? '处理中...' : item.action }} <span>→</span>
        </button>
      </article>
    </section>
    <section class="import-view__result">
      <div class="import-view__result-title"><div><h2>导入草稿</h2><p>检查启动路径，并调整下一阶段查询资料时使用的关键词。</p></div><span>{{ draft ? 1 : 0 }} 个项目</span></div>
      <div class="import-view__result-body">
        <p v-if="errorMessage" class="import-view__message import-view__message--error" role="alert">{{ errorMessage }}</p>

        <div v-if="draft" class="import-draft">
          <label class="import-draft__field import-draft__field--wide">
            <span>启动程序</span>
            <input :value="draft.ExecutablePath" type="text" readonly />
          </label>
          <label class="import-draft__field import-draft__field--wide">
            <span>工作目录</span>
            <input :value="draft.WorkingDirectory" type="text" readonly />
          </label>
          <label class="import-draft__field">
            <span>搜索关键词</span>
            <input v-model="draft.SearchKeyword" type="text" autocomplete="off" />
          </label>
          <label class="import-draft__field">
            <span>标题</span>
            <input v-model="draft.Title" type="text" autocomplete="off" />
          </label>
        </div>

        <div v-else-if="alreadyExists" class="import-view__existing">
          <span class="import-view__existing-icon">✓</span>
          <div>
            <h3>该启动程序已经存在于游戏库中</h3>
            <p v-if="existingGameId">游戏记录 ID：{{ existingGameId }}</p>
            <p v-else>无需再次导入，可以前往游戏库查看。</p>
          </div>
        </div>

        <div v-else class="import-view__empty"><span class="import-view__empty-icon">＋</span><h3>还没有待导入的游戏</h3><p>选择一个 EXE 开始导入。</p></div>
      </div>
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
.import-card button:disabled { opacity: .58; cursor: wait; }
.import-card button:disabled:hover { border-color: rgba(255, 255, 255, .1); background: rgba(255, 255, 255, .06); }
.import-view__result { border: 1px solid var(--line); border-radius: 22px; background: rgba(12, 19, 32, .72); }
.import-view__result-title { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px; border-bottom: 1px solid var(--line); }
.import-view__result-title h2 { margin: 0 0 4px; font-size: 1rem; }
.import-view__result-title p, .import-view__result-title span { margin: 0; color: var(--muted); font-size: .75rem; }
.import-view__result-body { padding: 24px; }
.import-view__message { margin: 0 0 18px; padding: 12px 14px; border-radius: 12px; font-size: .8rem; line-height: 1.5; }
.import-view__message--error { border: 1px solid rgba(255, 132, 132, .28); color: #ffc3c3; background: rgba(176, 55, 55, .12); }
.import-draft { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
.import-draft__field { display: grid; gap: 8px; min-width: 0; color: #b8c7d9; font-size: .76rem; }
.import-draft__field--wide { grid-column: 1 / -1; }
.import-draft__field input { width: 100%; min-width: 0; height: 44px; padding: 0 13px; border: 1px solid rgba(255, 255, 255, .11); border-radius: 11px; outline: 0; color: #f4f8fd; background: rgba(255, 255, 255, .055); }
.import-draft__field input:not([readonly]):focus { border-color: rgba(117, 222, 255, .55); box-shadow: 0 0 0 3px rgba(117, 222, 255, .08); }
.import-draft__field input[readonly] { color: #9caabd; cursor: default; }
.import-view__existing { display: flex; align-items: center; gap: 16px; min-height: 96px; padding: 18px; border: 1px solid rgba(113, 222, 181, .2); border-radius: 15px; background: rgba(65, 167, 127, .08); }
.import-view__existing-icon { display: grid; place-items: center; flex: 0 0 auto; width: 42px; height: 42px; border-radius: 50%; color: #0c241b; background: #83e1ba; font-weight: 900; }
.import-view__existing h3 { margin: 0 0 5px; font-size: .95rem; }
.import-view__existing p { margin: 0; color: var(--muted); font-size: .76rem; }
.import-view__empty { display: grid; place-items: center; padding: 34px 20px 42px; text-align: center; }
.import-view__empty-icon { display: grid; place-items: center; width: 58px; height: 58px; border: 1px dashed #52627a; border-radius: 18px; color: #8fa0b8; font-size: 1.8rem; }
.import-view__empty h3 { margin: 15px 0 5px; }
.import-view__empty p { margin: 0; color: var(--muted); font-size: .8rem; }
@container levabox-app (max-width: 900px) { .import-view__actions { grid-template-columns: 1fr; } .import-draft { grid-template-columns: 1fr; } }
@container levabox-app (max-height: 760px) { .import-view__actions { margin-top: 20px; } .import-card { min-height: 215px; } }
</style>
