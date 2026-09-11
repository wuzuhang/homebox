<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { showFailToast, showSuccessToast } from 'vant'
import { addMedicine, updateMedicine } from '@/api/medicine'
import { useMedicineStore } from '@/store/modules/medicine'
import { useUserStore } from '@/store/modules/user'
import 'vant/es/toast/style'

defineOptions({ name: 'MedicineForm' })

const route = useRoute()
const router = useRouter()
const medicineStore = useMedicineStore()
const userStore = useUserStore()
const { diseases } = storeToRefs(medicineStore)
const submitting = ref(false)
const isEdit = computed(() => route.name === 'MedicineEdit')
const medicineId = computed(() => Number(route.params.id || 0))

const form = reactive({
  name: '',
  manufacturer: '',
  disease_ids: [] as number[],
  stock: '',
  unit: '盒',
  min_stock_warn: '2',
  daily_dose: '',
  usage: '',
  photo: '',
  remark: '',
})

function fillForm() {
  const medicine = medicineStore.findMedicine(medicineId.value)
  if (!medicine)
    return false
  Object.assign(form, {
    name: medicine.name,
    manufacturer: medicine.manufacturer || '',
    disease_ids: [...(medicine.disease_ids || [])],
    stock: String(medicine.stock ?? ''),
    unit: medicine.unit || '盒',
    min_stock_warn: String(medicine.min_stock_warn ?? 2),
    daily_dose: medicine.daily_dose ? String(medicine.daily_dose) : '',
    usage: medicine.usage || '',
    photo: medicine.photo || '',
    remark: medicine.remark || '',
  })
  return true
}

function errorMessage(error: unknown) {
  const payload = error as { msg?: string, message?: string }
  return payload?.msg || payload?.message || '保存失败，请稍后重试'
}

async function initialize() {
  try {
    await medicineStore.loadDiseases()
    if (isEdit.value && !fillForm()) {
      await medicineStore.loadMedicines()
      if (!fillForm()) {
        showFailToast('没有找到该药品')
        await router.replace({ name: 'Medicines' })
      }
    }
  }
  catch (error) {
    showFailToast(errorMessage(error))
  }
}

async function submit() {
  submitting.value = true
  try {
    const userId = userStore.user?.ID
    if (!userId)
      throw new Error('未获取到当前用户 ID，请重新登录')

    const payload = {
      id: isEdit.value ? medicineId.value : undefined,
      user_id: userId,
      name: form.name.trim(),
      manufacturer: form.manufacturer.trim(),
      disease_ids: form.disease_ids,
      stock: Number(form.stock || 0),
      unit: form.unit,
      min_stock_warn: Number(form.min_stock_warn || 0),
      daily_dose: Number(form.daily_dose || 0),
      usage: form.usage.trim(),
      photo: form.photo.trim(),
      remark: form.remark.trim(),
    }
    if (isEdit.value)
      await updateMedicine(payload)
    else
      await addMedicine(payload)

    showSuccessToast(isEdit.value ? '更新成功' : '添加成功')
    await router.replace({ name: 'Medicines' })
  }
  catch (error) {
    showFailToast(errorMessage(error))
  }
  finally {
    submitting.value = false
  }
}

onMounted(initialize)
</script>

<template>
  <main class="form-page">
    <section class="form-card">
      <header class="form-heading">
        <div>
          <p>MEDICINE RECORD</p>
          <h1>{{ isEdit ? '编辑药品' : '新增药品' }}</h1>
          <span>记录库存与用法，低于预警值时会在列表中提醒。</span>
        </div>
        <van-image v-if="form.photo" round width="64" height="64" fit="cover" :src="form.photo" alt="药品图片预览" />
        <div v-else class="photo-preview"><van-icon name="medicines-o" /></div>
      </header>

      <van-form class="medicine-form" @submit="submit">
        <div class="field-grid">
          <van-field v-model.trim="form.name" name="name" label="药品名称" placeholder="如：布洛芬缓释胶囊" :rules="[{ required: true, message: '请输入药品名称' }]" />
          <van-field v-model.trim="form.manufacturer" name="manufacturer" label="生产企业" placeholder="选填" />
          <van-field v-model="form.stock" name="stock" label="当前库存" type="number" placeholder="0" :min="0" step="0.01" />
          <van-field name="unit" label="库存单位">
            <template #input>
              <select v-model="form.unit" class="unit-select" aria-label="库存单位">
                <option v-for="unit in ['粒', '片', '盒', '支', '瓶', '袋', 'ml']" :key="unit" :value="unit">{{ unit }}</option>
              </select>
            </template>
          </van-field>
          <van-field v-model="form.min_stock_warn" name="min_stock_warn" label="预警库存" type="number" placeholder="2" :min="0" step="0.01" />
          <van-field v-model="form.daily_dose" name="daily_dose" label="每日剂量" type="number" placeholder="0" :min="0" step="0.01" />
          <van-field v-model.trim="form.photo" name="photo" label="图片地址" type="url" placeholder="https://" />
          <van-field v-model.trim="form.usage" name="usage" label="用法用量" placeholder="如：饭后口服，一日 2 次" />
          <van-field class="full-field disease-field" name="disease_ids" label="适用疾病">
            <template #input>
              <van-checkbox-group v-if="diseases.length" v-model="form.disease_ids" direction="horizontal">
                <van-checkbox v-for="disease in diseases" :key="disease.id" :name="disease.id" shape="square">{{ disease.name }}</van-checkbox>
              </van-checkbox-group>
              <span v-else class="empty-hint">暂无疾病标签</span>
            </template>
          </van-field>
          <van-field v-model.trim="form.remark" class="full-field" name="remark" label="备注" type="textarea" rows="3" autosize maxlength="255" show-word-limit placeholder="其他需要留意的信息" />
        </div>

        <div class="form-actions">
          <van-button block round plain type="primary" native-type="button" @click="router.back()">取消</van-button>
          <van-button block round type="primary" native-type="submit" :loading="submitting" loading-text="正在保存…">{{ isEdit ? '保存修改' : '添加药品' }}</van-button>
        </div>
      </van-form>
    </section>
  </main>
</template>

<style lang="less" scoped>
.form-page { min-height: 100%; padding: 22px 18px 34px; color: var(--app-text); background: var(--app-page-background); }
.form-card { max-width: 860px; margin: 0 auto; padding: 30px; border: 1px solid var(--app-border); border-radius: 24px; background: var(--app-surface); box-shadow: 0 16px 42px var(--app-shadow); }
.form-heading { margin-bottom: 28px; display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.form-heading p, .form-heading h1, .form-heading span { margin: 0; }
.form-heading p { margin-bottom: 6px; font-size: 12px; font-weight: 800; letter-spacing: 0.16em; color: #0b8a7b; }
.form-heading h1 { margin-bottom: 8px; font-size: 28px; }
.form-heading span { font-size: 14px; color: var(--app-text-muted); }
.photo-preview { width: 64px; height: 64px; flex: 0 0 64px; display: grid; place-items: center; border-radius: 18px; font-size: 30px; color: #087568; background: var(--app-accent-soft); }
.medicine-form { --van-field-label-color: var(--app-text-strong); --van-field-input-text-color: var(--app-text); --van-field-label-width: 5.8em; }
.field-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0 14px; }
.medicine-form :deep(.van-cell) { margin-bottom: 14px; padding: 11px 14px; align-items: center; border: 1px solid var(--app-border); border-radius: 13px; background: var(--app-surface-muted); }
.medicine-form :deep(.van-cell::after) { display: none; }
.medicine-form :deep(.van-field__label), .medicine-form :deep(.van-field__control), .medicine-form :deep(.van-checkbox__label) { font-size: 14px; }
.medicine-form :deep(.van-field--error) { border-color: #df6c6c; }
.unit-select { width: 100%; border: 0; outline: 0; font-size: 14px; color: var(--app-text); background: transparent; }
.full-field { grid-column: 1 / -1; }
.disease-field :deep(.van-checkbox-group) { display: flex; flex-wrap: wrap; gap: 12px 16px; }
.empty-hint { font-size: 14px; color: var(--app-text-muted); }
.form-actions { margin-top: 10px; display: grid; grid-template-columns: 1fr 1.4fr; gap: 12px; }
.form-actions :deep(.van-button) { height: 48px; font-size: 15px; font-weight: 700; }
.form-actions :deep(.van-button--normal:not(.van-button--plain)) { border: 0; background: linear-gradient(115deg, #087f72, #0a9683); box-shadow: 0 10px 24px rgba(8, 127, 114, 0.2); }

@media (max-width: 680px) {
  .form-page { padding: 14px 12px 26px; }
  .form-card { padding: 22px 14px; border-radius: 20px; }
  .form-heading { padding: 0 4px; }
  .form-heading h1 { font-size: 24px; }
  .field-grid { grid-template-columns: 1fr; }
  .full-field { grid-column: auto; }
  .form-actions { grid-template-columns: 1fr 1.35fr; }
}
</style>
