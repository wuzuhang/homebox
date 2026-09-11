<script setup lang="ts">
import { showFailToast, showSuccessToast } from 'vant'
import { login, register } from '@/api/auth'
import { useUserStore } from '@/store/modules/user'
import 'vant/es/toast/style'

defineOptions({ name: 'Auth' })

type AuthMode = 'login' | 'register'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const mode = ref<AuthMode>(route.query.mode === 'register' ? 'register' : 'login')
const submitting = ref(false)
const showLoginPassword = ref(false)
const showRegisterPassword = ref(false)
const diseaseIdsText = ref('')

const loginForm = reactive({ username: '', password: '' })
const registerForm = reactive({
  username: '',
  password: '',
  phone: '',
  email: '',
  nickname: '',
  avatar: '',
  gender: 0,
  birthday: '',
  weight: '',
  remarks: '',
})

function switchMode(nextMode: AuthMode) {
  mode.value = nextMode
  router.replace({ query: nextMode === 'register' ? { mode: 'register' } : {} })
}

function getErrorMessage(error: unknown, fallback: string) {
  if (error && typeof error === 'object') {
    const payload = error as { message?: string, msg?: string, error?: string }
    return payload.message || payload.msg || payload.error || fallback
  }
  return fallback
}

async function submitLogin() {
  submitting.value = true
  try {
    const result = await login(loginForm)
    if (!result?.token || !result?.data)
      throw new Error('登录响应缺少用户信息')
    userStore.setSession(result.token, result.data)
    showSuccessToast('登录成功')
    await router.replace({ name: 'Home' })
  }
  catch (error) {
    showFailToast(getErrorMessage(error, '登录失败，请检查账号和密码'))
  }
  finally {
    submitting.value = false
  }
}

async function submitRegister() {
  submitting.value = true
  try {
    const diseaseIds = diseaseIdsText.value
      .split(/[，,\s]+/)
      .map(item => Number(item))
      .filter(item => Number.isInteger(item) && item > 0)

    await register({
      ...registerForm,
      weight: registerForm.weight ? Number(registerForm.weight) : 0,
      disease_ids: diseaseIds,
    })
    loginForm.username = registerForm.username
    loginForm.password = ''
    showSuccessToast({
      message: '注册成功，3 秒后返回登录',
      duration: 3000,
    })
    await new Promise(resolve => setTimeout(resolve, 3000))
    switchMode('login')
  }
  catch (error) {
    showFailToast(getErrorMessage(error, '注册失败，请稍后重试'))
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <section class="auth-shell" aria-label="家庭健康账号">
      <aside class="brand-panel">
        <div class="brand-mark" aria-hidden="true">♥</div>
        <div>
          <p class="brand-eyebrow">HOME BOX</p>
          <h1>把家人的健康，<br>稳稳放在心上。</h1>
          <p class="brand-copy">集中管理家庭成员的健康档案，让每一次记录都有迹可循。</p>
        </div>
        <div class="privacy-note">
          <span class="privacy-note__icon">✓</span>
          <span>你的健康资料将被妥善保护</span>
        </div>
      </aside>

      <section class="form-panel">
        <div class="mobile-brand">
          <div class="mobile-brand__mark">♥</div>
          <span>Home Box</span>
        </div>

        <div class="auth-heading">
          <p>{{ mode === 'login' ? '欢迎回来' : '创建家庭健康账号' }}</p>
          <h2>{{ mode === 'login' ? '登录' : '注册' }}</h2>
        </div>

        <div class="mode-switch" role="tablist" aria-label="选择登录或注册">
          <button type="button" role="tab" :aria-selected="mode === 'login'" :class="{ active: mode === 'login' }" @click="switchMode('login')">
            登录
          </button>
          <button type="button" role="tab" :aria-selected="mode === 'register'" :class="{ active: mode === 'register' }" @click="switchMode('register')">
            注册
          </button>
        </div>

        <van-form v-if="mode === 'login'" class="auth-form" @submit="submitLogin">
          <van-field
            v-model.trim="loginForm.username"
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            autocomplete="username"
            :rules="[{ required: true, message: '请输入用户名' }]"
          />
          <van-field
            v-model="loginForm.password"
            name="password"
            label="密码"
            placeholder="请输入密码"
            autocomplete="current-password"
            :type="showLoginPassword ? 'text' : 'password'"
            :right-icon="showLoginPassword ? 'closed-eye' : 'eye-o'"
            :rules="[{ required: true, message: '请输入密码' }]"
            @click-right-icon="showLoginPassword = !showLoginPassword"
          />
          <van-button class="submit-button" block round type="primary" native-type="submit" :loading="submitting" loading-text="正在登录…">
            登录
          </van-button>
          <p class="form-footnote">
            还没有账号？<button type="button" @click="switchMode('register')">立即注册</button>
          </p>
        </van-form>

        <van-form v-else class="auth-form register-form" @submit="submitRegister">
          <div class="field-grid">
            <van-field
              v-model.trim="registerForm.username"
              name="username"
              label="用户名"
              placeholder="至少 3 个字符"
              autocomplete="username"
              :rules="[
                { required: true, message: '请输入用户名' },
                { pattern: /^.{3,}$/, message: '用户名至少 3 个字符' },
              ]"
            />
            <van-field
              v-model="registerForm.password"
              name="password"
              label="密码"
              placeholder="至少 6 个字符"
              autocomplete="new-password"
              :type="showRegisterPassword ? 'text' : 'password'"
              :right-icon="showRegisterPassword ? 'closed-eye' : 'eye-o'"
              :rules="[
                { required: true, message: '请输入密码' },
                { pattern: /^.{6,}$/, message: '密码至少 6 个字符' },
              ]"
              @click-right-icon="showRegisterPassword = !showRegisterPassword"
            />
            <van-field v-model.trim="registerForm.phone" name="phone" label="手机号" type="tel" placeholder="选填" autocomplete="tel" />
            <van-field v-model.trim="registerForm.email" name="email" label="邮箱" type="email" placeholder="选填" autocomplete="email" />
          </div>

          <details class="profile-details">
            <summary>
              <span>
                <strong>完善健康资料</strong>
                <small>选填，可稍后补充</small>
              </span>
              <span class="summary-icon" aria-hidden="true">＋</span>
            </summary>
            <div class="field-grid profile-fields">
              <van-field v-model.trim="registerForm.nickname" name="nickname" label="昵称" placeholder="怎么称呼你" />
              <van-field v-model.trim="registerForm.avatar" name="avatar" label="头像地址" type="url" placeholder="https://" />
              <van-field name="gender" label="性别">
                <template #input>
                  <van-radio-group v-model="registerForm.gender" direction="horizontal">
                    <van-radio :name="0">未设置</van-radio>
                    <van-radio :name="1">男</van-radio>
                    <van-radio :name="2">女</van-radio>
                  </van-radio-group>
                </template>
              </van-field>
              <van-field v-model="registerForm.birthday" name="birthday" label="出生日期" type="date" />
              <van-field v-model="registerForm.weight" name="weight" label="体重" type="number" placeholder="kg" :min="0" step="0.1" />
              <van-field v-model.trim="diseaseIdsText" name="disease_ids" label="疾病标签" placeholder="ID 用逗号分隔" />
              <van-field v-model.trim="registerForm.remarks" class="remarks-field" name="remarks" label="备注" type="textarea" rows="2" autosize maxlength="200" show-word-limit placeholder="需要特别留意的健康信息" />
            </div>
          </details>

          <van-button class="submit-button" block round type="primary" native-type="submit" :loading="submitting" loading-text="正在注册…">
            创建账号
          </van-button>
          <p class="form-footnote">
            已有账号？<button type="button" @click="switchMode('login')">返回登录</button>
          </p>
        </van-form>
      </section>
    </section>
  </main>
</template>

<style lang="less" scoped>
.auth-page {
  --auth-primary: #087f72;
  --auth-primary-dark: #05655b;
  --auth-ink: var(--app-text);
  --auth-muted: var(--app-text-muted);
  min-height: 100%;
  padding: 32px;
  display: grid;
  place-items: center;
  background: var(--app-auth-background);
}

.auth-shell {
  width: min(1040px, 100%);
  min-height: 640px;
  display: grid;
  grid-template-columns: 0.92fr 1.08fr;
  overflow: hidden;
  border: 1px solid var(--app-border);
  border-radius: 32px;
  background: var(--app-surface);
  box-shadow: 0 28px 80px var(--app-shadow);
}

.brand-panel {
  position: relative;
  padding: 58px 52px 46px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
  color: #fff;
  background: linear-gradient(155deg, #0c8c7d 0%, #087568 58%, #075e59 100%);
}

.brand-panel::after {
  content: '';
  position: absolute;
  width: 340px;
  height: 340px;
  right: -170px;
  bottom: -85px;
  border: 42px solid rgba(255, 255, 255, 0.07);
  border-radius: 50%;
}

.brand-mark {
  width: 64px;
  height: 64px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.34);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.13);
  backdrop-filter: blur(8px);
  font-size: 30px;
  color: #d8fff6;
}

.brand-eyebrow {
  margin: 0 0 18px;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.24em;
  color: rgba(255, 255, 255, 0.66);
}

.brand-panel h1 {
  margin: 0;
  font-size: clamp(34px, 3.5vw, 48px);
  line-height: 1.34;
  letter-spacing: -0.04em;
}

.brand-copy {
  max-width: 340px;
  margin: 22px 0 0;
  font-size: 16px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.76);
}

.privacy-note {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.72);
}

.privacy-note__icon {
  width: 26px;
  height: 26px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.14);
  color: #dcfff8;
}

.form-panel {
  padding: 50px clamp(34px, 6vw, 76px) 36px;
  align-self: center;
}

.mobile-brand {
  display: none;
}

.auth-heading p {
  margin: 0 0 8px;
  font-size: 15px;
  color: var(--auth-muted);
}

.auth-heading h2 {
  margin: 0;
  font-size: 34px;
  line-height: 1.2;
  letter-spacing: -0.04em;
  color: var(--auth-ink);
}

.mode-switch {
  margin: 30px 0 24px;
  padding: 5px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 5px;
  border-radius: 14px;
  background: var(--app-surface-muted);
}

.mode-switch button {
  min-height: 42px;
  border: 0;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  color: var(--app-text-muted);
  background: transparent;
  cursor: pointer;
  transition: 180ms ease;
}

.mode-switch button.active {
  color: var(--auth-primary-dark);
  background: var(--app-surface);
  box-shadow: 0 4px 14px var(--app-shadow);
}

.auth-form {
  --van-cell-horizontal-padding: 0;
  --van-field-label-color: var(--app-text-strong);
  --van-field-label-width: 5.5em;
  --van-field-input-text-color: var(--app-text);
}

.auth-form :deep(.van-cell) {
  margin-bottom: 12px;
  padding: 10px 14px;
  align-items: center;
  border: 1px solid var(--app-border);
  border-radius: 13px;
  background: var(--app-surface-muted);
}

.auth-form :deep(.van-cell::after) {
  display: none;
}

.auth-form :deep(.van-field__label),
.auth-form :deep(.van-field__control),
.auth-form :deep(.van-radio__label) {
  font-size: 14px;
}

.auth-form :deep(.van-field--error) {
  border-color: #df6c6c;
}

.auth-form :deep(.van-field__error-message) {
  padding-top: 5px;
  font-size: 12px;
}

.field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 12px;
}

.profile-details {
  margin: 2px 0 16px;
  border-top: 1px solid var(--app-border-soft);
  border-bottom: 1px solid var(--app-border-soft);
}

.profile-details summary {
  padding: 15px 2px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  list-style: none;
  cursor: pointer;
  color: var(--auth-ink);
}

.profile-details summary::-webkit-details-marker {
  display: none;
}

.profile-details summary strong,
.profile-details summary small {
  display: block;
}

.profile-details summary strong {
  margin-bottom: 3px;
  font-size: 14px;
}

.profile-details summary small {
  font-size: 12px;
  color: var(--app-text-muted);
}

.summary-icon {
  font-size: 21px;
  color: var(--auth-primary);
  transition: transform 180ms ease;
}

.profile-details[open] .summary-icon {
  transform: rotate(45deg);
}

.profile-fields {
  padding-top: 2px;
}

.profile-fields .remarks-field {
  grid-column: 1 / -1;
}

.submit-button {
  height: 50px;
  margin-top: 8px;
  border: 0;
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(115deg, var(--auth-primary), #0a9683);
  box-shadow: 0 12px 25px rgba(8, 127, 114, 0.22);
}

.form-footnote {
  margin: 20px 0 0;
  text-align: center;
  font-size: 14px;
  color: var(--auth-muted);
}

.form-footnote button {
  padding: 4px;
  border: 0;
  font: inherit;
  font-weight: 700;
  color: var(--auth-primary);
  background: transparent;
  cursor: pointer;
}

@media (max-width: 760px) {
  .auth-page {
    padding: 0;
    place-items: stretch;
    background: var(--app-surface);
  }

  .auth-shell {
    min-height: 100vh;
    display: block;
    border: 0;
    border-radius: 0;
    box-shadow: none;
  }

  .brand-panel {
    display: none;
  }

  .form-panel {
    padding: 32px 22px 28px;
  }

  .mobile-brand {
    margin-bottom: 44px;
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 17px;
    font-weight: 800;
    color: var(--auth-ink);
  }

  .mobile-brand__mark {
    width: 36px;
    height: 36px;
    display: grid;
    place-items: center;
    border-radius: 11px;
    background: var(--auth-primary);
    color: #fff;
  }

  .auth-heading h2 {
    font-size: 30px;
  }

  .field-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 380px) {
  .form-panel {
    padding-inline: 16px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .mode-switch button,
  .summary-icon {
    transition: none;
  }
}
</style>
