import {
  PrepareImportMetadata,
  SelectExecutable,
  StartImport,
} from '../../wailsjs/go/main/App'
import type { metadata, service } from '../../wailsjs/go/models'

export async function selectExecutable(): Promise<string> {
  return SelectExecutable()
}

export async function startImport(
  executablePath: string,
): Promise<service.StartImportResult> {
  return StartImport(executablePath)
}

export async function prepareImportMetadata(
  draft: service.ImportDraft,
  sources: metadata.Source[],
): Promise<service.ImportMetadataResult> {
  return PrepareImportMetadata(draft, sources)
}
