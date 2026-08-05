<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import GlobalNavigation from './components/GlobalNavigation.vue'
import { mockGames } from './data/mockGames'
import GameDetailView from './views/GameDetailView.vue'
import HomeView from './views/HomeView.vue'
import ImportView from './views/ImportView.vue'
import LibraryView from './views/LibraryView.vue'
import SettingsView from './views/SettingsView.vue'
import { getBatteryLevel } from './services/deviceService'
import type { Game, QueueMode } from './types/game'
import type { MainPage, Page } from './types/navigation'

const currentPage = ref<Page>('home')
const previousPage = ref<MainPage>('home')
const originalImportOrder = new Map(mockGames.map((game, index) => [game.id, index]))
const homeQueue = ref<Game[]>(sortByDefaultRule(mockGames))
const selectedHomeGameId = ref<string | null>(homeQueue.value[0]?.id ?? null)
const detailGame = ref<Game | null>(null)
const queueMode = ref<QueueMode>('default')
const batteryLevel = ref(78)
const toastMessage = ref('')
let toastTimer: ReturnType<typeof setTimeout> | undefined

const activeMainPage = computed<MainPage>(() => currentPage.value === 'detail' ? previousPage.value : currentPage.value)
const selectedHomeGame = computed<Game | null>(() => {
  return homeQueue.value.find((game) => game.id === selectedHomeGameId.value) ?? homeQueue.value[0] ?? null
})

function navigate(page: MainPage) {
  currentPage.value = page
  previousPage.value = page
  if (page === 'home' && !selectedHomeGame.value) {
    selectedHomeGameId.value = homeQueue.value[0]?.id ?? null
  }
}

function selectGame(game: Game) {
  if (homeQueue.value.some((queuedGame) => queuedGame.id === game.id)) {
    selectedHomeGameId.value = game.id
  }
}

function openDetail(game: Game) {
  detailGame.value = game
  if (currentPage.value !== 'detail') previousPage.value = currentPage.value
  currentPage.value = 'detail'
}

function returnFromDetail() {
  currentPage.value = previousPage.value
}

function moveGameToFront(gameId: string) {
  const index = homeQueue.value.findIndex((game) => game.id === gameId)
  if (index <= 0) return
  const nextQueue = [...homeQueue.value]
  const [game] = nextQueue.splice(index, 1)
  if (!game) return
  nextQueue.unshift(game)
  homeQueue.value = nextQueue
  selectedHomeGameId.value = gameId
  queueMode.value = 'custom'
}

function commitQueueOrder(orderedIds: string[]) {
  const byId = new Map(homeQueue.value.map((game) => [game.id, game]))
  const nextQueue = orderedIds.map((id) => byId.get(id)).filter((game): game is Game => Boolean(game))
  if (nextQueue.length !== homeQueue.value.length) return
  const changed = nextQueue.some((game, index) => game.id !== homeQueue.value[index]?.id)
  if (!changed) return
  homeQueue.value = nextQueue
  queueMode.value = 'custom'
}

function removeGameFromQueue(gameId: string) {
  const index = homeQueue.value.findIndex((game) => game.id === gameId)
  if (index < 0) return
  const nextQueue = homeQueue.value.filter((game) => game.id !== gameId)
  homeQueue.value = nextQueue
  selectedHomeGameId.value = nextQueue[index]?.id ?? nextQueue[index - 1]?.id ?? null
}

function restoreDefaultQueue() {
  const selectedId = selectedHomeGameId.value
  homeQueue.value = sortByDefaultRule(homeQueue.value)
  queueMode.value = 'default'
  selectedHomeGameId.value = selectedId && homeQueue.value.some((game) => game.id === selectedId)
    ? selectedId
    : homeQueue.value[0]?.id ?? null
}

/** 后续游戏库的“加入队列”入口可直接调用此边界函数。 */
function addGameToQueue(gameId: string) {
  if (homeQueue.value.some((game) => game.id === gameId)) return
  const game = mockGames.find((item) => item.id === gameId)
  if (!game) return
  homeQueue.value = queueMode.value === 'custom'
    ? [...homeQueue.value, game]
    : sortByDefaultRule([...homeQueue.value, game])
  selectedHomeGameId.value ??= game.id
}

function sortByDefaultRule(games: Game[]): Game[] {
  return [...games].sort((a, b) => {
    const aPlayed = a.lastPlayedAt ? Date.parse(a.lastPlayedAt) : Number.NaN
    const bPlayed = b.lastPlayedAt ? Date.parse(b.lastPlayedAt) : Number.NaN
    const aHasPlayed = Number.isFinite(aPlayed)
    const bHasPlayed = Number.isFinite(bPlayed)
    if (aHasPlayed && bHasPlayed) return bPlayed - aPlayed
    if (aHasPlayed) return -1
    if (bHasPlayed) return 1
    return (originalImportOrder.get(a.id) ?? 0) - (originalImportOrder.get(b.id) ?? 0)
  })
}

function notify(message: string) {
  toastMessage.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toastMessage.value = '' }, 2600)
}

onMounted(async () => {
  batteryLevel.value = await getBatteryLevel()
})
</script>

<template>
  <div class="app-shell">
    <GlobalNavigation :current-page="activeMainPage" @navigate="navigate" />

    <HomeView
      v-if="currentPage === 'home'"
      :queue="homeQueue"
      :selected-game="selectedHomeGame"
      :battery-level="batteryLevel"
      @select-game="selectGame"
      @notify="notify"
      @navigate="navigate"
      @move-to-front="moveGameToFront"
      @commit-reorder="commitQueueOrder"
      @remove-from-queue="removeGameFromQueue"
      @restore-default="restoreDefaultQueue"
    />
    <LibraryView v-else-if="currentPage === 'library'" :games="mockGames" @open-detail="openDetail" />
    <ImportView v-else-if="currentPage === 'import'" @notify="notify" />
    <SettingsView v-else-if="currentPage === 'settings'" />
    <GameDetailView v-else-if="detailGame" :game="detailGame" @back="returnFromDetail" @notify="notify" />

    <Transition name="toast">
      <div v-if="toastMessage" class="app-toast" role="status"><span>i</span>{{ toastMessage }}</div>
    </Transition>
  </div>
</template>
