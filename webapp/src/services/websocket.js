import { BACKEND_HOST } from "@/utils"

export const shell = (path) => {
  return new WebSocket(`ws://${BACKEND_HOST}/ws/shell?path=${path}`);
}

export const ssh = (sessionId) => {
  return new WebSocket(`ws://${BACKEND_HOST}/ws/ssh?id=${sessionId}`);
}

export const sshQuickAccess = (sshConnection) => {
  return new WebSocket(`ws://${BACKEND_HOST}/ws/ssh?connection=${sshConnection.connection}&password=${sshConnection.password}`);
}
