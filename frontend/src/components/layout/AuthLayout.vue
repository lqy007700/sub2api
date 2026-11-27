<template>
  <div class="apple-mist-ambient relative flex min-h-screen items-center justify-center overflow-hidden bg-[#f4f5f8] p-4 dark:bg-dark-950">
    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-[400px]">
      <!-- Logo/Brand -->
      <div class="mb-7 text-center">
        <template v-if="settingsLoaded">
          <div
            v-if="siteLogo"
            class="mb-4 inline-flex h-12 w-12 items-center justify-center overflow-hidden rounded-2xl border border-slate-900/[0.08] bg-white shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800"
          >
            <img :src="siteLogo" alt="Logo" class="h-full w-full object-contain p-1.5" />
          </div>
          <span
            v-else
            class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-[#111827] text-white shadow-card dark:bg-white dark:text-[#0b0c0f]"
          >
            <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none">
              <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
              <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
            </svg>
          </span>
          <h1 class="mb-1.5 text-2xl font-semibold tracking-tight text-[#111827] dark:text-white">
            <template v-if="isDefaultBrand">
              <span>Thisai</span><span class="text-primary-500 dark:text-primary-400">.site</span>
            </template>
            <template v-else>{{ siteName }}</template>
          </h1>
          <p class="text-xs text-[#526071] dark:text-dark-300">
            {{ siteSubtitle }}
          </p>
        </template>
      </div>

      <!-- Card Container -->
      <div class="card rounded-[20px] p-7 sm:p-8">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-7 text-center text-xs text-[#8d99a8] dark:text-dark-400">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()

const rawSiteName = computed(() => (appStore.siteName || '').trim())
const isDefaultBrand = computed(() => !rawSiteName.value || rawSiteName.value === 'Sub2API' || rawSiteName.value.toLowerCase() === 'thisai.site')
const siteName = computed(() => (isDefaultBrand.value ? 'Thisai.site' : rawSiteName.value))
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform')
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>

<style scoped>
.text-gradient {
  @apply bg-gradient-to-r from-primary-600 to-primary-500 bg-clip-text text-transparent;
}
</style>
