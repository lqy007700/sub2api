import { mount } from '@vue/test-utils'
import { expect, it, vi } from 'vitest'
import MonitorTimeline from '../monitor/MonitorTimeline.vue'

vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key, te: () => true }),
}))

it('fills the selected range with green bars without status tooltips', async () => {
  const wrapper = mount(MonitorTimeline, {
    props: {
      days: 30,
      countdownSeconds: 10,
    },
  })

  const bars = () => wrapper.findAll('[role="img"] > .h-full.bg-emerald-500')
  expect(bars()).toHaveLength(30)
  expect(wrapper.find('[role="img"] [title]').exists()).toBe(false)
  await wrapper.setProps({ days: 15 })
  expect(bars()).toHaveLength(15)
  await wrapper.setProps({ days: 7 })
  expect(bars()).toHaveLength(7)
  wrapper.unmount()
})
