<template>
  <div class="card overflow-hidden">
    <!-- Toolbar -->
    <div
      class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-900/[0.06] bg-[#f8f9fb] px-4 py-3 dark:border-white/[0.07] dark:bg-dark-900/60 sm:px-5"
    >
      <div class="flex min-w-0 flex-wrap items-center gap-2.5">
        <span
          class="inline-flex shrink-0 items-center gap-1.5 rounded-full bg-emerald-500/10 px-2.5 py-1 text-xs font-semibold text-emerald-700 dark:bg-emerald-400/15 dark:text-emerald-300"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
          {{ activeItem?.label || t('payment.hub.storeBadge') }}
        </span>

        <!-- Multi-store sub-switcher (only shown if > 1 third-party recharge pages configured) -->
        <div
          v-if="items.length > 1"
          class="inline-flex items-center gap-1 rounded-full bg-slate-900/[0.05] p-0.5 dark:bg-white/[0.07]"
        >
          <button
            v-for="item in items"
            :key="item.id"
            type="button"
            class="rounded-full px-2.5 py-1 text-xs font-medium transition-all"
            :class="
              activeItem?.id === item.id
                ? 'bg-white text-[#111827] shadow-glass-sm dark:bg-dark-700 dark:text-white'
                : 'text-[#526071] hover:text-[#111827] dark:text-dark-300 dark:hover:text-white'
            "
            @click="selectedId = item.id"
          >
            {{ item.label }}
          </button>
        </div>

        <span class="text-xs text-[#526071] dark:text-dark-300">
          {{ t('payment.hub.storeToolbarHint') }}
        </span>
      </div>

      <div class="flex shrink-0 items-center gap-2">
        <button
          v-if="showRedeemAction"
          type="button"
          class="inline-flex items-center gap-1 rounded-full border border-primary-500/25 bg-primary-500/[0.09] px-3 py-1.5 text-xs font-semibold text-primary-600 transition-colors hover:bg-primary-500 hover:text-white dark:border-primary-400/30 dark:bg-primary-400/15 dark:text-primary-300 dark:hover:bg-primary-500 dark:hover:text-white"
          @click="emit('go-redeem')"
        >
          {{ t('payment.hub.goRedeemAction') }}
          <span aria-hidden="true">→</span>
        </button>

        <a
          v-if="activeItem && !activeItem.hide_open_button && embeddedUrl"
          :href="embeddedUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1.5 rounded-full border border-slate-900/[0.1] bg-white px-3 py-1.5 text-xs font-medium text-[#526071] shadow-glass-sm transition-colors hover:border-slate-900/20 hover:text-[#111827] dark:border-white/[0.1] dark:bg-dark-800 dark:text-dark-200 dark:hover:border-white/20 dark:hover:text-white"
        >
          <Icon name="externalLink" size="xs" :stroke-width="2" />
          {{ t('customPage.openInNewTab') }}
        </a>
      </div>
    </div>

    <!-- Embedded Iframe -->
    <div v-if="embeddedUrl" class="relative w-full bg-white dark:bg-dark-900">
      <iframe
        :key="embeddedUrl"
        :src="embeddedUrl"
        class="block h-[680px] min-h-[540px] w-full border-0"
        allowfullscreen
      ></iframe>
    </div>
    <div v-else class="px-6 py-16 text-center text-sm text-[#526071] dark:text-dark-300">
      {{ t('customPage.notConfiguredDesc') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { buildEmbeddedUrl, detectTheme } from '@/utils/embedded-url'
import type { CustomMenuItem } from '@/types'
import Icon from '@/components/icons/Icon.vue'

const props = withDefaults(
  defineProps<{
    items: CustomMenuItem[]
    showRedeemAction?: boolean
  }>(),
  {
    showRedeemAction: true,
  }
)

const emit = defineEmits<{
  'go-redeem': []
}>()

const { t, locale } = useI18n()
const authStore = useAuthStore()

const selectedId = ref<string>(props.items[0]?.id || '')
const pageTheme = ref<'light' | 'dark'>(detectTheme())
let themeObserver: MutationObserver | null = null

watch(
  () => props.items,
  (list) => {
    if (!list.some((item) => item.id === selectedId.value)) {
      selectedId.value = list[0]?.id || ''
    }
  },
  { immediate: true }
)

const activeItem = computed(() => {
  return props.items.find((item) => item.id === selectedId.value) ?? props.items[0] ?? null
})

const embeddedUrl = computed(() => {
  const item = activeItem.value
  if (!item?.url) return ''
  return buildEmbeddedUrl(
    item.url,
    authStore.user?.id,
    authStore.token,
    pageTheme.value,
    locale.value
  )
})

onMounted(() => {
  pageTheme.value = detectTheme()
  if (typeof MutationObserver !== 'undefined') {
    themeObserver = new MutationObserver(() => {
      pageTheme.value = detectTheme()
    })
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    })
  }
})

onUnmounted(() => {
  themeObserver?.disconnect()
})
</script>
