export interface LobbyInputHandlers {
  moveLeft: () => void
  moveRight: () => void
  confirm: () => void
  cancel: () => void
  openMenu: () => void
}

/**
 * Future boundary for Gamepad API or Wails/native controller events.
 * Keyboard and mouse input stay inside HomeView for the current prototype.
 */
export function registerLobbyGamepadInput(_handlers: LobbyInputHandlers): () => void {
  return () => undefined
}
