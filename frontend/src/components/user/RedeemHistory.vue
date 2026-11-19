<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        {{ title }}
      </h2>
    </div>
    <div class="p-6">
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
      <div v-else-if="history.length > 0" class="space-y-3">
        <div
          v-for="item in history"
          :key="item.id"
          class="flex items-center justify-between rounded-xl bg-gray-50 p-4 dark:bg-dark-800"
        >
          <div class="flex items-center gap-4">
            <div
              :class="[
                'flex h-10 w-10 items-center justify-center rounded-xl',
                isBalanceType(item.type)
                  ? item.value >= 0
                    ? 'bg-emerald-100 dark:bg-emerald-900/30'
                    : 'bg-red-100 dark:bg-red-900/30'
                  : isSubscriptionType(item.type)
                    ? 'bg-purple-100 dark:bg-purple-900/30'
                    : item.value >= 0
                      ? 'bg-blue-100 dark:bg-blue-900/30'
                      : 'bg-orange-100 dark:bg-orange-900/30'
              ]"
            >
              <!-- 余额类型图标 -->
              <Icon
                v-if="isBalanceType(item.type)"
                name="dollar"
                size="md"
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
                size="md"
                class="text-purple-600 dark:text-purple-400"
              />
              <!-- 并发类型图标 -->
              <Icon
                v-else
                name="bolt"
                size="md"
                :class="
                  item.value >= 0
                    ? 'text-blue-600 dark:text-blue-400'
                    : 'text-orange-600 dark:text-orange-400'
                "
              />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ getHistoryItemTitle(item) }}
              </p>
              <p class="text-xs text-gray-500 dark:text-dark-400">
                {{ formatDateTime(item.used_at) }}
              </p>
            </div>
          </div>
          <div class="text-right">
            <p
              :class="[
                'text-sm font-semibold',
                isBalanceType(item.type)
                  ? item.value >= 0
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-red-600 dark:text-red-400'
                  : isSubscriptionType(item.type)
                    ? 'text-purple-600 dark:text-purple-400'
                    : item.value >= 0
                      ? 'text-blue-600 dark:text-blue-400'
                      : 'text-orange-600 dark:text-orange-400'
              ]"
            >
              {{ formatHistoryValue(item) }}
            </p>
            <p
              v-if="!isAdminAdjustment(item.type)"
              class="font-mono text-xs text-gray-400 dark:text-dark-500"
            >
              {{ item.code.slice(0, 8) }}...
            </p>
            <p v-else class="text-xs text-gray-400 dark:text-dark-500">
              {{ t('redeem.adminAdjustment') }}
            </p>
            <!-- Display notes for admin adjustments -->
            <p
              v-if="item.notes"
              class="mt-1 text-xs text-gray-500 dark:text-dark-400 italic max-w-[200px] truncate"
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
          class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800"
        >
          <Icon name="clock" size="xl" class="text-gray-400 dark:text-dark-500" />
        </div>
        <p class="text-sm text-gray-500 dark:text-dark-400">
          {{ t('redeem.historyWillAppear') }}
        </p>
      </div>
      <div class="mt-4 flex flex-wrap items-center justify-between gap-3 text-sm">
        <span>{{ t('common.total') }}: {{ historyTotal }} {{ t('pagination.results') }}</span>
        <label>
          {{ t('pagination.perPage') }}
          <select
            v-model="historyPageSize"
            class="input w-20"
            :disabled="loadingHistory || submitting"
            @change="fetchHistory(1)"
          >
            <option v-for="size in [20, 50, 100]" :key="size" :value="size">{{ size }}</option>
          </select>
        </label>
        <button
          class="btn btn-secondary"
          :disabled="loadingHistory || submitting || historyPage <= 1"
          @click="fetchHistory(historyPage - 1)"
        >{{ t('pagination.previous') }}</button>
        <button
          class="btn btn-secondary"
          :disabled="loadingHistory || submitting || historyPage * historyPageSize >= historyTotal"
          @click="fetchHistory(historyPage + 1)"
        >{{ t('pagination.next') }}</button>
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
