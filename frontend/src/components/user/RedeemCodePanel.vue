<template>
  <div class="space-y-5">
    <!-- Redeem Form Card (Apple Mist) -->
    <div class="card p-5 sm:p-6">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
        <label for="code" class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">
          {{ t('redeem.redeemCodeLabel') }}
        </label>
        <span class="text-xs text-[#8d99a8] dark:text-dark-400">
          {{ t('redeem.redeemCodeHint') }}
        </span>
      </div>

      <form @submit.prevent="handleRedeem" class="space-y-4">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-stretch">
          <div
            class="flex flex-1 items-center rounded-xl border border-slate-900/[0.1] bg-[#f8f9fb] pl-3.5 pr-2 transition-all focus-within:border-primary-500 focus-within:bg-white focus-within:ring-2 focus-within:ring-primary-500/20 dark:border-white/[0.1] dark:bg-dark-900/60 dark:focus-within:border-primary-400 dark:focus-within:bg-dark-900"
          >
            <Icon name="gift" size="md" class="mr-2.5 shrink-0 text-[#8d99a8] dark:text-dark-400" />
            <input
              id="code"
              v-model="redeemCode"
              type="text"
              required
              :placeholder="t('redeem.redeemCodePlaceholder')"
              :disabled="submitting"
              class="h-12 w-full border-0 bg-transparent font-mono text-base font-semibold tracking-wide text-[#111827] placeholder:font-sans placeholder:font-normal placeholder:tracking-normal placeholder:text-[#8d99a8] focus:outline-none focus:ring-0 disabled:opacity-60 dark:text-white dark:placeholder:text-dark-500"
            />
            <button
              v-if="canPaste"
              type="button"
              :disabled="submitting"
              class="ml-2 shrink-0 rounded-lg border border-slate-900/[0.08] bg-white px-2.5 py-1.5 text-xs font-medium text-[#526071] transition-colors hover:text-[#111827] disabled:opacity-50 dark:border-white/[0.08] dark:bg-dark-800 dark:text-dark-300 dark:hover:text-white"
              @click="handlePaste"
            >
              {{ t('payment.hub.pasteCode') }}
            </button>
          </div>

          <button
            type="submit"
            :disabled="!redeemCode || submitting"
            class="btn btn-primary h-12 shrink-0 rounded-xl px-6 text-sm font-semibold sm:min-w-[136px]"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
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
            <Icon v-else name="checkCircle" size="sm" class="mr-1.5" />
            {{ submitting ? t('redeem.redeeming') : t('redeem.redeemButton') }}
          </button>
        </div>

        <div
          v-if="showStoreLink || showOnlineLink"
          class="flex flex-wrap items-center justify-between gap-2 pt-0.5 text-xs text-[#8d99a8] dark:text-dark-400"
        >
          <span>{{ t('redeem.codeRule1') }}</span>
          <div class="flex flex-wrap items-center gap-3">
            <button
              v-if="showStoreLink"
              type="button"
              class="font-semibold text-primary-600 transition-colors hover:underline dark:text-primary-400"
              @click="emit('open-store')"
            >
              {{ t('payment.hub.noCodeBuyStore') }} →
            </button>
            <button
              v-if="showOnlineLink"
              type="button"
              class="font-medium text-[#526071] transition-colors hover:text-[#111827] hover:underline dark:text-dark-300 dark:hover:text-white"
              @click="emit('open-online')"
            >
              {{ t('payment.hub.goOnlineRecharge') }} →
            </button>
          </div>
        </div>
      </form>
    </div>

    <!-- Success Message -->
    <transition name="fade">
      <div
        v-if="redeemResult"
        class="rounded-2xl border border-emerald-500/25 bg-emerald-50/80 p-5 dark:border-emerald-400/25 dark:bg-emerald-500/10"
      >
        <div class="flex items-start gap-3.5">
          <div
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-emerald-500/15 dark:bg-emerald-400/20"
          >
            <Icon name="checkCircle" size="md" class="text-emerald-600 dark:text-emerald-400" />
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-sm font-semibold text-emerald-900 dark:text-emerald-200">
              {{ t('redeem.redeemSuccess') }}
            </h3>
            <div class="mt-1.5 text-sm text-emerald-800 dark:text-emerald-300">
              <p>{{ redeemResult.message }}</p>
              <div class="mt-2.5 flex flex-wrap gap-x-5 gap-y-1 text-xs">
                <p v-if="redeemResult.type === 'balance'" class="font-semibold tabular-nums">
                  {{ t('redeem.added') }}: ${{ redeemResult.value.toFixed(2) }}
                </p>
                <p v-else-if="redeemResult.type === 'concurrency'" class="font-semibold tabular-nums">
                  {{ t('redeem.added') }}: {{ redeemResult.value }}
                  {{ t('redeem.concurrentRequests') }}
                </p>
                <p v-else-if="redeemResult.type === 'subscription'" class="font-semibold">
                  {{ t('redeem.subscriptionAssigned') }}
                  <span v-if="redeemResult.group_name"> - {{ redeemResult.group_name }}</span>
                  <span v-if="redeemResult.validity_days">
                    ({{
                      t('redeem.subscriptionDays', { days: redeemResult.validity_days })
                    }})</span
                  >
                </p>
                <p v-if="redeemResult.new_balance !== undefined" class="tabular-nums">
                  {{ t('redeem.newBalance') }}:
                  <span class="font-semibold">${{ redeemResult.new_balance.toFixed(2) }}</span>
                </p>
                <p v-if="redeemResult.new_concurrency !== undefined" class="tabular-nums">
                  {{ t('redeem.newConcurrency') }}:
                  <span class="font-semibold">
                    {{ redeemResult.new_concurrency }} {{ t('redeem.requests') }}
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Error Message -->
    <transition name="fade">
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-500/25 bg-red-50/80 p-5 dark:border-red-400/25 dark:bg-red-500/10"
      >
        <div class="flex items-start gap-3.5">
          <div
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-red-500/15 dark:bg-red-400/20"
          >
            <Icon name="exclamationCircle" size="md" class="text-red-600 dark:text-red-400" />
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-sm font-semibold text-red-900 dark:text-red-200">
              {{ t('redeem.redeemFailed') }}
            </h3>
            <p class="mt-1 text-sm text-red-700 dark:text-red-300">
              {{ errorMessage }}
            </p>
          </div>
        </div>
      </div>
    </transition>

    <!-- Redemption History -->
    <RedeemHistory ref="history" :title="t('redeem.recentActivity')" :submitting="submitting" />

    <!-- Subtle Information Card -->
    <div
      class="rounded-2xl border border-slate-900/[0.06] bg-[#f8f9fb] p-5 dark:border-white/[0.07] dark:bg-dark-900/50"
    >
      <div class="flex items-start gap-3.5">
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-xl bg-primary-500/10 text-primary-600 dark:bg-primary-400/15 dark:text-primary-400"
        >
          <Icon name="infoCircle" size="sm" />
        </div>
        <div class="min-w-0 flex-1">
          <h3 class="text-xs font-semibold text-[#111827] dark:text-white">
            {{ t('redeem.aboutCodes') }}
          </h3>
          <ul class="mt-2 list-inside list-disc space-y-1 text-xs leading-relaxed text-[#526071] dark:text-dark-300">
            <li>{{ t('redeem.codeRule1') }}</li>
            <li>{{ t('redeem.codeRule2') }}</li>
            <li>
              {{ t('redeem.codeRule3') }}
              <span
                v-if="contactInfo"
                class="ml-1.5 inline-flex items-center rounded-md border border-slate-900/[0.08] bg-white px-2 py-0.5 text-[11px] font-medium text-[#111827] dark:border-white/[0.1] dark:bg-dark-800 dark:text-dark-100"
              >
                {{ contactInfo }}
              </span>
            </li>
            <li>{{ t('redeem.codeRule4') }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useSubscriptionStore } from '@/stores/subscriptions'
import { redeemAPI, authAPI } from '@/api'
import Icon from '@/components/icons/Icon.vue'
import RedeemHistory from '@/components/user/RedeemHistory.vue'

withDefaults(
  defineProps<{
    showStoreLink?: boolean
    showOnlineLink?: boolean
  }>(),
  {
    showStoreLink: false,
    showOnlineLink: false,
  }
)

const emit = defineEmits<{
  'open-store': []
  'open-online': []
}>()

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const subscriptionStore = useSubscriptionStore()

const redeemCode = ref('')
const submitting = ref(false)
const redeemResult = ref<{
  message: string
  type: string
  value: number
  new_balance?: number
  new_concurrency?: number
  group_name?: string
  validity_days?: number
} | null>(null)
const errorMessage = ref('')

const history = ref<InstanceType<typeof RedeemHistory> | null>(null)
const contactInfo = ref('')

const canPaste = computed(
  () => typeof navigator !== 'undefined' && Boolean(navigator.clipboard?.readText)
)

const handlePaste = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      redeemCode.value = text.trim()
    }
  } catch {
    // Ignore clipboard permission denial silently
  }
}

const handleRedeem = async () => {
  if (!redeemCode.value.trim()) {
    appStore.showError(t('redeem.pleaseEnterCode'))
    return
  }

  submitting.value = true
  errorMessage.value = ''
  redeemResult.value = null

  try {
    const result = await redeemAPI.redeem(redeemCode.value.trim())

    redeemResult.value = result

    // Refresh user data to get updated balance/concurrency
    try {
      await authStore.refreshUser()
    } catch (error) {
      console.error('Failed to refresh user after redeem:', error)
      appStore.showWarning(t('redeem.userRefreshFailed'))
    }

    // If subscription type, immediately refresh subscription status
    if (result.type === 'subscription') {
      try {
        await subscriptionStore.fetchActiveSubscriptions(true) // force refresh
      } catch (error) {
        console.error('Failed to refresh subscriptions after redeem:', error)
        appStore.showWarning(t('redeem.subscriptionRefreshFailed'))
      }
    }

    // Clear the input
    redeemCode.value = ''

    // Refresh history
    await history.value?.refresh()

    // Show success toast
    appStore.showSuccess(t('redeem.codeRedeemSuccess'))
  } catch (error: any) {
    errorMessage.value = error.response?.data?.detail || t('redeem.failedToRedeem')

    appStore.showError(t('redeem.redeemFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  if (appStore.contactInfo) {
    contactInfo.value = appStore.contactInfo
    return
  }
  try {
    const settings = await authAPI.getPublicSettings()
    contactInfo.value = settings.contact_info || ''
  } catch (error) {
    console.error('Failed to load contact info:', error)
  }
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
