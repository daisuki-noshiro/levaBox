import {
  SelectExecutable,
  StartImport,
} from '../../wailsjs/go/main/App'
import type { service } from '../../wailsjs/go/models'

export async function selectExecutable(): Promise<string> {
  return SelectExecutable()
}

export async function startImport(
  executablePath: string,
): Promise<service.StartImportResult> {
  return StartImport(executablePath)
}
