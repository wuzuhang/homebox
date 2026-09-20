<script setup lang="ts">
import { storeToRefs } from 'pinia'
import type { UploaderAfterRead } from 'vant'
import { showFailToast, showSuccessToast } from 'vant'
import { updateProfile } from '@/api/profile'
import { uploadImage } from '@/api/upload'
import { useMedicineStore } from '@/store/modules/medicine'
import { useUserStore } from '@/store/modules/user'
import 'vant/es/toast/style'

defineOptions({ name: 'ProfileEdit' })

const router = useRouter()
const userStore = useUserStore()
const medicineStore = useMedicineStore()
const { diseases } = storeToRefs(medicineStore)
const submitting = ref(false)
const avatarUploading = ref(false)
const showPassword = ref(false)

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
  disease_ids: [...(userStore.user?.disease_ids || [])],
  remarks: userStore.user?.remarks || '',
})

function getErrorMessage(error: unknown) {
  if (error && typeof error === 'object') {
    const payload = error as { message?: string, msg?: string, error?: string }
    return payload.message || payload.msg || payload.error || '更新失败，请稍后重试'
  }
  return '更新失败，请稍后重试'
}

function createProfilePayload() {
  const userId = userStore.user?.ID
  if (!userId)
    throw new Error('未获取到当前用户 ID，请重新登录')

  return {
    id: userId,
    ...form,
    weight: form.weight ? Number(form.weight) : 0,
    disease_ids: form.disease_ids,
  }
}

const uploadAvatar: UploaderAfterRead = async (items) => {
  const item = Array.isArray(items) ? items[0] : items
  const file = item.file
  if (!file)
    return

  if (!file.type.startsWith('image/')) {
    showFailToast('请选择图片文件')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    showFailToast('图片大小不能超过 5MB')
    return
  }

  avatarUploading.value = true
  item.status = 'uploading'
  item.message = '上传中…'
  try {
    const result = await uploadImage(file)
    if (!result?.fileURL)
      throw new Error('上传接口未返回图片地址')
    form.avatar = result.fileURL
    item.status = 'done'
    showSuccessToast('头像上传成功，提交后保存')
  }
  catch (error) {
    item.status = 'failed'
    item.message = '上传失败'
    showFailToast(getErrorMessage(error))
  }
  finally {
    avatarUploading.value = false
  }
}

async function submitProfile() {
  submitting.value = true
  try {
    const updatedUser = await updateProfile(createProfilePayload())
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

onMounted(async () => {
  try {
    await medicineStore.loadDiseases()
  }
  catch (error) {
    showFailToast(getErrorMessage(error))
  }
})
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
      </header>

      <van-form class="profile-form" @submit="submitProfile">
        <div class="field-grid">
          <van-field v-model="form.username" name="username" label="用户名" readonly />
          <van-field v-model.trim="form.nickname" name="nickname" label="用户别名" placeholder="请输入用户别名" />
          <van-field v-model.trim="form.phone" name="phone" label="手机号" type="tel" placeholder="请输入手机号" autocomplete="tel" />
          <van-field v-model.trim="form.email" name="email" label="邮箱" type="email" placeholder="请输入邮箱" autocomplete="email" />
          <van-field name="avatar" label="头像">
            <template #input>
              <div class="avatar-field">
                <van-image v-if="form.avatar" round width="46" height="46" fit="cover" :src="form.avatar" alt="头像预览" />
                <div v-else class="avatar-preview">{{ (form.nickname || form.username || '家').slice(0, 1) }}</div>
                <van-uploader accept="image/*" :after-read="uploadAvatar" :disabled="avatarUploading || submitting" :preview-image="false">
                  <van-button size="small" round plain type="primary" native-type="button" :loading="avatarUploading">
                    {{ form.avatar ? '更换头像' : '上传头像' }}
                  </van-button>
                </van-uploader>
              </div>
            </template>
          </van-field>
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
          <van-field class="full-field disease-field" name="disease_ids" label="疾病标签">
            <template #input>
              <van-checkbox-group v-if="diseases.length" v-model="form.disease_ids" direction="horizontal">
                <van-checkbox v-for="disease in diseases" :key="disease.id" :name="disease.id" shape="square">{{ disease.name }}</van-checkbox>
              </van-checkbox-group>
              <span v-else class="empty-hint">暂无疾病标签</span>
            </template>
          </van-field>
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
          <van-button block round type="primary" native-type="submit" :loading="submitting" :disabled="avatarUploading" loading-text="正在保存…">保存修改</van-button>
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
.avatar-preview { width: 46px; height: 46px; flex: 0 0 46px; display: grid; place-items: center; border-radius: 50%; font-size: 18px; font-weight: 800; color: #087568; background: var(--app-accent-soft); }
.avatar-field { width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.profile-form { --van-field-label-color: var(--app-text-strong); --van-field-input-text-color: var(--app-text); --van-field-label-width: 5.5em; }
.field-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 14px; }
.profile-form :deep(.van-cell) { margin-bottom: 14px; padding: 11px 14px; align-items: center; border: 1px solid var(--app-border); border-radius: 13px; background: var(--app-surface-muted); }
.profile-form :deep(.van-cell::after) { display: none; }
.profile-form :deep(.van-field__label), .profile-form :deep(.van-field__control), .profile-form :deep(.van-radio__label) { font-size: 14px; }
.profile-form :deep(.van-field--error) { border-color: #df6c6c; }
.full-field { grid-column: 1 / -1; }
.disease-field :deep(.van-checkbox-group) { display: flex; flex-wrap: wrap; gap: 12px 16px; }
.empty-hint { font-size: 14px; color: var(--app-text-muted); }
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
