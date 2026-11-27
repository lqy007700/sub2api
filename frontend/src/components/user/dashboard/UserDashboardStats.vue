<template>
  <!-- Asymmetric Apple Mist Hero Row -->
  <div class="grid grid-cols-1 gap-4 lg:grid-cols-12">
    <!-- Left Focal Card: Balance (or API Keys in simple mode) -->
    <div class="card flex flex-col justify-between gap-5 p-6 lg:col-span-7">
      <template v-if="!isSimple">
        <div>
          <div class="flex items-center justify-between gap-3">
            <span class="text-xs font-medium text-[#526071] dark:text-dark-300">
              {{ t('dashboard.balance') }}
            </span>
            <span class="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-400/15 dark:text-emerald-300">
              <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
              {{ stats?.active_api_keys || 0 }} / {{ stats?.total_api_keys || 0 }} {{ t('common.active') }}
            </span>
          </div>

          <div class="mt-3 flex items-baseline gap-0.5 font-semibold tracking-tight tabular-nums">
            <span class="mr-0.5 text-2xl font-normal text-[#8d99a8] dark:text-dark-400">$</span>
            <span class="text-4xl text-[#111827] dark:text-white">{{ balanceParts.int }}</span>
            <span class="text-2xl font-normal text-[#8d99a8] dark:text-dark-400">.{{ balanceParts.dec }}</span>
          </div>

          <p class="mt-2.5 text-xs text-[#526071] dark:text-dark-300 tabular-nums">
            <span>{{ t('dashboard.todayCost') }}</span>
            <strong class="ml-1.5 font-semibold text-[#111827] dark:text-white" :title="t('dashboard.actual')">
              ${{ formatCost(stats?.today_actual_cost || 0) }}
            </strong>
            <span class="text-[#8d99a8] dark:text-dark-400" :title="t('dashboard.standard')">
              / ${{ formatCost(stats?.today_cost || 0) }}
            </span>
          </p>

          <div class="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-slate-900/[0.05] pt-3.5 text-xs text-[#8d99a8] dark:border-white/[0.06] dark:text-dark-400 tabular-nums">
            <span>
              {{ t('common.total') }}
              <strong class="ml-1 font-medium text-[#526071] dark:text-dark-200" :title="t('dashboard.actual')">
                ${{ formatCost(stats?.total_actual_cost || 0) }}
              </strong>
              <span :title="t('dashboard.standard')"> / ${{ formatCost(stats?.total_cost || 0) }}</span>
            </span>
            <span>{{ t('dashboard.totalTokens') }}: {{ formatTokens(stats?.total_tokens || 0) }}</span>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-2.5">
          <button
            v-if="canPurchase"
            type="button"
            class="btn btn-primary rounded-full px-4 py-2 text-xs"
            @click="emit('navigate', '/purchase')"
          >
            {{ t('nav.recharge') }}
          </button>
          <button
            type="button"
            class="btn rounded-full px-4 py-2 text-xs"
            :class="canPurchase ? 'btn-secondary' : 'btn-primary'"
            @click="emit('navigate', '/keys')"
          >
            {{ t('dashboard.apiKeys') }}
          </button>
          <button
            type="button"
            class="btn btn-secondary rounded-full px-4 py-2 text-xs"
            @click="emit('navigate', canPurchase ? '/purchase?tab=redeem' : '/redeem')"
          >
            {{ t('dashboard.redeemCode') }}
          </button>
        </div>
      </template>

      <template v-else>
        <div>
          <div class="flex items-center justify-between gap-3">
            <span class="text-xs font-medium text-[#526071] dark:text-dark-300">
              {{ t('dashboard.apiKeys') }}
            </span>
            <span class="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-400/15 dark:text-emerald-300">
              <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
              {{ stats?.active_api_keys || 0 }} {{ t('common.active') }}
            </span>
          </div>

          <div class="mt-3 text-4xl font-semibold tracking-tight text-[#111827] dark:text-white tabular-nums">
            {{ stats?.total_api_keys || 0 }}
          </div>

          <p class="mt-2.5 text-xs text-[#526071] dark:text-dark-300 tabular-nums">
            {{ t('dashboard.totalTokens') }}:
            <strong class="ml-1 font-semibold text-[#111827] dark:text-white">{{ formatTokens(stats?.total_tokens || 0) }}</strong>
          </p>
        </div>

        <div class="flex items-center gap-2.5">
          <button
            type="button"
            class="btn btn-primary rounded-full px-4 py-2 text-xs"
            @click="emit('navigate', '/keys')"
          >
            {{ t('dashboard.createApiKey') }}
          </button>
        </div>
      </template>
    </div>

    <!-- Right Serene Metrics Stack -->
    <div class="card flex flex-col justify-center divide-y divide-slate-900/[0.05] px-5 py-2 dark:divide-white/[0.06] lg:col-span-5">
      <!-- Today Requests -->
      <div class="flex items-center justify-between gap-4 py-3.5">
        <div class="min-w-0">
          <p class="text-xs font-medium text-[#526071] dark:text-dark-300">{{ t('dashboard.todayRequests') }}</p>
          <p class="mt-0.5 text-[11px] text-[#8d99a8] dark:text-dark-400 tabular-nums">
            {{ t('common.total') }}: {{ formatNumber(stats?.total_requests || 0) }}
          </p>
        </div>
        <div class="text-right tabular-nums">
          <p class="text-lg font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ formatNumber(stats?.today_requests || 0) }}
          </p>
          <p class="text-[11px] text-[#8d99a8] dark:text-dark-400">
            {{ formatTokens(stats?.rpm || 0) }} RPM
          </p>
        </div>
      </div>

      <!-- Today Tokens -->
      <div class="flex items-center justify-between gap-4 py-3.5">
        <div class="min-w-0">
          <p class="text-xs font-medium text-[#526071] dark:text-dark-300">{{ t('dashboard.todayTokens') }}</p>
          <p class="mt-0.5 truncate text-[11px] text-[#8d99a8] dark:text-dark-400 tabular-nums">
            {{ t('dashboard.input') }} {{ formatTokens(stats?.today_input_tokens || 0) }} ·
            {{ t('dashboard.output') }} {{ formatTokens(stats?.today_output_tokens || 0) }} ·
            {{ t('dashboard.cache') }} {{ formatTokens((stats?.today_cache_creation_tokens || 0) + (stats?.today_cache_read_tokens || 0)) }}
          </p>
        </div>
        <div class="shrink-0 text-right tabular-nums">
          <p class="text-lg font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ formatTokens(stats?.today_tokens || 0) }}
          </p>
          <p class="text-[11px] text-primary-600 dark:text-primary-400">
            {{ formatTokens(stats?.tpm || 0) }} TPM
          </p>
        </div>
      </div>

      <!-- Avg Response Time -->
      <div class="flex items-center justify-between gap-4 py-3.5">
        <div class="min-w-0">
          <p class="text-xs font-medium text-[#526071] dark:text-dark-300">{{ t('dashboard.avgResponse') }}</p>
          <p class="mt-0.5 text-[11px] text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.averageTime') }}</p>
        </div>
        <div class="text-right tabular-nums">
          <p class="text-lg font-semibold tracking-tight text-[#111827] dark:text-white">
            {{ formatDuration(stats?.average_duration_ms || 0) }}
          </p>
          <p class="text-[11px] text-emerald-600 dark:text-emerald-400">
            {{ t('dashboard.performance') }}
          </p>
        </div>
      </div>
    </div>
  </div>

  <!-- Row 3: Per-platform breakdown -->
  <div v-if="!isSimple && platformCards.length > 0" class="card p-5">
    <div class="mb-3.5 flex items-center justify-between">
      <h3 class="text-sm font-semibold tracking-tight text-[#111827] dark:text-white">{{ t('dashboard.platformBreakdown') }}</h3>
      <span class="rounded-full bg-[#f4f5f8] px-2.5 py-0.5 text-[11px] font-medium text-[#526071] dark:bg-white/[0.06] dark:text-dark-300">
        {{ t('dashboard.platformCount', { count: platformCount }) }}
      </span>
    </div>
    <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="item in platformCards"
        :key="item.platform"
        data-testid="platform-card"
        :data-platform="item.platform"
        :class="[
          'rounded-xl border p-3.5 transition-colors',
          item.isOther
            ? 'border-dashed border-slate-900/[0.12] bg-[#f8f9fb] dark:border-white/[0.12] dark:bg-dark-900/50'
            : 'border-slate-900/[0.07] bg-[#f8f9fb]/80 dark:border-white/[0.07] dark:bg-dark-900/60'
        ]"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-semibold text-[#111827] dark:text-white">
            {{ item.isOther ? t('dashboard.platformOther') : platformLabel(item.platform) }}
          </span>
          <span class="font-mono text-sm font-semibold tabular-nums text-primary-600 dark:text-primary-400" :title="t('dashboard.actual')">
            ${{ formatCost(item.total_actual_cost) }}
          </span>
        </div>
        <div class="mt-2.5 space-y-1 text-xs tabular-nums">
          <div class="flex items-center justify-between">
            <span class="text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.todayCost') }}</span>
            <span class="font-mono text-[#111827] dark:text-white">${{ formatCost(item.today_actual_cost) }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.requests') }}</span>
            <span class="font-mono text-[#526071] dark:text-dark-200">
              {{ item.total_requests > 0 ? formatNumber(item.total_requests) : '-' }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-[#8d99a8] dark:text-dark-400">{{ t('dashboard.tokens') }}</span>
            <span class="font-mono text-[#526071] dark:text-dark-200">
              {{ item.total_tokens > 0 ? formatTokens(item.total_tokens) : '-' }}
            </span>
          </div>
        </div>

        <!-- Quota 区：仅当 quota 配置存在、非 __other__ 且至少有一个窗口配了 limit 时显示 -->
        <div v-if="hasAnyLimit(item.quota) && !item.isOther" class="mt-3 space-y-1.5 border-t border-slate-900/[0.06] pt-2.5 dark:border-white/[0.07]">
          <p class="text-[10px] font-medium uppercase tracking-wider text-[#8d99a8] dark:text-dark-400">
            {{ t('dashboard.platformQuota.title') }}
          </p>
          <template v-for="w in (['daily', 'weekly', 'monthly'] as const)" :key="w">
            <div v-if="quotaVal(item.quota, `${w}_limit_usd`) != null" class="space-y-1">
              <!-- limit=0：完全禁用 -->
              <template v-if="(quotaVal(item.quota, `${w}_limit_usd`) as number) === 0">
                <div class="flex items-center justify-between text-xs">
                  <span class="text-[#526071] dark:text-gray-300">{{ t(`dashboard.platformQuota.${w}`) }}</span>
                  <span class="font-mono text-red-500">{{ t('dashboard.platformQuota.disabled') }}</span>
                </div>
                <div class="h-1.5 w-full overflow-hidden rounded-full bg-slate-900/[0.08] dark:bg-white/[0.08]">
                  <div class="h-full w-full rounded-full bg-red-500" />
                </div>
              </template>
              <!-- limit>0：正常用量进度条 -->
              <template v-else>
                <div class="flex items-center justify-between text-xs tabular-nums">
                  <span class="text-[#526071] dark:text-gray-300">{{ t(`dashboard.platformQuota.${w}`) }}</span>
                  <span class="font-mono text-[#111827] dark:text-gray-200">
                    ${{ formatUsd((quotaVal(item.quota, `${w}_usage_usd`) as number) ?? 0) }} / ${{ formatUsd(quotaVal(item.quota, `${w}_limit_usd`) as number) }}
                  </span>
                </div>
                <div class="h-1.5 w-full overflow-hidden rounded-full bg-slate-900/[0.08] dark:bg-white/[0.08]">
                  <div
                    class="h-full rounded-full transition-all"
                    :class="quotaBarClass(calcPercent((quotaVal(item.quota, `${w}_usage_usd`) as number) ?? 0, quotaVal(item.quota, `${w}_limit_usd`) as number))"
                    :style="{ width: calcPercent((quotaVal(item.quota, `${w}_usage_usd`) as number) ?? 0, quotaVal(item.quota, `${w}_limit_usd`) as number) + '%' }"
                  />
                </div>
                <p v-if="quotaVal(item.quota, `${w}_window_resets_at`)" class="text-[10px] text-[#8d99a8] dark:text-dark-400">
                  {{ t('dashboard.platformQuota.resetsAt', { time: formatResetTime(quotaVal(item.quota, `${w}_window_resets_at`) as string) }) }}
                </p>
              </template>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PlatformDashboardStats, UserDashboardStats as UserStatsType } from '@/api/usage'
import type { PlatformQuotaItem } from '@/types'

interface FusedPlatformCard {
  platform: string
  total_actual_cost: number
  today_actual_cost: number
  total_requests: number
  total_tokens: number
  isOther?: boolean
  quota?: PlatformQuotaItem
}

const props = defineProps<{
  stats: UserStatsType
  balance: number
  isSimple: boolean
  platformQuotas?: PlatformQuotaItem[] | null
  canPurchase?: boolean
}>()

const emit = defineEmits<{
  navigate: [path: string]
}>()

const { t } = useI18n()

const balanceParts = computed(() => {
  const formatted = formatBalance(props.balance || 0)
  const [int = '0', dec = '00'] = formatted.split('.')
  return { int, dec }
})

const PLATFORM_LABELS: Record<string, string> = {
  anthropic: 'Claude',
  openai: 'OpenAI',
  gemini: 'Gemini',
  antigravity: 'Antigravity',
  grok: 'Grok',
  kimi: 'Kimi',
  zhipu: 'Zhipu GLM',
  deepseek: 'DeepSeek',
  minimax: 'MiniMax',
  typesafe: 'TypeSafe / Jev',
}

const platformLabel = (p: string) => PLATFORM_LABELS[p] ?? p

// 处理"各平台之和 < 总值"的差值：后端按平台聚合时过滤了无法归属平台的行
// （group 与 account 都缺 platform）。这里把差值作为"其他"卡片显式展示，
// 避免 Row 1 总值与 Row 3 平台拆分加总对不上、用户困惑。
const OTHER_THRESHOLD = 0.0001
const platformCards = computed<FusedPlatformCard[]>(() => {
  // 建立 by_platform Map
  const byPlat = new Map<string, PlatformDashboardStats>()
  for (const item of props.stats?.by_platform ?? []) byPlat.set(item.platform, item)

  // 建立 quota Map。三档全空的记录不产生卡片，挂到卡片上也不渲染配额区。
  const byQuota = new Map<string, PlatformQuotaItem>()
  for (const q of props.platformQuotas ?? []) byQuota.set(q.platform, q)

  // 卡片集合 = 有用量的平台 ∪ 至少配置了一档限额的平台。
  // 三档全空的限额记录等价于不限额，不单独产生卡片。
  // 后端 by_platform / quota 接口均不会返回 platform='__other__'，
  // 无需显式排除；__other__ 由下方差值补差逻辑单独追加。
  const platforms = new Set<string>(byPlat.keys())
  for (const [platform, q] of byQuota) {
    if (hasAnyLimit(q)) platforms.add(platform)
  }

  const PLATFORM_ORDER = ['anthropic', 'openai', 'gemini', 'antigravity', 'grok', 'typesafe']
  const cards: FusedPlatformCard[] = []

  for (const p of platforms) {
    const stat = byPlat.get(p)
    cards.push({
      platform: p,
      total_actual_cost: stat?.total_actual_cost ?? 0,
      today_actual_cost: stat?.today_actual_cost ?? 0,
      total_requests: stat?.total_requests ?? 0,
      total_tokens: stat?.total_tokens ?? 0,
      quota: byQuota.get(p),
    })
  }

  // 排序：按 PLATFORM_ORDER，未知平台按名称排序
  cards.sort((a, b) => {
    const ai = PLATFORM_ORDER.indexOf(a.platform)
    const bi = PLATFORM_ORDER.indexOf(b.platform)
    if (ai === -1 && bi === -1) return a.platform.localeCompare(b.platform)
    if (ai === -1) return 1
    if (bi === -1) return -1
    return ai - bi
  })

  // __other__ 补差逻辑：只对 by_platform 有 usage 数据的总和计算
  const total = props.stats?.total_actual_cost ?? 0
  const today = props.stats?.today_actual_cost ?? 0
  const sumTotal = cards.reduce((s, c) => s + c.total_actual_cost, 0)
  const sumToday = cards.reduce((s, c) => s + c.today_actual_cost, 0)
  const diffTotal = Math.max(0, total - sumTotal)
  const diffToday = Math.max(0, today - sumToday)

  if (diffTotal > OTHER_THRESHOLD || diffToday > OTHER_THRESHOLD) {
    cards.push({
      platform: '__other__',
      total_actual_cost: diffTotal,
      today_actual_cost: diffToday,
      total_requests: 0,
      total_tokens: 0,
      isOther: true,
    })
  }

  return cards
})

// 标题右侧的平台计数 = 实际渲染的平台卡片数，不含"其他"差额卡。
const platformCount = computed(() => platformCards.value.filter((c) => !c.isOther).length)

// Quota helpers

type QuotaWindow = 'daily' | 'weekly' | 'monthly'
type QuotaField = `${QuotaWindow}_limit_usd` | `${QuotaWindow}_usage_usd` | `${QuotaWindow}_window_resets_at`

function quotaVal(q: PlatformQuotaItem | undefined, key: QuotaField): PlatformQuotaItem[QuotaField] {
  return q?.[key]
}

function hasAnyLimit(q: PlatformQuotaItem | undefined): boolean {
  if (!q) return false
  return q.daily_limit_usd != null || q.weekly_limit_usd != null || q.monthly_limit_usd != null
}

function calcPercent(usage: number, limit: number): number {
  if (!limit || limit <= 0) return 0
  return Math.min(100, Math.max(0, Math.round((usage / limit) * 100)))
}

function quotaBarClass(p: number): string {
  if (p >= 95) return 'bg-red-500'
  if (p >= 75) return 'bg-amber-500'
  return 'bg-green-500'
}

// 与 formatBalance 一致使用 Intl.NumberFormat 做半偶舍入，避免 toFixed 在不同 JS 引擎
// 下偶发截断而非四舍五入（与后端展示精度不一致）。
const usdFormatter = new Intl.NumberFormat('en-US', {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})
function formatUsd(n: number): string {
  if (!Number.isFinite(n)) return '0.00'
  return usdFormatter.format(n)
}

function formatResetTime(iso: string | null | undefined): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString(undefined, {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

const formatBalance = (b: number) =>
  new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(b)

const formatNumber = (n: number) => n.toLocaleString()
const formatCost = (c: number) => c.toFixed(4)
const formatTokens = (t: number) => {
  if (t >= 1_000_000) return `${(t / 1_000_000).toFixed(1)}M`
  if (t >= 1000) return `${(t / 1000).toFixed(1)}K`
  return t.toString()
}
const formatDuration = (ms: number) => ms >= 1000 ? `${(ms / 1000).toFixed(2)}s` : `${ms.toFixed(0)}ms`
</script>
