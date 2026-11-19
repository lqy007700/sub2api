import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { mount, RouterLinkStub } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { nextTick, reactive } from 'vue'
import { createI18n } from 'vue-i18n'
import { describe, expect, it, vi } from 'vitest'
import { useAdminSettingsStore, useAppStore, useAuthStore } from '@/stores'
import type { CustomMenuItem, PublicSettings, User } from '@/types'
import AppSidebar from '../AppSidebar.vue'

vi.mock('vue-router', () => ({
  useRoute: () => ({ path: '/dashboard' }),
  useRouter: () => ({ push: vi.fn() })
}))

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AppSidebar.vue')
const componentSource = readFileSync(componentPath, 'utf8')
const stylePath = resolve(dirname(fileURLToPath(import.meta.url)), '../../../style.css')
const styleSource = readFileSync(stylePath, 'utf8')

describe('AppSidebar custom SVG styles', () => {
  it('does not override uploaded SVG fill or stroke colors', () => {
    expect(componentSource).toContain('.sidebar-svg-icon {')
    expect(componentSource).toContain('color: currentColor;')
    expect(componentSource).toContain('display: block;')
    expect(componentSource).not.toContain('stroke: currentColor;')
    expect(componentSource).not.toContain('fill: none;')
  })
})

describe('AppSidebar custom menu visibility', () => {
  it.each(['user', 'admin'] as const)('hides and restores links in the %s sidebar', async (role) => {
    const pinia = createPinia()
    const appStore = useAppStore(pinia)
    const authStore = useAuthStore(pinia)
    const adminSettingsStore = useAdminSettingsStore(pinia)
    authStore.user = { role } as User
    adminSettingsStore.loaded = true

    const docs: CustomMenuItem = {
      id: 'docs', label: 'Docs', url: 'https://example.com/docs', icon_svg: '', visibility: role, sort_order: 0
    }
    const items = reactive([docs, { ...docs, id: 'help', label: 'Help', sort_order: 1, hide_in_sidebar: true }])
    appStore.cachedPublicSettings = { custom_menu_items: role === 'user' ? items : [] } as PublicSettings
    adminSettingsStore.customMenuItems = role === 'admin' ? items : []

    const wrapper = mount(AppSidebar, {
      global: {
        plugins: [pinia, createI18n({ legacy: false, missingWarn: false, fallbackWarn: false })],
        stubs: { RouterLink: RouterLinkStub, VersionBadge: true }
      }
    })
    const links = () => wrapper.findAllComponents(RouterLinkStub)
      .map(link => link.props('to'))
      .filter(to => String(to).startsWith('/custom/'))

    expect(links()).toEqual(['/custom/docs'])
    items[1].hide_in_sidebar = false
    await nextTick()
    expect(links()).toEqual(['/custom/docs', '/custom/help'])
    wrapper.unmount()
  })
})

describe('AppSidebar scroll position persistence', () => {
  it('binds a template ref to the sidebar nav element', () => {
    expect(componentSource).toContain('ref="sidebarNavRef"')
    expect(componentSource).toContain('sidebar-nav')
  })

  it('declares sidebarNavRef in script setup', () => {
    expect(componentSource).toContain("const sidebarNavRef = ref<HTMLElement | null>(null)")
  })

  it('saves scroll position on beforeUnmount', () => {
    expect(componentSource).toContain('onBeforeUnmount')
    expect(componentSource).toContain('appStore.sidebarScrollTop')
    expect(componentSource).toContain('sidebarNavRef.value.scrollTop')
  })

  it('restores scroll position on mount', () => {
    expect(componentSource).toContain('onMounted')
    expect(componentSource).toContain('appStore.sidebarScrollTop')
    expect(componentSource).toContain('nextTick')
  })
})

describe('AppSidebar collapsible groups', () => {
  it('lets the user collapse a group even while a child route is active', () => {
    // The expand state must come from the user's override first, falling back
    // to the active-route heuristic only when the user has not clicked yet.
    expect(componentSource).toContain('const groupExpandOverrides = ref<Map<string, boolean>>(new Map())')
    expect(componentSource).not.toContain('expandedGroups.value.has(item.path) || isGroupActive(item)')
  })
})

describe('AppSidebar header styles', () => {
  it('does not clip the version badge dropdown', () => {
    const sidebarHeaderBlockMatch = styleSource.match(/\.sidebar-header\s*\{[\s\S]*?\n {2}\}/)
    const sidebarBrandBlockMatch = componentSource.match(/\.sidebar-brand\s*\{[\s\S]*?\n\}/)

    expect(sidebarHeaderBlockMatch).not.toBeNull()
    expect(sidebarBrandBlockMatch).not.toBeNull()
    expect(sidebarHeaderBlockMatch?.[0]).not.toContain('@apply overflow-hidden;')
    expect(sidebarBrandBlockMatch?.[0]).not.toContain('overflow: hidden;')
  })
})

describe('AppSidebar subscription feature flag', () => {
  it('gates the My Subscriptions entry behind the subscription public-settings flag', () => {
    expect(componentSource).toContain('const flagSubscription = makeSidebarFlag(FeatureFlags.subscription)')
    expect(componentSource).toMatch(/path: '\/subscriptions'[^\n]*featureFlag: flagSubscription/)
  })

  it('also hides the admin Subscription Management entry on recharge-only sites', () => {
    expect(componentSource).toMatch(/path: '\/admin\/subscriptions'[^\n]*featureFlag: flagSubscription/)
  })

  it('derives the purchase entry label from the site billing mode', () => {
    expect(componentSource).toContain("import { resolveSiteBillingMode } from '@/utils/siteBillingMode'")
    expect(componentSource).toMatch(/case 'recharge_only':\s*return t\('nav\.recharge'\)/)
    expect(componentSource).toMatch(/case 'subscription_only':\s*return t\('nav\.subscribe'\)/)
    expect(componentSource).toMatch(/path: '\/purchase'[^\n]*label: purchaseNavLabel\.value/)
  })
})

describe('AppSidebar MCP entry', () => {
  it('uses the MCP logo and shows the entry only when the deployment enables MCP', () => {
    expect(componentSource).toContain("render: () => h(Icon, { name: 'mcp' })")
    expect(componentSource).toContain('const flagMCP = makeSidebarFlag(FeatureFlags.mcp)')
    expect(componentSource).toMatch(/path: '\/settings\/mcp'[^\n]*icon: MCPIcon, featureFlag: flagMCP/)
    expect(componentSource).toMatch(/if \(flagMCP\(\)\) filtered\.push\(\{ path: '\/settings\/mcp'[^\n]*icon: MCPIcon/)
    expect(componentSource).not.toMatch(/path: '\/settings\/mcp'[^\n]*PluginIcon/)
  })
})
