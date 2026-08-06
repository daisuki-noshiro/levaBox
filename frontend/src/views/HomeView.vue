<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import DeviceStatus from '../components/DeviceStatus.vue'
import GameCard from '../components/GameCard.vue'
import { playUiSound, unlockMediaAudio } from '../services/audioService'
import { launchGame } from '../services/gameLauncher'
import { registerLobbyGamepadInput } from '../services/inputService'
import type { Game } from '../types/game'
import type { MainPage } from '../types/navigation'

type MenuActionId = 'launch' | 'front' | 'reorder' | 'remove' | 'restore'

const props = defineProps<{
  queue: Game[]
  selectedGame: Game | null
  batteryLevel: number
}>()

const emit = defineEmits<{
  selectGame: [game: Game]
  notify: [message: string]
  navigate: [page: MainPage]
  moveToFront: [gameId: string]
  commitReorder: [orderedIds: string[]]
  removeFromQueue: [gameId: string]
  restoreDefault: []
}>()

const actionMenuOpen = ref(false)
const confirmRestoreOpen = ref(false)
const focusedActionIndex = ref(0)
const confirmChoice = ref<0 | 1>(0)
const actionButtons = ref<HTMLButtonElement[]>([])
const confirmButtons = ref<HTMLButtonElement[]>([])
const cardsViewport = ref<HTMLElement | null>(null)
const cardsTrack = ref<HTMLElement | null>(null)
const trackOffset = ref(0)
const trackMotionReady = ref(false)
const reorderGameId = ref<string | null>(null)
const reorderDraft = ref<Game[]>([])
const reorderOriginalIds = ref<string[]>([])
const pointerDragGameId = ref<string | null>(null)
const pointerDragId = ref<number | null>(null)
const pointerDragMoved = ref(false)
let unregisterGamepadInput: () => void = () => undefined
let trackResizeObserver: ResizeObserver | undefined
let trackReadyFrame = 0

const isReordering = computed(() => reorderGameId.value !== null)
const displayQueue = computed(() => isReordering.value ? reorderDraft.value : props.queue)
const selectedIndex = computed(() => props.selectedGame
  ? props.queue.findIndex((game) => game.id === props.selectedGame?.id)
  : -1)

const menuActions = computed<Array<{ id: MenuActionId; label: string; primary?: boolean }>>(() => {
  const actions: Array<{ id: MenuActionId; label: string; primary?: boolean }> = [
    { id: 'launch', label: '开始游戏', primary: true },
  ]
  if (selectedIndex.value > 0) actions.push({ id: 'front', label: '设为队首' })
  actions.push(
    { id: 'reorder', label: '顺序调整' },
    { id: 'remove', label: '移除队列' },
    { id: 'restore', label: '恢复默认队列' },
  )
  return actions
})

function setActionButton(element: unknown, index: number) {
  if (element instanceof HTMLButtonElement) actionButtons.value[index] = element
}

function setConfirmButton(element: unknown, index: number) {
  if (element instanceof HTMLButtonElement) confirmButtons.value[index] = element
}

function navigate(page: MainPage) {
  playUiSound('confirm')
  emit('navigate', page)
}

function handleGameClick(game: Game) {
  if (isReordering.value) return
  unlockMediaAudio()
  if (props.selectedGame?.id === game.id) {
    openActionMenu()
    return
  }
  playUiSound('select')
  emit('selectGame', game)
}

function selectAdjacent(direction: -1 | 1) {
  if (!props.selectedGame || !props.queue.length) return
  const index = props.queue.findIndex((game) => game.id === props.selectedGame?.id)
  const nextIndex = Math.min(props.queue.length - 1, Math.max(0, index + direction))
  const nextGame = props.queue[nextIndex]
  if (!nextGame || nextGame.id === props.selectedGame.id) return
  playUiSound('select')
  emit('selectGame', nextGame)
  focusGameCard(nextGame.id)
}

function focusGameCard(gameId: string) {
  nextTick(() => {
    cardsTrack.value
      ?.querySelector<HTMLElement>(`[data-game-id="${gameId}"]`)
      ?.focus({ preventScroll: true })
  })
}

function openActionMenu() {
  if (!props.selectedGame || isReordering.value || confirmRestoreOpen.value) return
  playUiSound('open')
  actionMenuOpen.value = true
  focusedActionIndex.value = 0
  actionButtons.value = []
  nextTick(() => actionButtons.value[0]?.focus())
}

function closeActionMenu() {
  if (!actionMenuOpen.value) return
  playUiSound('cancel')
  actionMenuOpen.value = false
}

function focusAction(index: number) {
  const count = menuActions.value.length
  if (!count) return
  focusedActionIndex.value = (index + count) % count
  nextTick(() => actionButtons.value[focusedActionIndex.value]?.focus())
}

async function executeAction(action: MenuActionId) {
  const game = props.selectedGame
  if (!game) return
  if (action === 'launch') {
    playUiSound('launch')
    const result = await launchGame(game)
    emit('notify', result.message)
    actionMenuOpen.value = false
    return
  }
  if (action === 'front') {
    playUiSound('confirm')
    emit('moveToFront', game.id)
    actionMenuOpen.value = false
    return
  }
  if (action === 'reorder') {
    beginReorder(game)
    return
  }
  if (action === 'remove') {
    playUiSound('confirm')
    emit('removeFromQueue', game.id)
    actionMenuOpen.value = false
    return
  }
  actionMenuOpen.value = false
  openRestoreConfirmation()
}

function beginReorder(game: Game) {
  playUiSound('confirm')
  reorderOriginalIds.value = props.queue.map((item) => item.id)
  reorderDraft.value = [...props.queue]
  reorderGameId.value = game.id
  actionMenuOpen.value = false
}

function moveReorderingGame(direction: -1 | 1) {
  if (!reorderGameId.value) return
  const currentIndex = reorderDraft.value.findIndex((game) => game.id === reorderGameId.value)
  const targetIndex = Math.min(reorderDraft.value.length - 1, Math.max(0, currentIndex + direction))
  if (currentIndex < 0 || currentIndex === targetIndex) return
  moveDraftItem(currentIndex, targetIndex)
  playUiSound('select')
}

function moveDraftItem(fromIndex: number, toIndex: number) {
  const nextQueue = [...reorderDraft.value]
  const [movingGame] = nextQueue.splice(fromIndex, 1)
  if (!movingGame) return
  nextQueue.splice(toIndex, 0, movingGame)
  reorderDraft.value = nextQueue
}

function confirmReorder() {
  if (!isReordering.value) return
  const nextIds = reorderDraft.value.map((game) => game.id)
  const changed = nextIds.some((id, index) => id !== reorderOriginalIds.value[index])
  if (changed) emit('commitReorder', nextIds)
  playUiSound('confirm')
  clearReorderState()
}

function cancelReorder() {
  if (!isReordering.value) return
  playUiSound('cancel')
  clearReorderState()
}

function clearReorderState() {
  reorderGameId.value = null
  reorderDraft.value = []
  reorderOriginalIds.value = []
  pointerDragGameId.value = null
  pointerDragId.value = null
  pointerDragMoved.value = false
}

function handlePointerDown(game: Game, event: PointerEvent) {
  if (game.id !== reorderGameId.value || event.button !== 0) return
  pointerDragGameId.value = game.id
  pointerDragId.value = event.pointerId
  pointerDragMoved.value = false
}

function handlePointerMove(event: PointerEvent) {
  if (pointerDragId.value !== event.pointerId || !pointerDragGameId.value) return
  const targetCard = document
    .elementFromPoint(event.clientX, event.clientY)
    ?.closest<HTMLElement>('.game-card--rail')
  const targetId = targetCard?.dataset.gameId
  if (!targetId || targetId === pointerDragGameId.value) return
  const fromIndex = reorderDraft.value.findIndex((item) => item.id === pointerDragGameId.value)
  const toIndex = reorderDraft.value.findIndex((item) => item.id === targetId)
  if (fromIndex < 0 || toIndex < 0 || fromIndex === toIndex) return
  moveDraftItem(fromIndex, toIndex)
  pointerDragMoved.value = true
}

function handlePointerUp(event: PointerEvent) {
  if (pointerDragId.value !== event.pointerId) return
  const shouldCommit = pointerDragMoved.value
  pointerDragGameId.value = null
  pointerDragId.value = null
  pointerDragMoved.value = false
  if (shouldCommit) confirmReorder()
}

function handlePointerCancel(event: PointerEvent) {
  if (pointerDragId.value !== event.pointerId) return
  pointerDragGameId.value = null
  pointerDragId.value = null
  pointerDragMoved.value = false
  cancelReorder()
}

function openRestoreConfirmation() {
  playUiSound('open')
  confirmRestoreOpen.value = true
  confirmChoice.value = 0
  confirmButtons.value = []
  nextTick(() => confirmButtons.value[0]?.focus())
}

function closeRestoreConfirmation() {
  if (!confirmRestoreOpen.value) return
  playUiSound('cancel')
  confirmRestoreOpen.value = false
}

function confirmRestore() {
  playUiSound('confirm')
  emit('restoreDefault')
  confirmRestoreOpen.value = false
}

function handleKeydown(event: KeyboardEvent) {
  const target = event.target instanceof HTMLElement ? event.target : null
  const isTextEntry = Boolean(target?.closest('input, select, textarea, [contenteditable="true"]'))
  const focusedButton = target?.closest<HTMLButtonElement>('button') ?? null
  const focusedGameCard = focusedButton?.matches('.game-card--rail') ? focusedButton : null
  const overlayOpen = actionMenuOpen.value || confirmRestoreOpen.value || isReordering.value

  if (isTextEntry) return
  if (!overlayOpen && focusedButton && !focusedGameCard) return

  if (event.key === 'Escape') {
    if (isReordering.value || confirmRestoreOpen.value || actionMenuOpen.value) event.preventDefault()
    if (isReordering.value) cancelReorder()
    else if (confirmRestoreOpen.value) closeRestoreConfirmation()
    else if (actionMenuOpen.value) closeActionMenu()
    return
  }

  if (actionMenuOpen.value && (event.key === 'ArrowUp' || event.key === 'ArrowDown')) {
    event.preventDefault()
    focusAction(focusedActionIndex.value + (event.key === 'ArrowDown' ? 1 : -1))
    return
  }

  if (confirmRestoreOpen.value && (event.key === 'ArrowLeft' || event.key === 'ArrowRight')) {
    event.preventDefault()
    confirmChoice.value = confirmChoice.value === 0 ? 1 : 0
    nextTick(() => confirmButtons.value[confirmChoice.value]?.focus())
    return
  }

  if (event.key === 'Enter') {
    if (!overlayOpen && focusedGameCard?.dataset.gameId !== props.selectedGame?.id) return
    event.preventDefault()
    if (confirmRestoreOpen.value) {
      confirmChoice.value === 0 ? confirmRestore() : closeRestoreConfirmation()
    } else if (isReordering.value) {
      confirmReorder()
    } else if (actionMenuOpen.value) {
      void executeAction(menuActions.value[focusedActionIndex.value]?.id ?? 'launch')
    } else {
      openActionMenu()
    }
    return
  }

  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight') return
  event.preventDefault()
  const direction = event.key === 'ArrowRight' ? 1 : -1
  if (isReordering.value) moveReorderingGame(direction)
  else if (!actionMenuOpen.value && !confirmRestoreOpen.value) selectAdjacent(direction)
}

function updateTrackPosition() {
  const selectedId = props.selectedGame?.id
  if (!selectedId || !cardsViewport.value || !cardsTrack.value) return
  nextTick(() => {
    const viewport = cardsViewport.value
    const track = cardsTrack.value
    const card = track?.querySelector<HTMLElement>(`[data-game-id="${selectedId}"]`)
    if (!viewport || !track || !card) return

    const focusRatio = Number.parseFloat(
      getComputedStyle(viewport).getPropertyValue('--rail-focus-ratio'),
    ) || .34
    const focusPoint = viewport.clientWidth * focusRatio
    const leadingSpace = Math.max(0, focusPoint - card.offsetWidth / 2)
    const trailingSpace = Math.max(0, viewport.clientWidth - focusPoint - card.offsetWidth / 2)
    const leadingValue = `${Math.round(leadingSpace)}px`
    const trailingValue = `${Math.round(trailingSpace)}px`

    if (track.style.getPropertyValue('--rail-leading-space') !== leadingValue) {
      track.style.setProperty('--rail-leading-space', leadingValue)
    }
    if (track.style.getPropertyValue('--rail-trailing-space') !== trailingValue) {
      track.style.setProperty('--rail-trailing-space', trailingValue)
    }

    trackOffset.value = Math.round(focusPoint - (card.offsetLeft + card.offsetWidth / 2))

    if (!trackMotionReady.value && !trackReadyFrame) {
      trackReadyFrame = requestAnimationFrame(() => {
        trackReadyFrame = requestAnimationFrame(() => {
          trackMotionReady.value = true
          trackReadyFrame = 0
        })
      })
    }
  })
}

watch(() => props.selectedGame?.id, updateTrackPosition)
watch(() => displayQueue.value.map((game) => game.id).join('|'), updateTrackPosition)

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('pointermove', handlePointerMove)
  window.addEventListener('pointerup', handlePointerUp)
  window.addEventListener('pointercancel', handlePointerCancel)
  unregisterGamepadInput = registerLobbyGamepadInput({
    moveLeft: () => isReordering.value ? moveReorderingGame(-1) : selectAdjacent(-1),
    moveRight: () => isReordering.value ? moveReorderingGame(1) : selectAdjacent(1),
    confirm: () => window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter' })),
    cancel: () => window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' })),
    openMenu: openActionMenu,
  })
  trackResizeObserver = new ResizeObserver(updateTrackPosition)
  if (cardsViewport.value) trackResizeObserver.observe(cardsViewport.value)
  if (cardsTrack.value) trackResizeObserver.observe(cardsTrack.value)
  updateTrackPosition()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('pointermove', handlePointerMove)
  window.removeEventListener('pointerup', handlePointerUp)
  window.removeEventListener('pointercancel', handlePointerCancel)
  unregisterGamepadInput()
  trackResizeObserver?.disconnect()
  if (trackReadyFrame) cancelAnimationFrame(trackReadyFrame)
})
</script>

<template>
  <section class="home" :class="{ 'home--modal-open': actionMenuOpen || confirmRestoreOpen }">
    <div v-if="selectedGame" class="home__media" aria-hidden="true">
      <Transition name="media-fade" mode="out-in">
        <video
          v-if="selectedGame.backgroundType === 'video' && selectedGame.backgroundVideo"
          :key="`${selectedGame.id}-video`"
          class="home__backdrop"
          :src="selectedGame.backgroundVideo"
          :poster="selectedGame.backgroundImage"
          :style="{ objectPosition: selectedGame.backgroundPosition ?? 'center' }"
          autoplay muted loop playsinline preload="metadata"
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

    <DeviceStatus :battery-level="batteryLevel" />

    <div v-if="queue.length" class="home__lower-content">
      <Transition name="hero" mode="out-in">
        <article v-if="selectedGame" :key="selectedGame.id" class="home__info">
          <img v-if="selectedGame.logo" class="home__logo" :src="selectedGame.logo" :alt="selectedGame.title" />
          <h1 v-else>{{ selectedGame.title }}</h1>
          <p class="home__company">{{ selectedGame.developer }} · {{ selectedGame.year }}</p>
          <div class="home__details">
            <span v-for="genre in selectedGame.genres" :key="genre" class="home__tag">{{ genre }}</span>
            <span class="home__status"><i></i>{{ selectedGame.status }}</span>
          </div>
        </article>
      </Transition>

      <section class="home__carousel" aria-label="游戏队列">
        <header class="home__carousel-header">
          <div><span></span><strong>游戏队列</strong></div>
          <p v-if="isReordering">← / → 调整位置　Enter 确认　Esc 取消</p>
        </header>
        <div ref="cardsViewport" class="home__cards-viewport">
          <div
            ref="cardsTrack"
            class="home__cards-track"
            :class="{ 'home__cards-track--ready': trackMotionReady }"
            :style="{ transform: `translate3d(${trackOffset}px, 0, 0)` }"
          >
            <TransitionGroup name="queue">
              <GameCard
                v-for="game in displayQueue"
                :key="game.id"
                :game="game"
                variant="rail"
                :selected="game.id === selectedGame?.id"
                :reordering="isReordering"
                :moving="game.id === reorderGameId"
                @select="handleGameClick"
                @pointer-down="handlePointerDown"
              />
            </TransitionGroup>
          </div>
        </div>
      </section>
    </div>

    <section v-else class="home__empty">
      <span>＋</span><h1>游戏队列为空</h1><p>游戏资料仍保留在游戏库中。</p>
      <button @click="navigate('library')">打开游戏库</button>
    </section>

    <Transition name="dialog-fade">
      <div v-if="actionMenuOpen && selectedGame" class="home-dialog-layer" role="presentation" @click.self="closeActionMenu">
        <section class="game-menu" role="dialog" aria-modal="true" :aria-label="`${selectedGame.title}操作菜单`">
          <header><small>队列操作</small><h2>{{ selectedGame.title }}</h2></header>
          <div class="game-menu__actions">
            <button
              v-for="(action, index) in menuActions"
              :key="action.id"
              :ref="(element) => setActionButton(element, index)"
              :class="{ 'game-menu__action--primary': action.primary }"
              @focus="focusedActionIndex = index"
              @click="executeAction(action.id)"
            >{{ action.label }}</button>
          </div>
          <footer>Enter 确认　Esc 返回</footer>
        </section>
      </div>
    </Transition>

    <Transition name="dialog-fade">
      <div v-if="confirmRestoreOpen" class="home-dialog-layer" role="presentation" @click.self="closeRestoreConfirmation">
        <section class="confirm-dialog" role="alertdialog" aria-modal="true" aria-label="恢复默认队列确认">
          <small>恢复默认队列</small>
          <p>恢复默认队列后，手动调整的顺序将被清除。</p>
          <div>
            <button :ref="(element) => setConfirmButton(element, 0)" class="confirm-dialog__primary" @focus="confirmChoice = 0" @click="confirmRestore">确认恢复</button>
            <button :ref="(element) => setConfirmButton(element, 1)" @focus="confirmChoice = 1" @click="closeRestoreConfirmation">取消</button>
          </div>
        </section>
      </div>
    </Transition>
  </section>
</template>

<style scoped>
.home {
  --home-content-left: clamp(88px, 6cqw, 150px);
  --rail-card-width: clamp(170px, min(13.75cqw, 24.5cqh, calc(8.333cqh + 160px)), 290px);
  --rail-card-gap: clamp(10px, .75cqw, 16px);
  --rail-selected-scale: 1.08;
  --rail-selected-lift: clamp(16px, 2cqh, 22px);
  --rail-reorder-lift: clamp(27px, 3.2cqh, 35px);
  --rail-bottom-space: clamp(10px, 2cqh, 24px);
  --rail-focus-ratio: .34;
  --rail-motion-duration: 360ms;
  --rail-motion-easing: cubic-bezier(.22, .8, .25, 1);
  --rail-edge-fade: clamp(16px, 1.8cqw, 36px);
  --home-info-queue-gap: clamp(12px, 1.5cqh, 22px);
  position: relative;
  height: 100%;
  min-height: 0;
  overflow: hidden;
  background: radial-gradient(circle at 70% 20%, #182438, #070b12 64%);
}
.home__media, .home__backdrop, .home__shade { position: absolute; inset: 0; }
.home__media { transition: filter 180ms ease, transform 180ms ease; }
.home__backdrop { width: 100%; height: 100%; object-fit: cover; }
.home__shade { background: linear-gradient(90deg, rgba(3, 7, 13, .78) 0%, rgba(3, 7, 13, .36) 30%, rgba(3, 7, 13, .04) 70%), linear-gradient(0deg, rgba(3, 7, 13, .94) 0%, rgba(3, 7, 13, .48) 39%, transparent 66%), linear-gradient(180deg, rgba(3, 7, 13, .23), transparent 20%); transition: background 180ms ease; }
.home--modal-open .home__media { filter: blur(6px) brightness(.66); transform: scale(1.015); }
.home--modal-open .home__shade { background: rgba(3, 7, 13, .44); }

.home__lower-content { position: absolute; z-index: 4; right: 0; bottom: 0; left: 0; display: grid; gap: var(--home-info-queue-gap); padding-bottom: var(--rail-bottom-space); }
.home__info { position: relative; z-index: 3; width: min(560px, 52cqw); margin-left: var(--home-content-left); text-shadow: 0 3px 18px rgba(0, 0, 0, .72); }
.home__info h1 { margin: 0; max-width: 540px; font-size: clamp(2.4rem, 4.5cqw, 5.1rem); line-height: 1.02; letter-spacing: -.06em; }
.home__logo { display: block; max-width: min(450px, 40cqw); max-height: clamp(82px, 10.5cqh, 128px); object-fit: contain; object-position: left center; }
.home__company { margin: 10px 0 13px; color: rgba(255, 255, 255, .74); font-size: .77rem; letter-spacing: .04em; }
.home__details { display: flex; align-items: center; flex-wrap: wrap; gap: 7px; }
.home__tag, .home__status { padding: 5px 9px; border: 1px solid rgba(255, 255, 255, .15); border-radius: 999px; color: rgba(255, 255, 255, .78); background: rgba(4, 9, 17, .2); font-size: .66rem; backdrop-filter: blur(6px); }
.home__status { display: flex; align-items: center; gap: 6px; margin-left: 2px; color: #fff; }
.home__status i { width: 5px; height: 5px; border-radius: 50%; background: #74e2ff; box-shadow: 0 0 8px #74e2ff; }

.home__carousel { min-width: 0; }
.home__carousel-header { display: flex; align-items: end; justify-content: space-between; width: calc(100% - var(--home-content-left) - clamp(20px, 2.2cqw, 46px)); min-height: 24px; margin-left: var(--home-content-left); padding: 0 8px 4px; }
.home__carousel-header div { display: flex; align-items: center; gap: 9px; }
.home__carousel-header div > span { width: 2px; height: 15px; border-radius: 2px; background: rgba(191, 235, 249, .82); }
.home__carousel-header strong { color: rgba(255, 255, 255, .72); font-size: clamp(.72rem, .7cqw, .88rem); font-weight: 600; letter-spacing: .04em; }
.home__carousel-header p { margin: 0; color: rgba(255, 255, 255, .62); font-size: clamp(.61rem, .58cqw, .72rem); }
.home__cards-viewport { width: 100%; overflow: hidden; mask-image: linear-gradient(90deg, transparent 0, #000 var(--rail-edge-fade), #000 calc(100% - var(--rail-edge-fade)), transparent 100%); }
.home__cards-track { display: flex; gap: var(--rail-card-gap); width: max-content; padding: calc(var(--rail-reorder-lift) + 18px) var(--rail-trailing-space, 24px) 8px var(--rail-leading-space, 8px); will-change: transform; }
.home__cards-track--ready { transition: transform var(--rail-motion-duration) var(--rail-motion-easing); }
.queue-move { transition: transform var(--rail-motion-duration) var(--rail-motion-easing); }

.home__empty { position: absolute; z-index: 4; top: 50%; left: 50%; display: grid; justify-items: center; transform: translate(-50%, -50%); color: rgba(255, 255, 255, .72); text-align: center; }
.home__empty > span { display: grid; place-items: center; width: 52px; height: 52px; border: 1px solid rgba(255, 255, 255, .18); border-radius: 16px; font-size: 1.7rem; }
.home__empty h1 { margin: 15px 0 4px; color: #fff; font-size: 1.35rem; }
.home__empty p { margin: 0 0 18px; font-size: .76rem; }
.home__empty button { min-height: 40px; padding: 0 16px; border: 1px solid rgba(116, 223, 255, .32); border-radius: 11px; color: #dff8ff; background: rgba(93, 201, 241, .11); cursor: pointer; }

.home-dialog-layer { position: absolute; inset: 0; z-index: 30; display: grid; place-items: center; background: rgba(2, 5, 10, .24); }
.game-menu { width: 324px; padding: 20px; border: 1px solid rgba(255, 255, 255, .14); border-radius: 19px; background: rgba(8, 14, 24, .9); box-shadow: 0 26px 70px rgba(0, 0, 0, .48); backdrop-filter: blur(24px); }
.game-menu header { margin-bottom: 15px; text-align: center; }
.game-menu header small, .confirm-dialog > small { color: #77defb; font-size: .62rem; font-weight: 700; letter-spacing: .16em; }
.game-menu h2 { margin: 5px 0 0; overflow: hidden; color: #fff; font-size: 1.05rem; white-space: nowrap; text-overflow: ellipsis; }
.game-menu__actions { display: grid; gap: 7px; }
.game-menu__actions button { width: 100%; height: 43px; padding: 0 16px; border: 1px solid rgba(255, 255, 255, .1); border-radius: 11px; color: rgba(255, 255, 255, .83); background: rgba(255, 255, 255, .055); cursor: pointer; }
.game-menu__actions button:hover, .game-menu__actions button:focus-visible { border-color: rgba(129, 225, 255, .58); color: #fff; outline: none; background: rgba(105, 205, 241, .13); box-shadow: 0 0 0 2px rgba(100, 214, 255, .1); }
.game-menu__actions .game-menu__action--primary { height: 49px; border-color: rgba(139, 230, 255, .5); color: #07121b; background: linear-gradient(135deg, #e7faff, #7edcf7); font-weight: 800; }
.game-menu__actions .game-menu__action--primary:hover, .game-menu__actions .game-menu__action--primary:focus-visible { color: #041018; background: linear-gradient(135deg, #fff, #98e9ff); box-shadow: 0 0 0 2px rgba(156, 235, 255, .35), 0 8px 22px rgba(79, 199, 235, .2); }
.game-menu footer { margin-top: 13px; color: rgba(255, 255, 255, .38); font-size: .58rem; text-align: center; }

.confirm-dialog { width: min(390px, calc(100cqw - 40px)); padding: 24px; border: 1px solid rgba(255, 255, 255, .14); border-radius: 18px; background: rgba(8, 14, 24, .94); box-shadow: 0 26px 70px rgba(0, 0, 0, .5); backdrop-filter: blur(24px); }
.confirm-dialog p { margin: 10px 0 20px; color: rgba(255, 255, 255, .78); font-size: .82rem; line-height: 1.6; }
.confirm-dialog div { display: grid; grid-template-columns: 1fr 1fr; gap: 9px; }
.confirm-dialog button { height: 42px; border: 1px solid rgba(255, 255, 255, .12); border-radius: 10px; color: #fff; background: rgba(255, 255, 255, .06); cursor: pointer; }
.confirm-dialog button:hover, .confirm-dialog button:focus-visible { border-color: #7adefb; outline: none; box-shadow: 0 0 0 2px rgba(122, 222, 251, .14); }
.confirm-dialog .confirm-dialog__primary { color: #07121b; background: #92e6fc; font-weight: 800; }

.media-fade-enter-active, .media-fade-leave-active { transition: opacity 280ms ease; }
.media-fade-enter-from, .media-fade-leave-to { opacity: 0; }
.hero-enter-active, .hero-leave-active { transition: opacity 160ms ease, transform 160ms ease; }
.hero-enter-from { opacity: 0; transform: translateY(8px); }
.hero-leave-to { opacity: 0; transform: translateY(-6px); }
.dialog-fade-enter-active, .dialog-fade-leave-active { transition: opacity 150ms ease; }
.dialog-fade-enter-from, .dialog-fade-leave-to { opacity: 0; }

@container levabox-app (max-width: 1100px) {
  .home { --home-content-left: 72px; }
  .home__info { width: 58cqw; }
  .home__info h1 { font-size: clamp(2.2rem, 4.6cqw, 3.8rem); }
}

@container levabox-app (max-width: 1300px) and (min-height: 780px) and (min-aspect-ratio: 3 / 2) {
  .home { --rail-card-width: clamp(185px, 23.75cqh, 195px); }
}

@container levabox-app (min-width: 2300px) and (min-height: 1300px) {
  .home {
    --rail-card-gap: 18px;
    --rail-selected-lift: 26px;
    --rail-reorder-lift: 40px;
    --rail-bottom-space: 30px;
  }
  .home__info { width: 650px; }
  .home__logo { max-width: 540px; max-height: 150px; }
  .home__company { font-size: .9rem; }
  .home__tag, .home__status { font-size: .76rem; }
}

@container levabox-app (max-height: 760px) {
  .home {
    --rail-bottom-space: 8px;
    --rail-selected-lift: 16px;
    --rail-reorder-lift: 27px;
  }
  .home__info h1 { font-size: clamp(2.1rem, 4cqw, 3.7rem); }
  .home__logo { max-height: 78px; }
  .home__company { margin: 7px 0 9px; }
  .home__cards-track { padding-top: calc(var(--rail-reorder-lift) + 12px); }
}

@media (prefers-reduced-motion: reduce) {
  .home__cards-track,
  .queue-move { transition-duration: 1ms; }
}
</style>
