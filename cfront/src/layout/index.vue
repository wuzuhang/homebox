<script setup lang="ts">
import NavBar from '@/components/nav-bar/index.vue'
import { useCachedViewStore } from '@/store/modules/cached-view'
import { useDarkModeStore } from '@/store/modules/dark-mode'

const cachedViewStore = useCachedViewStore()
const { cachedViewList } = storeToRefs(cachedViewStore)

const darkModeStore = useDarkModeStore()
const { theme } = storeToRefs(darkModeStore)
const route = useRoute()
const appContentRef = ref<HTMLElement | null>(null)

watch(
  () => route.fullPath,
  async () => {
    await nextTick()
    if (!appContentRef.value)
      return
    appContentRef.value.scrollTop = 0
    appContentRef.value.scrollLeft = 0
  },
  { immediate: true, flush: 'post' },
)
</script>

<template>
  <div class="app-wrapper">
    <van-config-provider class="app-shell" :theme="theme">
      <NavBar />
      <div ref="appContentRef" class="app-content">
        <router-view v-slot="{ Component }">
          <keep-alive :include="cachedViewList">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </div>
    </van-config-provider>
  </div>
</template>

<style lang="less" scoped>
@import "@/styles/mixin.less";

.app-wrapper {
  .clearfix();
  position: relative;
  height: 100%;
  height: 100dvh;
  width: 100%;
  overflow: hidden;
  background-color: var(--app-page-background);
}

.app-shell {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.app-content {
  flex: 1;
  min-height: 0;
  overflow-x: hidden;
  overflow-y: auto;
  overscroll-behavior-y: contain;
  -webkit-overflow-scrolling: touch;
}
</style>
