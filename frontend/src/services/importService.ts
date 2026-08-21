import {
  PrepareImportMetadata,
  SaveImport,
  SelectExecutable,
  StartImport,
} from '../../wailsjs/go/main/App'
import type { metadata, model, service } from '../../wailsjs/go/models'

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

export async function saveImport(
  request: service.SaveImportRequest,
): Promise<model.Game> {
  return SaveImport(request)
}
