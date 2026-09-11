<script setup lang="ts">
import { showFailToast, showSuccessToast } from 'vant'
import { updateProfile } from '@/api/profile'
import { useUserStore } from '@/store/modules/user'
import 'vant/es/toast/style'

defineOptions({ name: 'ProfileEdit' })

const router = useRouter()
const userStore = useUserStore()
const submitting = ref(false)
const showPassword = ref(false)
const diseaseIdsText = ref(userStore.user?.disease_ids?.join(', ') || '')

const form = reactive({
  username: userStore.user?.username || '',
  password: '',
  phone: userStore.user?.phone || '',
  email: userStore.user?.email || '',
  nickname: userStore.user?.nickname || '',
  avatar: userStore.user?.avatar || '',
  gender: userStore.user?.gender ?? 0,
  birthday: userStore.user?.birthday?.slice(0, 10) || '',
  weight: userStore.user?.weight ? String(userStore.user.weight) : '',
  remarks: userStore.user?.remarks || '',
})

function getErrorMessage(error: unknown) {
  if (error && typeof error === 'object') {
    const payload = error as { message?: string, msg?: string, error?: string }
    return payload.message || payload.msg || payload.error || '更新失败，请稍后重试'
  }
  return '更新失败，请稍后重试'
}

async function submitProfile() {
  submitting.value = true
  try {
    const diseaseIds = diseaseIdsText.value
      .split(/[，,\s]+/)
      .map(item => Number(item))
      .filter(item => Number.isInteger(item) && item > 0)

    const updatedUser = await updateProfile({
      ...form,
      weight: form.weight ? Number(form.weight) : 0,
      disease_ids: diseaseIds,
    })
    userStore.setUser(updatedUser)
    showSuccessToast('更新成功')
    await router.replace({ name: 'Home' })
  }
  catch (error) {
    showFailToast(getErrorMessage(error))
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="edit-page">
    <section class="edit-card">
      <header class="edit-heading">
        <div>
          <p>PERSONAL PROFILE</p>
          <h1>编辑基础信息</h1>
          <span>保存后，首页将同步显示最新资料。</span>
        </div>
        <van-image v-if="form.avatar" round width="64" height="64" fit="cover" :src="form.avatar" alt="头像预览" />
        <div v-else class="avatar-preview">{{ (form.nickname || form.username || '家').slice(0, 1) }}</div>
      </header>

      <van-form class="profile-form" @submit="submitProfile">
        <div class="field-grid">
          <van-field v-model="form.username" name="username" label="用户名" readonly />
          <van-field v-model.trim="form.nickname" name="nickname" label="用户别名" placeholder="请输入用户别名" />
          <van-field v-model.trim="form.phone" name="phone" label="手机号" type="tel" placeholder="请输入手机号" autocomplete="tel" />
          <van-field v-model.trim="form.email" name="email" label="邮箱" type="email" placeholder="请输入邮箱" autocomplete="email" />
          <van-field v-model.trim="form.avatar" name="avatar" label="头像地址" type="url" placeholder="https://" />
          <van-field v-model="form.birthday" name="birthday" label="出生日期" type="date" />
          <van-field name="gender" label="性别">
            <template #input>
              <van-radio-group v-model="form.gender" direction="horizontal">
                <van-radio :name="0">未设置</van-radio>
                <van-radio :name="1">男</van-radio>
                <van-radio :name="2">女</van-radio>
              </van-radio-group>
            </template>
          </van-field>
          <van-field v-model="form.weight" name="weight" label="体重" type="number" placeholder="kg" :min="0" step="0.1" />
          <van-field v-model.trim="diseaseIdsText" name="disease_ids" label="疾病标签" placeholder="ID 用逗号分隔" />
          <van-field
            v-model="form.password"
            name="password"
            label="当前密码"
            placeholder="用于确认本次修改"
            autocomplete="current-password"
            :type="showPassword ? 'text' : 'password'"
            :right-icon="showPassword ? 'closed-eye' : 'eye-o'"
            :rules="[
              { required: true, message: '请输入当前密码' },
              { pattern: /^.{6,}$/, message: '密码至少 6 个字符' },
            ]"
            @click-right-icon="showPassword = !showPassword"
          />
          <van-field v-model.trim="form.remarks" class="full-field" name="remarks" label="备注" type="textarea" rows="3" autosize maxlength="200" show-word-limit placeholder="需要特别留意的健康信息" />
        </div>

        <div class="form-actions">
          <van-button block round plain type="primary" native-type="button" @click="router.back()">取消</van-button>
          <van-button block round type="primary" native-type="submit" :loading="submitting" loading-text="正在保存…">保存修改</van-button>
        </div>
      </van-form>
    </section>
  </main>
</template>

<style lang="less" scoped>
.edit-page { min-height: 100%; padding: 22px 18px 32px; color: var(--app-text); background: var(--app-page-background); }
.edit-card { max-width: 820px; margin: 0 auto; padding: 30px; border: 1px solid var(--app-border); border-radius: 24px; background: var(--app-surface); box-shadow: 0 16px 42px var(--app-shadow); }
.edit-heading { margin-bottom: 28px; display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.edit-heading p, .edit-heading h1, .edit-heading span { margin: 0; }
.edit-heading p { margin-bottom: 6px; font-size: 12px; font-weight: 800; letter-spacing: 0.16em; color: #0b8a7b; }
.edit-heading h1 { margin-bottom: 8px; font-size: 28px; }
.edit-heading span { font-size: 14px; color: var(--app-text-muted); }
.avatar-preview { width: 64px; height: 64px; flex: 0 0 64px; display: grid; place-items: center; border-radius: 50%; font-size: 24px; font-weight: 800; color: #087568; background: var(--app-accent-soft); }
.profile-form { --van-field-label-color: var(--app-text-strong); --van-field-input-text-color: var(--app-text); --van-field-label-width: 5.5em; }
.field-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 14px; }
.profile-form :deep(.van-cell) { margin-bottom: 14px; padding: 11px 14px; align-items: center; border: 1px solid var(--app-border); border-radius: 13px; background: var(--app-surface-muted); }
.profile-form :deep(.van-cell::after) { display: none; }
.profile-form :deep(.van-field__label), .profile-form :deep(.van-field__control), .profile-form :deep(.van-radio__label) { font-size: 14px; }
.profile-form :deep(.van-field--error) { border-color: #df6c6c; }
.full-field { grid-column: 1 / -1; }
.form-actions { margin-top: 10px; display: grid; grid-template-columns: 1fr 1.4fr; gap: 12px; }
.form-actions :deep(.van-button) { height: 48px; font-size: 15px; font-weight: 700; }
.form-actions :deep(.van-button--normal:not(.van-button--plain)) { border: 0; background: linear-gradient(115deg, #087f72, #0a9683); box-shadow: 0 10px 24px rgba(8, 127, 114, 0.2); }

@media (max-width: 680px) {
  .edit-page { padding: 14px 12px 24px; }
  .edit-card { padding: 22px 14px; border-radius: 20px; }
  .edit-heading { padding: 0 4px; }
  .edit-heading h1 { font-size: 24px; }
  .field-grid { grid-template-columns: 1fr; }
  .full-field { grid-column: auto; }
  .form-actions { grid-template-columns: 1fr 1.35fr; }
}
</style>
