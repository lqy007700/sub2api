import { afterEach, expect, it, vi } from 'vitest'

afterEach(() => {
  localStorage.removeItem('sub2api_locale')
  document.documentElement.removeAttribute('lang')
  vi.restoreAllMocks()
})

it.each([null, 'en'])('starts in Chinese on an English browser with saved locale %s', async (savedLocale) => {
  vi.resetModules()
  vi.spyOn(navigator, 'language', 'get').mockReturnValue('en-US')
  if (savedLocale !== null) localStorage.setItem('sub2api_locale', savedLocale)

  const { initI18n, getLocale, i18n } = await import('@/i18n')
  await initI18n()

  expect(getLocale()).toBe('zh')
  expect(document.documentElement.lang).toBe('zh')
  expect(i18n.global.t('nav.profile')).toBe('个人资料')
})
