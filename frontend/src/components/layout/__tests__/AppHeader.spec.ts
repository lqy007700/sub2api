import { enableAutoUnmount, mount, RouterLinkStub } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { afterEach, expect, it, vi } from 'vitest'
import { createI18n } from 'vue-i18n'
import zhCommon from '@/i18n/locales/zh/common'
import { useAuthStore } from '@/stores'
import type { User } from '@/types'
import AppHeader from '../AppHeader.vue'

vi.mock('vue-router', () => ({
  useRoute: () => ({ name: 'Profile', meta: {}, params: {} }),
  useRouter: () => ({ push: vi.fn() })
}))

enableAutoUnmount(afterEach)

function mountHeader() {
  const pinia = createPinia()
  useAuthStore(pinia).user = { role: 'user', username: 'Tester' } as User
  const i18n = createI18n({
    legacy: false, locale: 'zh', messages: { zh: zhCommon }, missingWarn: false, fallbackWarn: false
  })
  const wrapper = mount(AppHeader, {
    attachTo: document.body,
    global: {
      plugins: [pinia, i18n],
      stubs: { RouterLink: RouterLinkStub, SubscriptionProgressMini: true, AnnouncementBell: true }
    }
  })
  const trigger = wrapper.get<HTMLButtonElement>('button[aria-expanded]')
  return { wrapper, trigger, area: wrapper.get({ ref: 'dropdownRef' }) }
}

it('opens on mouse hover and keeps the profile link in the dropdown', async () => {
  const { wrapper, trigger, area } = mountHeader()
  expect(trigger.attributes('aria-expanded')).toBe('false')
  await area.trigger('pointerenter', { pointerType: 'mouse' })
  expect(trigger.attributes('aria-expanded')).toBe('true')
  expect(wrapper.findAllComponents(RouterLinkStub).some(link =>
    link.props('to') === '/profile' && link.text() === '个人资料'
  )).toBe(true)
  await area.trigger('pointerleave', { pointerType: 'mouse' })
  expect(trigger.attributes('aria-expanded')).toBe('false')
})

it('keeps tap and keyboard access and returns focus on Escape', async () => {
  const { trigger, area } = mountHeader()
  await area.trigger('pointerenter', { pointerType: 'touch' })
  expect(trigger.attributes('aria-expanded')).toBe('false')
  await trigger.trigger('click')
  expect(trigger.attributes('aria-expanded')).toBe('true')
  await area.get('.dropdown button').trigger('keydown', { key: 'Escape' })
  expect(trigger.attributes('aria-expanded')).toBe('false')
  expect(document.activeElement).toBe(trigger.element)
})
