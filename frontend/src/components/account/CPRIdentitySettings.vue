<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const enabled = defineModel<boolean>('enabled', { required: true })
const sendTokenName = defineModel<boolean>('sendTokenName', { required: true })
const { t } = useI18n()
</script>

<template>
  <div class="space-y-2">
    <label class="flex items-start gap-3">
      <input
        v-model="enabled"
        type="checkbox"
        role="switch"
        class="mt-1 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-500 dark:bg-dark-700"
        data-testid="cpr-trusted-identity"
      />
      <span>
        <span class="block text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.accounts.openai.cprTrustedIdentity') }}
        </span>
        <span class="input-hint block">{{ t('admin.accounts.openai.cprTrustedIdentityDesc') }}</span>
      </span>
    </label>
    <label v-if="enabled" class="flex items-start gap-3 pl-7 text-sm text-gray-700 dark:text-gray-300">
      <input
        v-model="sendTokenName"
        type="checkbox"
        class="mt-1 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-500 dark:bg-dark-700"
        data-testid="cpr-send-token-name"
      />
      {{ t('admin.accounts.openai.cprSendTokenName') }}
    </label>
  </div>
</template>
