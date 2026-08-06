import { reactive, watch } from 'vue'

const STORAGE_KEY = 'levabox.settings.v1'

export interface SettingsState {
  backgroundMusic: boolean
  musicVolume: number
  soundEffects: boolean
  fullscreen: boolean
  gamepadEnabled: boolean
}

const defaults: SettingsState = {
  backgroundMusic: true,
  musicVolume: 65,
  soundEffects: true,
  fullscreen: false,
  gamepadEnabled: true,
}

export const settingsState = reactive<SettingsState>({
  ...defaults,
  ...readStoredSettings(),
})

watch(settingsState, persistSettings, { deep: true })

function readStoredSettings(): Partial<SettingsState> {
  if (typeof window === 'undefined') return {}
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return {}
    const stored = JSON.parse(raw) as Partial<SettingsState>
    return {
      backgroundMusic: typeof stored.backgroundMusic === 'boolean' ? stored.backgroundMusic : defaults.backgroundMusic,
      musicVolume: typeof stored.musicVolume === 'number'
        ? Math.min(100, Math.max(0, Math.round(stored.musicVolume)))
        : defaults.musicVolume,
      soundEffects: typeof stored.soundEffects === 'boolean' ? stored.soundEffects : defaults.soundEffects,
      fullscreen: typeof stored.fullscreen === 'boolean' ? stored.fullscreen : defaults.fullscreen,
      gamepadEnabled: typeof stored.gamepadEnabled === 'boolean' ? stored.gamepadEnabled : defaults.gamepadEnabled,
    }
  } catch {
    return {}
  }
}

function persistSettings(): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(settingsState))
  } catch {
    // localStorage may be unavailable in restricted WebViews; settings remain valid for this session.
  }
}
