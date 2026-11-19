import { mount } from '@vue/test-utils'
import { afterEach, expect, it, vi } from 'vitest'
import MonitorTimeline from '../monitor/MonitorTimeline.vue'

vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key, te: () => true }),
}))

afterEach(() => vi.useRealTimers())

it('shows degraded checks as available and keeps failures and missing days distinct', async () => {
  vi.useFakeTimers({ toFake: ['Date'] })
  vi.setSystemTime(new Date('2026-11-19T10:00:00Z'))
  const wrapper = mount(MonitorTimeline, {
    props: {
      days: 30,
      monitorId: 1,
      countdownSeconds: 10,
      buckets: [
        { date: '2026-11-18', status: 'failed' },
        { date: '2026-11-19', status: 'degraded' },
      ],
    },
  })

  const bars = () => wrapper.findAll('[role="listitem"]')
  expect(bars()).toHaveLength(30)
  expect(wrapper.findAll('.bg-amber-100')).toHaveLength(2)
  expect(wrapper.findAll('.bg-emerald-100')).toHaveLength(26)
  expect(wrapper.text()).toContain('channelStatus.historyPlaceholder')
  expect(bars()[0].attributes('aria-label')).toContain('channelStatus.noHistorySample')
  expect(bars()[28].classes()).toContain('bg-red-500')
  expect(bars()[29].classes()).toContain('bg-emerald-500')
  expect(bars()[29].attributes('aria-label')).toContain('monitorCommon.status.operational')

  const recentClasses = bars().slice(-15).map(bar => bar.classes())
  await wrapper.setProps({ days: 15 })
  expect(bars().map(bar => bar.classes())).toEqual(recentClasses)
  await wrapper.setProps({ days: 7 })
  expect(bars().map(bar => bar.classes())).toEqual(recentClasses.slice(-7))
  wrapper.unmount()
})
