import request from './request'

export function getStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  })
}
