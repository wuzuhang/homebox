<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useDarkModeStore } from '@/store/modules/dark-mode'

const darkModeStore = useDarkModeStore()
const { isDark } = storeToRefs(darkModeStore)
const route = useRoute()
const router = useRouter()

const iconName = computed(() => isDark.value ? 'light' : 'dark')
const showBack = computed(() => [
  'ProfileDetail',
  'ProfileEdit',
  'Medicines',
  'MedicineCreate',
  'MedicineEdit',
  'MedicineDetail',
].includes(String(route.name)))

function onClickRight(event: TouchEvent | MouseEvent) {
  darkModeStore.toggleDarkMode(event)
}
</script>

<template>
  <van-nav-bar :title="route.meta.title" :left-arrow="showBack" fixed placeholder @click-left="router.back()" @click-right="onClickRight">
    <template #right>
      <svg-icon class="text-[18px]" :name="iconName" />
    </template>
  </van-nav-bar>
</template>

<style scoped></style>
