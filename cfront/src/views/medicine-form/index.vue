<script setup lang="ts">
import { storeToRefs } from "pinia";
import { showFailToast, showSuccessToast } from "vant";
import type { MedicinePayload } from "@/api/medicine";
import { addMedicine, updateMedicine } from "@/api/medicine";
import { useMedicineStore } from "@/store/modules/medicine";
import { useUserStore } from "@/store/modules/user";
import "vant/es/toast/style";

defineOptions({ name: "MedicineForm" });

const route = useRoute();
const router = useRouter();
const medicineStore = useMedicineStore();
const userStore = useUserStore();
const { diseases } = storeToRefs(medicineStore);
const submitting = ref(false);
const stateSubmitting = ref(false);
const isEdit = computed(() => route.name === "MedicineEdit");
const medicineId = computed(() => Number(route.params.id || 0));
const medicineUnits = ["粒", "片", "盒", "支", "瓶", "袋", "ml"];

const form = reactive({
  name: "",
  manufacturer: "",
  disease_ids: [] as number[],
  stock: "",
  package_unit: "盒",
  dose_unit: "粒",
  specifications: "0",
  price: "0",
  min_stock_warn: "2",
  daily_dose: "",
  usage: "",
  photo: "",
  state: 1,
  remark: "",
});

function fillForm() {
  const medicine = medicineStore.findMedicine(medicineId.value);
  if (!medicine) return false;
  Object.assign(form, {
    name: medicine.name,
    manufacturer: medicine.manufacturer || "",
    disease_ids: [...(medicine.disease_ids || [])],
    stock: String(medicine.stock ?? ""),
    package_unit: medicine.package_unit || "盒",
    dose_unit: medicine.dose_unit || "粒",
    specifications: String(medicine.specifications ?? 0),
    price: String(medicine.price ?? 0),
    min_stock_warn: String(medicine.min_stock_warn ?? 2),
    daily_dose: medicine.daily_dose ? String(medicine.daily_dose) : "",
    usage: medicine.usage || "",
    photo: medicine.photo || "",
    state: medicine.state ?? 1,
    remark: medicine.remark || "",
  });
  return true;
}

function errorMessage(error: unknown) {
  const payload = error as { msg?: string; message?: string };
  return payload?.msg || payload?.message || "保存失败，请稍后重试";
}

function validateSpecifications(value: string) {
  return /^\d+$/.test(value) && Number(value) >= 0;
}

function validatePrice(value: string) {
  return /^\d+(\.\d{1,2})?$/.test(value) && Number(value) <= 999999.99;
}

async function initialize() {
  try {
    await medicineStore.loadDiseases();
    if (isEdit.value && !fillForm()) {
      await medicineStore.loadMedicines();
      if (!fillForm()) {
        showFailToast("没有找到该药品");
        await router.replace({ name: "Medicines" });
      }
    }
  } catch (error) {
    showFailToast(errorMessage(error));
  }
}

function createPayload(state = form.state): MedicinePayload {
  const userId = userStore.user?.ID;
  if (!userId) throw new Error("未获取到当前用户 ID，请重新登录");

  return {
    id: isEdit.value ? medicineId.value : undefined,
    user_id: userId,
    name: form.name.trim(),
    manufacturer: form.manufacturer.trim(),
    disease_ids: form.disease_ids,
    stock: Number(form.stock || 0),
    package_unit: form.package_unit,
    dose_unit: form.dose_unit,
    specifications: Number(form.specifications || 0),
    price: Number(form.price || 0),
    min_stock_warn: Number(form.min_stock_warn || 0),
    daily_dose: Number(form.daily_dose || 0),
    usage: form.usage.trim(),
    photo: form.photo.trim(),
    state,
    remark: form.remark.trim(),
  };
}

async function submit() {
  submitting.value = true;
  try {
    const payload = createPayload();
    if (isEdit.value) await updateMedicine(payload);
    else await addMedicine(payload);

    showSuccessToast(isEdit.value ? "更新成功" : "添加成功");
    await router.replace({ name: "Medicines" });
  } catch (error) {
    showFailToast(errorMessage(error));
  } finally {
    submitting.value = false;
  }
}

async function toggleMedicineState() {
  if (!isEdit.value) return;

  stateSubmitting.value = true;
  const nextState = form.state === 1 ? 0 : 1;
  try {
    await updateMedicine(createPayload(nextState));
    form.state = nextState;
    const medicine = medicineStore.findMedicine(medicineId.value);
    if (medicine) medicine.state = nextState;
    showSuccessToast({
      message: nextState === 1 ? "药品已上架" : "药品已下架",
      duration: 3000,
    });
    await new Promise((resolve) => setTimeout(resolve, 3000));
    await router.replace({ name: "Medicines" });
  } catch (error) {
    showFailToast(errorMessage(error));
  } finally {
    stateSubmitting.value = false;
  }
}

onMounted(initialize);
</script>

<template>
  <main class="form-page">
    <section class="form-card">
      <header class="form-heading">
        <div>
          <p>MEDICINE RECORD</p>
          <h1>{{ isEdit ? "编辑药品" : "新增药品" }}</h1>
          <span>记录库存与用法，低于预警值时会在列表中提醒。</span>
        </div>
        <van-image
          v-if="form.photo"
          round
          width="64"
          height="64"
          fit="cover"
          :src="form.photo"
          alt="药品图片预览"
        />
        <div v-else class="photo-preview"><van-icon name="medicines-o" /></div>
      </header>

      <van-form class="medicine-form" @submit="submit">
        <div class="field-grid">
          <van-field
            v-model.trim="form.name"
            name="name"
            label="药品名称"
            placeholder="如：布洛芬缓释胶囊"
            :rules="[{ required: true, message: '请输入药品名称' }]"
          />
          <van-field
            v-model.trim="form.manufacturer"
            name="manufacturer"
            label="生产企业"
            placeholder="选填"
          />
          <div class="unit-row">
            <van-field name="package_unit" label="包装单位">
              <template #input>
                <select
                  v-model="form.package_unit"
                  class="unit-select"
                  aria-label="包装单位"
                >
                  <option
                    v-for="unit in medicineUnits"
                    :key="unit"
                    :value="unit"
                  >
                    {{ unit }}
                  </option>
                </select>
              </template>
            </van-field>
            <van-field name="dose_unit" label="剂量单位">
              <template #input>
                <select
                  v-model="form.dose_unit"
                  class="unit-select"
                  aria-label="剂量单位"
                >
                  <option
                    v-for="unit in medicineUnits"
                    :key="unit"
                    :value="unit"
                  >
                    {{ unit }}
                  </option>
                </select>
              </template>
            </van-field>
          </div>
          <van-field
            v-model="form.specifications"
            name="specifications"
            label="药品规格"
            type="digit"
            placeholder="如：12"
            :rules="[
              {
                validator: validateSpecifications,
                message: '请输入非负整数规格',
              },
            ]"
          >
            <template #button>
              <span class="field-unit">{{ form.dose_unit }}/{{ form.package_unit }}</span>
            </template>
          </van-field>
          <van-field
            v-model="form.price"
            name="price"
            label="药品单价"
            type="number"
            placeholder="如：29.90"
            :min="0"
            step="0.01"
            :rules="[
              {
                validator: validatePrice,
                message: '请输入非负金额，最多两位小数',
              },
            ]"
          >
            <template #button>
              <span class="field-unit">元</span>
            </template>
          </van-field>

          <van-field
            v-model="form.stock"
            name="stock"
            label="当前库存"
            type="number"
            placeholder="0"
            :min="0"
            step="0.01"
          >
            <template #button>
              <span class="field-unit">{{ form.dose_unit }}</span>
            </template>
          </van-field>
          <van-field
            v-model="form.min_stock_warn"
            name="min_stock_warn"
            label="预警库存"
            type="number"
            placeholder="2"
            :min="0"
            step="0.01"
          >
            <template #button>
              <span class="field-unit">{{ form.dose_unit }}</span>
            </template>
          </van-field>
          <van-field
            v-model="form.daily_dose"
            name="daily_dose"
            label="每日剂量"
            type="number"
            placeholder="0"
            :min="0"
            step="0.01"
          >
            <template #button>
              <span class="field-unit">{{ form.dose_unit }}</span>
            </template>
          </van-field>
          <van-field
            v-model.trim="form.photo"
            name="photo"
            label="图片地址"
            type="url"
            placeholder="https://"
          />
          <van-field
            v-model.trim="form.usage"
            name="usage"
            label="用法用量"
            placeholder="如：饭后口服，一日 2 次"
          />
          <van-field
            class="full-field disease-field"
            name="disease_ids"
            label="适用疾病"
          >
            <template #input>
              <van-checkbox-group
                v-if="diseases.length"
                v-model="form.disease_ids"
                direction="horizontal"
              >
                <van-checkbox
                  v-for="disease in diseases"
                  :key="disease.id"
                  :name="disease.id"
                  shape="square"
                  >{{ disease.name }}</van-checkbox
                >
              </van-checkbox-group>
              <span v-else class="empty-hint">暂无疾病标签</span>
            </template>
          </van-field>
          <van-field
            v-model.trim="form.remark"
            class="full-field"
            name="remark"
            label="备注"
            type="textarea"
            rows="3"
            autosize
            maxlength="255"
            show-word-limit
            placeholder="其他需要留意的信息"
          />
        </div>

        <div v-if="isEdit" class="state-action">
          <div>
            <strong>药品状态</strong>
            <span>{{
              form.state === 1 ? "当前为上架状态" : "当前为下架状态"
            }}</span>
          </div>
          <van-button
            round
            plain
            :type="form.state === 1 ? 'danger' : 'primary'"
            native-type="button"
            :loading="stateSubmitting"
            :disabled="submitting"
            @click="toggleMedicineState"
          >
            {{ form.state === 1 ? "下架药品" : "上架药品" }}
          </van-button>
        </div>

        <div class="form-actions">
          <van-button
            block
            round
            plain
            type="primary"
            native-type="button"
            @click="router.back()"
            >取消</van-button
          >
          <van-button
            block
            round
            type="primary"
            native-type="submit"
            :loading="submitting"
            loading-text="正在保存…"
            >{{ isEdit ? "保存修改" : "添加药品" }}</van-button
          >
        </div>
      </van-form>
    </section>
  </main>
</template>

<style lang="less" scoped>
.form-page {
  min-height: 100%;
  padding: 22px 18px 34px;
  color: var(--app-text);
  background: var(--app-page-background);
}
.form-card {
  max-width: 860px;
  margin: 0 auto;
  padding: 30px;
  border: 1px solid var(--app-border);
  border-radius: 24px;
  background: var(--app-surface);
  box-shadow: 0 16px 42px var(--app-shadow);
}
.form-heading {
  margin-bottom: 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}
.form-heading p,
.form-heading h1,
.form-heading span {
  margin: 0;
}
.form-heading p {
  margin-bottom: 6px;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  color: #0b8a7b;
}
.form-heading h1 {
  margin-bottom: 8px;
  font-size: 28px;
}
.form-heading span {
  font-size: 14px;
  color: var(--app-text-muted);
}
.photo-preview {
  width: 64px;
  height: 64px;
  flex: 0 0 64px;
  display: grid;
  place-items: center;
  border-radius: 18px;
  font-size: 30px;
  color: #087568;
  background: var(--app-accent-soft);
}
.medicine-form {
  --van-field-label-color: var(--app-text-strong);
  --van-field-input-text-color: var(--app-text);
  --van-field-label-width: 5.8em;
}
.field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 14px;
}
.unit-row {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 14px;
}
.unit-row :deep(.van-field__label) {
  width: 4.5em;
}
.medicine-form :deep(.van-cell) {
  margin-bottom: 14px;
  padding: 11px 14px;
  align-items: center;
  border: 1px solid var(--app-border);
  border-radius: 13px;
  background: var(--app-surface-muted);
}
.medicine-form :deep(.van-cell::after) {
  display: none;
}
.medicine-form :deep(.van-field__label),
.medicine-form :deep(.van-field__control),
.medicine-form :deep(.van-checkbox__label) {
  font-size: 14px;
}
.medicine-form :deep(.van-field--error) {
  border-color: #df6c6c;
}
.unit-select {
  width: 100%;
  border: 0;
  outline: 0;
  font-size: 14px;
  color: var(--app-text);
  background: transparent;
}
.field-unit {
  color: var(--app-text-muted);
  font-size: 14px;
}
.full-field {
  grid-column: 1 / -1;
}
.disease-field :deep(.van-checkbox-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
}
.empty-hint {
  font-size: 14px;
  color: var(--app-text-muted);
}
.state-action {
  margin-top: 4px;
  padding: 16px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border: 1px solid var(--app-border);
  border-radius: 14px;
  background: var(--app-surface-muted);
}
.state-action strong,
.state-action span {
  display: block;
}
.state-action strong {
  margin-bottom: 5px;
  font-size: 15px;
  color: var(--app-text-strong);
}
.state-action span {
  font-size: 13px;
  color: var(--app-text-muted);
}
.state-action :deep(.van-button) {
  min-width: 106px;
}
.form-actions {
  margin-top: 10px;
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 12px;
}
.form-actions :deep(.van-button) {
  height: 48px;
  font-size: 15px;
  font-weight: 700;
}
.form-actions :deep(.van-button--normal:not(.van-button--plain)) {
  border: 0;
  background: linear-gradient(115deg, #087f72, #0a9683);
  box-shadow: 0 10px 24px rgba(8, 127, 114, 0.2);
}

@media (max-width: 680px) {
  .form-page {
    padding: 14px 12px 26px;
  }
  .form-card {
    padding: 22px 14px;
    border-radius: 20px;
  }
  .form-heading {
    padding: 0 4px;
  }
  .form-heading h1 {
    font-size: 24px;
  }
  .field-grid {
    grid-template-columns: 1fr;
  }
  .full-field {
    grid-column: auto;
  }
  .unit-row {
    grid-column: auto;
    gap: 0 8px;
  }
  .unit-row :deep(.van-cell) {
    padding-inline: 10px;
  }
  .form-actions {
    grid-template-columns: 1fr 1.35fr;
  }
}
</style>
