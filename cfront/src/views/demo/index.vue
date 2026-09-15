<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/modules/user'

defineOptions({ name: 'Home' })

const userStore = useUserStore()
const { user, displayName } = storeToRefs(userStore)

const genderText = computed(() => ({ 0: '未设置', 1: '男', 2: '女' }[user.value?.gender ?? 0] || '未设置'))
const ageText = computed(() => {
  if (!user.value?.birthday)
    return '未设置'
  const birthday = new Date(user.value.birthday)
  if (Number.isNaN(birthday.getTime()))
    return '未设置'
  const today = new Date()
  let age = today.getFullYear() - birthday.getFullYear()
  const offset = today.getMonth() - birthday.getMonth()
  if (offset < 0 || (offset === 0 && today.getDate() < birthday.getDate()))
    age--
  return age >= 0 ? `${age} 岁` : '未设置'
})
const avatarText = computed(() => displayName.value.trim().slice(0, 1).toUpperCase())
const profileItems = computed(() => [
  { label: '用户名', value: user.value?.username || '未设置', icon: 'contact-o' },
  { label: '年龄', value: ageText.value, icon: 'underway-o' },
  { label: '性别', value: genderText.value, icon: 'friends-o' },
  { label: '体重', value: user.value?.weight ? `${user.value.weight} kg` : '未设置', icon: 'balance-o' },
])

</script>

<template>
  <main class="profile-page">
    <section class="welcome-card">
      <div class="welcome-card__glow" />
      <router-link class="identity" :to="{ name: 'ProfileDetail' }" aria-label="查看个人信息详情">
        <van-image v-if="user?.avatar" round width="76" height="76" fit="cover" :src="user.avatar" alt="用户头像" />
        <div v-else class="avatar-fallback" aria-label="默认头像">{{ avatarText }}</div>
        <div class="identity__text">
          <p>欢迎回来</p>
          <h1>{{ displayName }}</h1>
          <span>健康档案已连接</span>
        </div>
      </router-link>
      <div class="card-actions">
        <router-link class="edit-button primary-action" :to="{ name: 'Medicines' }">
          <van-icon name="medicines-o" />
          我的药品
        </router-link>
        <router-link class="edit-button" :to="{ name: 'ProfileEdit' }">
          <van-icon name="edit" />
          编辑基础信息
        </router-link>
      </div>
    </section>

    <section class="profile-section">
      <div class="section-heading">
        <div><p>PROFILE</p><h2>基础信息</h2></div>
        <span>用于健康记录与用药参考</span>
      </div>
      <div class="info-grid">
        <article v-for="item in profileItems" :key="item.label" class="info-item">
          <div class="info-item__icon"><van-icon :name="item.icon" /></div>
          <div><p>{{ item.label }}</p><strong>{{ item.value }}</strong></div>
        </article>
      </div>
      <article v-if="user?.remarks" class="remarks-card">
        <div class="remarks-card__title"><van-icon name="notes-o" /><span>健康备注</span></div>
        <p>{{ user.remarks }}</p>
      </article>
    </section>
  </main>
</template>

<style lang="less" scoped>
.profile-page { min-height: 100%; padding: 22px 18px 30px; color: var(--app-text); background: var(--app-page-background); }
.welcome-card { position: relative; max-width: 920px; min-height: 190px; margin: 0 auto; padding: 30px; display: flex; align-items: center; justify-content: space-between; gap: 24px; overflow: hidden; border-radius: 26px; color: #fff; background: linear-gradient(135deg, #087f72, #09665e); box-shadow: 0 18px 40px rgba(8, 101, 91, 0.18); }
.welcome-card__glow { position: absolute; width: 230px; height: 230px; right: -80px; top: -120px; border: 34px solid rgba(255, 255, 255, 0.08); border-radius: 50%; }
.identity { position: relative; z-index: 1; padding: 8px; margin: -8px; display: flex; align-items: center; gap: 18px; border-radius: 18px; cursor: pointer; transition: background 180ms ease; }
.identity:hover, .identity:focus-visible { background: rgba(255, 255, 255, 0.1); outline: none; }
.avatar-fallback { width: 76px; height: 76px; flex: 0 0 76px; display: grid; place-items: center; border: 3px solid rgba(255, 255, 255, 0.45); border-radius: 50%; font-size: 28px; font-weight: 800; color: #087568; background: var(--app-accent-soft); }
.identity__text p, .identity__text h1 { margin: 0; }
.identity__text p { margin-bottom: 5px; font-size: 14px; color: rgba(255, 255, 255, 0.7); }
.identity__text h1 { margin-bottom: 9px; font-size: 28px; line-height: 1.2; }
.identity__text span { font-size: 13px; color: #c8fff3; }
.edit-button { position: relative; z-index: 1; min-height: 44px; padding: 0 18px; display: inline-flex; align-items: center; justify-content: center; gap: 7px; border: 1px solid rgba(255, 255, 255, 0.32); border-radius: 22px; font-size: 14px; font-weight: 700; background: rgba(255, 255, 255, 0.13); backdrop-filter: blur(8px); }
.card-actions { position: relative; z-index: 1; display: flex; align-items: center; gap: 10px; }
.primary-action { border-color: #d9fff6; color: #087568; background: #d9fff6; }
.profile-section { max-width: 920px; margin: 20px auto 0; padding: 28px; border: 1px solid var(--app-border); border-radius: 24px; background: var(--app-surface); }
.section-heading { margin-bottom: 22px; display: flex; align-items: flex-end; justify-content: space-between; gap: 16px; }
.section-heading p, .section-heading h2 { margin: 0; }
.section-heading p { margin-bottom: 5px; font-size: 12px; font-weight: 800; letter-spacing: 0.16em; color: #0b8a7b; }
.section-heading h2 { font-size: 22px; color: var(--app-text-strong); }
.section-heading > span { font-size: 13px; color: var(--app-text-muted); }
.info-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.info-item { min-height: 98px; padding: 18px; display: flex; align-items: center; gap: 13px; border: 1px solid var(--app-border); border-radius: 17px; background: var(--app-surface-muted); }
.info-item__icon { width: 38px; height: 38px; flex: 0 0 38px; display: grid; place-items: center; border-radius: 12px; font-size: 20px; color: #087f72; background: var(--app-accent-soft); }
.info-item p, .info-item strong { display: block; margin: 0; overflow-wrap: anywhere; }
.info-item p { margin-bottom: 6px; font-size: 13px; color: var(--app-text-muted); }
.info-item strong { font-size: 15px; color: var(--app-text-strong); }
.remarks-card { margin-top: 12px; padding: 18px; border-radius: 17px; background: var(--app-surface-muted); }
.remarks-card__title { display: flex; align-items: center; gap: 7px; font-size: 14px; font-weight: 700; color: #087f72; }
.remarks-card p { margin: 10px 0 0; font-size: 14px; line-height: 1.7; color: var(--app-text-muted); }

@media (max-width: 680px) {
  .profile-page { padding: 14px 14px 24px; }
  .welcome-card { min-height: 230px; padding: 25px 22px; align-items: flex-start; flex-direction: column; justify-content: center; }
  .card-actions { width: 100%; display: grid; grid-template-columns: 1fr 1fr; }
  .edit-button { width: 100%; padding-inline: 10px; }
  .profile-section { padding: 22px 16px; }
  .section-heading > span { display: none; }
  .info-grid { grid-template-columns: 1fr 1fr; }
  .info-item { padding: 15px 12px; align-items: flex-start; flex-direction: column; }
}
</style>
