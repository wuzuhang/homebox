<script setup lang="ts">
import type { Medicine } from '@/api/medicine'
import { storeToRefs } from 'pinia'
import { showFailToast } from 'vant'
import { useMedicineStore } from '@/store/modules/medicine'
import 'vant/es/toast/style'

defineOptions({ name: 'Medicines' })

const router = useRouter()
const medicineStore = useMedicineStore()
const { medicines, diseases, loading, lowStockMedicines } = storeToRefs(medicineStore)
const selectedDiseaseId = ref<number | null>(null)
const costPeriods = [
  { key: 'day', label: '按天', days: 1, suffix: '元/天' },
  { key: 'month', label: '按月 ×30', days: 30, suffix: '元/月' },
  { key: 'year', label: '按年 ×365', days: 365, suffix: '元/年' },
] as const
type CostPeriod = typeof costPeriods[number]['key']
const costPeriod = ref<CostPeriod>('day')
const totalDailyCost = computed(() => medicines.value.reduce((summary, medicine) => {
  if (medicine.state !== 1)
    return summary

  const price = Number(medicine.price)
  const specifications = Number(medicine.specifications)
  const dailyDose = Number(medicine.daily_dose)
  if (price <= 0 || specifications <= 0 || dailyDose <= 0)
    return summary

  summary.amount += price / specifications * dailyDose
  summary.medicineCount += 1
  return summary
}, { amount: 0, medicineCount: 0 }))
const selectedTotalCost = computed(() => {
  const period = costPeriods.find(item => item.key === costPeriod.value) || costPeriods[0]
  return {
    ...period,
    amount: (totalDailyCost.value.amount * period.days).toFixed(2),
  }
})
const diseaseFilters = computed(() => {
  const usedIds = new Set(medicines.value.flatMap(item => item.disease_ids || []))
  return diseases.value
    .filter(disease => usedIds.has(disease.id))
    .map(disease => ({
      ...disease,
      medicineCount: medicines.value.filter(item => item.disease_ids?.includes(disease.id)).length,
    }))
})
const filteredMedicines = computed(() => {
  if (selectedDiseaseId.value === null)
    return medicines.value
  return medicines.value.filter(item => item.disease_ids?.includes(selectedDiseaseId.value as number))
})

watch(diseaseFilters, (filters) => {
  if (selectedDiseaseId.value !== null && !filters.some(item => item.id === selectedDiseaseId.value))
    selectedDiseaseId.value = null
})

function remainingDays(item: Medicine) {
  const dailyDose = Number(item.daily_dose)
  if (dailyDose <= 0)
    return null
  return Math.floor(Number(item.stock) / dailyDose)
}

function stockWarningClass(item: Medicine) {
  const days = remainingDays(item)
  if (days === null)
    return ''
  if (days <= 3)
    return 'stock-warning--critical'
  if (days <= 5)
    return 'stock-warning--danger'
  if (days <= 7)
    return 'stock-warning--yellow'
  return ''
}

function remainingDaysText(item: Medicine) {
  const days = remainingDays(item)
  if (days === null)
    return '预计可用天数未知'
  return `预计可用 ${days} 天`
}

async function loadPage() {
  try {
    await Promise.all([medicineStore.loadMedicines(), medicineStore.loadDiseases()])
  }
  catch (error) {
    const payload = error as { msg?: string }
    showFailToast(payload?.msg || '药品列表加载失败')
  }
}

onMounted(loadPage)
</script>

<template>
  <main class="medicine-page">
    <header class="page-heading">
      <div>
        <p>MEDICINE CABINET</p>
        <h1>我的药品</h1>
        <span>管理库存、剂量和用药说明</span>
      </div>
      <router-link class="add-button" :to="{ name: 'MedicineCreate' }">
        <van-icon name="plus" />
        新增药品
      </router-link>
    </header>

    <section v-if="medicines.length" class="summary-row" aria-label="药品库存概览">
      <div><strong>{{ medicines.length }}</strong><span>种药品</span></div>
      <div :class="{ warning: lowStockMedicines.length }"><strong>{{ lowStockMedicines.length }}</strong><span>种库存偏低</span></div>
    </section>

    <section v-if="medicines.length" class="total-cost-card" aria-label="全部药品预计消耗价格">
      <div class="total-cost-copy">
        <div>
          <p>全部药品预计消耗</p>
          <span>已计算 {{ totalDailyCost.medicineCount }}/{{ medicines.length }} 种药品（仅上架参与）</span>
        </div>
        <strong>{{ selectedTotalCost.amount }} <small>{{ selectedTotalCost.suffix }}</small></strong>
      </div>
      <div class="cost-periods" role="group" aria-label="消耗价格周期">
        <button
          v-for="period in costPeriods"
          :key="period.key"
          type="button"
          :class="{ active: costPeriod === period.key }"
          :aria-pressed="costPeriod === period.key"
          @click="costPeriod = period.key"
        >
          {{ period.label }}
        </button>
      </div>
    </section>

    <section v-if="medicines.length" class="category-filter" aria-label="按适用疾病筛选药品">
      <button
        type="button"
        :class="{ active: selectedDiseaseId === null }"
        :aria-pressed="selectedDiseaseId === null"
        @click="selectedDiseaseId = null"
      >
        <span>全部</span>
        <strong>{{ medicines.length }}</strong>
      </button>
      <button
        v-for="disease in diseaseFilters"
        :key="disease.id"
        type="button"
        :class="{ active: selectedDiseaseId === disease.id }"
        :aria-pressed="selectedDiseaseId === disease.id"
        @click="selectedDiseaseId = disease.id"
      >
        <span>{{ disease.name }}</span>
        <strong>{{ disease.medicineCount }}</strong>
      </button>
    </section>

    <div v-if="loading && !medicines.length" class="loading-list">
      <van-skeleton v-for="item in 3" :key="item" title avatar :row="2" />
    </div>

    <section v-else-if="medicines.length" class="medicine-list">
      <article
        v-for="item in filteredMedicines"
        :key="item.ID"
        class="medicine-card"
        :class="{ 'medicine-card--offline': item.state === 0 }"
        tabindex="0"
        role="link"
        @click="router.push({ name: 'MedicineDetail', params: { id: item.ID } })"
        @keydown.enter="router.push({ name: 'MedicineDetail', params: { id: item.ID } })"
      >
        <van-image v-if="item.photo" class="medicine-photo" width="76" height="76" radius="16" fit="cover" :src="item.photo" :alt="item.name" />
        <div v-else class="medicine-photo medicine-photo--empty"><van-icon name="medicines-o" /></div>
        <div class="medicine-main">
          <div class="medicine-title">
            <div>
              <h2>{{ item.name }}</h2>
              <p>{{ item.manufacturer || '未填写生产企业' }}</p>
            </div>
            <div class="medicine-statuses">
              <span v-if="item.state === 0" class="offline-status">已下架</span>
              <span v-if="stockWarningClass(item)" class="stock-warning" :class="stockWarningClass(item)">低库存</span>
            </div>
          </div>
          <div class="medicine-meta">
            <span><strong>{{ item.stock }}</strong> {{ item.dose_unit || '件' }}库存</span>
            <span>{{ remainingDaysText(item) }}</span>
          </div>
          <div v-if="item.disease_ids?.length" class="tag-list">
            <span v-for="id in item.disease_ids.slice(0, 3)" :key="id">{{ medicineStore.diseaseName(id) }}</span>
          </div>
        </div>
        <van-icon class="detail-arrow" name="arrow" />
      </article>
      <div v-if="!filteredMedicines.length" class="filter-empty">
        <van-icon name="search" />
        <span>该分类下暂无药品</span>
      </div>
    </section>

    <van-empty v-else image="search" description="还没有药品记录">
      <van-button round type="primary" @click="router.push({ name: 'MedicineCreate' })">添加第一种药品</van-button>
    </van-empty>
  </main>
</template>

<style lang="less" scoped>
.medicine-page { min-height: 100%; padding: 24px 18px 34px; color: var(--app-text); background: var(--app-page-background); }
.page-heading { max-width: 920px; margin: 0 auto 18px; display: flex; align-items: flex-end; justify-content: space-between; gap: 18px; }
.page-heading p, .page-heading h1, .page-heading span { margin: 0; }
.page-heading p { margin-bottom: 6px; font-size: 12px; font-weight: 800; letter-spacing: 0.16em; color: #0b8a7b; }
.page-heading h1 { margin-bottom: 7px; font-size: 28px; }
.page-heading span { font-size: 14px; color: var(--app-text-muted); }
.add-button { min-height: 44px; padding: 0 18px; display: inline-flex; align-items: center; gap: 7px; border-radius: 22px; font-size: 14px; font-weight: 700; color: #fff; background: linear-gradient(115deg, #087f72, #0a9683); box-shadow: 0 10px 24px rgba(8, 127, 114, 0.2); }
.summary-row { max-width: 920px; margin: 0 auto 14px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.summary-row > div { padding: 14px 18px; display: flex; align-items: baseline; gap: 7px; border: 1px solid var(--app-border); border-radius: 16px; background: var(--app-surface); }
.summary-row strong { font-size: 24px; color: #087f72; }
.summary-row span { font-size: 14px; color: var(--app-text-muted); }
.summary-row .warning strong { color: #d57928; }
.total-cost-card { max-width: 920px; margin: 0 auto 14px; padding: 18px; border: 1px solid var(--app-border); border-radius: 18px; background: var(--app-surface); }
.total-cost-copy { display: flex; align-items: center; justify-content: space-between; gap: 18px; }
.total-cost-copy p, .total-cost-copy span, .total-cost-copy strong { margin: 0; }
.total-cost-copy p { margin-bottom: 5px; font-size: 15px; font-weight: 700; color: var(--app-text-strong); }
.total-cost-copy span { font-size: 12px; color: var(--app-text-muted); }
.total-cost-copy strong { flex: 0 0 auto; font-size: 25px; color: #087f72; }
.total-cost-copy small { font-size: 12px; font-weight: 500; }
.cost-periods { margin-top: 14px; display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.cost-periods button { min-height: 36px; padding: 0 8px; border: 1px solid var(--app-border); border-radius: 18px; font-size: 13px; color: var(--app-text-muted); background: var(--app-surface-muted); cursor: pointer; }
.cost-periods button.active { border-color: #0a9683; font-weight: 700; color: #fff; background: #0a9683; }
.category-filter { max-width: 920px; margin: 0 auto 14px; padding: 12px 14px; display: flex; align-items: center; gap: 8px; overflow-x: auto; border: 1px solid var(--app-border); border-radius: 16px; background: var(--app-surface); scrollbar-width: none; }
.category-filter::-webkit-scrollbar { display: none; }
.category-filter button { min-height: 34px; padding: 0 10px 0 14px; flex: 0 0 auto; display: inline-flex; align-items: center; gap: 7px; border: 1px solid var(--app-border); border-radius: 17px; font-size: 13px; color: var(--app-text-muted); background: var(--app-surface-muted); cursor: pointer; transition: color 160ms ease, border-color 160ms ease, background 160ms ease; }
.category-filter button strong { min-width: 21px; height: 21px; padding: 0 6px; display: inline-flex; align-items: center; justify-content: center; border-radius: 11px; font-size: 11px; color: #087f72; background: var(--app-accent-soft); }
.category-filter button.active { border-color: #0a9683; font-weight: 700; color: #fff; background: #0a9683; }
.category-filter button.active strong { color: #087f72; background: #fff; }
.medicine-list, .loading-list { max-width: 920px; margin: 0 auto; display: grid; gap: 12px; }
.loading-list { padding: 24px; border-radius: 20px; background: var(--app-surface); }
.filter-empty { min-height: 150px; display: flex; align-items: center; justify-content: center; flex-direction: column; gap: 10px; border: 1px dashed var(--app-border); border-radius: 20px; font-size: 14px; color: var(--app-text-muted); background: var(--app-surface); }
.filter-empty :deep(.van-icon) { font-size: 30px; color: #72b9b0; }
.medicine-card { position: relative; padding: 18px; display: grid; grid-template-columns: 76px 1fr auto; align-items: center; gap: 18px; border: 1px solid var(--app-border); border-radius: 20px; background: var(--app-surface); cursor: pointer; transition: transform 180ms ease, box-shadow 180ms ease; }
.medicine-card:hover, .medicine-card:focus-visible { transform: translateY(-2px); outline: none; box-shadow: 0 12px 28px var(--app-shadow); }
.medicine-card--offline { border-color: var(--app-border-soft); background: var(--app-surface-muted); filter: grayscale(1); opacity: 0.68; }
.medicine-card--offline:hover, .medicine-card--offline:focus-visible { opacity: 0.82; }
.medicine-photo--empty { display: grid; place-items: center; border-radius: 16px; font-size: 32px; color: #087f72; background: var(--app-accent-soft); }
.medicine-main { min-width: 0; }
.medicine-title { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; }
.medicine-title h2, .medicine-title p { margin: 0; }
.medicine-title h2 { margin-bottom: 5px; font-size: 18px; color: var(--app-text-strong); }
.medicine-title p { font-size: 13px; color: var(--app-text-muted); }
.medicine-statuses { flex: 0 0 auto; display: flex; flex-wrap: wrap; justify-content: flex-end; gap: 6px; }
.offline-status { padding: 5px 9px; border-radius: 12px; font-size: 12px; font-weight: 700; color: #626a69; background: #e2e5e4; }
.stock-warning { flex: 0 0 auto; padding: 5px 9px; border-radius: 12px; font-size: 12px; font-weight: 700; color: #fff; }
.stock-warning--yellow { color: #6f4b00; background: #f4c542; }
.stock-warning--danger { background: #df5252; }
.stock-warning--critical { background: #991f32; }
.medicine-meta { margin-top: 11px; display: grid; grid-template-columns: auto 1fr; align-items: baseline; gap: 8px; font-size: 13px; color: var(--app-text-muted); }
.medicine-meta span { white-space: nowrap; }
.medicine-meta span:last-child { text-align: right; }
.medicine-meta strong { font-size: 16px; color: #087f72; }
.tag-list { margin-top: 10px; display: flex; flex-wrap: wrap; gap: 6px; }
.tag-list span { padding: 4px 8px; border-radius: 10px; font-size: 12px; color: #0a9683; background: var(--app-accent-soft); }
.detail-arrow { margin-left: 4px; font-size: 18px; color: #91a4a1; }

@media (max-width: 640px) {
  .medicine-page { padding: 18px 12px 28px; }
  .page-heading { align-items: flex-start; }
  .page-heading h1 { font-size: 25px; }
  .add-button { padding: 0 14px; }
  .total-cost-card { padding: 15px; }
  .total-cost-copy { align-items: flex-start; }
  .total-cost-copy strong { font-size: 21px; }
  .cost-periods button { padding-inline: 4px; font-size: 12px; }
  .category-filter { margin-bottom: 12px; padding: 10px; }
  .medicine-card { padding: 14px; grid-template-columns: 62px 1fr; gap: 13px; }
  .medicine-photo { width: 62px !important; height: 62px !important; }
  .medicine-title { display: block; }
  .medicine-statuses { position: absolute; top: 12px; right: 12px; }
  .medicine-title h2 { padding-right: 62px; font-size: 16px; }
  .medicine-meta { font-size: 12px; }
  .detail-arrow { position: absolute; right: 15px; bottom: 16px; margin: 0; }
}
</style>
