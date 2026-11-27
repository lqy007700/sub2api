<template>
  <div class="space-y-4">
    <!-- Compact Inline Controls Strip -->
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div class="flex flex-wrap items-center gap-2.5">
        <DateRangePicker
          :start-date="startDate"
          :end-date="endDate"
          @update:startDate="$emit('update:startDate', $event)"
          @update:endDate="$emit('update:endDate', $event)"
          @change="$emit('dateRangeChange', $event)"
        />
        <button @click="$emit('refresh')" :disabled="loading" class="btn btn-secondary btn-sm rounded-full px-3.5">
          {{ t('common.refresh') }}
        </button>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-xs font-medium text-[#526071] dark:text-dark-300">{{ t('dashboard.granularity') }}</span>
        <div class="w-28">
          <Select
            :model-value="granularity"
            :options="[{ value: 'day', label: t('dashboard.day') }, { value: 'hour', label: t('dashboard.hour') }]"
            @update:model-value="$emit('update:granularity', $event)"
            @change="$emit('granularityChange')"
          />
        </div>
      </div>
    </div>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-12">
      <!-- Token Usage Trend Chart (Primary Wide Card) -->
      <div class="lg:col-span-7">
        <TokenUsageTrend :trend-data="trend" :loading="loading" />
      </div>

      <!-- Model Distribution Card (Apple Mist Segmented Bar + Tabular Breakdown) -->
      <div class="card relative flex flex-col justify-between overflow-hidden p-5 lg:col-span-5">
        <div v-if="loading" class="absolute inset-0 z-10 flex items-center justify-center bg-white/60 backdrop-blur-sm dark:bg-dark-800/60">
          <LoadingSpinner size="md" />
        </div>
        <div class="mb-3.5 flex items-center justify-between">
          <h3 class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">{{ t('dashboard.modelDistribution') }}</h3>
          <span v-if="models.length > 0" class="text-[11px] text-[#8d99a8] dark:text-dark-400">
            {{ models.length }} {{ t('dashboard.model') }}
          </span>
        </div>

        <div v-if="modelSegments.length > 0" class="flex flex-1 flex-col">
          <!-- Proportional Segmented Bar -->
          <div class="mb-4 flex h-2 w-full gap-1 overflow-hidden rounded-full">
            <div
              v-for="seg in modelSegments"
              :key="seg.model"
              class="h-full rounded-full transition-all duration-300"
              :style="{ width: `${Math.max(seg.percent, 2)}%`, backgroundColor: seg.color }"
              :title="`${seg.model}: ${seg.percent.toFixed(1)}%`"
            />
          </div>

          <!-- Clean Breakdown Table -->
          <div class="max-h-48 w-full min-w-0 flex-1 overflow-auto">
            <table class="w-full text-xs tabular-nums">
              <thead>
                <tr class="text-[#8d99a8] dark:text-dark-400">
                  <th class="pb-2 text-left font-medium">{{ t('dashboard.model') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.tokens') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.actual') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.standard') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="seg in modelSegments"
                  :key="seg.model"
                  class="border-t border-slate-900/[0.045] dark:border-white/[0.05]"
                >
                  <td class="max-w-[140px] py-2 font-medium text-[#111827] dark:text-white" :title="seg.model">
                    <div class="flex items-center gap-2">
                      <span class="h-2 w-2 shrink-0 rounded-full" :style="{ backgroundColor: seg.color }"></span>
                      <span class="truncate font-mono text-[11.5px]">{{ seg.model }}</span>
                    </div>
                  </td>
                  <td class="py-2 text-right text-[#526071] dark:text-dark-300">{{ formatTokens(seg.total_tokens) }}</td>
                  <td class="py-2 text-right font-medium text-[#111827] dark:text-white">${{ formatCost(seg.actual_cost) }}</td>
                  <td class="py-2 text-right text-[#8d99a8] dark:text-dark-400">${{ formatCost(seg.cost) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div v-else class="flex h-48 items-center justify-center text-sm text-[#8d99a8] dark:text-dark-400">
          {{ t('dashboard.noDataAvailable') }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Select from '@/components/common/Select.vue'
import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import type { TrendDataPoint, ModelStat } from '@/types'
import { formatCostFixed as formatCost, formatTokensK as formatTokens } from '@/utils/format'

const props = defineProps<{ loading: boolean, startDate: string, endDate: string, granularity: string, trend: TrendDataPoint[], models: ModelStat[] }>()
defineEmits(['update:startDate', 'update:endDate', 'update:granularity', 'dateRangeChange', 'granularityChange', 'refresh'])
const { t } = useI18n()

const MIST_PALETTE = ['#0071e3', '#5e5ce6', '#64748b', '#0ea5e9', '#16a34a', '#94a3b8', '#cbd5e1', '#38bdf8']

const modelSegments = computed(() => {
  if (!props.models?.length) return []
  const total = props.models.reduce((sum, m) => sum + (m.total_tokens || 0), 0)
  return props.models.map((m, idx) => ({
    ...m,
    percent: total > 0 ? ((m.total_tokens || 0) / total) * 100 : 0,
    color: MIST_PALETTE[idx % MIST_PALETTE.length]
  }))
})
</script>
