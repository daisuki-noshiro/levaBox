export interface LobbyInputHandlers {
  moveLeft: () => void
  moveRight: () => void
  moveUp: () => void
  moveDown: () => void
  confirm: () => void
  cancel: () => void
}

export interface GlobalGamepadInputHandlers {
  isMenuOpen: () => boolean
  toggleMenu: () => void
  moveUp: () => void
  moveDown: () => void
  confirm: () => void
  cancel: () => void
}

interface InputSnapshot {
  left: boolean
  right: boolean
  up: boolean
  down: boolean
  confirm: boolean
  cancel: boolean
  menu: boolean
}

let lobbyHandlers: LobbyInputHandlers | null = null
let globalHandlers: GlobalGamepadInputHandlers | null = null
let pollingFrame = 0
const previousInput = new Map<number, InputSnapshot>()

const emptySnapshot = (): InputSnapshot => ({
  left: false,
  right: false,
  up: false,
  down: false,
  confirm: false,
  cancel: false,
  menu: false,
})

function buttonPressed(gamepad: Gamepad, index: number): boolean {
  return gamepad.buttons[index]?.pressed === true
}

function readInput(gamepad: Gamepad): InputSnapshot {
  const horizontalAxis = gamepad.axes[0] ?? 0
  const verticalAxis = gamepad.axes[1] ?? 0
  return {
    left: buttonPressed(gamepad, 14) || horizontalAxis < -.6,
    right: buttonPressed(gamepad, 15) || horizontalAxis > .6,
    up: buttonPressed(gamepad, 12) || verticalAxis < -.6,
    down: buttonPressed(gamepad, 13) || verticalAxis > .6,
    confirm: buttonPressed(gamepad, 0),
    cancel: buttonPressed(gamepad, 1),
    menu: buttonPressed(gamepad, 9),
  }
}

function justPressed(current: InputSnapshot, previous: InputSnapshot, key: keyof InputSnapshot): boolean {
  return current[key] && !previous[key]
}

function dispatchInput(current: InputSnapshot, previous: InputSnapshot): void {
  if (justPressed(current, previous, 'menu')) {
    globalHandlers?.toggleMenu()
    return
  }

  if (globalHandlers?.isMenuOpen()) {
    if (justPressed(current, previous, 'up')) globalHandlers.moveUp()
    else if (justPressed(current, previous, 'down')) globalHandlers.moveDown()
    else if (justPressed(current, previous, 'confirm')) globalHandlers.confirm()
    else if (justPressed(current, previous, 'cancel')) globalHandlers.cancel()
    return
  }

  if (!lobbyHandlers) return
  if (justPressed(current, previous, 'left')) lobbyHandlers.moveLeft()
  else if (justPressed(current, previous, 'right')) lobbyHandlers.moveRight()
  else if (justPressed(current, previous, 'up')) lobbyHandlers.moveUp()
  else if (justPressed(current, previous, 'down')) lobbyHandlers.moveDown()
  else if (justPressed(current, previous, 'confirm')) lobbyHandlers.confirm()
  else if (justPressed(current, previous, 'cancel')) lobbyHandlers.cancel()
}

function pollGamepads(): void {
  if (!lobbyHandlers && !globalHandlers) {
    pollingFrame = 0
    previousInput.clear()
    return
  }

  const gamepads = typeof navigator.getGamepads === 'function' ? navigator.getGamepads() : []
  for (const gamepad of gamepads) {
    if (!gamepad) continue
    const current = readInput(gamepad)
    const previous = previousInput.get(gamepad.index) ?? emptySnapshot()
    dispatchInput(current, previous)
    previousInput.set(gamepad.index, current)
  }
  pollingFrame = requestAnimationFrame(pollGamepads)
}

function ensurePolling(): void {
  if (!pollingFrame && typeof requestAnimationFrame === 'function') {
    pollingFrame = requestAnimationFrame(pollGamepads)
  }
}

function stopPollingIfUnused(): void {
  if (lobbyHandlers || globalHandlers || !pollingFrame) return
  cancelAnimationFrame(pollingFrame)
  pollingFrame = 0
  previousInput.clear()
}

export function registerLobbyGamepadInput(handlers: LobbyInputHandlers): () => void {
  lobbyHandlers = handlers
  ensurePolling()
  return () => {
    if (lobbyHandlers === handlers) lobbyHandlers = null
    stopPollingIfUnused()
  }
}

export function registerGlobalGamepadInput(handlers: GlobalGamepadInputHandlers): () => void {
  globalHandlers = handlers
  ensurePolling()
  return () => {
    if (globalHandlers === handlers) globalHandlers = null
    stopPollingIfUnused()
  }
}
