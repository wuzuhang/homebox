<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/modules/user'

defineOptions({ name: 'ProfileDetail' })

const userStore = useUserStore()
const { user, displayName } = storeToRefs(userStore)

function formatDate(value?: string | null) {
  if (!value)
    return '未设置'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()))
    return '未设置'
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date)
}

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
const details = computed(() => [
  { label: '用户名', value: user.value?.username || '未设置' },
  { label: '用户别名', value: user.value?.nickname || '未设置' },
  { label: '年龄', value: ageText.value },
  { label: '性别', value: genderText.value },
  { label: '出生日期', value: formatDate(user.value?.birthday) },
  { label: '体重', value: user.value?.weight ? `${user.value.weight} kg` : '未设置' },
  { label: '手机号', value: user.value?.phone || '未设置' },
  { label: '邮箱', value: user.value?.email || '未设置' },
])
</script>

<template>
  <main class="detail-page">
    <section class="profile-card">
      <header class="profile-hero">
        <van-image v-if="user?.avatar" round width="88" height="88" fit="cover" :src="user.avatar" alt="用户头像" />
        <div v-else class="avatar-fallback">{{ avatarText }}</div>
        <div class="profile-hero__text">
          <p>个人健康档案</p>
          <h1>{{ displayName }}</h1>
          <span>@{{ user?.username }}</span>
        </div>
        <router-link class="edit-link" :to="{ name: 'ProfileEdit' }">
          <van-icon name="edit" />
          编辑资料
        </router-link>
      </header>

      <section class="detail-section">
        <div class="section-title">
          <van-icon name="contact-o" />
          <h2>基础信息</h2>
        </div>
        <dl class="detail-list">
          <div v-for="item in details" :key="item.label" class="detail-row">
            <dt>{{ item.label }}</dt>
            <dd>{{ item.value }}</dd>
          </div>
        </dl>
      </section>

      <section class="detail-section health-section">
        <div class="section-title">
          <van-icon name="records-o" />
          <h2>健康信息</h2>
        </div>
        <div class="health-block">
          <p class="health-block__label">疾病标签</p>
          <div v-if="user?.disease_ids?.length" class="tag-list">
            <span v-for="id in user.disease_ids" :key="id">标签 {{ id }}</span>
          </div>
          <p v-else class="empty-value">暂未添加</p>
        </div>
        <div class="health-block">
          <p class="health-block__label">健康备注</p>
          <p class="remarks">{{ user?.remarks || '暂未填写健康备注' }}</p>
        </div>
      </section>
    </section>
  </main>
</template>

<style lang="less" scoped>
.detail-page { min-height: 100%; padding: 22px 18px 34px; color: var(--app-text); background: var(--app-page-background); }
.profile-card { max-width: 820px; margin: 0 auto; overflow: hidden; border: 1px solid var(--app-border); border-radius: 26px; background: var(--app-surface); box-shadow: 0 16px 42px var(--app-shadow); }
.profile-hero { position: relative; min-height: 190px; padding: 34px; display: flex; align-items: center; gap: 20px; color: #fff; background: linear-gradient(135deg, #087f72, #09665e); }
.profile-hero::after { content: ''; position: absolute; width: 210px; height: 210px; right: -90px; top: -110px; border: 32px solid rgba(255, 255, 255, 0.07); border-radius: 50%; }
.avatar-fallback { width: 88px; height: 88px; flex: 0 0 88px; display: grid; place-items: center; border: 3px solid rgba(255, 255, 255, 0.45); border-radius: 50%; font-size: 32px; font-weight: 800; color: #087568; background: var(--app-accent-soft); }
.profile-hero__text { position: relative; z-index: 1; flex: 1; min-width: 0; }
.profile-hero__text p, .profile-hero__text h1 { margin: 0; }
.profile-hero__text p { margin-bottom: 5px; font-size: 14px; color: rgba(255, 255, 255, 0.7); }
.profile-hero__text h1 { margin-bottom: 8px; font-size: 30px; overflow-wrap: anywhere; }
.profile-hero__text span { font-size: 14px; color: #c8fff3; }
.edit-link { position: relative; z-index: 1; min-height: 42px; padding: 0 16px; display: inline-flex; align-items: center; justify-content: center; gap: 7px; border: 1px solid rgba(255, 255, 255, 0.32); border-radius: 22px; font-size: 14px; font-weight: 700; background: rgba(255, 255, 255, 0.12); }
.detail-section { padding: 28px 32px; }
.detail-section + .detail-section { border-top: 1px solid var(--app-border-soft); }
.section-title { margin-bottom: 18px; display: flex; align-items: center; gap: 9px; color: #087f72; }
.section-title :deep(.van-icon) { font-size: 20px; }
.section-title h2 { margin: 0; font-size: 18px; color: var(--app-text-strong); }
.detail-list { margin: 0; display: grid; grid-template-columns: 1fr 1fr; gap: 0 34px; }
.detail-row { min-height: 52px; display: flex; align-items: center; justify-content: space-between; gap: 16px; border-bottom: 1px solid var(--app-border-soft); }
.detail-row dt { flex: 0 0 auto; font-size: 14px; color: var(--app-text-muted); }
.detail-row dd { margin: 0; text-align: right; font-size: 14px; font-weight: 700; color: var(--app-text-strong); overflow-wrap: anywhere; }
.health-section { display: grid; grid-template-columns: 1fr 1fr; gap: 18px 28px; }
.health-section .section-title { grid-column: 1 / -1; margin-bottom: 0; }
.health-block { padding: 17px; border-radius: 15px; background: var(--app-surface-muted); }
.health-block__label { margin: 0 0 10px; font-size: 13px; font-weight: 700; color: var(--app-text-muted); }
.tag-list { display: flex; flex-wrap: wrap; gap: 7px; }
.tag-list span { padding: 6px 10px; border-radius: 16px; font-size: 13px; color: #0a9683; background: var(--app-accent-soft); }
.empty-value, .remarks { margin: 0; font-size: 14px; line-height: 1.7; color: var(--app-text-muted); }
.empty-value { color: var(--app-text-muted); }

@media (max-width: 680px) {
  .detail-page { padding: 14px 12px 24px; }
  .profile-card { border-radius: 20px; }
  .profile-hero { padding: 28px 22px; align-items: flex-start; flex-wrap: wrap; }
  .profile-hero__text { padding-top: 8px; }
  .edit-link { width: 100%; }
  .detail-section { padding: 24px 20px; }
  .detail-list, .health-section { grid-template-columns: 1fr; }
  .health-section .section-title { grid-column: auto; }
}
</style>
