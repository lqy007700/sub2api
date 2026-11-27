<template>
  <AppLayout>
    <div class="mx-auto max-w-[820px] space-y-5">
      <!-- Compact Apple Mist Balance Strip (replaces legacy loud blue gradient hero) -->
      <div class="card flex flex-wrap items-center justify-between gap-4 px-5 py-4 sm:px-6">
        <div>
          <h1 class="text-lg font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ t('redeem.title') }}
          </h1>
          <p class="mt-0.5 text-xs text-[#526071] dark:text-dark-300">
            {{ t('redeem.concurrency') }}: {{ user?.concurrency || 0 }} {{ t('redeem.requests') }}
          </p>
        </div>

        <div class="flex items-center gap-4">
          <div class="text-right">
            <p class="flex items-center justify-end gap-1.5 text-[11px] font-medium text-[#8d99a8] dark:text-dark-400">
              <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
              {{ t('redeem.currentBalance') }}
            </p>
            <p class="font-mono text-2xl font-bold tracking-tight tabular-nums text-[#111827] dark:text-white">
              ${{ user?.balance?.toFixed(2) || '0.00' }}
            </p>
          </div>
        </div>
      </div>

      <RedeemCodePanel />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import RedeemCodePanel from '@/components/user/RedeemCodePanel.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const user = computed(() => authStore.user)
</script>
