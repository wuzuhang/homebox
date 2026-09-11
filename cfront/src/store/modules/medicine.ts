import type { Disease, Medicine } from '@/api/medicine'
import { getDiseases, getMedicines } from '@/api/medicine'

export const useMedicineStore = defineStore('medicine', () => {
  const medicines = ref<Medicine[]>([])
  const diseases = ref<Disease[]>([])
  const loading = ref(false)

  const lowStockMedicines = computed(() => medicines.value.filter(item => item.stock <= item.min_stock_warn))

  async function loadMedicines() {
    loading.value = true
    try {
      medicines.value = (await getMedicines()) || []
    }
    finally {
      loading.value = false
    }
  }

  async function loadDiseases() {
    if (diseases.value.length)
      return
    diseases.value = (await getDiseases()) || []
  }

  function findMedicine(id: number) {
    return medicines.value.find(item => item.ID === id)
  }

  function diseaseName(id: number) {
    return diseases.value.find(item => item.id === id)?.name || `标签 ${id}`
  }

  return {
    medicines,
    diseases,
    loading,
    lowStockMedicines,
    loadMedicines,
    loadDiseases,
    findMedicine,
    diseaseName,
  }
})
