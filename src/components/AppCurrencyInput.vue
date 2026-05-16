<template>
  <div class="relative">
    <span class="absolute z-1 left-3 top-1/2 -translate-y-1/2 text-sm text-base-content/40 font-medium">
      Rp
    </span>
    <input
      type="text"
      :value="displayValue"
      @input="handleInput"
      :placeholder="placeholder"
      :disabled="disabled"
      class="input input-bordered w-full pl-10"
      :class="{ 'opacity-40 cursor-not-allowed': disabled }"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [Number, String],
    default: ''
  },
  placeholder: {
    type: String,
    default: '0'
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const displayValue = computed(() => {
  if (props.modelValue === '' || props.modelValue === null || props.modelValue === undefined) {
    return ''
  }
  
  // Format with dot as thousand separator (id-ID style)
  return props.modelValue.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')
})

const handleInput = (e) => {
  // Remove all non-digits
  const rawValue = e.target.value.replace(/\D/g, '')
  
  if (rawValue === '') {
    emit('update:modelValue', '')
    return
  }
  
  // Parse to integer
  const numValue = parseInt(rawValue, 10)
  
  // Emit the raw number
  emit('update:modelValue', numValue)
}
</script>
