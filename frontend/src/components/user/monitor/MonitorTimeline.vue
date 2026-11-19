<template>
  <div class="mt-4 pt-3 border-t border-gray-100 dark:border-dark-700/60">
    <div class="mb-2 flex justify-between text-[10px] font-semibold uppercase tracking-widest text-gray-400">
      <span>{{ t('channelStatus.dailyHistory', { n: days }) }}</span>
      <span class="tabular-nums">{{ t('monitorCommon.nextUpdateIn', { n: countdownSeconds }) }}</span>
    </div>

    <div class="flex h-5 w-full items-end gap-[2px]" role="img" :aria-label="t('channelStatus.dailyHistory', { n: days })">
      <div
        v-for="day in days"
        :key="day"
        class="h-full min-w-0 flex-1 rounded-sm bg-emerald-500"
      ></div>
    </div>

    <div class="mt-1 flex justify-between text-[10px] text-gray-500 dark:text-gray-400">
      <span>{{ t('monitorCommon.past') }}</span>
      <span>{{ t('monitorCommon.now') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

defineProps<{
  days: number
  countdownSeconds: number
}>()

const { t } = useI18n()
</script>
