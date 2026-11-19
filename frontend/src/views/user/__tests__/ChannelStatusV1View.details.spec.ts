import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, shallowMount } from '@vue/test-utils'
import ChannelStatusV1View from '../ChannelStatusV1View.vue'

const mocks = vi.hoisted(() => ({ list: vi.fn(), status: vi.fn() }))
vi.mock('@/api/channelMonitor', () => mocks)
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ cachedPublicSettings: { channel_monitor_enabled: true }, showError: vi.fn() }) }))
vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key })
}))
enableAutoUnmount(afterEach)
afterEach(() => { vi.useRealTimers(); localStorage.clear() })

function detail(availability: number) {
  return { id: 1, models: [{ model: 'model', availability_15d: availability, availability_30d: availability }] }
}
beforeEach(() => {
  vi.useFakeTimers()
  vi.resetAllMocks()
  localStorage.clear()
  localStorage.setItem('channel-status-auto-refresh', JSON.stringify({ enabled: true, interval_seconds: 30 }))
  mocks.list.mockImplementation(async () => ({ items: [{ id: 1, name: 'Monitor', primary_model: 'model', primary_status: 'operational', availability_7d: 99 }] }))
  mocks.status.mockResolvedValue(detail(90))
})

function mountView() {
  return shallowMount(ChannelStatusV1View, {
    global: { stubs: {
      AppLayout: { template: '<div><slot /></div>' },
      MonitorHero: false,
      MonitorCardGrid: false,
      MonitorCard: { props: ['availabilityValue', 'item'], template: '<output :data-status="item.primary_status">{{ availabilityValue }}</output>' }
    } }
  })
}

it('does not show an overall badge beside the tabs when a channel is degraded', async () => {
  mocks.list.mockResolvedValue({ items: [{ id: 1, primary_status: 'degraded' }] })
  const wrapper = mountView()
  await flushPromises()
  expect(wrapper.findAll('[role="tab"]')).toHaveLength(3)
  expect(wrapper.text()).not.toContain('channelStatus.overall')
  expect(wrapper.get('output').attributes('data-status')).toBe('degraded')
})

describe('channel monitor automatic detail refresh', () => {
  it.each([15, 30])('refreshes the cached %s-day availability on timer ticks', async (days) => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.findAll('[role="tab"]')[days === 15 ? 1 : 2].trigger('click')
    await flushPromises()
    expect(wrapper.get('output').text()).toBe('90')
    mocks.status.mockResolvedValue(detail(98))
    await vi.advanceTimersByTimeAsync(30000)
    await flushPromises()
    expect(wrapper.get('output').text()).toBe('98')
    expect(mocks.status).toHaveBeenCalledTimes(2)
    expect(wrapper.get('button[title="common.refresh"]').attributes('disabled')).toBeUndefined()
  })

  it('keeps the seven-day view on the list endpoint without extra detail requests', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('output').trigger('click')
    await vi.advanceTimersByTimeAsync(30000)
    await flushPromises()
    expect(mocks.list).toHaveBeenCalledTimes(2)
    expect(mocks.status).not.toHaveBeenCalled()
    expect(wrapper.get('output').text()).toBe('99')
  })

  it('still refreshes cached availability on manual reload', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.findAll('[role="tab"]')[1].trigger('click')
    await flushPromises()
    mocks.status.mockResolvedValue(detail(98))
    await wrapper.get('button[title="common.refresh"]').trigger('click')
    await flushPromises()
    expect(wrapper.get('output').text()).toBe('98')
  })
})
