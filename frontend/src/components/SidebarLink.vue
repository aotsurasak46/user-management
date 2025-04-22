<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  label: { type: String, required: true },
  routeName: { type: String, required: true },
  icon: { type: String, required: true },
})
</script>

<template>
  <RouterLink
    :to="{ name: props.routeName }"
    exact-active-class="bg-gray-200 text-black lg:border-l-8"  
    class="w-full h-10 flex items-center justify-center rounded-md mt-2  text-black hover:bg-gray-200 hover:text-[#333333] "
  >
    <i :class="icon"></i>
    <span class="pl-3 hidden lg:inline">{{ props.label }}</span>
  </RouterLink>
</template>
