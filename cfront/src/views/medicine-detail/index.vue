<script setup lang="ts">
import { showConfirmDialog, showFailToast, showSuccessToast } from 'vant'
import { deleteMedicine } from '@/api/medicine'
import { useMedicineStore } from '@/store/modules/medicine'
import 'vant/es/dialog/style'
import 'vant/es/toast/style'

defineOptions({ name: 'MedicineDetail' })

const route = useRoute()
const router = useRouter()
const medicineStore = useMedicineStore()
const medicineId = computed(() => Number(route.params.id))
const medicine = computed(() => medicineStore.findMedicine(medicineId.value))
const isLowStock = computed(() => Boolean(medicine.value && medicine.value.stock <= medicine.value.min_stock_warn))

const detailRows = computed(() => medicine.value ? [
  { label: '生产企业', value: medicine.value.manufacturer || '未设置' },
  { label: '当前库存', value: `${medicine.value.stock} ${medicine.value.unit || '件'}` },
  { label: '低库存阈值', value: `${medicine.value.min_stock_warn} ${medicine.value.unit || '件'}` },
  { label: '每日剂量', value: medicine.value.daily_dose ? `${medicine.value.daily_dose} ${medicine.value.unit || '件'}` : '未设置' },
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
        <van-image v-if="medicine.photo" class="medicine-photo" width="104" height="104" radius="22" fit="cover" :src="medicine.photo" :alt="medicine.name" />
        <div v-else class="medicine-photo medicine-photo--empty"><van-icon name="medicines-o" /></div>
        <div class="hero-copy">
          <div class="hero-labels">
            <span>药品档案</span>
            <span v-if="isLowStock" class="warning-label">低库存预警</span>
          </div>
          <h1>{{ medicine.name }}</h1>
          <p>{{ medicine.manufacturer || '未填写生产企业' }}</p>
        </div>
      </header>

      <div class="detail-content">
        <section class="stock-panel" :class="{ warning: isLowStock }">
          <div>
            <p>当前库存</p>
            <strong>{{ medicine.stock }} <small>{{ medicine.unit || '件' }}</small></strong>
          </div>
          <van-icon :name="isLowStock ? 'warning-o' : 'passed'" />
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
.hero-copy { min-width: 0; }
.hero-copy h1, .hero-copy p { margin: 0; }
.hero-copy h1 { margin: 8px 0 7px; font-size: 28px; overflow-wrap: anywhere; }
.hero-copy p { font-size: 14px; color: rgba(255, 255, 255, 0.72); }
.hero-labels { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.hero-labels span { font-size: 12px; font-weight: 700; color: #c8fff3; }
.hero-labels .warning-label { padding: 5px 9px; border-radius: 12px; color: #d9772c; background: var(--app-warning-soft); }
.detail-content { padding: 28px 32px 32px; }
.stock-panel { padding: 18px 20px; display: flex; align-items: center; justify-content: space-between; border-radius: 17px; color: #0a9683; background: var(--app-accent-soft); }
.stock-panel.warning { color: #d9772c; background: var(--app-warning-soft); }
.stock-panel p, .stock-panel strong { margin: 0; }
.stock-panel p { margin-bottom: 5px; font-size: 13px; }
.stock-panel strong { font-size: 26px; }
.stock-panel small { font-size: 14px; }
.stock-panel :deep(.van-icon) { font-size: 30px; }
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
  .info-section dl, .text-section { grid-template-columns: 1fr; }
  .text-section { gap: 22px; }
}
</style>
