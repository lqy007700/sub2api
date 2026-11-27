<template>
  <div class="apple-mist-ambient min-h-screen bg-[#f4f5f8] dark:bg-dark-950">
    <!-- Sidebar (always mounted for admin & mobile drawer; hidden on desktop for regular user top-nav) -->
    <AppSidebar :desktop-hidden="useTopNavMode" />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="[
        useTopNavMode
          ? 'lg:ml-0'
          : sidebarCollapsed
            ? 'lg:ml-[72px]'
            : 'lg:ml-[232px]'
      ]"
    >
      <!-- Header -->
      <AppHeader />

      <!-- Main Content -->
      <main
        class="px-4 py-5 md:px-6 md:py-6 lg:px-8"
        :class="useTopNavMode ? 'mx-auto w-full max-w-[1200px]' : ''"
      >
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')
const useTopNavMode = computed(() => !isAdmin.value && !appStore.backendModeEnabled)

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
})

defineExpose({ replayTour })
</script>
