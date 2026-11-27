import type { CustomMenuItem } from '@/types'

const THIRD_PARTY_RECHARGE_KEYWORD_RE =
  /(充值|购卡|买卡|卡密|发卡|兑换码|买码|购买|商店|商城|recharge|top[-_\s]?up|buy|shop|store|card|faka|redeem|purchase|pay)/i

function isHttpUrl(rawUrl?: string): boolean {
  if (!rawUrl) return false
  const trimmed = rawUrl.trim()
  return /^https?:\/\//i.test(trimmed)
}

function isMarkdownCustomMenuItem(item: CustomMenuItem): boolean {
  return Boolean(item.page_slug) || Boolean(item.url?.trim().startsWith('md:'))
}

/**
 * 判断某个用户侧自定义菜单项是否为「三方购买兑换码 / 充值发卡」页面。
 * 命中项会被收口到统一「充值中心」的「购买兑换码」标签中，并从用户顶栏移除避免重复展示。
 */
export function isThirdPartyRechargeCustomMenuItem(item: CustomMenuItem): boolean {
  if (!item || item.visibility !== 'user' || item.hide_in_sidebar) return false
  if (isMarkdownCustomMenuItem(item)) return false
  if (!isHttpUrl(item.url)) return false

  const haystack = `${item.label || ''} ${item.id || ''} ${item.url || ''}`
  return THIRD_PARTY_RECHARGE_KEYWORD_RE.test(haystack)
}

export function getThirdPartyRechargeCustomMenuItems(items?: CustomMenuItem[] | null): CustomMenuItem[] {
  if (!Array.isArray(items) || items.length === 0) return []
  return items
    .filter((item) => isThirdPartyRechargeCustomMenuItem(item))
    .sort((a, b) => a.sort_order - b.sort_order)
}

export function getStandaloneUserCustomMenuItems(
  items: CustomMenuItem[] | null | undefined,
  mergeRechargeIntoHub: boolean
): CustomMenuItem[] {
  if (!Array.isArray(items) || items.length === 0) return []
  return items
    .filter((item) => {
      if (item.visibility !== 'user' || item.hide_in_sidebar) return false
      if (mergeRechargeIntoHub && isThirdPartyRechargeCustomMenuItem(item)) return false
      return true
    })
    .sort((a, b) => a.sort_order - b.sort_order)
}
