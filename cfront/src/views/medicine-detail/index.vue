<script setup lang="ts">
import { showConfirmDialog, showFailToast, showImagePreview, showSuccessToast } from 'vant'
import { deleteMedicine } from '@/api/medicine'
import { useMedicineStore } from '@/store/modules/medicine'
import 'vant/es/dialog/style'
import 'vant/es/image-preview/style'
import 'vant/es/toast/style'

defineOptions({ name: 'MedicineDetail' })

const route = useRoute()
const router = useRouter()
const medicineStore = useMedicineStore()
const medicineId = computed(() => Number(route.params.id))
const medicine = computed(() => medicineStore.findMedicine(medicineId.value))
const remainingDays = computed(() => {
  const dailyDose = Number(medicine.value?.daily_dose)
  if (!medicine.value || dailyDose <= 0)
    return null
  return Number(medicine.value.stock) / dailyDose
})
const stockWarningClass = computed(() => {
  if (remainingDays.value === null || remainingDays.value > 7)
    return ''
  if (remainingDays.value <= 3)
    return 'stock-warning--critical'
  if (remainingDays.value <= 5)
    return 'stock-warning--danger'
  return 'stock-warning--yellow'
})
const isLowStock = computed(() => Boolean(stockWarningClass.value))
const remainingDaysText = computed(() => {
  if (remainingDays.value === null)
    return '未设置每日剂量'
  const days = Number.isInteger(remainingDays.value) ? remainingDays.value : remainingDays.value.toFixed(1)
  return `${days} 天`
})
const costPeriods = [
  { key: 'day', label: '按天', button: '按天', days: 1, suffix: '元/天' },
  { key: 'month', label: '按月（30 天）', button: '按月 ×30', days: 30, suffix: '元/月' },
  { key: 'year', label: '按年（365 天）', button: '按年 ×365', days: 365, suffix: '元/年' },
] as const
type CostPeriod = typeof costPeriods[number]['key']
const costPeriod = ref<CostPeriod>('day')
const dailyConsumptionCost = computed(() => {
  if (!medicine.value)
    return null

  const price = Number(medicine.value.price)
  const specifications = Number(medicine.value.specifications)
  const dailyDose = Number(medicine.value.daily_dose)
  if (price <= 0 || specifications <= 0 || dailyDose <= 0)
    return null

  return price / specifications * dailyDose
})
const selectedConsumptionCost = computed(() => {
  if (dailyConsumptionCost.value === null)
    return null

  const period = costPeriods.find(item => item.key === costPeriod.value) || costPeriods[0]
  return {
    ...period,
    amount: (dailyConsumptionCost.value * period.days).toFixed(2),
  }
})

function previewPhoto() {
  if (!medicine.value?.photo)
    return
  showImagePreview({
    images: [medicine.value.photo],
    closeable: true,
  })
}

const detailRows = computed(() => medicine.value ? [
  { label: '生产企业', value: medicine.value.manufacturer || '未设置' },
  { label: '当前库存', value: `${medicine.value.stock} ${medicine.value.dose_unit || '件'}` },
  { label: '药品规格', value: medicine.value.specifications ? `${medicine.value.specifications} ${medicine.value.dose_unit || '件'}/${medicine.value.package_unit || '盒'}` : '未设置' },
  { label: '药品单价', value: medicine.value.price ? `${Number(medicine.value.price).toFixed(2)} 元/${medicine.value.package_unit || '盒'}` : '未设置' },
  { label: '低库存阈值', value: `${medicine.value.min_stock_warn} ${medicine.value.dose_unit || '件'}` },
  { label: '每日剂量', value: medicine.value.daily_dose ? `${medicine.value.daily_dose} ${medicine.value.dose_unit || '件'}` : '未设置' },
  { label: '预计可用', value: remainingDaysText.value },
] : [])

async function initialize() {
  try {
    await Promise.all([
      medicineStore.loadDiseases(),
      medicine.value ? Promise.resolve() : medicineStore.loadMedicines(),
    ])
    if (!medicine.value) {
      showFailToast('没有找到该药品')
      await router.replace({ name: 'Medicines' })
    }
  }
  catch (error) {
    const payload = error as { msg?: string }
    showFailToast(payload?.msg || '药品详情加载失败')
  }
}

async function removeCurrent() {
  if (!medicine.value)
    return
  try {
    await showConfirmDialog({
      title: '删除药品',
      message: `确定删除“${medicine.value.name}”吗？删除后无法恢复。`,
      confirmButtonText: '确认删除',
      confirmButtonColor: '#d24b4b',
    })
    await deleteMedicine(medicine.value.ID)
    await medicineStore.loadMedicines()
    showSuccessToast('删除成功')
    await router.replace({ name: 'Medicines' })
  }
  catch (error) {
    if (error === 'cancel' || error === 'close')
      return
    const payload = error as { msg?: string }
    showFailToast(payload?.msg || '删除失败，请稍后重试')
  }
}

onMounted(initialize)
</script>

<template>
  <main class="detail-page">
    <section v-if="medicine" class="medicine-detail">
      <header class="medicine-hero">
        <van-image
          v-if="medicine.photo"
          class="medicine-photo medicine-photo--clickable"
          width="104"
          height="104"
          radius="22"
          fit="cover"
          :src="medicine.photo"
          :alt="medicine.name"
          role="button"
          tabindex="0"
          aria-label="查看药品大图"
          @click="previewPhoto"
          @keydown.enter="previewPhoto"
          @keydown.space.prevent="previewPhoto"
        />
        <div v-else class="medicine-photo medicine-photo--empty"><van-icon name="medicines-o" /></div>
        <div class="hero-copy">
          <div class="hero-labels">
            <span>药品档案</span>
            <span v-if="medicine.state === 0" class="offline-label">已下架</span>
            <span v-if="isLowStock" class="warning-label" :class="stockWarningClass">低库存预警</span>
          </div>
          <h1>{{ medicine.name }}</h1>
          <p>{{ medicine.manufacturer || '未填写生产企业' }}</p>
        </div>
      </header>

      <div class="detail-content">
        <section class="stock-panel" :class="stockWarningClass">
          <div>
            <p>当前库存</p>
            <strong>{{ medicine.stock }} <small>{{ medicine.dose_unit || '件' }}</small></strong>
          </div>
          <van-icon :name="isLowStock ? 'warning-o' : 'passed'" />
        </section>

        <section class="cost-section">
          <div class="cost-heading">
            <div>
              <p>预计消耗价格</p>
              <span>单价 ÷ 每包装规格 × 每日剂量</span>
            </div>
            <van-icon name="gold-coin-o" />
          </div>
          <div v-if="selectedConsumptionCost" class="cost-display">
            <span>{{ selectedConsumptionCost.label }}</span>
            <strong>{{ selectedConsumptionCost.amount }} <small>{{ selectedConsumptionCost.suffix }}</small></strong>
          </div>
          <p v-else class="cost-empty">请先完善药品单价、规格和每日剂量</p>
          <div class="cost-periods" role="group" aria-label="消耗价格周期">
            <button
              v-for="period in costPeriods"
              :key="period.key"
              type="button"
              :class="{ active: costPeriod === period.key }"
              :aria-pressed="costPeriod === period.key"
              @click="costPeriod = period.key"
            >
              {{ period.button }}
            </button>
          </div>
        </section>

        <section class="info-section">
          <h2>基础信息</h2>
          <dl>
            <div v-for="row in detailRows" :key="row.label">
              <dt>{{ row.label }}</dt>
              <dd>{{ row.value }}</dd>
            </div>
          </dl>
        </section>

        <section class="info-section">
          <h2>适用疾病</h2>
          <div v-if="medicine.disease_ids?.length" class="tag-list">
            <span v-for="id in medicine.disease_ids" :key="id">{{ medicineStore.diseaseName(id) }}</span>
          </div>
          <p v-else class="empty-text">暂未关联疾病标签</p>
        </section>

        <section class="info-section text-section">
          <div>
            <h2>用法用量</h2>
            <p>{{ medicine.usage || '暂未填写用法用量' }}</p>
          </div>
          <div>
            <h2>备注</h2>
            <p>{{ medicine.remark || '暂未填写备注' }}</p>
          </div>
        </section>

        <div class="page-actions">
          <button class="delete-button" type="button" @click="removeCurrent"><van-icon name="delete-o" />删除</button>
          <router-link class="edit-button" :to="{ name: 'MedicineEdit', params: { id: medicine.ID } }"><van-icon name="edit" />编辑药品</router-link>
        </div>
      </div>
    </section>
  </main>
</template>

<style lang="less" scoped>
.detail-page { min-height: 100%; padding: 22px 18px 34px; color: var(--app-text); background: var(--app-page-background); }
.medicine-detail { max-width: 820px; margin: 0 auto; overflow: hidden; border: 1px solid var(--app-border); border-radius: 26px; background: var(--app-surface); box-shadow: 0 16px 42px var(--app-shadow); }
.medicine-hero { padding: 34px; display: flex; align-items: center; gap: 22px; color: #fff; background: linear-gradient(135deg, #087f72, #09665e); }
.medicine-photo--empty { display: grid; place-items: center; flex: 0 0 104px; border-radius: 22px; font-size: 44px; color: #087568; background: var(--app-accent-soft); }
.medicine-photo--clickable { cursor: zoom-in; transition: transform 180ms ease, box-shadow 180ms ease; }
.medicine-photo--clickable:hover, .medicine-photo--clickable:focus-visible { transform: scale(1.03); outline: none; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.22); }
.hero-copy { min-width: 0; }
.hero-copy h1, .hero-copy p { margin: 0; }
.hero-copy h1 { margin: 8px 0 7px; font-size: 28px; overflow-wrap: anywhere; }
.hero-copy p { font-size: 14px; color: rgba(255, 255, 255, 0.72); }
.hero-labels { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.hero-labels span { font-size: 12px; font-weight: 700; color: #c8fff3; }
.hero-labels .offline-label { padding: 5px 9px; border-radius: 12px; color: #555d5c; background: #e2e5e4; }
.hero-labels .warning-label { padding: 5px 9px; border-radius: 12px; color: #fff; }
.hero-labels .stock-warning--yellow { color: #6f4b00; background: #f4c542; }
.hero-labels .stock-warning--danger { background: #df5252; }
.hero-labels .stock-warning--critical { background: #991f32; }
.detail-content { padding: 28px 32px 32px; }
.stock-panel { padding: 18px 20px; display: flex; align-items: center; justify-content: space-between; border-radius: 17px; color: #0a9683; background: var(--app-accent-soft); }
.stock-panel.stock-warning--yellow { color: #6f4b00; background: #f4c542; }
.stock-panel.stock-warning--danger { color: #fff; background: #df5252; }
.stock-panel.stock-warning--critical { color: #fff; background: #991f32; }
.stock-panel p, .stock-panel strong { margin: 0; }
.stock-panel p { margin-bottom: 5px; font-size: 13px; }
.stock-panel strong { font-size: 26px; }
.stock-panel small { font-size: 14px; }
.stock-panel :deep(.van-icon) { font-size: 30px; }
.cost-section { margin-top: 14px; padding: 18px 20px; border: 1px solid var(--app-border); border-radius: 17px; background: var(--app-surface-muted); }
.cost-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.cost-heading p, .cost-heading span { margin: 0; }
.cost-heading p { margin-bottom: 5px; font-size: 15px; font-weight: 700; color: var(--app-text-strong); }
.cost-heading span { font-size: 12px; color: var(--app-text-muted); }
.cost-heading :deep(.van-icon) { font-size: 28px; color: #d6962f; }
.cost-display { margin-top: 16px; padding: 16px 18px; border-radius: 14px; background: var(--app-surface); }
.cost-display span, .cost-display strong { display: block; }
.cost-display span { margin-bottom: 7px; font-size: 13px; color: var(--app-text-muted); }
.cost-display strong { font-size: 26px; color: #087f72; }
.cost-display small { font-size: 12px; font-weight: 500; }
.cost-empty { margin: 14px 0 0; font-size: 13px; color: var(--app-text-muted); }
.cost-periods { margin-top: 12px; display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.cost-periods button { min-height: 38px; padding: 0 10px; border: 1px solid var(--app-border); border-radius: 19px; font-size: 13px; color: var(--app-text-muted); background: var(--app-surface); cursor: pointer; transition: color 160ms ease, border-color 160ms ease, background 160ms ease; }
.cost-periods button.active { border-color: #0a9683; font-weight: 700; color: #fff; background: #0a9683; }
.info-section { padding: 24px 0; border-bottom: 1px solid var(--app-border); }
.info-section h2 { margin: 0 0 15px; font-size: 17px; color: var(--app-text-strong); }
.info-section dl { margin: 0; display: grid; grid-template-columns: 1fr 1fr; gap: 0 30px; }
.info-section dl > div { min-height: 48px; display: flex; align-items: center; justify-content: space-between; gap: 15px; border-bottom: 1px solid var(--app-border-soft); }
.info-section dt { font-size: 14px; color: var(--app-text-muted); }
.info-section dd { margin: 0; text-align: right; font-size: 14px; font-weight: 700; color: var(--app-text-strong); }
.tag-list { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-list span { padding: 7px 11px; border-radius: 14px; font-size: 13px; color: #0a9683; background: var(--app-accent-soft); }
.empty-text, .text-section p { margin: 0; font-size: 14px; line-height: 1.75; color: var(--app-text-muted); }
.text-section { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; }
.page-actions { padding-top: 24px; display: grid; grid-template-columns: 1fr 1.5fr; gap: 12px; }
.page-actions > * { min-height: 46px; display: flex; align-items: center; justify-content: center; gap: 7px; border-radius: 23px; font-size: 14px; font-weight: 700; cursor: pointer; }
.delete-button { border: 1px solid #eccaca; color: #bd4949; background: var(--app-surface); }
.edit-button { color: #fff; background: linear-gradient(115deg, #087f72, #0a9683); box-shadow: 0 10px 24px rgba(8, 127, 114, 0.2); }

@media (max-width: 620px) {
  .detail-page { padding: 14px 12px 28px; }
  .medicine-hero { padding: 26px 20px; align-items: flex-start; }
  .medicine-photo { width: 76px !important; height: 76px !important; flex-basis: 76px; }
  .hero-copy h1 { font-size: 22px; }
  .detail-content { padding: 20px 16px 24px; }
  .cost-section { padding: 16px; }
  .cost-display { padding: 14px; }
  .cost-display strong { font-size: 22px; }
  .cost-periods button { padding-inline: 5px; font-size: 12px; }
  .info-section dl, .text-section { grid-template-columns: 1fr; }
  .text-section { gap: 22px; }
}
</style>
