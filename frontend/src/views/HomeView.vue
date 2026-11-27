<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="apple-mist-ambient flex min-h-screen flex-col bg-[#f4f5f8] text-[#111827] dark:bg-dark-950 dark:text-white"
  >
    <header class="glass sticky top-0 z-20 border-b border-slate-900/[0.06] px-4 py-3 sm:px-6 dark:border-white/[0.07]">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-2.5">
          <img
            v-if="siteLogo"
            :src="siteLogo"
            alt="Logo"
            class="h-8 w-8 shrink-0 rounded-lg border border-slate-900/[0.08] bg-white object-contain p-0.5 shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800"
          />
          <span
            v-else
            class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-[#111827] text-white shadow-sm dark:bg-white dark:text-[#0b0c0f]"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none">
              <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
              <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
            </svg>
          </span>
          <span v-if="isDefaultBrand" class="min-w-0 truncate text-[15px] tracking-tight">
            <span class="font-bold text-[#111827] dark:text-white">Thisai</span><span class="font-medium text-primary-500 dark:text-primary-400">.site</span>
          </span>
          <span v-else class="min-w-0 truncate text-[15px] font-semibold tracking-tight">{{ displaySiteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="sm" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-9 shrink-0 items-center gap-1.5 rounded-full px-3 text-xs font-medium text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="sm" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-9 shrink-0 items-center justify-center rounded-full bg-[#111827] px-4 py-2 text-xs font-medium text-white transition-colors hover:bg-[#1f2937] dark:bg-white dark:text-[#0b0c0f] dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          v-if="siteLogo"
          :src="siteLogo"
          alt="Logo"
          class="mx-auto mb-6 h-16 w-16 rounded-2xl border border-slate-900/[0.08] bg-white object-contain p-2 shadow-card dark:border-white/[0.1] dark:bg-dark-800"
        />
        <span
          v-else
          class="mx-auto mb-6 inline-flex h-14 w-14 items-center justify-center rounded-2xl bg-[#111827] text-white shadow-card dark:bg-white dark:text-[#0b0c0f]"
        >
          <svg class="h-7 w-7" viewBox="0 0 24 24" fill="none">
            <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
            <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
          </svg>
        </span>
        <h1 class="[overflow-wrap:anywhere] text-3xl font-semibold tracking-tight md:text-4xl">
          <template v-if="isDefaultBrand">
            <span>Thisai</span><span class="text-primary-500 dark:text-primary-400">.site</span>
          </template>
          <template v-else>{{ displaySiteName }}</template>
        </h1>
        <p class="mt-3 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-[#526071] dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="btn btn-primary mt-8 inline-flex min-h-10 rounded-full px-6 py-2.5 text-sm"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-slate-900/[0.06] px-4 py-5 text-center text-xs text-[#8d99a8] [overflow-wrap:anywhere] sm:px-6 dark:border-white/[0.07] dark:text-dark-400">
      &copy; {{ currentYear }} {{ displaySiteName }}
    </footer>
  </div>

  <!-- Default Home Page (Apple Mist · Thisai.site) -->
  <div
    v-else
    class="apple-mist-ambient relative flex min-h-screen flex-col overflow-hidden bg-[#f4f5f8] text-[#111827] dark:bg-dark-950 dark:text-white"
  >
    <!-- Header -->
    <header class="glass sticky top-0 z-20 border-b border-slate-900/[0.06] px-5 py-3 sm:px-6 dark:border-white/[0.07]">
      <nav class="mx-auto flex max-w-5xl items-center justify-between">
        <!-- Brand Lockup -->
        <div class="flex items-center gap-2.5">
          <div
            v-if="siteLogo"
            class="h-8 w-8 overflow-hidden rounded-lg border border-slate-900/[0.08] bg-white p-0.5 shadow-glass-sm dark:border-white/[0.1] dark:bg-dark-800"
          >
            <img :src="siteLogo" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span
            v-else
            class="inline-flex h-7 w-7 items-center justify-center rounded-lg bg-[#111827] text-white shadow-sm dark:bg-white dark:text-[#0b0c0f]"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none">
              <path d="M5 7h14M12 7v11" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
              <circle cx="17.5" cy="16.5" r="2.2" fill="#0071e3" />
            </svg>
          </span>
          <span v-if="isDefaultBrand" class="text-[15px] tracking-tight">
            <span class="font-bold text-[#111827] dark:text-white">Thisai</span><span class="font-medium text-primary-500 dark:text-primary-400">.site</span>
          </span>
          <span v-else class="text-[15px] font-semibold tracking-tight text-[#111827] dark:text-white">{{ displaySiteName }}</span>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-2">
          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-full p-2 text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="sm" />
          </a>

          <!-- Model Plaza Link -->
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-medium text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="sm" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-full p-2 text-[#526071] transition-colors hover:bg-slate-900/[0.05] hover:text-[#111827] dark:text-dark-300 dark:hover:bg-white/[0.07] dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-[#111827] py-1.5 pl-1.5 pr-3.5 text-xs font-medium text-white transition-colors hover:bg-[#1f2937] dark:bg-white dark:text-[#0b0c0f] dark:hover:bg-gray-200"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-white/15 text-[10px] font-semibold dark:bg-black/10"
            >
              {{ userInitial }}
            </span>
            <span>{{ t('home.dashboard') }}</span>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-[#111827] px-4 py-1.5 text-xs font-medium text-white transition-colors hover:bg-[#1f2937] dark:bg-white dark:text-[#0b0c0f] dark:hover:bg-gray-200"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-5 py-12 sm:px-6 md:py-16">
      <div class="mx-auto max-w-5xl">
        <!-- Centered Editorial Hero -->
        <section class="mx-auto mb-10 max-w-2xl text-center">
          <div class="mb-5 inline-flex items-center gap-2 rounded-full border border-slate-900/[0.08] bg-white px-3.5 py-1.5 text-xs font-medium text-[#526071] shadow-pill dark:border-white/[0.1] dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-emerald-500 ring-2 ring-emerald-500/20"></span>
            <span class="font-mono font-semibold text-[#111827] dark:text-white">{{ domainBadge }}</span>
            <span class="h-2.5 w-px bg-slate-900/[0.12] dark:bg-white/[0.14]"></span>
            <span>{{ t('home.heroKicker') }}</span>
          </div>

          <h1 class="mb-4 text-4xl font-semibold tracking-tight text-[#111827] dark:text-white md:text-5xl">
            <span>{{ t('home.heroTitleLead') }}</span><span class="text-primary-500 dark:text-primary-400">{{ t('home.heroTitleAccent') }}</span><span>。</span>
          </h1>

          <p class="mx-auto mb-7 max-w-xl text-sm leading-relaxed text-[#526071] dark:text-dark-300 md:text-base">
            {{ t('home.heroDescription') }}
          </p>

          <!-- CTA Buttons -->
          <div class="flex flex-wrap items-center justify-center gap-3">
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="btn btn-primary rounded-full px-7 py-2.5 text-sm"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
              <Icon name="arrowRight" size="sm" :stroke-width="2" />
            </router-link>
            <router-link
              v-if="showModelPlazaEntry"
              to="/model-plaza"
              class="btn btn-secondary rounded-full px-6 py-2.5 text-sm"
            >
              {{ t('home.welcome.browseModels') }}
            </router-link>
          </div>
        </section>

        <!-- Core Dual Bento: 01 不降智 + 02 速度快 -->
        <section class="mb-5 grid grid-cols-1 gap-4 md:grid-cols-2">
          <!-- Card 01: Native Fidelity (不降智) -->
          <div class="card flex flex-col justify-between gap-4 p-6">
            <div>
              <div class="flex min-h-[26px] items-center justify-between gap-2">
                <span class="inline-flex items-center rounded-full bg-emerald-500/10 px-2.5 py-0.5 text-[11px] font-semibold text-emerald-700 dark:bg-emerald-400/15 dark:text-emerald-300">
                  {{ t('home.welcome.fidelityTag') }}
                </span>
                <span class="font-mono text-[11px] text-[#8d99a8] dark:text-dark-400">01</span>
              </div>
              <h2 class="mt-3 text-base font-semibold tracking-tight text-[#111827] dark:text-white">
                {{ t('home.welcome.fidelityTitle') }}
              </h2>
              <p class="mt-1.5 text-xs leading-relaxed text-[#526071] dark:text-dark-300">
                {{ t('home.welcome.fidelityDesc') }}
              </p>
            </div>

            <div class="flex flex-1 flex-col justify-between rounded-xl border border-slate-900/[0.05] bg-[#f8f9fb] px-4 py-3 dark:border-white/[0.06] dark:bg-dark-900/70">
              <div class="divide-y divide-slate-900/[0.05] dark:divide-white/[0.06]">
                <div class="flex items-center justify-between gap-3 py-2 first:pt-0.5 text-xs">
                  <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.fidelityModelLabel') }}</span>
                  <span class="font-medium text-[#111827] dark:text-white">
                    <span class="mr-1 text-emerald-600 dark:text-emerald-400">✓</span>{{ t('home.welcome.fidelityModelValue') }}
                  </span>
                </div>
                <div class="flex items-center justify-between gap-3 py-2 text-xs">
                  <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.fidelityThinkingLabel') }}</span>
                  <span class="font-medium text-[#111827] dark:text-white">
                    <span class="mr-1 text-emerald-600 dark:text-emerald-400">✓</span>{{ t('home.welcome.fidelityThinkingValue') }}
                  </span>
                </div>
                <div class="flex items-center justify-between gap-3 py-2 text-xs">
                  <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.fidelityCacheLabel') }}</span>
                  <span class="font-medium text-[#111827] dark:text-white">
                    <span class="mr-1 text-emerald-600 dark:text-emerald-400">✓</span>{{ t('home.welcome.fidelityCacheValue') }}
                  </span>
                </div>
                <div class="flex items-center justify-between gap-3 py-2 text-xs">
                  <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.fidelityToolsLabel') }}</span>
                  <span class="font-medium text-[#111827] dark:text-white">
                    <span class="mr-1 text-emerald-600 dark:text-emerald-400">✓</span>{{ t('home.welcome.fidelityToolsValue') }}
                  </span>
                </div>
              </div>
              <p class="mt-2 border-t border-slate-900/[0.05] pt-2.5 text-[11px] text-[#8d99a8] dark:border-white/[0.06] dark:text-dark-400">
                {{ t('home.welcome.fidelityNote') }}
              </p>
            </div>
          </div>

          <!-- Card 02: Fast Low-Latency Streaming (速度快 · 动态适配后台已开启模型) -->
          <div class="terminal-container card flex flex-col justify-between gap-4 p-6">
            <div>
              <div class="flex min-h-[26px] items-center justify-between gap-2">
                <span class="inline-flex items-center rounded-full bg-primary-500/10 px-2.5 py-0.5 text-[11px] font-semibold text-primary-600 dark:bg-primary-400/15 dark:text-primary-300">
                  {{ t('home.welcome.speedTag') }}
                </span>

                <!-- Only show protocol tabs when 2+ platforms are enabled in backend -->
                <div
                  v-if="speedTabs.length > 1"
                  class="inline-flex items-center gap-0.5 rounded-full bg-[#f1f3f7] p-0.5 dark:bg-dark-900"
                >
                  <button
                    v-for="tab in speedTabs"
                    :key="tab.key"
                    type="button"
                    @click="activeProtoKey = tab.key"
                    class="rounded-full px-2.5 py-0.5 text-[11px] font-medium transition-all"
                    :class="
                      activeProto.key === tab.key
                        ? 'bg-white text-[#111827] shadow-pill dark:bg-dark-700 dark:text-white dark:shadow-none'
                        : 'text-[#526071] hover:text-[#111827] dark:text-dark-300 dark:hover:text-white'
                    "
                  >
                    {{ tab.label }}
                  </button>
                </div>
                <span v-else class="font-mono text-[11px] text-[#8d99a8] dark:text-dark-400">02</span>
              </div>

              <h2 class="mt-3 text-base font-semibold tracking-tight text-[#111827] dark:text-white">
                {{ t('home.welcome.speedTitle') }}
              </h2>
              <p class="mt-1.5 text-xs leading-relaxed text-[#526071] dark:text-dark-300">
                {{ t('home.welcome.speedDesc') }}
              </p>
            </div>

            <div class="flex flex-1 flex-col justify-between gap-3 rounded-xl border border-slate-900/[0.05] bg-[#f8f9fb] px-4 py-3 dark:border-white/[0.06] dark:bg-dark-900/70">
              <div class="flex items-center justify-between border-b border-slate-900/[0.05] pb-2.5 font-mono text-[11.5px] dark:border-white/[0.06]">
                <span class="font-medium text-[#111827] dark:text-white">{{ activeProto.path }}</span>
                <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ t('home.terminal.success') }}</span>
              </div>

              <div class="space-y-3">
                <div class="space-y-1.5">
                  <div class="flex items-center justify-between text-[11.5px]">
                    <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.speedRouteLabel') }}</span>
                    <span class="font-mono font-semibold tabular-nums text-[#111827] dark:text-white">{{ activeProto.route }}</span>
                  </div>
                  <div class="h-1.5 w-full overflow-hidden rounded-full bg-[#e8ebf0] dark:bg-white/[0.08]">
                    <div class="h-full rounded-full bg-primary-500 transition-all duration-300 dark:bg-primary-400" :style="{ width: activeProto.routeWidth }"></div>
                  </div>
                </div>

                <div class="space-y-1.5">
                  <div class="flex items-center justify-between text-[11.5px]">
                    <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.speedTtftLabel') }}</span>
                    <span class="font-mono font-semibold tabular-nums text-[#111827] dark:text-white">{{ activeProto.ttft }}</span>
                  </div>
                  <div class="h-1.5 w-full overflow-hidden rounded-full bg-[#e8ebf0] dark:bg-white/[0.08]">
                    <div class="h-full rounded-full bg-primary-500 transition-all duration-300 dark:bg-primary-400" :style="{ width: activeProto.ttftWidth }"></div>
                  </div>
                </div>

                <div class="space-y-1.5">
                  <div class="flex items-center justify-between text-[11.5px]">
                    <span class="text-[#526071] dark:text-dark-300">{{ t('home.welcome.speedFailoverLabel') }}</span>
                    <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ t('home.welcome.speedFailoverValue') }}</span>
                  </div>
                  <div class="h-1.5 w-full overflow-hidden rounded-full bg-[#e8ebf0] dark:bg-white/[0.08]">
                    <div class="h-full w-[92%] rounded-full bg-emerald-500 dark:bg-emerald-400"></div>
                  </div>
                </div>
              </div>

              <p class="border-t border-slate-900/[0.05] pt-2.5 text-[11px] text-[#8d99a8] dark:border-white/[0.06] dark:text-dark-400">
                <template v-if="enabledPlatforms.length > 0">
                  {{ t('home.welcome.speedChannelPrefix') }}<span class="font-mono text-[#526071] dark:text-dark-200">{{ activeProto.label }}</span> · {{ activeProto.protocolHint }}
                </template>
                <template v-else>
                  {{ t('home.welcome.speedGenericNote') }}
                </template>
              </p>
            </div>
          </div>
        </section>

        <!-- ModelTrace Public Model Attribution Verification Banner -->
        <section class="mb-5">
          <div class="card flex flex-wrap items-center justify-between gap-4 px-5 py-4">
            <div class="flex min-w-0 items-start gap-3 sm:items-center">
              <span class="mt-0.5 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-600 dark:bg-emerald-400/15 dark:text-emerald-300 sm:mt-0">
                <Icon name="shield" size="sm" />
              </span>
              <div class="min-w-0">
                <p class="text-xs font-semibold text-[#111827] dark:text-white sm:text-[13px]">
                  {{ t('home.welcome.modelTraceTitle') }}
                </p>
                <p class="mt-0.5 text-xs text-[#526071] dark:text-dark-300">
                  {{ t('home.welcome.modelTraceDesc') }}
                </p>
              </div>
            </div>
            <a
              href="https://xqy2006.github.io/ModelTrace/"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex shrink-0 items-center gap-1.5 rounded-full border border-emerald-600/20 bg-emerald-500/10 px-4 py-1.5 text-xs font-semibold text-emerald-700 transition-colors hover:bg-emerald-500/15 dark:border-emerald-400/25 dark:bg-emerald-400/15 dark:text-emerald-300 dark:hover:bg-emerald-400/20"
            >
              <span>{{ t('home.welcome.modelTraceAction') }}</span>
              <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M7 17L17 7M9 7h8v8" />
              </svg>
            </a>
          </div>
        </section>

        <!-- One-Line Endpoint Bar -->
        <section class="mb-6">
          <div class="card flex flex-wrap items-center justify-between gap-3 px-5 py-3.5">
            <div class="flex min-w-0 flex-wrap items-center gap-2.5 text-xs text-[#526071] dark:text-dark-300">
              <span class="font-semibold text-[#111827] dark:text-white">{{ t('home.welcome.endpointLabel') }}</span>
              <code class="rounded-lg bg-[#f1f3f7] px-2.5 py-1 font-mono text-xs font-medium text-[#111827] dark:bg-dark-900 dark:text-gray-100">
                {{ apiEndpoint }}
              </code>
              <span class="truncate">{{ compatDescription }}</span>
            </div>
            <button
              type="button"
              @click="copyEndpoint"
              class="shrink-0 rounded-full bg-primary-500/10 px-3.5 py-1.5 text-xs font-semibold text-primary-600 transition-colors hover:bg-primary-500/15 dark:bg-primary-400/15 dark:text-primary-300 dark:hover:bg-primary-400/20"
            >
              {{ copiedEndpoint ? t('home.welcome.copiedEndpoint') : t('home.welcome.copyEndpoint') }}
            </button>
          </div>
        </section>

        <!-- Dynamic Enabled Platforms Strip (never shows disabled Claude/Gemini) -->
        <section class="flex flex-wrap items-center justify-between gap-3 pt-2">
          <div class="flex flex-wrap items-center gap-2">
            <template v-if="enabledPlatforms.length > 0">
              <div
                v-for="p in enabledPlatforms"
                :key="p"
                class="inline-flex items-center gap-2 rounded-full border border-slate-900/[0.08] bg-white px-3.5 py-1.5 text-xs font-medium text-[#111827] shadow-glass-sm dark:border-white/[0.09] dark:bg-dark-800 dark:text-white"
              >
                <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
                <span>{{ platformDisplayName(p) }}</span>
                <span
                  v-if="enabledPlatforms.length === 1"
                  class="text-[11px] text-[#8d99a8] dark:text-dark-400"
                >
                  · {{ t('home.welcome.enabledBadge') }}
                </span>
              </div>
            </template>
            <div
              v-else
              class="inline-flex items-center gap-2 rounded-full border border-slate-900/[0.08] bg-white px-3.5 py-1.5 text-xs font-medium text-[#526071] shadow-glass-sm dark:border-white/[0.09] dark:bg-dark-800 dark:text-dark-200"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
              <span>{{ t('home.welcome.privatePlatformsNote') }}</span>
            </div>
          </div>

          <span class="text-xs text-[#8d99a8] dark:text-dark-400">
            {{ t('home.welcome.transparentBilling') }}
          </span>
        </section>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-slate-900/[0.06] px-6 py-5 dark:border-white/[0.07]">
      <div
        class="mx-auto flex max-w-5xl flex-col items-center justify-between gap-3 text-center sm:flex-row sm:text-left"
      >
        <p class="text-xs text-[#8d99a8] dark:text-dark-400">
          &copy; {{ currentYear }} {{ displaySiteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-5">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-xs text-[#526071] transition-colors hover:text-[#111827] dark:text-dark-300 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="text-xs text-[#526071] transition-colors hover:text-[#111827] dark:text-dark-300 dark:hover:text-white"
          >
            {{ t('nav.modelPlaza') }}
          </router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import { getModelPlaza } from '@/api/modelPlaza'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

interface ProtoPreview {
  key: string
  label: string
  path: string
  route: string
  routeWidth: string
  ttft: string
  ttftWidth: string
  protocolHint: string
}

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const rawSiteName = computed(() => (appStore.cachedPublicSettings?.site_name || appStore.siteName || '').trim())
const isDefaultBrand = computed(() => {
  const n = rawSiteName.value
  return !n || n === 'Sub2API' || n.toLowerCase() === 'thisai.site'
})
const displaySiteName = computed(() => (isDefaultBrand.value ? 'Thisai.site' : rawSiteName.value))
const domainBadge = computed(() => {
  if (!isDefaultBrand.value) return displaySiteName.value
  if (typeof window !== 'undefined' && window.location?.hostname) {
    const host = window.location.hostname
    if (host && host !== 'localhost' && host !== '127.0.0.1') {
      return host
    }
  }
  return 'thisai.site'
})

const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Endpoint bar
const apiEndpoint = computed(() => {
  const configured = (appStore.cachedPublicSettings?.api_base_url || '').trim().replace(/\/+$/, '')
  if (configured) {
    return configured.endsWith('/v1') ? configured : `${configured}/v1`
  }
  if (typeof window !== 'undefined' && window.location?.origin) {
    const origin = window.location.origin
    if (!origin.startsWith('http://localhost') && !origin.startsWith('http://127.0.0.1')) {
      return `${origin}/v1`
    }
  }
  return 'https://thisai.site/v1'
})

const copiedEndpoint = ref(false)
let copyTimer: ReturnType<typeof setTimeout> | null = null

async function copyEndpoint() {
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(apiEndpoint.value)
    }
  } catch {
    // Ignore clipboard permission errors in restricted contexts
  }
  copiedEndpoint.value = true
  if (copyTimer) clearTimeout(copyTimer)
  copyTimer = setTimeout(() => {
    copiedEndpoint.value = false
  }, 1500)
}

// Dynamic enabled platforms loaded from /model-plaza when publicly available
const PLATFORM_ORDER = ['openai', 'anthropic', 'gemini', 'antigravity', 'grok', 'deepseek', 'kimi', 'zhipu', 'minimax']
const PLATFORM_NAMES: Record<string, string> = {
  openai: 'OpenAI',
  anthropic: 'Claude',
  gemini: 'Gemini',
  antigravity: 'Antigravity',
  grok: 'Grok',
  deepseek: 'DeepSeek',
  kimi: 'Kimi',
  zhipu: 'GLM',
  minimax: 'MiniMax'
}

const PROTO_PRESETS: Record<string, ProtoPreview> = {
  openai: {
    key: 'openai',
    label: 'OpenAI',
    path: 'POST /v1/responses',
    route: '~19ms',
    routeWidth: '16%',
    ttft: '~0.5s',
    ttftWidth: '38%',
    protocolHint: 'Responses / Chat'
  },
  anthropic: {
    key: 'anthropic',
    label: 'Claude',
    path: 'POST /v1/messages',
    route: '~22ms',
    routeWidth: '19%',
    ttft: '~0.6s',
    ttftWidth: '44%',
    protocolHint: 'Messages · Thinking & Cache'
  },
  gemini: {
    key: 'gemini',
    label: 'Gemini',
    path: 'POST /v1beta/models',
    route: '~25ms',
    routeWidth: '22%',
    ttft: '~0.5s',
    ttftWidth: '40%',
    protocolHint: 'GenerateContent · Stream'
  },
  antigravity: {
    key: 'antigravity',
    label: 'Antigravity',
    path: 'POST /v1/messages',
    route: '~24ms',
    routeWidth: '20%',
    ttft: '~0.6s',
    ttftWidth: '42%',
    protocolHint: 'Native Passthrough'
  },
  generic: {
    key: 'generic',
    label: 'API',
    path: 'POST /v1/chat/completions',
    route: '~20ms',
    routeWidth: '17%',
    ttft: '~0.5s',
    ttftWidth: '40%',
    protocolHint: 'Native Stream'
  }
}

const enabledPlatforms = ref<string[]>([])
const activeProtoKey = ref<string>('generic')
let plazaAbortController: AbortController | null = null

function platformDisplayName(p: string): string {
  return PLATFORM_NAMES[p] || p.toUpperCase()
}

const speedTabs = computed<ProtoPreview[]>(() => {
  const seen = new Set<string>()
  const list: ProtoPreview[] = []
  for (const p of enabledPlatforms.value) {
    const preset = PROTO_PRESETS[p]
    if (preset && !seen.has(preset.key)) {
      seen.add(preset.key)
      list.push(preset)
    }
    if (list.length >= 3) break
  }
  return list
})

const activeProto = computed<ProtoPreview>(() => {
  if (speedTabs.value.length === 0) return PROTO_PRESETS.generic
  return speedTabs.value.find((item) => item.key === activeProtoKey.value) || speedTabs.value[0]
})

const compatDescription = computed(() => {
  if (enabledPlatforms.value.length === 0) {
    return t('home.welcome.compatGeneric')
  }
  const clients: string[] = []
  if (enabledPlatforms.value.includes('anthropic')) clients.push('Claude Code')
  if (enabledPlatforms.value.includes('openai')) clients.push('Codex')
  clients.push('Cursor')
  return `${t('home.welcome.compatPrefix')}${clients.join('、')}${t('home.welcome.compatSuffix')}`
})

async function loadEnabledPlatforms() {
  if (!showModelPlazaEntry.value || hasHomeContent.value || compactHomeEnabled.value) {
    enabledPlatforms.value = []
    return
  }
  plazaAbortController?.abort()
  plazaAbortController = new AbortController()
  try {
    const res = await getModelPlaza({ signal: plazaAbortController.signal })
    const set = new Set<string>()
    for (const group of res?.groups || []) {
      const p = (group.platform || '').trim().toLowerCase()
      if (p) set.add(p)
    }
    const sorted = Array.from(set).sort((a, b) => {
      const ai = PLATFORM_ORDER.indexOf(a)
      const bi = PLATFORM_ORDER.indexOf(b)
      if (ai === -1 && bi === -1) return a.localeCompare(b)
      if (ai === -1) return 1
      if (bi === -1) return -1
      return ai - bi
    })
    enabledPlatforms.value = sorted
    if (sorted.length > 0 && PROTO_PRESETS[sorted[0]]) {
      activeProtoKey.value = sorted[0]
    }
  } catch {
    // Gracefully stay in neutral generic mode if plaza is disabled or unreachable
  }
}

watch(showModelPlazaEntry, (enabled) => {
  if (enabled && enabledPlatforms.value.length === 0) {
    void loadEnabledPlatforms()
  }
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  if (showModelPlazaEntry.value) {
    void loadEnabledPlatforms()
  }
})

onBeforeUnmount(() => {
  plazaAbortController?.abort()
  if (copyTimer) clearTimeout(copyTimer)
})
</script>

<style scoped>
.terminal-container {
  position: relative;
}
</style>
