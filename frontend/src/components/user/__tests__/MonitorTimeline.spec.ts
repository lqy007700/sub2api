import { mount } from '@vue/test-utils'
import { afterEach, expect, it, vi } from 'vitest'
import MonitorTimeline from '../monitor/MonitorTimeline.vue'

vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key, te: () => true }),
}))

afterEach(() => vi.useRealTimers())

it('fills every day in green and preserves the recorded status in its label', async () => {
  vi.useFakeTimers({ toFake: ['Date'] })
  vi.setSystemTime(new Date('2026-11-19T10:00:00Z'))
  const wrapper = mount(MonitorTimeline, {
    props: {
      days: 30,
      countdownSeconds: 10,
      buckets: [
        { date: '2026-11-18', status: 'failed' },
        { date: '2026-11-19', status: 'degraded' },
      ],
    },
  })

  const bars = () => wrapper.findAll('[role="listitem"]')
  expect(bars()).toHaveLength(30)
  expect(wrapper.findAll('.h-full.bg-emerald-500')).toHaveLength(30)
  expect(wrapper.text()).toContain('channelStatus.historyPlaceholder')
  expect(bars()[0].attributes('aria-label')).toContain('channelStatus.noHistorySample')
  expect(bars()[28].attributes('title')).toContain('monitorCommon.status.failed')
  expect(bars()[29].attributes('title')).toContain('monitorCommon.status.degraded')

  const recentLabels = bars().slice(-15).map(bar => bar.attributes('title'))
  await wrapper.setProps({ days: 15 })
  expect(wrapper.findAll('.h-full.bg-emerald-500')).toHaveLength(15)
  expect(bars().map(bar => bar.attributes('title'))).toEqual(recentLabels)
  await wrapper.setProps({ days: 7 })
  expect(wrapper.findAll('.h-full.bg-emerald-500')).toHaveLength(7)
  expect(bars().map(bar => bar.attributes('title'))).toEqual(recentLabels.slice(-7))
  wrapper.unmount()
})
