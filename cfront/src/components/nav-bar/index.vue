<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useDarkModeStore } from '@/store/modules/dark-mode'
import { useUserStore } from '@/store/modules/user'

const darkModeStore = useDarkModeStore()
const userStore = useUserStore()
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

async function logout() {
  userStore.clearSession()
  await router.replace({ name: 'Auth', query: { mode: 'login' } })
}

function onClickLeft() {
  const routeName = String(route.name)

  if (routeName === 'ProfileDetail' || routeName === 'Medicines') {
    void router.replace({ name: 'Home' })
    return
  }

  if (routeName === 'ProfileEdit') {
    void router.replace({ name: 'ProfileDetail' })
    return
  }

  if (routeName === 'MedicineDetail' || routeName === 'MedicineCreate') {
    void router.replace({ name: 'Medicines' })
    return
  }

  if (routeName === 'MedicineEdit') {
    const id = Number(route.params.id)
    void router.replace(id
      ? { name: 'MedicineDetail', params: { id } }
      : { name: 'Medicines' })
  }
}
</script>

<template>
  <van-nav-bar :title="route.meta.title" :left-arrow="showBack" :z-index="100" fixed placeholder @click-left="onClickLeft">
    <template #right>
      <div class="nav-actions">
        <button type="button" aria-label="切换黑白主题" @click="onClickRight">
          <svg-icon :name="iconName" />
        </button>
        <button type="button" aria-label="退出登录" @click="logout">
          <svg-icon name="logout" />
        </button>
      </div>
    </template>
  </van-nav-bar>
</template>

<style scoped>
.nav-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-actions button {
  width: 32px;
  height: 32px;
  padding: 0;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 50%;
  font-size: 18px;
  color: var(--van-nav-bar-icon-color);
  background: transparent;
  cursor: pointer;
}

.nav-actions button:active {
  background: var(--app-surface-muted);
}
</style>
