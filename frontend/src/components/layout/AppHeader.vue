<template>
  <header class="glass sticky top-0 z-30 border-b border-slate-900/[0.06] dark:border-white/[0.07]">
    <div
      class="flex h-[58px] items-center justify-between gap-3 px-3 sm:px-5 md:px-6"
      :class="useTopNavMode ? 'mx-auto max-w-[1280px]' : ''"
    >
      <!-- Left: Mobile Menu Toggle + Brand (User Top-Nav) or Page Title (Admin) -->
      <div class="flex shrink-0 items-center gap-2 sm:gap-3">
        <button
          @click="toggleMobileSidebar"
          class="btn-ghost btn-icon lg:hidden"
          :aria-label="t('common.toggleMenu')"
        >
          <Icon name="menu" size="md" />
        </button>

        <!-- Regular User Desktop Brand -->
        <router-link
          v-if="useTopNavMode"
          to="/dashboard"
          class="hidden items-center gap-2.5 transition-opacity hover:opacity-80 lg:flex"
        >
          <div
            v-if="siteLogo"
            class="flex h-7 w-7 items-center justify-center overflow-hidden rounded-lg border border-slate-900/[0.08] bg-white p-0.5 shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800"
          >
            <img
              :src="siteLogo"
              alt="Logo"
              class="h-full w-full object-contain"
            />
          </div>
          <span
            v-else
            class="inline-flex h-7 w-7 items-center justify-center rounded-lg bg-[#111827] text-white shadow-sm dark:bg-white dark:text-[#0b0c0f]"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none">
              <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
              <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
            </svg>
          </span>
          <span v-if="isDefaultBrand" class="text-[15px] tracking-tight">
            <span class="font-bold text-[#111827] dark:text-white">Thisai</span><span class="font-medium text-primary-500 dark:text-primary-400">.site</span>
          </span>
          <span v-else class="text-[15px] font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ siteName }}
          </span>
        </router-link>

        <!-- Admin / Backend Mode Page Title -->
        <div v-else class="hidden lg:block">
          <h1 class="text-[15px] font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ pageTitle }}
          </h1>
        </div>
      </div>

      <!-- Center: Floating Pill Navigation (Regular User Desktop Only — Show All Directly) -->
      <nav
        v-if="useTopNavMode && userTopNavItems.length > 0"
        class="scrollbar-hide hidden min-w-0 items-center gap-0.5 overflow-x-auto rounded-full bg-slate-900/[0.045] p-1 dark:bg-white/[0.06] lg:flex"
      >
        <router-link
          v-for="item in userTopNavItems"
          :key="item.path"
          :to="item.path"
          :data-tour="item.path === '/keys' ? 'sidebar-my-keys' : undefined"
          @click="handleTopNavClick(item.path)"
          class="shrink-0 whitespace-nowrap rounded-full px-3 py-1.5 text-[13px] font-medium transition-all duration-150"
          :class="
            isNavActive(item.path)
              ? 'bg-white text-[#111827] shadow-pill dark:bg-dark-700 dark:text-white dark:shadow-none'
              : 'text-[#526071] hover:text-[#111827] dark:text-dark-300 dark:hover:text-white'
          "
        >
          {{ item.label }}
        </router-link>
      </nav>

      <!-- Right: Announcements + Theme + Docs + Model Plaza + Subscriptions + Balance + User Dropdown -->
      <div class="flex min-w-0 items-center gap-1 sm:gap-2">
        <!-- Theme Toggle (visible on desktop when user top-nav hides sidebar footer) -->
        <button
          v-if="useTopNavMode"
          type="button"
          @click="toggleTheme"
          class="hidden h-8 w-8 items-center justify-center rounded-full text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white lg:inline-flex"
          :title="isDark ? t('nav.lightMode') : t('nav.darkMode')"
        >
          <Icon v-if="isDark" name="sun" size="sm" />
          <Icon v-else name="moon" size="sm" />
        </button>

        <!-- Announcement Bell -->
        <AnnouncementBell v-if="user" />

        <!-- Docs Link -->
        <a
          v-if="docUrl"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="hidden items-center gap-1.5 rounded-full px-2.5 py-1.5 text-xs font-medium text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white sm:flex"
          :title="t('nav.docs')"
        >
          <Icon name="book" size="sm" />
          <span class="hidden xl:inline">{{ t('nav.docs') }}</span>
        </a>

        <!-- Model Plaza Entry (icon only below xl in top-nav mode to keep header airy) -->
        <router-link
          v-if="user && modelPlazaEnabled"
          :to="{ path: '/model-plaza', query: { embedded: '1' } }"
          :title="t('nav.modelPlaza')"
          :aria-label="t('nav.modelPlaza')"
          class="flex items-center gap-1.5 rounded-full px-2.5 py-1.5 text-xs font-medium text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
        >
          <Icon name="grid" size="sm" />
          <span class="hidden xl:inline">{{ t('nav.modelPlaza') }}</span>
        </router-link>

        <!-- Subscription Progress (for users with active subscriptions; not mounted at all when the feature is off) -->
        <SubscriptionProgressMini v-if="user && subscriptionFeatureEnabled" />

        <!-- Balance Display (Apple Mist Capsule) -->
        <div
          v-if="user"
          class="group relative hidden items-center gap-2 rounded-full border border-slate-900/[0.08] bg-white px-3 py-1 shadow-glass-sm dark:border-white/[0.09] dark:bg-dark-800 sm:flex"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-emerald-500 ring-2 ring-emerald-500/20"></span>
          <span class="font-mono text-xs font-semibold tabular-nums text-[#111827] dark:text-white">
            {{ formatHeaderMoney(availableBalance) }}
          </span>
          <span
            v-if="frozenBalance > 0"
            class="rounded-full bg-amber-100 px-1.5 py-0.5 text-[10px] font-medium text-amber-700 dark:bg-amber-900/40 dark:text-amber-200"
          >
            {{ balanceFrozenLabel }}
          </span>
          <div
            class="pointer-events-none absolute right-0 top-full mt-2 hidden w-56 rounded-2xl border border-slate-900/[0.08] bg-white/95 p-3.5 text-xs shadow-glass backdrop-blur-xl group-hover:block dark:border-white/[0.1] dark:bg-dark-800/95"
          >
            <div class="flex items-center justify-between">
              <span class="text-[#526071] dark:text-dark-300">{{ balanceAvailableText }}</span>
              <span class="font-mono font-medium tabular-nums text-[#111827] dark:text-white">{{ formatHeaderMoney(availableBalance) }}</span>
            </div>
            <div class="mt-2 flex items-center justify-between">
              <span class="text-[#526071] dark:text-dark-300">{{ balanceFrozenText }}</span>
              <span class="font-mono font-medium tabular-nums text-amber-700 dark:text-amber-200">{{ formatHeaderMoney(frozenBalance) }}</span>
            </div>
            <div class="mt-2.5 border-t border-slate-900/[0.06] pt-2.5 dark:border-white/[0.07]">
              <div class="flex items-center justify-between">
                <span class="text-[#526071] dark:text-dark-300">{{ balanceTotalText }}</span>
                <span class="font-mono font-semibold tabular-nums text-[#111827] dark:text-white">{{ formatHeaderMoney(totalBalance) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- User Dropdown -->
        <div
          v-if="user"
          ref="dropdownRef"
          class="relative"
          :class="{ 'user-menu-open': dropdownOpen }"
          @pointerenter="openDropdownOnHover"
          @pointerleave="closeDropdownOnHover"
          @keydown.esc="dismissDropdown"
        >
          <button
            ref="dropdownTriggerRef"
            type="button"
            @click="toggleDropdown"
            class="flex items-center gap-2 rounded-full p-1 transition-colors hover:bg-slate-900/[0.05] dark:hover:bg-white/[0.07]"
            :aria-label="t('common.userMenu')"
            :aria-expanded="dropdownOpen"
          >
            <div class="flex h-8 w-8 items-center justify-center overflow-hidden rounded-full border border-slate-900/[0.08] bg-white text-xs font-semibold text-[#111827] shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800 dark:text-white">
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                :alt="displayName"
                class="h-full w-full object-cover"
              >
              <span v-else>{{ userInitials }}</span>
            </div>
            <div v-if="!useTopNavMode" class="hidden text-left md:block">
              <div class="text-xs font-semibold text-[#111827] dark:text-white">
                {{ displayName }}
              </div>
              <div class="text-[11px] text-[#8d99a8] dark:text-dark-400">
                {{ t('admin.users.roles.' + user.role) }}
              </div>
            </div>
          </button>

          <!-- Dropdown Menu -->
          <transition name="dropdown">
            <div v-if="dropdownOpen" class="dropdown right-0 mt-2 w-56">
              <!-- User Info -->
              <div class="border-b border-gray-100 px-4 py-3 dark:border-dark-700">
                <div class="text-sm font-medium text-gray-900 dark:text-white">
                  {{ displayName }}
                </div>
                <div class="text-xs text-gray-500 dark:text-dark-400">{{ user.email }}</div>
              </div>

              <!-- Balance (mobile only) -->
              <div class="border-b border-gray-100 px-4 py-2 dark:border-dark-700 sm:hidden">
                <div class="text-xs text-gray-500 dark:text-dark-400">
                  {{ t('common.balance') }}
                </div>
                <div class="text-sm font-semibold text-primary-600 dark:text-primary-400">
                  {{ formatHeaderMoney(availableBalance) }}
                </div>
                <div v-if="frozenBalance > 0" class="mt-1 text-xs text-amber-600 dark:text-amber-300">
                  {{ balanceFrozenText }} {{ formatHeaderMoney(frozenBalance) }}
                </div>
              </div>

              <div class="py-1">
                <router-link to="/profile" @click="closeDropdown" class="dropdown-item">
                  <Icon name="user" size="sm" />
                  {{ t('nav.profile') }}
                </router-link>

                <router-link to="/keys" @click="closeDropdown" class="dropdown-item">
                  <Icon name="key" size="sm" />
                  {{ t('nav.apiKeys') }}
                </router-link>

                <a
                  v-if="authStore.isAdmin"
                  href="https://github.com/Wei-Shaw/sub2api"
                  target="_blank"
                  rel="noopener noreferrer"
                  @click="closeDropdown"
                  class="dropdown-item"
                >
                  <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                    <path
                      fill-rule="evenodd"
                      clip-rule="evenodd"
                      d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.464-1.11-1.464-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z"
                    />
                  </svg>
                  {{ t('nav.github') }}
                </a>

              </div>

              <!-- Contact Support (only show if configured) -->
              <div
                v-if="contactInfo"
                class="border-t border-gray-100 px-4 py-2.5 dark:border-dark-700"
              >
                <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                  <svg
                    class="h-3.5 w-3.5 flex-shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 01-.825-.242m9.345-8.334a2.126 2.126 0 00-.476-.095 48.64 48.64 0 00-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0011.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155"
                    />
                  </svg>
                  <span>{{ t('common.contactSupport') }}:</span>
                  <span class="font-medium text-gray-700 dark:text-gray-300">{{
                    contactInfo
                  }}</span>
                </div>
              </div>

              <div v-if="showOnboardingButton" class="border-t border-gray-100 py-1 dark:border-dark-700">
                <button @click="handleReplayGuide" class="dropdown-item w-full">
                  <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                    <path
                      d="M12 2a10 10 0 100 20 10 10 0 000-20zm0 14a1 1 0 110 2 1 1 0 010-2zm1.07-7.75c0-.6-.49-1.25-1.32-1.25-.7 0-1.22.4-1.43 1.02a1 1 0 11-1.9-.62A3.41 3.41 0 0111.8 5c2.02 0 3.25 1.4 3.25 2.9 0 2-1.83 2.55-2.43 3.12-.43.4-.47.75-.47 1.23a1 1 0 01-2 0c0-1 .16-1.82 1.1-2.7.69-.64 1.82-1.05 1.82-2.06z"
                    />
                  </svg>
                  {{ $t('onboarding.restartTour') }}
                </button>
              </div>

              <div class="border-t border-gray-100 py-1 dark:border-dark-700">
                <button
                  @click="handleLogout"
                  class="dropdown-item w-full text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                >
                  <svg
                    class="h-4 w-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75"
                    />
                  </svg>
                  {{ t('nav.logout') }}
                </button>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore, useOnboardingStore } from '@/stores'
import { useAdminSettingsStore } from '@/stores/adminSettings'
import SubscriptionProgressMini from '@/components/common/SubscriptionProgressMini.vue'
import AnnouncementBell from '@/components/common/AnnouncementBell.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled, makeSidebarFlag } from '@/utils/featureFlags'
import { resolveRouteMetaKeys } from '@/router/title'
import { resolveSiteBillingMode } from '@/utils/siteBillingMode'
import { useBatchImageAccess } from '@/composables/useBatchImageAccess'

interface TopNavItem {
  path: string
  label: string
  hideInSimpleMode?: boolean
  featureFlag?: () => boolean | undefined
}

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const adminSettingsStore = useAdminSettingsStore()
const onboardingStore = useOnboardingStore()
const { canUseBatchImage } = useBatchImageAccess()

const user = computed(() => authStore.user)
const isAdmin = computed(() => user.value?.role === 'admin')
const useTopNavMode = computed(() => Boolean(user.value) && !isAdmin.value && !appStore.backendModeEnabled)
const rawSiteName = computed(() => (appStore.siteName || '').trim())
const isDefaultBrand = computed(() => !rawSiteName.value || rawSiteName.value === 'Sub2API' || rawSiteName.value.toLowerCase() === 'thisai.site')
const siteName = computed(() => (isDefaultBrand.value ? 'Thisai.site' : rawSiteName.value))
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const dropdownTriggerRef = ref<HTMLButtonElement | null>(null)
const isDark = ref(document.documentElement.classList.contains('dark'))
const contactInfo = computed(() => appStore.contactInfo)
const docUrl = computed(() => sanitizeUrl(appStore.docUrl))
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))
const avatarUrl = computed(() => user.value?.avatar_url?.trim() || '')
const availableBalance = computed(() => Number(user.value?.balance || 0))
const frozenBalance = computed(() => Number(user.value?.frozen_balance || 0))
const totalBalance = computed(() => availableBalance.value + frozenBalance.value)
const balanceAvailableText = computed(() => t('common.availableBalance') === 'common.availableBalance' ? '可用余额' : t('common.availableBalance'))
const balanceFrozenText = computed(() => t('common.frozenBalance') === 'common.frozenBalance' ? '冻结金额' : t('common.frozenBalance'))
const balanceTotalText = computed(() => t('common.totalBalance') === 'common.totalBalance' ? '总余额' : t('common.totalBalance'))
const balanceFrozenLabel = computed(() => `${balanceFrozenText.value} ${formatHeaderMoney(frozenBalance.value)}`)

const flagChannelMonitor = makeSidebarFlag(FeatureFlags.channelMonitor)
const flagPayment = makeSidebarFlag(FeatureFlags.payment)
const flagAvailableChannels = makeSidebarFlag(FeatureFlags.availableChannels)
const flagSubscription = makeSidebarFlag(FeatureFlags.subscription)
const flagAffiliate = makeSidebarFlag(FeatureFlags.affiliate)
const flagMCP = makeSidebarFlag(FeatureFlags.mcp)
const flagBatchImageAccess = () => canUseBatchImage.value

const purchaseNavLabel = computed(() => {
  switch (resolveSiteBillingMode(appStore.cachedPublicSettings)) {
    case 'recharge_only':
      return t('nav.recharge')
    case 'subscription_only':
      return t('nav.subscribe')
    default:
      return t('nav.buySubscription')
  }
})

const customMenuItemsForUser = computed(() => {
  const items = appStore.cachedPublicSettings?.custom_menu_items ?? []
  return items
    .filter((item) => item.visibility === 'user' && !item.hide_in_sidebar)
    .sort((a, b) => a.sort_order - b.sort_order)
})

const userTopNavItems = computed<TopNavItem[]>(() => {
  const raw: TopNavItem[] = [
    { path: '/dashboard', label: t('nav.sidebar.overview') },
    { path: '/keys', label: t('nav.sidebar.keys') },
    { path: '/usage', label: t('nav.usage'), hideInSimpleMode: true },
    { path: '/purchase', label: purchaseNavLabel.value, hideInSimpleMode: true, featureFlag: flagPayment },
    { path: '/monitor', label: t('nav.sidebar.channels'), featureFlag: flagChannelMonitor },
    { path: '/available-channels', label: t('nav.availableChannels'), hideInSimpleMode: true, featureFlag: flagAvailableChannels },
    { path: '/subscriptions', label: t('nav.mySubscriptions'), hideInSimpleMode: true, featureFlag: flagSubscription },
    { path: '/orders', label: t('nav.myOrders') },
    { path: '/redeem', label: t('nav.redeem'), hideInSimpleMode: true },
    { path: '/affiliate', label: t('nav.affiliate'), hideInSimpleMode: true, featureFlag: flagAffiliate },
    { path: '/batch-image', label: t('nav.batchImage'), hideInSimpleMode: true, featureFlag: flagBatchImageAccess },
    { path: '/settings/mcp', label: t('mcp.title'), featureFlag: flagMCP },
    ...customMenuItemsForUser.value.map((item) => ({
      path: `/custom/${item.id}`,
      label: item.label
    }))
  ]

  return raw.filter((item) => {
    if (item.featureFlag && item.featureFlag() === false) return false
    if (authStore.isSimpleMode && item.hideInSimpleMode) return false
    return true
  })
})

function isNavActive(path: string): boolean {
  const current = route.path || ''
  return current === path || current.startsWith(path + '/')
}

function handleTopNavClick(itemPath: string) {
  if (itemPath === '/keys' && onboardingStore.isCurrentStep('[data-tour="sidebar-my-keys"]')) {
    onboardingStore.nextStep(500)
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// 只在标准模式的管理员下显示新手引导按钮
const showOnboardingButton = computed(() => {
  return !authStore.isSimpleMode && user.value?.role === 'admin'
})

const userInitials = computed(() => {
  if (!user.value) return ''
  // Prefer username, fallback to email
  if (user.value.username) {
    return user.value.username.substring(0, 2).toUpperCase()
  }
  if (user.value.email) {
    // Get the part before @ and take first 2 chars
    const localPart = user.value.email.split('@')[0]
    return localPart.substring(0, 2).toUpperCase()
  }
  return ''
})

const displayName = computed(() => {
  if (!user.value) return ''
  return user.value.username || user.value.email?.split('@')[0] || ''
})

// 订阅功能关闭时不挂载顶栏订阅徽章（组件 onMounted 会拉取订阅接口）。
const subscriptionFeatureEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.subscription))

// /purchase 的标题/描述随站点计费模式切换，与 document.title 共用同一解析。
const routeMetaKeys = computed(() => resolveRouteMetaKeys(route, {
  billingMode: resolveSiteBillingMode(appStore.cachedPublicSettings),
}))

const pageTitle = computed(() => {
  // For custom pages, use the menu item's label instead of generic "自定义页面"
  if (route.name === 'CustomPage') {
    const id = route.params.id as string
    const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []
    const menuItem = publicItems.find((item) => item.id === id)
      ?? (authStore.isAdmin ? adminSettingsStore.customMenuItems.find((item) => item.id === id) : undefined)
    if (menuItem?.label) return menuItem.label
  }
  const titleKey = routeMetaKeys.value.titleKey
  if (titleKey) {
    return t(titleKey)
  }
  return (route.meta.title as string) || ''
})

function toggleMobileSidebar() {
  appStore.toggleMobileSidebar()
}

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function openDropdownOnHover(event: PointerEvent) {
  if (event.pointerType === 'mouse') dropdownOpen.value = true
}

function closeDropdownOnHover(event: PointerEvent) {
  if (event.pointerType === 'mouse') closeDropdown()
}

function closeDropdown() {
  dropdownOpen.value = false
}

function dismissDropdown() {
  if (!dropdownOpen.value) return
  closeDropdown()
  dropdownTriggerRef.value?.focus()
}

async function handleLogout() {
  closeDropdown()
  try {
    await authStore.logout()
  } catch (error) {
    // Ignore logout errors - still redirect to login
    console.error('Logout error:', error)
  }
  await router.push('/login')
}

function handleReplayGuide() {
  closeDropdown()
  onboardingStore.replay()
}

function formatHeaderMoney(value: number) {
  if (!Number.isFinite(value)) return '$0.00'
  return `$${value.toFixed(2)}`
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    closeDropdown()
  }
}

let themeObserver: MutationObserver | null = null

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  if (typeof MutationObserver !== 'undefined') {
    themeObserver = new MutationObserver(() => {
      isDark.value = document.documentElement.classList.contains('dark')
    })
    themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
  themeObserver?.disconnect()
})
</script>

<style scoped>
/* Bridge the gap so the pointer can move from the trigger into the menu. */
.user-menu-open::after {
  @apply absolute right-0 top-full h-2 w-56;
  content: '';
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
