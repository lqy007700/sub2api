import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import UserOrdersView from '../UserOrdersView.vue'
import Select from '@/components/common/Select.vue'
import Pagination from '@/components/common/Pagination.vue'
import OrderTable from '@/components/payment/OrderTable.vue'

const api = vi.hoisted(() => ({ getMyOrders: vi.fn(), getRefundEligibleProviders: vi.fn(), getHistory: vi.fn() }))
const store = vi.hoisted(() => ({ showError: vi.fn(), cachedPublicSettings: { payment_enabled: true } }))
vi.mock('@/api/payment', () => ({ paymentAPI: api }))
vi.mock('@/api', () => ({ redeemAPI: { getHistory: api.getHistory } }))
vi.mock('@/stores', () => ({ useAppStore: () => store }))
vi.mock('@/stores/app', () => ({ useAppStore: () => store }))
vi.mock('vue-router', () => ({ useRouter: () => ({ push: vi.fn() }) }))
vi.mock('vue-i18n', async (importOriginal) => ({
  ...await importOriginal<typeof import('vue-i18n')>(),
  useI18n: () => ({ t: (key: string) => key })
}))
enableAutoUnmount(afterEach)

beforeEach(() => {
  vi.clearAllMocks()
  api.getMyOrders.mockResolvedValue({ data: { items: [], total: 100 } })
  api.getRefundEligibleProviders.mockResolvedValue({ data: { provider_instance_ids: [] } })
  api.getHistory.mockResolvedValue({ items: [], total: 0 })
  store.cachedPublicSettings.payment_enabled = true
})

async function mountOrders() {
  const wrapper = mount(UserOrdersView, {
    global: { stubs: {
      AppLayout: { template: '<div><slot /></div>' },
      OrderTable: true, BaseDialog: true, Icon: true, Pagination: true, teleport: true
    } }
  })
  await flushPromises()
  return wrapper
}

async function openOrders() {
  const wrapper = await mountOrders()
  wrapper.getComponent(Pagination).vm.$emit('update:page', 4)
  await flushPromises()
  return wrapper
}

describe('order status filtering', () => {
  it('loads the first page with the selected status', async () => {
    const wrapper = await openOrders()
    const select = wrapper.getComponent(Select)
    await select.get('button').trigger('click')
    await select.findAll('[role="option"]').find(option => option.text() === 'payment.status.pending')!.trigger('click')
    await flushPromises()
    expect(api.getMyOrders).toHaveBeenLastCalledWith({ page: 1, page_size: 20, status: 'PENDING' })
    expect(wrapper.getComponent(Pagination).props('page')).toBe(1)
    expect(api.getMyOrders).toHaveBeenCalledTimes(3)
  })

  it('keeps the current page on manual refresh', async () => {
    const wrapper = await openOrders()
    await wrapper.get('[title="common.refresh"]').trigger('click')
    await flushPromises()
    expect(api.getMyOrders).toHaveBeenLastCalledWith({ page: 4, page_size: 20, status: undefined })
  })
})

describe('order categories', () => {
  it('loads recharge and purchase records with server-side filters starting on page one', async () => {
    const wrapper = await openOrders()
    const select = (kind: string) => wrapper.findAll('.tab').find(tab => tab.text() === 'payment.orders.tabs.' + kind)!
    await select('balance').trigger('click')
    await flushPromises()
    expect(api.getMyOrders).toHaveBeenLastCalledWith({ page: 1, page_size: 20, order_type: 'balance', status: undefined })
    await select('subscription').trigger('click')
    await flushPromises()
    expect(api.getMyOrders).toHaveBeenLastCalledWith({ page: 1, page_size: 20, order_type: 'subscription', status: undefined })
    expect(api.getHistory).not.toHaveBeenCalled()
  })

  it('paginates redemptions with online payments disabled', async () => {
    store.cachedPublicSettings.payment_enabled = false
    api.getHistory.mockResolvedValue({
      items: [{ id: 4, code: 'REDEEM-CODE', type: 'balance', value: 10, used_at: '2026-01-01T00:00:00Z' }],
      total: 31
    })
    const wrapper = await mountOrders()
    expect(wrapper.text()).not.toContain('payment.result.backToRecharge')
    await wrapper.findAll('.tab').find(tab => tab.text() === 'payment.orders.tabs.redemption')!.trigger('click')
    await flushPromises()
    expect(api.getHistory).toHaveBeenLastCalledWith(1, 20, { codesOnly: true })
    expect(wrapper.text()).toContain('REDEEM-C...')
    expect(wrapper.findComponent(OrderTable).exists()).toBe(false)
    await wrapper.findAll('button').find(button => button.text() === 'pagination.next')!.trigger('click')
    await flushPromises()
    expect(api.getHistory).toHaveBeenLastCalledWith(2, 20, { codesOnly: true })
    expect(api.getMyOrders).toHaveBeenCalledTimes(1)
  })

  it('ignores a payment response from the previous category', async () => {
    const wrapper = await mountOrders()
    let finishRecharge!: (value: unknown) => void
    api.getMyOrders.mockImplementationOnce(() => new Promise(resolve => { finishRecharge = resolve }))
    api.getMyOrders.mockResolvedValueOnce({ data: { items: [{ id: 2 }], total: 1 } })
    const select = (kind: string) => wrapper.findAll('.tab').find(tab => tab.text() === 'payment.orders.tabs.' + kind)!
    await select('balance').trigger('click')
    await select('subscription').trigger('click')
    await flushPromises()
    finishRecharge({ data: { items: [{ id: 1 }], total: 50 } })
    await flushPromises()
    expect(wrapper.getComponent(OrderTable).props('orders')).toEqual([{ id: 2 }])
    expect(wrapper.getComponent(Pagination).props('total')).toBe(1)
  })
})
