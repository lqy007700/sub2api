<template>
  <div class="card overflow-hidden">
    <div class="flex items-center justify-between border-b border-slate-900/[0.06] px-5 py-4 dark:border-white/[0.07] sm:px-6">
      <h2 class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">
        {{ title }}
      </h2>
      <span class="text-xs tabular-nums text-[#8d99a8] dark:text-dark-400">
        {{ t('common.total') }}: {{ historyTotal }} {{ t('pagination.results') }}
      </span>
    </div>
    <div class="p-5 sm:p-6">
      <!-- Loading State -->
      <div v-if="loadingHistory" class="flex items-center justify-center py-8">
        <svg class="h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          ></circle>
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          ></path>
        </svg>
      </div>

      <!-- History List -->
      <div v-else-if="history.length > 0" class="divide-y divide-slate-900/[0.06] dark:divide-white/[0.07]">
        <div
          v-for="item in history"
          :key="item.id"
          class="flex items-center justify-between gap-4 py-3.5 first:pt-0 last:pb-0"
        >
          <div class="flex min-w-0 items-center gap-3.5">
            <div
              :class="[
                'flex h-9 w-9 shrink-0 items-center justify-center rounded-xl',
                isBalanceType(item.type)
                  ? item.value >= 0
                    ? 'bg-emerald-500/10 dark:bg-emerald-400/15'
                    : 'bg-red-500/10 dark:bg-red-400/15'
                  : isSubscriptionType(item.type)
                    ? 'bg-purple-500/10 dark:bg-purple-400/15'
                    : item.value >= 0
                      ? 'bg-primary-500/10 dark:bg-primary-400/15'
                      : 'bg-amber-500/10 dark:bg-amber-400/15'
              ]"
            >
              <!-- 余额类型图标 -->
              <Icon
                v-if="isBalanceType(item.type)"
                name="dollar"
                size="sm"
                :class="
                  item.value >= 0
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-red-600 dark:text-red-400'
                "
              />
              <!-- 订阅类型图标 -->
              <Icon
                v-else-if="isSubscriptionType(item.type)"
                name="badge"
                size="sm"
                class="text-purple-600 dark:text-purple-400"
              />
              <!-- 并发类型图标 -->
              <Icon
                v-else
                name="bolt"
                size="sm"
                :class="
                  item.value >= 0
                    ? 'text-primary-600 dark:text-primary-400'
                    : 'text-amber-600 dark:text-amber-400'
                "
              />
            </div>
            <div class="min-w-0">
              <p class="truncate text-sm font-medium text-[#111827] dark:text-white">
                {{ getHistoryItemTitle(item) }}
              </p>
              <p class="mt-0.5 text-xs tabular-nums text-[#8d99a8] dark:text-dark-400">
                {{ formatDateTime(item.used_at) }}
              </p>
            </div>
          </div>
          <div class="shrink-0 text-right">
            <p
              :class="[
                'font-mono text-sm font-semibold tabular-nums',
                isBalanceType(item.type)
                  ? item.value >= 0
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-red-600 dark:text-red-400'
                  : isSubscriptionType(item.type)
                    ? 'text-purple-600 dark:text-purple-400'
                    : item.value >= 0
                      ? 'text-primary-600 dark:text-primary-400'
                      : 'text-amber-600 dark:text-amber-400'
              ]"
            >
              {{ formatHistoryValue(item) }}
            </p>
            <p
              v-if="!isAdminAdjustment(item.type)"
              class="mt-0.5 font-mono text-xs text-[#8d99a8] dark:text-dark-500"
            >
              {{ item.code.slice(0, 8) }}...
            </p>
            <p v-else class="mt-0.5 text-xs text-[#8d99a8] dark:text-dark-500">
              {{ t('redeem.adminAdjustment') }}
            </p>
            <!-- Display notes for admin adjustments -->
            <p
              v-if="item.notes"
              class="mt-1 max-w-[200px] truncate text-xs italic text-[#526071] dark:text-dark-400"
              :title="item.notes"
            >
              {{ item.notes }}
            </p>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state py-8">
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#f4f5f8] dark:bg-dark-800"
        >
          <Icon name="clock" size="lg" class="text-[#8d99a8] dark:text-dark-500" />
        </div>
        <p class="text-sm text-[#526071] dark:text-dark-400">
          {{ t('redeem.historyWillAppear') }}
        </p>
      </div>

      <div class="mt-5 flex flex-wrap items-center justify-between gap-3 border-t border-slate-900/[0.06] pt-4 text-xs text-[#526071] dark:border-white/[0.07] dark:text-dark-300">
        <label class="inline-flex items-center gap-2">
          <span>{{ t('pagination.perPage') }}</span>
          <select
            v-model="historyPageSize"
            class="input h-8 w-20 py-1 text-xs"
            :disabled="loadingHistory || submitting"
            @change="fetchHistory(1)"
          >
            <option v-for="size in [20, 50, 100]" :key="size" :value="size">{{ size }}</option>
          </select>
        </label>
        <div class="flex items-center gap-2">
          <button
            class="btn btn-secondary btn-sm"
            :disabled="loadingHistory || submitting || historyPage <= 1"
            @click="fetchHistory(historyPage - 1)"
          >{{ t('pagination.previous') }}</button>
          <button
            class="btn btn-secondary btn-sm"
            :disabled="loadingHistory || submitting || historyPage * historyPageSize >= historyTotal"
            @click="fetchHistory(historyPage + 1)"
          >{{ t('pagination.next') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { redeemAPI, type RedeemHistoryItem } from '@/api'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ title: string; submitting?: boolean; codesOnly?: boolean }>()
const { t } = useI18n()
const appStore = useAppStore()

// History data
const history = ref<RedeemHistoryItem[]>([])
const loadingHistory = ref(false)
const historyPage = ref(1)
const historyPageSize = ref(20)
const historyTotal = ref(0)
let historyRequest = 0
let loadedHistoryPageSize = 20

// Helper functions for history display
const isBalanceType = (type: string) => {
  return type === 'balance' || type === 'admin_balance'
}

const isSubscriptionType = (type: string) => {
  return type === 'subscription'
}

const isAdminAdjustment = (type: string) => {
  return type === 'admin_balance' || type === 'admin_concurrency'
}

const getHistoryItemTitle = (item: RedeemHistoryItem) => {
  if (item.type === 'balance') {
    return t('redeem.balanceAddedRedeem')
  } else if (item.type === 'admin_balance') {
    return item.value >= 0 ? t('redeem.balanceAddedAdmin') : t('redeem.balanceDeductedAdmin')
  } else if (item.type === 'concurrency') {
    return t('redeem.concurrencyAddedRedeem')
  } else if (item.type === 'admin_concurrency') {
    return item.value >= 0 ? t('redeem.concurrencyAddedAdmin') : t('redeem.concurrencyReducedAdmin')
  } else if (item.type === 'subscription') {
    return t('redeem.subscriptionAssigned')
  }
  return t('common.unknown')
}

const formatHistoryValue = (item: RedeemHistoryItem) => {
  if (isBalanceType(item.type)) {
    const sign = item.value >= 0 ? '+' : ''
    return `${sign}$${item.value.toFixed(2)}`
  } else if (isSubscriptionType(item.type)) {
    // 订阅类型显示有效天数和分组名称
    const days = item.validity_days || Math.round(item.value)
    const groupName = item.group?.name || ''
    return groupName ? `${days}${t('redeem.days')} - ${groupName}` : `${days}${t('redeem.days')}`
  } else {
    const sign = item.value >= 0 ? '+' : ''
    return `${sign}${item.value} ${t('redeem.requests')}`
  }
}

const fetchHistory = async (page = 1) => {
  const request = ++historyRequest
  const pageSize = historyPageSize.value
  loadingHistory.value = true
  try {
    const result = props.codesOnly
      ? await redeemAPI.getHistory(page, pageSize, { codesOnly: true })
      : await redeemAPI.getHistory(page, pageSize)
    if (request !== historyRequest) return
    history.value = result.items
    historyTotal.value = result.total
    historyPage.value = page
    historyPageSize.value = pageSize
    loadedHistoryPageSize = pageSize
  } catch (error) {
    if (request !== historyRequest) return
    historyPageSize.value = loadedHistoryPageSize
    appStore.showError(t('redeem.historyLoadFailed'))
    console.error('Failed to fetch history:', error)
  } finally {
    if (request === historyRequest) loadingHistory.value = false
  }
}

defineExpose({ refresh: fetchHistory })

onMounted(() => { void fetchHistory() })
onBeforeUnmount(() => { historyRequest++ })
</script>
