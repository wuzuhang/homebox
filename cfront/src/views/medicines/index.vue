<script setup lang="ts">
import type { Medicine } from '@/api/medicine'
import { storeToRefs } from 'pinia'
import { showConfirmDialog, showFailToast, showSuccessToast } from 'vant'
import { deleteMedicine } from '@/api/medicine'
import { useMedicineStore } from '@/store/modules/medicine'
import 'vant/es/dialog/style'
import 'vant/es/toast/style'

defineOptions({ name: 'Medicines' })

const router = useRouter()
const medicineStore = useMedicineStore()
const { medicines, loading, lowStockMedicines } = storeToRefs(medicineStore)

async function loadPage() {
  try {
    await Promise.all([medicineStore.loadMedicines(), medicineStore.loadDiseases()])
  }
  catch (error) {
    const payload = error as { msg?: string }
    showFailToast(payload?.msg || '药品列表加载失败')
  }
}

async function removeMedicine(item: Medicine) {
  try {
    await showConfirmDialog({
      title: '删除药品',
      message: `确定删除“${item.name}”吗？删除后无法恢复。`,
      confirmButtonText: '确认删除',
      confirmButtonColor: '#d24b4b',
    })
    await deleteMedicine(item.ID)
    showSuccessToast('删除成功')
    await medicineStore.loadMedicines()
  }
  catch (error) {
    if (error === 'cancel' || error === 'close')
      return
    const payload = error as { msg?: string }
    showFailToast(payload?.msg || '删除失败，请稍后重试')
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

    <div v-if="loading && !medicines.length" class="loading-list">
      <van-skeleton v-for="item in 3" :key="item" title avatar :row="2" />
    </div>

    <section v-else-if="medicines.length" class="medicine-list">
      <article
        v-for="item in medicines"
        :key="item.ID"
        class="medicine-card"
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
            <span v-if="item.stock <= item.min_stock_warn" class="stock-warning">低库存</span>
          </div>
          <div class="medicine-meta">
            <span><strong>{{ item.stock }}</strong> {{ item.unit || '件' }}库存</span>
            <span v-if="item.daily_dose">每日 {{ item.daily_dose }} {{ item.unit || '件' }}</span>
          </div>
          <div v-if="item.disease_ids?.length" class="tag-list">
            <span v-for="id in item.disease_ids.slice(0, 3)" :key="id">{{ medicineStore.diseaseName(id) }}</span>
          </div>
        </div>
        <div class="card-actions">
          <button type="button" aria-label="编辑药品" @click.stop="router.push({ name: 'MedicineEdit', params: { id: item.ID } })"><van-icon name="edit" /></button>
          <button class="delete-action" type="button" aria-label="删除药品" @click.stop="removeMedicine(item)"><van-icon name="delete-o" /></button>
          <van-icon class="detail-arrow" name="arrow" />
        </div>
      </article>
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
.medicine-list, .loading-list { max-width: 920px; margin: 0 auto; display: grid; gap: 12px; }
.loading-list { padding: 24px; border-radius: 20px; background: var(--app-surface); }
.medicine-card { position: relative; padding: 18px; display: grid; grid-template-columns: 76px 1fr auto; align-items: center; gap: 18px; border: 1px solid var(--app-border); border-radius: 20px; background: var(--app-surface); cursor: pointer; transition: transform 180ms ease, box-shadow 180ms ease; }
.medicine-card:hover, .medicine-card:focus-visible { transform: translateY(-2px); outline: none; box-shadow: 0 12px 28px var(--app-shadow); }
.medicine-photo--empty { display: grid; place-items: center; border-radius: 16px; font-size: 32px; color: #087f72; background: var(--app-accent-soft); }
.medicine-main { min-width: 0; }
.medicine-title { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; }
.medicine-title h2, .medicine-title p { margin: 0; }
.medicine-title h2 { margin-bottom: 5px; font-size: 18px; color: var(--app-text-strong); }
.medicine-title p { font-size: 13px; color: var(--app-text-muted); }
.stock-warning { flex: 0 0 auto; padding: 5px 9px; border-radius: 12px; font-size: 12px; font-weight: 700; color: #d9772c; background: var(--app-warning-soft); }
.medicine-meta { margin-top: 11px; display: flex; flex-wrap: wrap; gap: 8px 18px; font-size: 13px; color: var(--app-text-muted); }
.medicine-meta strong { font-size: 16px; color: #087f72; }
.tag-list { margin-top: 10px; display: flex; flex-wrap: wrap; gap: 6px; }
.tag-list span { padding: 4px 8px; border-radius: 10px; font-size: 12px; color: #0a9683; background: var(--app-accent-soft); }
.card-actions { display: flex; align-items: center; gap: 5px; }
.card-actions button { width: 38px; height: 38px; display: grid; place-items: center; border: 0; border-radius: 50%; font-size: 18px; color: #45a398; background: var(--app-accent-soft); cursor: pointer; }
.card-actions .delete-action { color: #dc6a6a; background: var(--app-danger-soft); }
.detail-arrow { margin-left: 4px; color: #91a4a1; }

@media (max-width: 640px) {
  .medicine-page { padding: 18px 12px 28px; }
  .page-heading { align-items: flex-start; }
  .page-heading h1 { font-size: 25px; }
  .add-button { padding: 0 14px; }
  .medicine-card { padding: 14px; grid-template-columns: 62px 1fr; gap: 13px; }
  .medicine-photo { width: 62px !important; height: 62px !important; }
  .medicine-title { display: block; }
  .stock-warning { position: absolute; top: 12px; right: 12px; }
  .medicine-title h2 { padding-right: 62px; font-size: 16px; }
  .card-actions { grid-column: 1 / -1; justify-content: flex-end; padding-top: 10px; border-top: 1px solid var(--app-border-soft); }
  .detail-arrow { margin-left: auto; }
}
</style>
