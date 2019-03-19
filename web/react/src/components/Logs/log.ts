export interface LogEntry {
  timestamp: string
  message: string
}

export interface LogResponse {
  entries: LogEntry[]
}
