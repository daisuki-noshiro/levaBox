import type { Game } from '../types/game'

export interface LaunchResult {
  launched: boolean
  message: string
}

/**
 * 游戏启动接口占位。
 * 后续将函数体替换为 Wails Go 后端调用，HomeView 无需改变交互结构。
 */
export async function launchGame(game: Game): Promise<LaunchResult> {
  return {
    launched: false,
    message: `“${game.title}”的启动功能将在后续接入`,
  }
}
