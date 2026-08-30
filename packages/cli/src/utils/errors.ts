import status from 'statuses'

/**
 * Thrown when a request to AgentBox API occurs.
 */
export class AgentBoxRequestError extends Error {
  constructor(message: any) {
    super(message)
    this.name = 'AgentBoxRequestError'
  }
}

type AgentBoxResponseError = { code?: number; message?: string }

type AgentBoxResponse<TData> =
  | {
      data: TData
      error?: undefined
    }
  | {
      data?: undefined
      error: AgentBoxResponseError
    }

function throwAgentBoxRequestError(
  error: AgentBoxResponseError,
  errMsg?: string
): never {
  let message: string
  const code = error.code ?? 0
  switch (code) {
    case 400:
      message = 'bad request'
      break
    case 401:
      message = 'unauthorized'
      break
    case 403:
      message = 'forbidden'
      break
    case 404:
      message = 'not found'
      break
    case 500:
      message = 'internal server error'
      break
    default:
      message = status.message[code] || 'unknown error'
      break
  }

  throw new AgentBoxRequestError(
    `${errMsg && `${errMsg}: `}[${code}] ${message && `${message}: `}${
      error.message ?? 'no message'
    }`
  )
}

export function handleAgentBoxRequestError(
  res: { error: AgentBoxResponseError },
  errMsg?: string
): never
export function handleAgentBoxRequestError<TData>(
  res: AgentBoxResponse<TData>,
  errMsg?: string
): asserts res is { data: TData; error?: undefined }
export function handleAgentBoxRequestError(
  res: AgentBoxResponse<unknown>,
  errMsg?: string
) {
  if (!res.error) {
    return
  }
  throwAgentBoxRequestError(res.error, errMsg)
}
