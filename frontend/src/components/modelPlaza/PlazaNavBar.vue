<template>
  <header
    class="glass sticky top-0 z-30 border-b border-slate-900/[0.06] dark:border-white/[0.07]"
  >
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
      <!-- 左:站点 logo + 名称 -->
      <RouterLink to="/" class="flex min-w-0 items-center gap-2.5 transition-opacity hover:opacity-85">
        <template v-if="settings">
          <span
            v-if="siteLogo"
            class="flex h-8 w-8 flex-shrink-0 items-center justify-center overflow-hidden rounded-lg border border-slate-900/[0.08] bg-white p-0.5 shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800"
          >
            <img :src="siteLogo" alt="Logo" class="h-full w-full object-contain" />
          </span>
          <span
            v-else
            class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-[#111827] text-white shadow-sm dark:bg-white dark:text-[#0b0c0f]"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none">
              <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
              <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
            </svg>
          </span>
          <span v-if="isDefaultBrand" class="truncate text-[15px] tracking-tight">
            <span class="font-bold text-[#111827] dark:text-white">Thisai</span><span class="font-medium text-primary-500 dark:text-primary-400">.site</span>
          </span>
          <span v-else class="truncate text-[15px] font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ siteName }}
          </span>
        </template>
        <template v-else>
          <span class="h-8 w-8 flex-shrink-0 animate-pulse rounded-lg bg-slate-900/[0.08] dark:bg-white/[0.08]" aria-hidden="true"></span>
          <span class="h-5 w-24 animate-pulse rounded bg-slate-900/[0.08] dark:bg-white/[0.08]" aria-hidden="true"></span>
        </template>
      </RouterLink>

      <!-- 右:登录 / 回到后台 -->
      <RouterLink
        v-if="isAuthenticated"
        :to="backTarget"
        class="btn btn-primary flex-shrink-0 rounded-full px-4 py-1.5 text-xs"
      >
        {{ t('modelPlaza.nav.backToDashboard') }}
      </RouterLink>
      <RouterLink
        v-else
        :to="{ path: '/login', query: { redirect: '/model-plaza' } }"
        class="btn btn-primary flex-shrink-0 rounded-full px-4 py-1.5 text-xs"
      >
        {{ t('modelPlaza.nav.login') }}
      </RouterLink>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const settings = computed(() => appStore.cachedPublicSettings)
const rawSiteName = computed(() => (settings.value?.site_name || '').trim())
const isDefaultBrand = computed(
  () => !rawSiteName.value || rawSiteName.value === 'Sub2API' || rawSiteName.value.toLowerCase() === 'thisai.site'
)
const siteName = computed(() => (isDefaultBrand.value ? 'Thisai.site' : rawSiteName.value))
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || '', { allowRelative: true, allowDataUrl: true })
)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const backTarget = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
</script>
