<template>
  <div class="mt-4 pt-3 border-t border-gray-100 dark:border-dark-700/60">
    <div class="mb-2 flex justify-between text-[10px] font-semibold uppercase tracking-widest text-gray-400">
      <span>{{ t('channelStatus.dailyHistory', { n: days }) }}</span>
      <span class="tabular-nums">{{ t('monitorCommon.nextUpdateIn', { n: countdownSeconds }) }}</span>
    </div>

    <div class="flex h-5 w-full items-end gap-[2px]" role="list" :aria-label="t('channelStatus.dailyHistory', { n: days })">
      <div
        v-for="bar in bars"
        :key="bar.date"
        role="listitem"
        class="min-w-0 flex-1 rounded-sm"
        :class="bar.classes"
        :title="bar.label"
        :aria-label="bar.label"
      ></div>
    </div>

    <div class="mt-1 flex justify-between text-[10px] text-gray-500 dark:text-gray-400">
      <span v-if="bars.some(bar => bar.placeholder)">{{ t('channelStatus.historyPlaceholder') }}</span>
      <span v-else>{{ t('monitorCommon.past') }}</span>
      <span>{{ t('monitorCommon.now') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { UserMonitorView } from '@/api/channelMonitor'
import { useChannelMonitorFormat, userMonitorDisplayStatus } from '@/composables/useChannelMonitorFormat'

const props = defineProps<{
  buckets?: UserMonitorView['daily_timeline']
  days: number
  monitorId: number
  countdownSeconds: number
}>()

const { t } = useI18n()
const { statusLabel } = useChannelMonitorFormat()

const statusClasses = {
  operational: 'h-full bg-emerald-500',
  failed: 'h-1/3 bg-red-500',
  error: 'h-1/3 bg-red-500',
}

const millisecondsPerDay = 24 * 60 * 60 * 1000

const bars = computed(() => {
  const statusesByDate = new Map(props.buckets?.map(point => [point.date, point.status]))
  const today = Math.floor(Date.now() / millisecondsPerDay)
  return Array.from({ length: props.days }, (_, index) => {
    const day = today - props.days + 1 + index
    const date = new Date(day * millisecondsPerDay).toISOString().slice(0, 10)
    const status = userMonitorDisplayStatus(statusesByDate.get(date) ?? '')
    const recordedClasses = status && statusClasses[status]
    if (recordedClasses) {
      return { date, classes: recordedClasses, label: `${date} UTC · ${statusLabel(status)}`, placeholder: false }
    }
    // Anchor the decorative accents to UTC dates so switching ranges keeps overlapping days identical.
    const accent = (day + props.monitorId) % 15 === 0
    return {
      date,
      classes: accent
        ? 'h-2/3 border border-dashed border-amber-300 bg-amber-100 dark:border-amber-700 dark:bg-amber-950/40'
        : 'h-full border border-dashed border-emerald-300 bg-emerald-100 dark:border-emerald-700 dark:bg-emerald-950/40',
      label: `${date} UTC · ${t('channelStatus.noHistorySample')}`,
      placeholder: true,
    }
  })
})
</script>
