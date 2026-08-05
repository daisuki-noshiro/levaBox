export type GameStatus = '未开始' | '游玩中' | '已通关'
export type BackgroundType = 'image' | 'video'

export interface Game {
  id: string
  title: string
  subtitle: string
  developer: string
  year: number
  genres: string[]
  shortDescription: string
  description: string
  status: GameStatus
  playtime: number
  favorite: boolean
  recentlyPlayed: boolean
  backgroundType: BackgroundType
  backgroundImage: string
  backgroundVideo?: string
  backgroundPosition?: string
  cover: string
  coverPosition?: string
  logo?: string
  backgroundGradient: string
  coverGradient: string
  accent: string
}
