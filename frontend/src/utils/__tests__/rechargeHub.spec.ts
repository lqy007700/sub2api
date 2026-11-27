import { describe, expect, it } from 'vitest'
import type { CustomMenuItem } from '@/types'
import {
  getStandaloneUserCustomMenuItems,
  getThirdPartyRechargeCustomMenuItems,
  isThirdPartyRechargeCustomMenuItem,
} from '../rechargeHub'

function makeItem(overrides: Partial<CustomMenuItem>): CustomMenuItem {
  return {
    id: 'item-1',
    label: '自定义页面',
    icon_svg: '',
    url: 'https://example.com/page',
    visibility: 'user',
    sort_order: 1,
    ...overrides,
  }
}

describe('rechargeHub custom menu classification', () => {
  it('identifies third-party recharge and voucher store iframe items', () => {
    const items: CustomMenuItem[] = [
      makeItem({ id: 'faka', label: '购买兑换码', url: 'https://shop.example.com/cards', sort_order: 2 }),
      makeItem({ id: 'topup', label: '三方充值', url: 'https://pay.example.com', sort_order: 1 }),
      makeItem({ id: 'docs', label: '使用文档', url: 'https://docs.example.com', sort_order: 3 }),
      makeItem({ id: 'md-guide', label: '充值指南', url: 'md:recharge-guide', page_slug: 'recharge-guide', sort_order: 4 }),
      makeItem({ id: 'admin-pay', label: '充值后台', url: 'https://pay.example.com', visibility: 'admin', sort_order: 5 }),
      makeItem({ id: 'hidden-shop', label: '发卡页', url: 'https://shop.example.com', hide_in_sidebar: true, sort_order: 6 }),
    ]

    expect(isThirdPartyRechargeCustomMenuItem(items[0])).toBe(true)
    expect(isThirdPartyRechargeCustomMenuItem(items[1])).toBe(true)
    expect(isThirdPartyRechargeCustomMenuItem(items[2])).toBe(false)
    expect(isThirdPartyRechargeCustomMenuItem(items[3])).toBe(false)
    expect(isThirdPartyRechargeCustomMenuItem(items[4])).toBe(false)
    expect(isThirdPartyRechargeCustomMenuItem(items[5])).toBe(false)

    expect(getThirdPartyRechargeCustomMenuItems(items).map((i) => i.id)).toEqual(['topup', 'faka'])
    expect(getStandaloneUserCustomMenuItems(items, true).map((i) => i.id)).toEqual(['docs', 'md-guide'])
    expect(getStandaloneUserCustomMenuItems(items, false).map((i) => i.id)).toEqual(['topup', 'faka', 'docs', 'md-guide'])
  })
})
