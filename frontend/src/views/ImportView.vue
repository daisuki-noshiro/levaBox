<script setup lang="ts">
import { ref } from 'vue'

import {
  prepareImportMetadata,
  saveImport,
  selectExecutable,
  startImport,
} from '../services/importService'
import { metadata, service } from '../../wailsjs/go/models'

const emit = defineEmits<{ notify: [message: string] }>()

const draft = ref<service.ImportDraft | null>(null)
const loading = ref(false)
const metadataLoading = ref(false)
const saveLoading = ref(false)
const metadataQueried = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const alreadyExists = ref(false)
const existingGameId = ref('')
const issues = ref<service.MetadataSourceIssue[]>([])
const selectedCover = ref<metadata.ImageCandidate | null>(null)
const selectedBackground = ref<metadata.ImageCandidate | null>(null)
const selectedTags = ref<string[]>([])
const customTagInput = ref('')
const failedImages = ref<Set<string>>(new Set())

const metadataSources: metadata.Source[] = [
  metadata.Source.VNDB,
  metadata.Source.Bangumi,
]

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

function resetMetadataState(): void {
  metadataQueried.value = false
  issues.value = []
  selectedCover.value = null
  selectedBackground.value = null
  selectedTags.value = []
  customTagInput.value = ''
  failedImages.value = new Set()
}

function tagKey(tag: string): string {
  return tag.trim().toLocaleLowerCase()
}

function isTagSelected(tag: string): boolean {
  const key = tagKey(tag)
  return selectedTags.value.some((selected) => tagKey(selected) === key)
}

function toggleTag(tag: string): void {
  const value = tag.trim()
  if (!value) return

  if (isTagSelected(value)) {
    selectedTags.value = selectedTags.value.filter(
      (selected) => tagKey(selected) !== tagKey(value),
    )
    return
  }

  selectedTags.value = [...selectedTags.value, value]
}

function addCustomTag(): void {
  const value = customTagInput.value.trim()
  if (!value) return

  if (!isTagSelected(value)) {
    selectedTags.value = [...selectedTags.value, value]
  }
  customTagInput.value = ''
}

function sourceLabel(source: metadata.Source): string {
  if (source === metadata.Source.VNDB) return 'VNDB'
  if (source === metadata.Source.Bangumi) return 'Bangumi'
  return source
}

function imageSource(candidate: metadata.ImageCandidate): string {
  return candidate.Thumbnail?.trim() || candidate.URL
}

function imageResolution(candidate: metadata.ImageCandidate): string {
  return candidate.Width !== undefined && candidate.Height !== undefined
    ? `${candidate.Width} × ${candidate.Height}`
    : '分辨率未知'
}

function isSelected(
  candidate: metadata.ImageCandidate,
  selected: metadata.ImageCandidate | null,
): boolean {
  return selected?.URL === candidate.URL
}

function handleImageError(
  event: Event,
  candidate: metadata.ImageCandidate,
): void {
  const image = event.currentTarget as HTMLImageElement
  if (candidate.Thumbnail && image.dataset.fallbackApplied !== 'true') {
    image.dataset.fallbackApplied = 'true'
    image.src = candidate.URL
    return
  }

  const next = new Set(failedImages.value)
  next.add(candidate.URL)
  failedImages.value = next
}

function updateYear(event: Event): void {
  if (!draft.value) return
  const value = (event.target as HTMLInputElement).value
  draft.value.Year = value === '' ? undefined : Number(value)
}

function updateDescription(event: Event): void {
  if (!draft.value) return
  const value = (event.target as HTMLTextAreaElement).value
  draft.value.Description = value === '' ? undefined : value
}

async function handleSelectExecutable(): Promise<void> {
  if (loading.value || metadataLoading.value || saveLoading.value) return

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
    resetMetadataState()
    successMessage.value = ''
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

async function handlePrepareMetadata(): Promise<void> {
  if (!draft.value || metadataLoading.value || saveLoading.value) return

  metadataLoading.value = true
  errorMessage.value = ''

  try {
    const result = await prepareImportMetadata(draft.value, metadataSources)
    draft.value = result.Draft
    issues.value = result.Issues ?? []
    selectedCover.value = null
    selectedBackground.value = null
    failedImages.value = new Set()
    metadataQueried.value = true
  } catch (error) {
    errorMessage.value = `查询游戏资料失败：${errorReason(error)}`
  } finally {
    metadataLoading.value = false
  }
}

async function handleSaveImport(): Promise<void> {
  if (!draft.value || saveLoading.value || metadataLoading.value) return

  saveLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const request = new service.SaveImportRequest({
      ExecutablePath: draft.value.ExecutablePath,
      WorkingDirectory: draft.value.WorkingDirectory,
      Title: draft.value.Title,
      Company: draft.value.Company,
      Year: draft.value.Year ?? null,
      Description: draft.value.Description ?? null,
      Tags: [...selectedTags.value],
      Cover: selectedCover.value,
      Background: selectedBackground.value,
      Sources: draft.value.Sources ?? [],
    })

    const game = await saveImport(request)
    successMessage.value = `“${game.Title}”已成功导入游戏库。`
    emit('notify', `“${game.Title}”导入成功`)

    draft.value = null
    resetMetadataState()
    alreadyExists.value = false
    existingGameId.value = ''
  } catch (error) {
    errorMessage.value = `保存游戏失败：${errorReason(error)}`
  } finally {
    saveLoading.value = false
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
          :disabled="item.executable && (loading || metadataLoading || saveLoading)"
          @click="item.executable ? handleSelectExecutable() : notifyDeveloping(item.title)"
        >
          {{ item.executable && loading ? '处理中...' : item.action }} <span>→</span>
        </button>
      </article>
    </section>
    <section class="import-view__result">
      <div class="import-view__result-title"><div><h2>导入草稿</h2><p>检查启动路径，并调整下一阶段查询资料时使用的关键词。</p></div><span>{{ draft ? 1 : 0 }} 个项目</span></div>
      <div class="import-view__result-body">
        <p v-if="successMessage" class="import-view__message import-view__message--success" role="status">{{ successMessage }}</p>
        <p v-if="errorMessage" class="import-view__message import-view__message--error" role="alert">{{ errorMessage }}</p>

        <div v-if="draft" class="import-workspace">
          <div class="import-draft">
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

            <template v-if="metadataQueried">
              <label class="import-draft__field">
                <span>开发商</span>
                <input v-model="draft.Company" type="text" autocomplete="off" />
              </label>
              <label class="import-draft__field">
                <span>年份</span>
                <input :value="draft.Year ?? ''" type="number" inputmode="numeric" @input="updateYear" />
              </label>
              <label class="import-draft__field import-draft__field--wide">
                <span>简介</span>
                <textarea :value="draft.Description ?? ''" rows="6" @input="updateDescription"></textarea>
              </label>
            </template>

            <div class="import-query import-draft__field--wide">
              <button
                type="button"
                :disabled="metadataLoading || saveLoading || !draft.SearchKeyword.trim()"
                @click="handlePrepareMetadata"
              >
                {{ metadataLoading ? '正在查询...' : metadataQueried ? '重新查询资料' : '查询游戏资料' }}
              </button>
              <p>使用当前搜索关键词查询 VNDB 和 Bangumi，不会自动保存到游戏库。</p>
            </div>
          </div>

          <section v-if="issues.length" class="metadata-issues" aria-label="来源查询警告">
            <div><strong>部分来源查询失败</strong><span>其他来源的结果仍已保留。</span></div>
            <ul>
              <li v-for="(issue, index) in issues" :key="`${issue.Source}-${index}`">
                <strong>{{ sourceLabel(issue.Source) }}</strong>：{{ issue.Message }}
              </li>
            </ul>
          </section>

          <div class="metadata-results">
            <section class="metadata-panel">
              <div class="metadata-panel__title"><div><h3>标签</h3><p>候选标签默认不选择，也可以添加自定义标签。</p></div><span>已选择 {{ selectedTags.length }}</span></div>
              <div v-if="draft.TagCandidates?.length" class="tag-candidates">
                <button
                  v-for="tag in draft.TagCandidates"
                  :key="tag"
                  type="button"
                  :class="{ 'tag-candidate--selected': isTagSelected(tag) }"
                  :aria-pressed="isTagSelected(tag)"
                  @click="toggleTag(tag)"
                >
                  {{ tag }} <i v-if="isTagSelected(tag)">✓</i>
                </button>
              </div>
              <p v-else class="metadata-panel__empty">暂无标签候选</p>

              <form class="custom-tag" @submit.prevent="addCustomTag">
                <input v-model="customTagInput" type="text" autocomplete="off" placeholder="输入自定义标签" aria-label="自定义标签" />
                <button type="submit" :disabled="!customTagInput.trim()">添加</button>
              </form>

              <div v-if="selectedTags.length" class="selected-tags" aria-label="已选择标签">
                <button v-for="tag in selectedTags" :key="tagKey(tag)" type="button" @click="toggleTag(tag)">{{ tag }} <span>×</span></button>
              </div>
            </section>

            <section v-if="metadataQueried" class="metadata-panel">
              <div class="metadata-panel__title"><div><h3>封面候选</h3><p>可选择一张封面，也可以不使用封面。</p></div><button type="button" :class="{ 'selection-clear--active': selectedCover === null }" @click="selectedCover = null">不使用封面</button></div>
              <div v-if="draft.CoverCandidates?.length" class="image-candidates image-candidates--cover">
                <button
                  v-for="candidate in draft.CoverCandidates"
                  :key="candidate.URL"
                  type="button"
                  class="image-candidate image-candidate--cover"
                  :class="{ 'image-candidate--selected': isSelected(candidate, selectedCover) }"
                  :aria-pressed="isSelected(candidate, selectedCover)"
                  @click="selectedCover = candidate"
                >
                  <span class="image-candidate__preview">
                    <img v-if="!failedImages.has(candidate.URL)" :src="imageSource(candidate)" :alt="`${sourceLabel(candidate.Source)} 封面候选`" loading="lazy" @error="handleImageError($event, candidate)" />
                    <span v-else class="image-candidate__fallback">图片加载失败</span>
                    <i v-if="isSelected(candidate, selectedCover)">✓</i>
                  </span>
                  <span class="image-candidate__meta"><strong>{{ sourceLabel(candidate.Source) }}</strong><small>{{ imageResolution(candidate) }}</small></span>
                </button>
              </div>
              <p v-else class="metadata-panel__empty">暂无封面候选</p>
            </section>

            <section v-if="metadataQueried" class="metadata-panel">
              <div class="metadata-panel__title"><div><h3>背景候选</h3><p>可选择一张横向背景，也可以不使用背景。</p></div><button type="button" :class="{ 'selection-clear--active': selectedBackground === null }" @click="selectedBackground = null">不使用背景</button></div>
              <div v-if="draft.BackgroundCandidates?.length" class="image-candidates image-candidates--background">
                <button
                  v-for="candidate in draft.BackgroundCandidates"
                  :key="candidate.URL"
                  type="button"
                  class="image-candidate image-candidate--background"
                  :class="{ 'image-candidate--selected': isSelected(candidate, selectedBackground) }"
                  :aria-pressed="isSelected(candidate, selectedBackground)"
                  @click="selectedBackground = candidate"
                >
                  <span class="image-candidate__preview">
                    <img v-if="!failedImages.has(candidate.URL)" :src="imageSource(candidate)" :alt="`${sourceLabel(candidate.Source)} 背景候选`" loading="lazy" @error="handleImageError($event, candidate)" />
                    <span v-else class="image-candidate__fallback">图片加载失败</span>
                    <i v-if="isSelected(candidate, selectedBackground)">✓</i>
                  </span>
                  <span class="image-candidate__meta"><strong>{{ sourceLabel(candidate.Source) }}</strong><small>{{ imageResolution(candidate) }}</small></span>
                </button>
              </div>
              <p v-else class="metadata-panel__empty">暂无背景候选</p>
            </section>
          </div>

          <section class="import-confirm">
            <div>
              <h3>确认导入</h3>
              <p>将按当前资料创建 Game，并保存已选择的标签和媒体；不会加入游戏大厅或队列。</p>
              <span>{{ selectedTags.length }} 个标签 · {{ selectedCover ? '使用封面' : '无封面' }} · {{ selectedBackground ? '使用背景' : '无背景' }}</span>
            </div>
            <button type="button" :disabled="saveLoading || metadataLoading" @click="handleSaveImport">
              {{ saveLoading ? '正在导入...' : '确认导入' }}
            </button>
          </section>
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
.import-view__message--success { border: 1px solid rgba(113, 222, 181, .25); color: #bff3dc; background: rgba(65, 167, 127, .1); }
.import-workspace { display: grid; gap: 22px; }
.import-draft { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
.import-draft__field { display: grid; gap: 8px; min-width: 0; color: #b8c7d9; font-size: .76rem; }
.import-draft__field--wide { grid-column: 1 / -1; }
.import-draft__field input, .import-draft__field textarea { width: 100%; min-width: 0; padding: 0 13px; border: 1px solid rgba(255, 255, 255, .11); border-radius: 11px; outline: 0; color: #f4f8fd; background: rgba(255, 255, 255, .055); font: inherit; }
.import-draft__field input { height: 44px; }
.import-draft__field textarea { padding-top: 11px; padding-bottom: 11px; resize: vertical; line-height: 1.55; }
.import-draft__field input:not([readonly]):focus, .import-draft__field textarea:focus { border-color: rgba(117, 222, 255, .55); box-shadow: 0 0 0 3px rgba(117, 222, 255, .08); }
.import-draft__field input[readonly] { color: #9caabd; cursor: default; }
.import-query { display: flex; align-items: center; gap: 14px; padding-top: 2px; }
.import-query button { flex: 0 0 auto; min-height: 44px; padding: 0 18px; border: 1px solid rgba(117, 222, 255, .32); border-radius: 11px; color: #dff8ff; background: rgba(74, 177, 219, .13); cursor: pointer; }
.import-query button:hover:not(:disabled) { border-color: rgba(117, 222, 255, .58); background: rgba(74, 177, 219, .2); }
.import-query button:disabled { opacity: .55; cursor: wait; }
.import-query p { margin: 0; color: var(--muted); font-size: .73rem; line-height: 1.5; }
.metadata-issues { display: grid; gap: 12px; padding: 16px 18px; border: 1px solid rgba(238, 190, 92, .24); border-radius: 14px; color: #f3d99e; background: rgba(194, 137, 40, .09); }
.metadata-issues > div { display: flex; flex-wrap: wrap; align-items: baseline; gap: 8px; }
.metadata-issues > div strong { font-size: .86rem; }
.metadata-issues > div span { color: #bfae8b; font-size: .72rem; }
.metadata-issues ul { display: grid; gap: 6px; margin: 0; padding-left: 19px; color: #d7c8aa; font-size: .75rem; line-height: 1.5; }
.metadata-results { display: grid; gap: 18px; }
.metadata-panel { min-width: 0; padding: 18px; border: 1px solid rgba(255, 255, 255, .08); border-radius: 16px; background: rgba(255, 255, 255, .025); }
.metadata-panel__title { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 15px; }
.metadata-panel__title h3 { margin: 0 0 4px; font-size: .9rem; }
.metadata-panel__title p { margin: 0; color: var(--muted); font-size: .72rem; }
.metadata-panel__title > span { color: #9eacbe; font-size: .72rem; }
.metadata-panel__title > button { flex: 0 0 auto; min-height: 34px; padding: 0 11px; border: 1px solid rgba(255, 255, 255, .11); border-radius: 9px; color: #aeb9c8; background: rgba(255, 255, 255, .045); cursor: pointer; }
.metadata-panel__title > button:hover, .metadata-panel__title > .selection-clear--active { border-color: rgba(117, 222, 255, .35); color: #dff8ff; background: rgba(74, 177, 219, .11); }
.metadata-panel__empty { margin: 0; padding: 18px 0; color: var(--muted); font-size: .76rem; text-align: center; }
.tag-candidates { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-candidates button { min-height: 31px; padding: 0 10px; border: 1px solid rgba(126, 214, 241, .17); border-radius: 999px; color: #badde8; background: rgba(72, 162, 190, .08); cursor: pointer; font-size: .72rem; }
.tag-candidates button:hover, .tag-candidates .tag-candidate--selected { border-color: rgba(126, 226, 255, .5); color: #e2f9ff; background: rgba(72, 178, 214, .18); }
.tag-candidates i { margin-left: 4px; color: #8ce9c3; font-style: normal; }
.custom-tag { display: flex; gap: 9px; margin-top: 16px; }
.custom-tag input { flex: 1 1 auto; min-width: 0; height: 40px; padding: 0 12px; border: 1px solid rgba(255, 255, 255, .11); border-radius: 10px; outline: 0; color: #f4f8fd; background: rgba(255, 255, 255, .05); }
.custom-tag input:focus { border-color: rgba(117, 222, 255, .5); }
.custom-tag > button { flex: 0 0 auto; min-width: 66px; border: 1px solid rgba(117, 222, 255, .25); border-radius: 10px; color: #d9f6ff; background: rgba(74, 177, 219, .11); cursor: pointer; }
.custom-tag > button:disabled { opacity: .48; cursor: default; }
.selected-tags { display: flex; flex-wrap: wrap; gap: 7px; margin-top: 12px; padding-top: 12px; border-top: 1px solid rgba(255, 255, 255, .065); }
.selected-tags button { min-height: 29px; padding: 0 9px; border: 1px solid rgba(139, 231, 195, .22); border-radius: 8px; color: #c7f0df; background: rgba(70, 166, 128, .09); cursor: pointer; font-size: .7rem; }
.selected-tags button:hover { border-color: rgba(139, 231, 195, .46); }
.selected-tags span { margin-left: 4px; color: #8ba99d; }
.image-candidates { display: grid; grid-auto-flow: column; gap: 12px; padding: 2px 2px 10px; overflow-x: auto; overscroll-behavior-inline: contain; }
.image-candidates--cover { grid-auto-columns: minmax(145px, 180px); }
.image-candidates--background { grid-auto-columns: minmax(245px, 300px); }
.image-candidate { min-width: 0; padding: 7px; border: 1px solid rgba(255, 255, 255, .09); border-radius: 13px; color: #edf4fa; background: rgba(8, 14, 24, .7); cursor: pointer; text-align: left; transition: border-color 150ms ease, background 150ms ease, transform 150ms ease; }
.image-candidate:hover { border-color: rgba(117, 222, 255, .34); transform: translateY(-2px); }
.image-candidate--selected { border-color: rgba(117, 222, 255, .72); background: rgba(65, 162, 198, .13); box-shadow: 0 0 0 2px rgba(117, 222, 255, .08); }
.image-candidate__preview { position: relative; display: grid; place-items: center; width: 100%; overflow: hidden; border-radius: 9px; color: #8e9bae; background: rgba(255, 255, 255, .045); }
.image-candidate--cover .image-candidate__preview { aspect-ratio: 2 / 3; }
.image-candidate--background .image-candidate__preview { aspect-ratio: 16 / 9; }
.image-candidate__preview img { width: 100%; height: 100%; object-fit: cover; }
.image-candidate__preview i { position: absolute; top: 7px; right: 7px; display: grid; place-items: center; width: 25px; height: 25px; border-radius: 50%; color: #09201a; background: #8ce9c3; box-shadow: 0 3px 12px rgba(0, 0, 0, .34); font-style: normal; font-weight: 900; }
.image-candidate__fallback { padding: 12px; font-size: .7rem; text-align: center; }
.image-candidate__meta { display: flex; align-items: center; justify-content: space-between; gap: 8px; padding: 9px 3px 2px; }
.image-candidate__meta strong { font-size: .72rem; }
.image-candidate__meta small { overflow: hidden; color: var(--muted); font-size: .65rem; text-overflow: ellipsis; white-space: nowrap; }
.import-confirm { display: flex; align-items: center; justify-content: space-between; gap: 24px; padding: 20px; border: 1px solid rgba(117, 222, 255, .19); border-radius: 16px; background: linear-gradient(135deg, rgba(51, 130, 164, .1), rgba(91, 75, 162, .08)); }
.import-confirm h3 { margin: 0 0 6px; font-size: .95rem; }
.import-confirm p { margin: 0 0 7px; color: var(--muted); font-size: .75rem; line-height: 1.55; }
.import-confirm span { color: #a8bacd; font-size: .7rem; }
.import-confirm > button { flex: 0 0 auto; min-width: 130px; min-height: 46px; border: 1px solid rgba(255, 255, 255, .25); border-radius: 12px; color: #07131d; background: linear-gradient(135deg, #e9fbff, #91e5ff); cursor: pointer; font-weight: 800; }
.import-confirm > button:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 9px 24px rgba(87, 210, 255, .14); }
.import-confirm > button:disabled { opacity: .55; cursor: wait; }
.import-view__existing { display: flex; align-items: center; gap: 16px; min-height: 96px; padding: 18px; border: 1px solid rgba(113, 222, 181, .2); border-radius: 15px; background: rgba(65, 167, 127, .08); }
.import-view__existing-icon { display: grid; place-items: center; flex: 0 0 auto; width: 42px; height: 42px; border-radius: 50%; color: #0c241b; background: #83e1ba; font-weight: 900; }
.import-view__existing h3 { margin: 0 0 5px; font-size: .95rem; }
.import-view__existing p { margin: 0; color: var(--muted); font-size: .76rem; }
.import-view__empty { display: grid; place-items: center; padding: 34px 20px 42px; text-align: center; }
.import-view__empty-icon { display: grid; place-items: center; width: 58px; height: 58px; border: 1px dashed #52627a; border-radius: 18px; color: #8fa0b8; font-size: 1.8rem; }
.import-view__empty h3 { margin: 15px 0 5px; }
.import-view__empty p { margin: 0; color: var(--muted); font-size: .8rem; }
@container levabox-app (max-width: 900px) { .import-view__actions { grid-template-columns: 1fr; } .import-draft { grid-template-columns: 1fr; } .import-query, .import-confirm { align-items: stretch; flex-direction: column; } .metadata-panel__title { align-items: flex-start; } .import-confirm > button { width: 100%; } }
@container levabox-app (max-height: 760px) { .import-view__actions { margin-top: 20px; } .import-card { min-height: 215px; } }
</style>
