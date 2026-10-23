import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import type { AdminGroup } from '@/types'
import GroupRateMultipliersModal from '../GroupRateMultipliersModal.vue'

const mocks = vi.hoisted(() => ({ getGroupRateMultipliers: vi.fn(), batchSetGroupRateMultipliers: vi.fn() }))
vi.mock('@/api/admin', () => ({ adminAPI: { groups: mocks } }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showSuccess: vi.fn(), showError: vi.fn() }) }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (key: string) => key }) }))
enableAutoUnmount(afterEach)
beforeEach(() => {
  vi.resetAllMocks()
  mocks.getGroupRateMultipliers.mockResolvedValue([
    { user_id: 7, user_email: 'user@example.com', user_name: '', user_status: 'active', rate_multiplier: 2, rpm_override: null }
  ])
})

async function openEditor() {
  const wrapper = mount(GroupRateMultipliersModal, {
    props: { show: false, group: { id: 1, name: 'Group', platform: 'openai', rate_multiplier: 1 } as AdminGroup },
    global: { stubs: {
      BaseDialog: { props: ['show'], template: '<div v-if="show"><slot /></div>' },
      Icon: true, PlatformIcon: true, Pagination: true
    } }
  })
  await wrapper.setProps({ show: true })
  await flushPromises()
  return wrapper
}

describe('group rate multiplier edits', () => {
  it.each(['0', '-1'])('does not stage invalid custom rate %s', async (value) => {
    const wrapper = await openEditor()
    await wrapper.get('tbody input[type="number"]').setValue(value)
    expect(wrapper.findAll('button').some(button => button.text() === 'common.save')).toBe(false)
    expect(mocks.batchSetGroupRateMultipliers).not.toHaveBeenCalled()
  })

  it.each([['0.5', 0.5], ['1e-3', 0.001]])('saves positive rate %s', async (value, expected) => {
    const wrapper = await openEditor()
    await wrapper.get('tbody input[type="number"]').setValue(value)
    await wrapper.findAll('button').find(button => button.text() === 'common.save')!.trigger('click')
    expect(mocks.batchSetGroupRateMultipliers).toHaveBeenCalledWith(1, [{ user_id: 7, rate_multiplier: expected }])
  })

  it('previews the hidden actual base without applying user-rate batch adjustments to it', async () => {
    const wrapper = await openEditor()
    await wrapper.setProps({ group: {
      id: 1, name: 'Group', platform: 'openai', rate_multiplier: 1,
      hidden_rate_multiplier: 0.25,
    } as AdminGroup })
    await flushPromises()
    expect(wrapper.text()).toContain('admin.groups.actualBaseRate')
    expect(wrapper.get('tbody tr').text()).toContain('0.25')
    await wrapper.get('input[placeholder="0.5"]').setValue('0.5')
    expect(wrapper.get('tbody tr').text()).toContain('0.25')
    expect(wrapper.get('tbody input').element.value).toBe('2')
  })

  it('still removes an override when its input is cleared', async () => {
    const wrapper = await openEditor()
    await wrapper.get('tbody input[type="number"]').setValue('')
    await wrapper.findAll('button').find(button => button.text() === 'common.save')!.trigger('click')
    expect(mocks.batchSetGroupRateMultipliers).toHaveBeenCalledWith(1, [])
  })
})
