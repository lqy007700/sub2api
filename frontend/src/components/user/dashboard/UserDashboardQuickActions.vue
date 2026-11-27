<template>
  <div class="card overflow-hidden">
    <div class="border-b border-slate-900/[0.05] px-5 py-4 dark:border-white/[0.06]">
      <h2 class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">{{ t('dashboard.quickActions') }}</h2>
    </div>
    <div class="space-y-2 p-4">
      <button
        @click="router.push('/keys')"
        class="group flex w-full items-center gap-3.5 rounded-xl bg-[#f8f9fb] px-3.5 py-3 text-left transition-colors hover:bg-[#f1f3f7] dark:bg-dark-900/60 dark:hover:bg-white/[0.06]"
      >
        <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl border border-slate-900/[0.06] bg-white text-[#111827] shadow-glass-sm dark:border-white/[0.08] dark:bg-dark-800 dark:text-white">
          <Icon name="key" size="sm" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-[#111827] dark:text-white">{{ t('dashboard.createApiKey') }}</p>
          <p class="mt-0.5 truncate text-[11px] text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.generateNewKey') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="sm"
          class="text-[#8d99a8] transition-transform group-hover:translate-x-0.5 group-hover:text-[#111827] dark:text-dark-400 dark:group-hover:text-white"
        />
      </button>

      <button
        @click="router.push('/usage')"
        class="group flex w-full items-center gap-3.5 rounded-xl bg-[#f8f9fb] px-3.5 py-3 text-left transition-colors hover:bg-[#f1f3f7] dark:bg-dark-900/60 dark:hover:bg-white/[0.06]"
      >
        <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl border border-slate-900/[0.06] bg-white text-[#111827] shadow-glass-sm dark:border-white/[0.08] dark:bg-dark-800 dark:text-white">
          <Icon name="chart" size="sm" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-[#111827] dark:text-white">{{ t('dashboard.viewUsage') }}</p>
          <p class="mt-0.5 truncate text-[11px] text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.checkDetailedLogs') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="sm"
          class="text-[#8d99a8] transition-transform group-hover:translate-x-0.5 group-hover:text-[#111827] dark:text-dark-400 dark:group-hover:text-white"
        />
      </button>

      <button
        v-if="canUseBatchImage"
        @click="router.push('/batch-image')"
        class="group flex w-full items-center gap-3.5 rounded-xl bg-[#f8f9fb] px-3.5 py-3 text-left transition-colors hover:bg-[#f1f3f7] dark:bg-dark-900/60 dark:hover:bg-white/[0.06]"
      >
        <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl border border-slate-900/[0.06] bg-white text-[#111827] shadow-glass-sm dark:border-white/[0.08] dark:bg-dark-800 dark:text-white">
          <Icon name="sparkles" size="sm" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-[#111827] dark:text-white">{{ t('dashboard.batchImageAgent') }}</p>
          <p class="mt-0.5 truncate text-[11px] text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.batchImageAgentDesc') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="sm"
          class="text-[#8d99a8] transition-transform group-hover:translate-x-0.5 group-hover:text-[#111827] dark:text-dark-400 dark:group-hover:text-white"
        />
      </button>

      <button
        @click="handleRedeemNavigate"
        class="group flex w-full items-center gap-3.5 rounded-xl bg-[#f8f9fb] px-3.5 py-3 text-left transition-colors hover:bg-[#f1f3f7] dark:bg-dark-900/60 dark:hover:bg-white/[0.06]"
      >
        <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl border border-slate-900/[0.06] bg-white text-[#111827] shadow-glass-sm dark:border-white/[0.08] dark:bg-dark-800 dark:text-white">
          <Icon name="gift" size="sm" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-[#111827] dark:text-white">{{ t('dashboard.redeemCode') }}</p>
          <p class="mt-0.5 truncate text-[11px] text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.addBalanceWithCode') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="sm"
          class="text-[#8d99a8] transition-transform group-hover:translate-x-0.5 group-hover:text-[#111827] dark:text-dark-400 dark:group-hover:text-white"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useBatchImageAccess } from '@/composables/useBatchImageAccess'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
const router = useRouter()
const { t } = useI18n()
const { canUseBatchImage, refreshBatchImageAccess } = useBatchImageAccess()

function handleRedeemNavigate() {
  if (isFeatureFlagEnabled(FeatureFlags.payment)) {
    router.push({ path: '/purchase', query: { tab: 'redeem' } })
  } else {
    router.push('/redeem')
  }
}

onMounted(() => {
  void refreshBatchImageAccess()
})
</script>
