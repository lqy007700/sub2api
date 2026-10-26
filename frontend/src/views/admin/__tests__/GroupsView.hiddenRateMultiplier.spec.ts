import { defineComponent } from 'vue'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { AdminGroup, CreateGroupRequest, Group, UpdateGroupRequest } from '@/types'
import GroupsView from '@/views/admin/GroupsView.vue'

const mocks = vi.hoisted(() => ({
  list: vi.fn(), create: vi.fn(), update: vi.fn(), getAll: vi.fn(),
  getModelAllowlistCandidates: vi.fn(), getUsageSummary: vi.fn(),
  getCapacitySummary: vi.fn(), getLiveCapability: vi.fn(), showError: vi.fn(),
}))
const auth = vi.hoisted(() => ({ isSimpleMode: false }))
vi.mock('@/api/admin', () => ({ adminAPI: { groups: mocks, accounts: { list: vi.fn(), getById: vi.fn() } } }))
vi.mock('@/stores/auth', () => ({ useAuthStore: () => auth }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showError: mocks.showError, showSuccess: vi.fn() }) }))
vi.mock('@/stores/onboarding', () => ({ useOnboardingStore: () => ({ isCurrentStep: () => false }) }))
vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string, values?: Record<string, unknown>) => values?.price ? `${key}: ${values.price}` : key }),
}))
enableAutoUnmount(afterEach)

const group = {
  id: 42, name: 'Group', platform: 'grok', rate_multiplier: 5, hidden_rate_multiplier: 2,
  status: 'active', is_exclusive: false, subscription_type: 'standard', model_pricing: [],
  allow_image_generation: true, image_rate_independent: false, image_rate_multiplier: 4,
  image_price_1k: 0.1, image_price_2k: 0.2, image_price_4k: 0.3,
  video_rate_independent: false, video_rate_multiplier: 6,
  video_price_480p: 0.15, video_price_720p: 0.2, video_price_1080p: 0.3,
  web_search_price_per_call: 0.02, model_routing: null,
  peak_rate_enabled: true, peak_start: '00:00', peak_end: '23:59', peak_rate_multiplier: 7,
} as AdminGroup

const Dialog = defineComponent({ props: ['show'], template: '<div v-if="show"><slot /><slot name="footer" /></div>' })
function mountView() {
  return mount(GroupsView, { global: { stubs: {
    AppLayout: { template: '<main><slot /></main>' },
    TablePageLayout: { template: '<section><slot name="filters" /><slot name="table" /></section>' },
    DataTable: { props: ['data'], template: '<div><div v-for="row in data" :key="row.id"><slot name="cell-actions" :row="row" /><slot name="cell-rate_multiplier" :value="row.rate_multiplier" /></div></div>' },
    BaseDialog: Dialog, Pagination: true, ConfirmDialog: true, EmptyState: true,
    Select: true, PlatformIcon: true, Icon: true, GroupCapacityBadge: true,
    GroupRateMultipliersModal: true, GroupRPMOverridesModal: true,
    ReasoningEffortPolicyFields: defineComponent({
      setup(_, { expose }) { expose({ validate: () => true, resetValidation: () => undefined }) },
      template: '<div />',
    }), CodexManifestAccountsField: true,
    PricingEntryCard: true, VueDraggable: true,
  } } })
}
async function editGroup() {
  const wrapper = mountView()
  await flushPromises()
  await wrapper.findAll('button').find(button => button.text() === 'common.edit')!.trigger('click')
  await flushPromises()
  return wrapper
}
async function createGroup() {
  const wrapper = mountView()
  await flushPromises()
  await wrapper.get('[data-tour="groups-create-btn"]').trigger('click')
  await wrapper.get('[data-tour="group-form-name"]').setValue('New group')
  return wrapper
}

beforeEach(() => {
  vi.resetAllMocks()
  localStorage.clear()
  auth.isSimpleMode = false
  mocks.list.mockResolvedValue({ items: [group], total: 1, page: 1, page_size: 20, pages: 1 })
  mocks.getAll.mockResolvedValue([])
  mocks.getModelAllowlistCandidates.mockResolvedValue([])
  mocks.getUsageSummary.mockResolvedValue([])
  mocks.getCapacitySummary.mockResolvedValue([])
  mocks.getLiveCapability.mockResolvedValue({ supported: false })
})

describe('admin hidden group multiplier', () => {
  it.each(['openai', 'composite'])('saves hidden rates together with per-image billing for %s', async platform => {
    mocks.list.mockResolvedValue({ items: [{ ...group, platform, image_count_billing: true }], total: 1 })
    const wrapper = await editGroup()
    const label = wrapper.findAll('label').find(item => item.text() === 'admin.groups.imagePricing.perImageBilling')!
    expect(label.get<HTMLInputElement>('input').element.checked).toBe(true)
    await wrapper.get('#edit-hidden-rate-multiplier').setValue('0.125')
    await wrapper.get('#edit-group-form').trigger('submit')
    await flushPromises()
    expect(mocks.update).toHaveBeenCalledWith(42, expect.objectContaining({
      hidden_rate_multiplier: 0.125, image_count_billing: true, rate_multiplier: 5,
    }))
  })

  it('clears per-image billing when image generation is disabled while preserving the hidden rate', async () => {
    mocks.list.mockResolvedValue({ items: [{ ...group, platform: 'openai', allow_image_generation: false, image_count_billing: true }], total: 1 })
    const wrapper = await editGroup()
    expect(wrapper.findAll('label').some(item => item.text() === 'admin.groups.imagePricing.perImageBilling')).toBe(false)
    await wrapper.get('#edit-group-form').trigger('submit')
    await flushPromises()
    expect(mocks.update).toHaveBeenCalledWith(42, expect.objectContaining({
      hidden_rate_multiplier: 2, image_count_billing: false,
    }))
  })

  it('keeps the field out of public Group and accepts nullable admin requests', () => {
    const publicHasHidden: 'hidden_rate_multiplier' extends keyof Group ? true : false = false
    const create: Pick<CreateGroupRequest, 'hidden_rate_multiplier'> = { hidden_rate_multiplier: null }
    const update: Pick<UpdateGroupRequest, 'hidden_rate_multiplier'> = { hidden_rate_multiplier: null }
    expect(publicHasHidden).toBe(false)
    expect(create).toEqual(update)
  })

  it.each([null, 0.125])('creates with an explicit hidden multiplier %s', async value => {
    const wrapper = await createGroup()
    const input = wrapper.get<HTMLInputElement>('#create-hidden-rate-multiplier')
    expect(input.element.value).toBe('')
    expect(input.attributes('required')).toBeUndefined()
    if (value !== null) await input.setValue(String(value))
    await wrapper.get('#create-group-form').trigger('submit')
    await flushPromises()
    expect(mocks.create).toHaveBeenCalledWith(expect.objectContaining({ hidden_rate_multiplier: value }))
  })

  it('resets the hidden multiplier when the create dialog is reopened', async () => {
    const wrapper = await createGroup()
    await wrapper.get('#create-hidden-rate-multiplier').setValue('0.5')
    await wrapper.get('#create-group-form').trigger('submit')
    await flushPromises()
    await wrapper.get('[data-tour="groups-create-btn"]').trigger('click')
    expect(wrapper.get<HTMLInputElement>('#create-hidden-rate-multiplier').element.value).toBe('')
  })

  it.each(['2', '', '0.125'])('loads and updates hidden multiplier input %j', async value => {
    const wrapper = await editGroup()
    const input = wrapper.get<HTMLInputElement>('#edit-hidden-rate-multiplier')
    expect(input.element.value).toBe('2')
    await input.setValue(value)
    await wrapper.get('#edit-group-form').trigger('submit')
    await flushPromises()
    expect(mocks.update).toHaveBeenCalledWith(42, expect.objectContaining({
      hidden_rate_multiplier: value === '' ? null : Number(value), rate_multiplier: 5,
    }))
  })

  it.each(['0', '-1', '1000000'])('blocks invalid hidden multiplier %s', async value => {
    const wrapper = await editGroup()
    await wrapper.get('#edit-hidden-rate-multiplier').setValue(value)
    await wrapper.get('#edit-group-form').trigger('submit')
    expect(mocks.update).not.toHaveBeenCalled()
    expect(mocks.showError).toHaveBeenCalledWith('admin.groups.form.hiddenRateMultiplierInvalid')
  })

  it.each([
    [2, false, '$0.2', '$0.3'],
    [2, true, '$0.4', '$0.9'],
    [null, false, '$0.5', '$0.75'],
  ])('previews media with hidden=%s and independent=%s', async (hidden, independent, imagePrice, videoPrice) => {
    mocks.list.mockResolvedValue({ items: [{ ...group, hidden_rate_multiplier: hidden, image_rate_independent: independent, video_rate_independent: independent }], total: 1 })
    const wrapper = await editGroup()
    expect(wrapper.text()).toContain(`1K: ${imagePrice}`)
    expect(wrapper.text()).toContain(`480p: ${videoPrice}`)
    expect(wrapper.text()).toContain('5x')
  })

  it('uses hidden base for the web search preview and normal base after clearing', async () => {
    mocks.list.mockResolvedValue({ items: [{ ...group, platform: 'openai' }], total: 1 })
    const wrapper = await editGroup()
    expect(wrapper.text()).toContain('admin.groups.webSearchPricing.finalPricePreview: $0.04')
    await wrapper.get('#edit-hidden-rate-multiplier').setValue('')
    expect(wrapper.text()).toContain('admin.groups.webSearchPricing.finalPricePreview: $0.1')
  })

  it('leaves an existing hidden multiplier untouched during a simple-mode edit', async () => {
    auth.isSimpleMode = true
    const wrapper = await editGroup()
    expect(wrapper.find('#edit-hidden-rate-multiplier').exists()).toBe(false)
    await wrapper.get('#edit-group-form').trigger('submit')
    await flushPromises()
    expect(mocks.update).toHaveBeenCalledWith(42, expect.not.objectContaining({ hidden_rate_multiplier: expect.anything() }))
    expect(mocks.update.mock.calls[0][1]).not.toHaveProperty('hidden_rate_multiplier')
  })
})
