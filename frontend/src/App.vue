<script setup lang="ts">
import { computed, ref } from 'vue'
import SidebarMenu from './components/SidebarMenu.vue'
import { mockGames } from './data/mockGames'
import GameDetailView from './views/GameDetailView.vue'
import HomeView from './views/HomeView.vue'
import ImportView from './views/ImportView.vue'
import LibraryView from './views/LibraryView.vue'
import SettingsView from './views/SettingsView.vue'
import type { Game } from './types/game'
import type { MainPage, Page } from './types/navigation'

const currentPage = ref<Page>('home')
const previousPage = ref<MainPage>('home')
const selectedGame = ref<Game>(mockGames[0]!)
const menuOpen = ref(false)
const toastMessage = ref('')
let toastTimer: ReturnType<typeof setTimeout> | undefined

const activeMainPage = computed<MainPage>(() => currentPage.value === 'detail' ? previousPage.value : currentPage.value)

function navigate(page: MainPage) {
  currentPage.value = page
  previousPage.value = page
  menuOpen.value = false
}

function selectGame(game: Game) {
  selectedGame.value = game
}

function openDetail(game: Game) {
  selectedGame.value = game
  if (currentPage.value !== 'detail') previousPage.value = currentPage.value
  currentPage.value = 'detail'
  menuOpen.value = false
}

function returnFromDetail() {
  currentPage.value = previousPage.value
}

function notify(message: string) {
  toastMessage.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toastMessage.value = '' }, 2600)
}
</script>

<template>
  <div class="app-shell">
    <button v-if="currentPage !== 'home'" class="menu-trigger" :class="{ 'menu-trigger--hidden': menuOpen }" aria-label="打开菜单" @click="menuOpen = true">
      <span></span><span></span><span></span>
    </button>

    <SidebarMenu :open="menuOpen" :current-page="activeMainPage" @close="menuOpen = false" @navigate="navigate" />

    <HomeView
      v-if="currentPage === 'home'"
      :games="mockGames"
      :selected-game="selectedGame"
      @select-game="selectGame"
      @open-detail="openDetail"
      @notify="notify"
      @open-menu="menuOpen = true"
      @navigate="navigate"
    />
    <LibraryView v-else-if="currentPage === 'library'" :games="mockGames" @open-detail="openDetail" />
    <ImportView v-else-if="currentPage === 'import'" @notify="notify" />
    <SettingsView v-else-if="currentPage === 'settings'" />
    <GameDetailView v-else :game="selectedGame" @back="returnFromDetail" @notify="notify" />

    <Transition name="toast">
      <div v-if="toastMessage" class="app-toast" role="status"><span>i</span>{{ toastMessage }}</div>
    </Transition>
  </div>
</template>
