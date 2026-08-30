import { describe, expect, test } from 'vitest'

import {
  handleAgentBoxRequestError,
  AgentBoxRequestError,
} from '../../src/utils/errors'

describe('handleAgentBoxRequestError', () => {
  test('does not throw when there is no error', () => {
    const res = { data: { id: '123' } }
    expect(() => handleAgentBoxRequestError(res)).not.toThrow()
  })

  test('throws AgentBoxRequestError for known status codes', () => {
    const res = { error: { code: 401, message: 'invalid token' } }
    expect(() => handleAgentBoxRequestError(res, 'Auth failed')).toThrow(
      AgentBoxRequestError
    )
    expect(() => handleAgentBoxRequestError(res, 'Auth failed')).toThrow(
      'Auth failed: [401] unauthorized: invalid token'
    )
  })

  test('throws AgentBoxRequestError with message for status code 0', () => {
    const res = { error: { code: 0, message: 'connection reset' } }
    expect(() => handleAgentBoxRequestError(res, 'Request failed')).toThrow(
      AgentBoxRequestError
    )
    expect(() => handleAgentBoxRequestError(res, 'Request failed')).toThrow(
      'Request failed: [0] unknown error: connection reset'
    )
  })

  test('throws AgentBoxRequestError when error code is missing', () => {
    const res = { error: { message: 'something went wrong' } } as any
    expect(() => handleAgentBoxRequestError(res, 'Request failed')).toThrow(
      AgentBoxRequestError
    )
    expect(() => handleAgentBoxRequestError(res, 'Request failed')).toThrow(
      'Request failed: [0] unknown error: something went wrong'
    )
  })

  test('handles valid but unlisted HTTP status codes via statuses package', () => {
    const res = { error: { code: 502, message: 'upstream down' } }
    expect(() => handleAgentBoxRequestError(res)).toThrow(AgentBoxRequestError)
    expect(() => handleAgentBoxRequestError(res)).toThrow(
      '[502] Bad Gateway: upstream down'
    )
  })
})
