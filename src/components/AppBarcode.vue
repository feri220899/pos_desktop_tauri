<template>
    <svg ref="svgRef"></svg>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import JsBarcode from 'jsbarcode'

const props = defineProps({
    value:       { type: String,  required: true },
    format:      { type: String,  default: 'CODE128' },
    width:       { type: Number,  default: 1.5 },
    height:      { type: Number,  default: 32 },
    displayValue:{ type: Boolean, default: false },
    fontSize:    { type: Number,  default: 10 },
})

const svgRef = ref(null)

function render() {
    if (!svgRef.value || !props.value) return
    try {
        JsBarcode(svgRef.value, props.value, {
            format:       props.format,
            width:        props.width,
            height:       props.height,
            displayValue: props.displayValue,
            fontSize:     props.fontSize,
            margin:       0,
            background:   'transparent',
            lineColor:    'currentColor',
        })
    } catch {
        svgRef.value.innerHTML = ''
    }
}

onMounted(render)
watch(() => props.value, render)
</script>
