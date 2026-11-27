<template>
  <div class="card overflow-hidden">
    <div class="flex items-center justify-between border-b border-slate-900/[0.05] px-5 py-4 dark:border-white/[0.06]">
      <h2 class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">{{ t('dashboard.recentUsage') }}</h2>
      <span class="rounded-full bg-[#f4f5f8] px-2.5 py-0.5 text-[11px] font-medium text-[#526071] dark:bg-white/[0.06] dark:text-dark-300">
        {{ t('dashboard.last7Days') }}
      </span>
    </div>
    <div class="p-4">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
      <div v-else-if="data.length === 0" class="py-8">
        <EmptyState :title="t('dashboard.noUsageRecords')" :description="t('dashboard.startUsingApi')" />
      </div>
      <div v-else class="divide-y divide-slate-900/[0.045] dark:divide-white/[0.05]">
        <div
          v-for="log in data"
          :key="log.id"
          class="flex items-center justify-between gap-4 rounded-xl px-3 py-3 transition-colors hover:bg-[#f8f9fb] dark:hover:bg-white/[0.03]"
        >
          <div class="flex min-w-0 items-center gap-3">
            <span class="h-2 w-2 shrink-0 rounded-full bg-emerald-500"></span>
            <div class="min-w-0">
              <p class="truncate font-mono text-xs font-medium text-[#111827] dark:text-white">{{ log.model }}</p>
              <p class="mt-0.5 text-[11px] text-[#8d99a8] dark:text-dark-400 tabular-nums">{{ formatDateTime(log.created_at) }}</p>
            </div>
          </div>
          <div class="shrink-0 text-right tabular-nums">
            <p class="font-mono text-xs font-semibold">
              <span class="text-[#111827] dark:text-white" :title="t('dashboard.actual')">${{ formatCost(log.actual_cost) }}</span>
              <span class="font-normal text-[#8d99a8] dark:text-dark-400" :title="t('dashboard.standard')"> / ${{ formatCost(log.total_cost) }}</span>
            </p>
            <p class="mt-0.5 text-[11px] text-[#8d99a8] dark:text-dark-400">{{ (log.input_tokens + log.output_tokens).toLocaleString() }} tokens</p>
          </div>
        </div>

        <router-link
          to="/usage"
          class="mt-1 flex items-center justify-center gap-1.5 pt-3.5 text-xs font-medium text-primary-600 transition-colors hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
        >
          {{ t('dashboard.viewAllUsage') }}
          <Icon name="arrowRight" size="xs" />
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'
import type { UsageLog } from '@/types'

defineProps<{
  data: UsageLog[]
  loading: boolean
}>()
const { t } = useI18n()
const formatCost = (c: number) => c.toFixed(4)
</script>
